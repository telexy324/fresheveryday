package main

import (
	_ "fresheveryday/routers"
	"github.com/astaxie/beego"
	_ "fresheveryday/models"
)

func main() {
	beego.AddFuncMap("ShowPrePage",PrePageIndex)
	beego.AddFuncMap("NextPage",NextPageIndex)
	beego.Run()
}

func PrePageIndex(pageIndex int) int {
	prePage := pageIndex - 1
	if prePage < 1 {
		prePage = 1
	}
	return prePage
}

func NextPageIndex(pageIndex int, pages int64) int {

	if pageIndex+1 > int(pages) {
		return pageIndex
	}
	return pageIndex + 1
}
