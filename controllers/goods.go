package controllers

import (
	"fresheveryday/models"
	"fresheveryday/utils"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"github.com/gomodule/redigo/redis"
	"math"
	"strconv"
)

type GoodsController struct {
	beego.Controller
}

func (this *GoodsController) ShowGoods() {
	userName:=utils.GetUserName(&this.Controller)

	fdfshost := beego.AppConfig.String("fdfshost")
	this.Data["fdfshost"] = fdfshost

	utils.GetCart(&this.Controller,userName)

	o := orm.NewOrm()
	//var goodsTypes []*models.GoodsType
	//o.QueryTable("GoodsType").All(&goodsTypes)
	//this.Data["goodsTypes"] = goodsTypes

	goodsTypes := utils.GetType(&this.Controller)

	var indexGoodsBanners []*models.IndexGoodsBanner
	o.QueryTable("IndexGoodsBanner").All(&indexGoodsBanners)
	this.Data["indexGoodsBanners"] = indexGoodsBanners

	var indexPromotionBanners []*models.IndexPromotionBanner
	o.QueryTable("IndexPromotionBanner").All(&indexPromotionBanners)
	this.Data["indexPromotionBanners"] = indexPromotionBanners

	var goodsSKUs = make([]map[string]interface{}, len(goodsTypes))

	for typeIndex, _ := range goodsSKUs {
		temp := make(map[string]interface{})
		temp["type"] = goodsTypes[typeIndex]
		goodsSKUs[typeIndex] = temp
	}
	for _, goodsMap := range goodsSKUs {
		var goodsImage []models.IndexTypeGoodsBanner
		var goodsText []models.IndexTypeGoodsBanner
		o.QueryTable("IndexTypeGoodsBanner").RelatedSel("GoodsType", "GoodsSKU").Filter("GoodsType", goodsMap["type"]).Filter("DisplayType", 0).All(&goodsText)
		o.QueryTable("IndexTypeGoodsBanner").RelatedSel("GoodsType", "GoodsSKU").Filter("GoodsType", goodsMap["type"]).Filter("DisplayType", 1).All(&goodsImage)
		goodsMap["goodsImage"] = goodsImage
		goodsMap["goodsText"] = goodsText

	}
	this.Data["goodsSKUs"] = goodsSKUs

	this.Data["title"] = "天天生鲜-首页"
	this.Layout = "topBar.html"
	this.TplName = "index.html"
}

func (this *GoodsController) ShowGoodsDetails() {
	userName := utils.GetUserName(&this.Controller)

	fdfshost := beego.AppConfig.String("fdfshost")
	this.Data["fdfshost"] = fdfshost

	utils.GetCart(&this.Controller,userName)

	id, err := this.GetInt("id")
	if err != nil {
		beego.Error("未获得商品ID")
		this.Redirect("/", 302)
	}
	o := orm.NewOrm()
	utils.GetType(&this.Controller)

	var goodsSKU models.GoodsSKU
	o.QueryTable("GoodsSKU").RelatedSel("Goods").Filter("ID", id).One(&goodsSKU)
	this.Data["goodsSKU"] = goodsSKU

	goodsTypeId := goodsSKU.GoodsType.Id
	utils.GetNewGoods(&this.Controller, goodsTypeId)

	conn, err := redis.Dial("tcp", "192.168.248.129:6379")
	if err != nil {
		beego.Error("redis连接失败", err)
	}
	defer conn.Close()

	var user models.User
	if userName != "" {
		user.Name = userName
		o.Read(&user, "Name")
		conn.Do("lrem", "history_"+strconv.Itoa(user.Id), 0, id)
		conn.Do("lpush", "history_"+strconv.Itoa(user.Id), id)
	}

	this.Data["title"] = "天天生鲜-商品详情"
	this.Layout = "topBar.html"
	this.TplName = "detail.html"
}

func (this *GoodsController) ShowGoodsList() {
	userName:=utils.GetUserName(&this.Controller)

	fdfshost := beego.AppConfig.String("fdfshost")
	this.Data["fdfshost"] = fdfshost

	utils.GetCart(&this.Controller,userName)

	id, err := this.GetInt("id")

	if err != nil {
		beego.Error("未获得商品ID")
		this.Redirect("/", 302)
	}
	this.Data["id"] = id

	o := orm.NewOrm()
	utils.GetType(&this.Controller)

	goodsPerPage := 1
	count, _ := o.QueryTable("GoodsSKU").RelatedSel("GoodsType").Filter("GoodsType__Id", id).Count()
	pagesRaw := float64(count) / float64(goodsPerPage)
	pages := int64(math.Ceil(pagesRaw))

	var pageIndex int
	pageIndex, err = this.GetInt("pageIndex")
	if err != nil {
		pageIndex = 1
	}
	start := int64(goodsPerPage * (pageIndex - 1))

	pageDisp := 3
	//var pageDisplay []int
	//if int(pages)<pageDisp {
	//	pageDisplay = make([]int,0)
	//	for i := 1;i<=int(pages);i++ {
	//		pageDisplay=append(pageDisplay, i)
	//	}
	//} else if pageIndex < pageDisp+1/2 {
	//	pageDisplay = make([]int,0)
	//	for i := 1;i<=pageDisp;i++ {
	//		pageDisplay=append(pageDisplay, i)
	//	}
	//} else if pageIndex > int(pages)-(pageDisp+1)/2+1 {
	//	pageDisplay = make([]int,0)
	//	for i := int(pages)-pageDisp+1;i<=int(pages);i++ {
	//		pageDisplay=append(pageDisplay, i)
	//	}
	//} else {
	//	pageDisplay = make([]int,0)
	//	for i := pageIndex-(pageDisp+1)/2+1;i<=pageIndex+(pageDisp+1)/2-1;i++ {
	//		pageDisplay=append(pageDisplay, i)
	//	}
	//}
	pageDisplay := getPages(pageDisp, int(pages), pageIndex)
	this.Data["pageDisplay"] = pageDisplay

	var goodsSKUs []*models.GoodsSKU
	sort := this.GetString("sort")

	if sort == "price" {
		o.QueryTable("GoodsSKU").RelatedSel("GoodsType").Filter("GoodsType__Id", id).OrderBy("Price").Limit(goodsPerPage, start).All(&goodsSKUs)
	} else if sort == "sale" {
		o.QueryTable("GoodsSKU").RelatedSel("GoodsType").Filter("GoodsType__Id", id).OrderBy("Sales").Limit(goodsPerPage, start).All(&goodsSKUs)
	} else {
		o.QueryTable("GoodsSKU").RelatedSel("GoodsType").Filter("GoodsType__Id", id).Limit(goodsPerPage, start).All(&goodsSKUs)
	}

	this.Data["sort"] = sort
	//o.QueryTable("GoodsSKU").RelatedSel("GoodsType").Filter("GoodsType__Id",id).Limit(goodsPerPage,start).All(&goodsSKUs)
	this.Data["goodsSKU"] = goodsSKUs
	this.Data["pages"] = pages
	this.Data["pageIndex"] = pageIndex

	utils.GetNewGoods(&this.Controller, id)

	this.Data["title"] = "天天生鲜-商品列表"
	this.Layout = "topBar.html"
	this.TplName = "list.html"
}

