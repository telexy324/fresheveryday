package routers

import (
	"fresheveryday/controllers"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/context"
)

func init() {
	beego.InsertFilter("/user/*",beego.BeforeExec,funcFilter)
    //beego.Router("/", &controllers.MainController{})
	beego.Router("/", &controllers.GoodsController{},"get:ShowGoods")
	beego.Router("/goodsDetails", &controllers.GoodsController{},"get:ShowGoodsDetails")
	beego.Router("/goodsList",&controllers.GoodsController{},"get:ShowGoodsList")
	beego.Router("/searchGoods" ,&controllers.GoodsController{},"*:HandleSearchGoods")

    beego.Router("/register", &controllers.UserController{},"get:ShowRegister;post:HandleRegister")
	beego.Router("/activate", &controllers.UserController{},"get:HandleActivation")
	beego.Router("/login", &controllers.UserController{},"get:ShowLogin;post:HandleLogin")
	beego.Router("/logout", &controllers.UserController{},"get:HandleLogout")

	beego.Router("/user/center_info", &controllers.ProfileController{},"get:ShowUserCenter")
	beego.Router("/user/address_info", &controllers.ProfileController{},"get:ShowUserAddress;post:HandleModifyUserAddress")

	beego.Router("/user/addCart", &controllers.CartController{},"post:HandleAddCart")
	beego.Router("/user/showCart", &controllers.CartController{},"get:ShowCart")
	beego.Router("/user/delCart", &controllers.CartController{},"*:DeleteCart")

	beego.Router("/user/confirmOrder", &controllers.OrderController{},"*:ConfirmOrder")
	beego.Router("/user/dealOrder", &controllers.OrderController{},"*:DealOrder")
}

func funcFilter(ctx *context.Context)  {
	if userName := ctx.Input.Session("userName");userName==nil {
		ctx.Redirect(302,"/login")
	}
}