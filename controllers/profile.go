package controllers

import (
	"fresheveryday/models"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"github.com/gomodule/redigo/redis"
	"regexp"
	"fresheveryday/utils"
	"strconv"
)

type ProfileController struct {
	beego.Controller
}

func (this *ProfileController) ShowUserCenter()  {
	//userName:=this.GetSession("userName").(string)
	userName:=utils.GetUserName(&this.Controller)

	fdfshost:=beego.AppConfig.String("fdfshost")
	this.Data["fdfshost"] = fdfshost

	o:=orm.NewOrm()
	//var user models.User
	//user.Name=userName
	//o.Read(&user,"Name")
	var address models.Address
	address.Isdefault=true
	qs := o.QueryTable("Address")
	qs.RelatedSel("User").Filter("IsDefault",address.Isdefault).Filter("User__Name",userName).One(&address)
	//if err !=nil {
	//	this.Data["errmsg"] = "用户名错误，请重新登陆！"
	//	this.TplName = "login.html"
	//	return
	//}
	utils.GetUserName(&this.Controller)
	this.Data["addr"] = address.Addr
	this.Data["phone"] = address.Phone

	conn,err:=redis.Dial("tcp","192.168.248.129:6379")
	if err!=nil {
		beego.Error("redis连接失败",err)
	}
	defer conn.Close()

	var user models.User
	user.Name = userName
	o.Read(&user,"Name")
	repl,err:=conn.Do("lrange","history_"+strconv.Itoa(user.Id),0,4)
	ids,err:=redis.Ints(repl,err)
	if err!=nil {
		beego.Error("获取最近浏览商品失败")
		this.Layout = "user_center_layout.html"
		this.TplName = "user_center_info.html"
		return
	}
	var goodsSKUs []*models.GoodsSKU
	for _,id := range ids {
		var goodsSKU models.GoodsSKU
		goodsSKU.Id=id
		o.Read(&goodsSKU)
		goodsSKUs=append(goodsSKUs, &goodsSKU)
	}
	this.Data["goodsSKUs"] = goodsSKUs
	this.Data["title"] = "用户中心"
	this.Layout = "user_center_layout.html"
	this.TplName = "user_center_info.html"
}

func (this *ProfileController) ShowUserAddress()  {
	//userName:=this.GetSession("userName").(string)
	//this.Data["userName"] =userName
	userName:=utils.GetUserName(&this.Controller)
	o:=orm.NewOrm()
	var user models.User
	user.Name=userName
	o.Read(&user,"Name")

	var address models.Address
	address.Isdefault=true
	qs := o.QueryTable("Address")
	qs.RelatedSel("User").Filter("IsDefault",address.Isdefault).Filter("User__Name",userName).One(&address)
	this.Data["addr"] = address.Addr
	this.Data["receiver"] = address.Receiver
	this.Data["phone"] = address.Phone

	this.Data["title"] = "用户中心"
	this.Layout = "user_center_layout.html"
	this.TplName = "user_center_site.html"
}

func (this *ProfileController) HandleModifyUserAddress()  {
	//userName:=this.GetSession("userName").(string)
	//this.Data["userName"] =userName
	userName:=utils.GetUserName(&this.Controller)

	receiver:=this.GetString("receiver")
	addr:=this.GetString("addr")
	zipcode:=this.GetString("zipcode")
	phone:=this.GetString("phone")

	if receiver == "" || receiver == "" || receiver == "" || receiver == "" {
		this.Redirect("/user/address_info?errmsg="+"地址填写有误",302)
		return
	}

	zipcodereg,_:=regexp.Compile(`\d{6}`)
	phonereg,_:=regexp.Compile(`1\d{10}`)
	zipres:=zipcodereg.FindString(zipcode)
	phres:=phonereg.FindString(phone)
	if zipres=="" {
		this.Redirect("/user/address_info?errmsg="+"邮编格式不正确",302)
		return
	}
	if phres=="" {
		this.Redirect("/user/address_info?errmsg="+"电话格式不正确",302)
		return
	}


	o:=orm.NewOrm()
	var user models.User
	user.Name=userName
	o.Read(&user,"Name")

	var addressold models.Address
	addressold.Isdefault=true
	err:=o.Read(&addressold,"Isdefault")
	if err!=nil {
		beego.Info("第一次添加地址")
	} else {
		addressold.Isdefault=false
		o.Update(&addressold)
	}

	var address models.Address
	address.Receiver=receiver
	address.Addr=addr
	address.Zipcode=zipcode
	address.Phone=phone
	address.Isdefault=true
	address.User=&user
	_,err=o.Insert(&address)
	if err!=nil {
		this.Redirect("/user/address_info?errmsg="+"插入地址失败",302)
		return
	}
	this.Redirect("/user/address_info",302)
}
