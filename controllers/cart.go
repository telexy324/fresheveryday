package controllers

import (
	"fresheveryday/models"
	"fresheveryday/utils"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"github.com/gomodule/redigo/redis"
	"strconv"
)

type CartController struct {
	beego.Controller
}

func (this *CartController) HandleAddCart()  {
	resp := make(map[string]interface{})
	defer this.ServeJSON()

	userName:=utils.GetUserName(&this.Controller)
	//if userName == "" {
	//	resp["code"] = 1
	//	resp["errmsg"] = "请您先登陆"
	//	this.Data["json"] = resp
	//	return
	//}


	fdfshost := beego.AppConfig.String("fdfshost")
	this.Data["fdfshost"] = fdfshost


	goodsId,errgoods:=this.GetInt("goodsId")
	count,errcount:=this.GetInt("count")


	if errgoods != nil || errcount != nil{
		resp["code"] = 1
		resp["errmsg"] = "ajax数据传输错误"
		this.Data["json"] = resp
		return
	}

	o:=orm.NewOrm()
	var user models.User
	user.Name=userName
	o.Read(&user,"Name")

	conn,err:=redis.Dial("tcp","192.168.248.129:6379")
	if err!=nil {
		resp["code"] = 1
		resp["errmsg"] = "获取数据失败"
		this.Data["json"] = resp
		return
	}

	repl,err:=conn.Do("hget","cart_"+strconv.Itoa(user.Id),goodsId)
	preCount,_:=redis.Int(repl,err)
	if preCount+count<0 {
		resp["code"] = 1
		resp["errmsg"] = "添加购物车失败"
		this.Data["json"] = resp
		return
	}
	conn.Do("hset","cart_"+strconv.Itoa(user.Id),goodsId,count+preCount)

	resp["code"] = 5
	this.Data["json"] = resp
	this.ServeJSON()
}

func (this *CartController) ShowCart()  {
	userName:=utils.GetUserName(&this.Controller)
	fdfshost := beego.AppConfig.String("fdfshost")
	this.Data["fdfshost"] = fdfshost


	o:=orm.NewOrm()
	var user models.User
	user.Name=userName
	o.Read(&user,"Name")

	conn,err:=redis.Dial("tcp","192.168.248.129:6379")
	if err!=nil {
		beego.Error("redis连接失败")
	}
	resp,err:=conn.Do("hgetall","cart_"+strconv.Itoa(user.Id))
	cartMap,_:=redis.IntMap(resp,err)
	//totalCount:=len(cartMap)
	goods:=make([]map[string]interface{},0)
	for goodsId,count:=range cartMap {
		var goodsSKU models.GoodsSKU
		goodsSKU.Id,_=strconv.Atoi(goodsId)
		o.Read(&goodsSKU)
		temp:=make(map[string]interface{})
		temp["goodsSKU"] = goodsSKU
		temp["count"] = count
		temp["perPrice"] = goodsSKU.Price*count
		goods=append(goods, temp)
	}

	this.Data["goods"] = goods

	this.Data["title"] = "购物车"
	this.Layout = "user_center_layout.html"
	this.TplName = "cart.html"
}

//func (this *CartController) DeleteCart()  {
//	userName:=utils.GetUserName(&this.Controller)
//	goodsId,_:=this.GetInt("goodsId")
//	o:=orm.NewOrm()
//	var user models.User
//	user.Name=userName
//	o.Read(&user,"Name")
//
//	conn,err:=redis.Dial("tcp","192.168.248.129:6379")
//	if err!=nil {
//		beego.Error("redis连接失败")
//	}
//	conn.Do("hdel","cart_"+strconv.Itoa(user.Id),goodsId)
//	this.Redirect("/user/showCart",302)
//}

func (this *CartController) DeleteCart()  {
	resp := make(map[string]interface{})
	defer this.ServeJSON()

	userName:=utils.GetUserName(&this.Controller)
	goodsId,_:=this.GetInt("goodsId")
	o:=orm.NewOrm()
	var user models.User
	user.Name=userName
	o.Read(&user,"Name")

	conn,err:=redis.Dial("tcp","192.168.248.129:6379")
	if err!=nil {
		beego.Error("redis连接失败")
	}
	_,err=conn.Do("hdel","cart_"+strconv.Itoa(user.Id),goodsId)
	if err!=nil {
		resp["code"] = 1
		resp["errmsg"] = "删除失败"
		return
	}
	resp["code"] = 5
	this.Data["json"] = resp
	this.ServeJSON()
}