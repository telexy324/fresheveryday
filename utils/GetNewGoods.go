package utils

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"fresheveryday/models"
)

func GetNewGoods(this *beego.Controller,id int)  {
	o:=orm.NewOrm()
	var newGoods []*models.GoodsSKU
	if id == -1 {
		o.QueryTable("GoodsSKU").OrderBy("Time").Limit(2,0).All(&newGoods)
	} else {
		o.QueryTable("GoodsSKU").RelatedSel("GoodsType").Filter("GoodsType__Id", id).OrderBy("Time").Limit(2, 0).All(&newGoods)
	}
	this.Data["newGoods"] = newGoods
}
