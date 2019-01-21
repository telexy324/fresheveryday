package controllers

import (
	"fresheveryday/models"
	"fresheveryday/utils"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"github.com/gomodule/redigo/redis"
	"strconv"
	"strings"
	"time"
)

type OrderController struct {
	beego.Controller
}

func (this *OrderController) ConfirmOrder()  {
	userName:=utils.GetUserName(&this.Controller)
	fdfshost := beego.AppConfig.String("fdfshost")
	this.Data["fdfshost"] = fdfshost

	o:=orm.NewOrm()
	var user models.User
	user.Name=userName
	o.Read(&user,"Name")

	ids:=this.GetStrings("id")
	if len(ids) == 0 {
		beego.Error("获取商品失败")
	}

	orderGoods := make([]map[string]interface{},0)
	conn,err:=redis.Dial("tcp","192.168.248.129:6379")
	if err != nil {
		beego.Error("连接redis失败")
	}
	totalCount:=0
	totalPrice:=0
	for _,id:=range ids {
		temp:=make(map[string]interface{})
		repl,err:=conn.Do("hget","cart_"+strconv.Itoa(user.Id),id)
		count,_:=redis.Int(repl,err)
		temp["goodsNum"] = count
		var goodsSKU models.GoodsSKU
		goodsSKU.Id,_=strconv.Atoi(id)
		o.Read(&goodsSKU)
		temp["goodsSKU"] = goodsSKU
		sumPrice := goodsSKU.Price * count
		temp["sumPrice"] = sumPrice

		totalCount += count
		totalPrice += sumPrice
		orderGoods = append(orderGoods ,temp)
	}
	this.Data["orderGoods"] = orderGoods
	var addresses []*models.Address
	o.QueryTable("Address").RelatedSel("User").Filter("User__Id",user.Id).All(&addresses)
	this.Data["addresses"] = addresses
	this.Data["totalPrice"]= totalPrice
	this.Data["totalCount"]= totalCount
	this.Data["deliverFee"]= 10
	this.Data["pricePlusDeliver"] = totalPrice+10
	this.Data["ids"] = ids

	this.Data["title"] = "提交订单"
	this.Layout="user_center_layout.html"
	this.TplName="place_order.html"
}

func (this *OrderController) DealOrder()  {
	userName:=utils.GetUserName(&this.Controller)
	addId,err1 :=this.GetInt("addId")
	payId,err2 :=this.GetInt("payId")
	//js获取页面数据都是以字符串类型获取
	goodsId:=this.GetString("goodsId")
	totalPrice,err3 := this.GetInt("totalPrice")
	totalCount,err4 :=this.GetInt("totalCount")

	//校验数据
	if err1 != nil || err2 != nil || err3 != nil ||err4 != nil || len(goodsId) == 0{
		beego.Error("获取数据失败")
	}
	ids :=strings.Split(goodsId[1:len(goodsId)-1]," ")

	//处理数据
	//向订单表和订单商品表插入数据
	o := orm.NewOrm()
	var order models.OrderInfo
	order.TransitPrice = 10
	order.TotalPrice = totalPrice
	order.TotalCount = totalCount
	order.PayMethod = payId

	//获取用户数据
	var user models.User
	user.Name = userName
	o.Read(&user,"Name")

	order.OrderId = time.Now().Format("20060102150405")+strconv.Itoa(user.Id)
	order.User = &user

	//获取地址信息
	var addr models.Address
	addr.Id = addId
	o.Read(&addr)
	order.Address = &addr

	//插入操作
	o.Insert(&order)
	//插入数据到订单商品表
	conn,_ :=redis.Dial("tcp","192.168.248.129:6379")
	for _,value := range ids{
		id,_ :=strconv.Atoi(value)
		//获取商品信息
		var goodsSku models.GoodsSKU
		goodsSku.Id = id
		o.Read(&goodsSku)
		//获取商品数量
		resp ,err :=conn.Do("hget","cart_"+strconv.Itoa(user.Id),id)
		count,_ :=redis.Int(resp,err)

		var orderGoods models.OrderGoods
		orderGoods.GoodsSKU = &goodsSku
		orderGoods.Price = goodsSku.Price * count
		orderGoods.OrderInfo = &order
		orderGoods.Count = count

		//插入操作
		o.Insert(&orderGoods)

	}


	//返回数据
	re := make(map[string]interface{})
	re["code"] = 5
	re["errmsg"] = "OK"
	this.Data["json"] = re
	this.ServeJSON()
}