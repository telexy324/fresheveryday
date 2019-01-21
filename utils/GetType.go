package utils

import (
	"fresheveryday/models"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
)

func GetType(this *beego.Controller) (goodsTypes []*models.GoodsType) {
	o:=orm.NewOrm()
	o.QueryTable("GoodsType").All(&goodsTypes)
	this.Data["goodsTypes"] = goodsTypes
	return goodsTypes
}