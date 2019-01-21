package utils

import (
	"fresheveryday/models"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"github.com/gomodule/redigo/redis"
	"strconv"
)

func GetCart(this *beego.Controller,userName string)  {
	var cartCount int
	if userName=="" {
		cartCount=0
	} else {
		o:=orm.NewOrm()
		var user models.User
		user.Name=userName
		o.Read(&user,"Name")

		conn,err:=redis.Dial("tcp","192.168.248.129:6379")
		if err!=nil {
			beego.Error("redis连接失败")
		}
		resp,err:=conn.Do("hlen","cart_"+strconv.Itoa(user.Id))
		cartCount,_=redis.Int(resp,err)
	}
	this.Data["cartCount"] = cartCount
}