func (this *GoodsController) HandleSearchGoods() {
	userName:=utils.GetUserName(&this.Controller)

	fdfshost := beego.AppConfig.String("fdfshost")
	this.Data["fdfshost"] = fdfshost

	utils.GetCart(&this.Controller,userName)

	o := orm.NewOrm()
	var goodsTypes []*models.GoodsType
	o.QueryTable("GoodsType").All(&goodsTypes)
	this.Data["goodsTypes"] = goodsTypes

	utils.GetNewGoods(&this.Controller, -1)

	searchName := this.GetString("searchName")
	this.Data["searchName"] = searchName
	pageIndex, err := this.GetInt("pageIndex")
	if err != nil {
		pageIndex = 1
	}
	this.Data["pageIndex"] = pageIndex
	goodsPerPage := 3
	pageDisp := 3
	var goodsSKUs []*models.GoodsSKU
	var count int64
	var pageDisplay []int
	var start int64
	var pages int64
	sort := this.GetString("sort")
	this.Data["sort"] = sort
	if searchName == "" {
		count, _ = o.QueryTable("goodsSKU").Count()
		pagesRaw := float64(count) / float64(goodsPerPage)
		pages = int64(math.Ceil(pagesRaw))
		pageDisplay = getPages(pageDisp, int(pages), pageIndex)
		start = int64(goodsPerPage * (pageIndex - 1))
		if sort == "price" {
			o.QueryTable("GoodsSKU").OrderBy("Price").Limit(goodsPerPage, start).All(&goodsSKUs)
		} else if sort == "sale" {
			o.QueryTable("GoodsSKU").OrderBy("Sales").Limit(goodsPerPage, start).All(&goodsSKUs)
		} else {
			o.QueryTable("GoodsSKU").Limit(goodsPerPage, start).All(&goodsSKUs)
		}
	} else {
		count, _ = o.QueryTable("goodsSKU").Filter("Name__icontains", searchName).Count()
		pagesRaw := float64(count) / float64(goodsPerPage)
		pages = int64(math.Ceil(pagesRaw))
		pageDisplay = getPages(pageDisp, int(pages), pageIndex)
		start = int64(goodsPerPage * (pageIndex - 1))
		if sort == "price" {
			o.QueryTable("GoodsSKU").Filter("Name__icontains", searchName).OrderBy("Price").Limit(goodsPerPage, start).All(&goodsSKUs)
		} else if sort == "sale" {
			o.QueryTable("GoodsSKU").Filter("Name__icontains", searchName).OrderBy("Sales").Limit(goodsPerPage, start).All(&goodsSKUs)
		} else {
			o.QueryTable("GoodsSKU").Filter("Name__icontains", searchName).Limit(goodsPerPage, start).All(&goodsSKUs)
		}
	}

	this.Data["pages"] = pages
	this.Data["goodsSKUs"] = goodsSKUs
	this.Data["pageDisplay"] = pageDisplay
	this.Data["title"] = "天天生鲜-查询结果"
	this.Layout = "topBar.html"
	this.TplName = "search.html"
}

func getPages(pageDisp, pages, pageIndex int) []int {
	var pageDisplay []int
	if pages < pageDisp {
		pageDisplay = make([]int, 0)
		for i := 1; i <= pages; i++ {
			pageDisplay = append(pageDisplay, i)
		}
	} else if pageIndex < pageDisp+1/2 {
		pageDisplay = make([]int, 0)
		for i := 1; i <= pageDisp; i++ {
			pageDisplay = append(pageDisplay, i)
		}
	} else if pageIndex > pages-(pageDisp+1)/2+1 {
		pageDisplay = make([]int, 0)
		for i := pages - pageDisp + 1; i <= pages; i++ {
			pageDisplay = append(pageDisplay, i)
		}
	} else {
		pageDisplay = make([]int, 0)
		for i := pageIndex - (pageDisp+1)/2 + 1; i <= pageIndex+(pageDisp+1)/2-1; i++ {
			pageDisplay = append(pageDisplay, i)
		}
	}
	return pageDisplay
}
