package controllers

import (
	"encoding/base64"
	"fresheveryday/models"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"github.com/astaxie/beego/utils"
	"regexp"
	"strconv"
)

type UserController struct {
	beego.Controller
}

func (this *UserController) ShowRegister() {
	this.TplName = "register.html"
}

func (this *UserController) HandleRegister() {
	userName := this.GetString("user_name")
	pwd := this.GetString("pwd")
	email := this.GetString("email")
	cpwd:= this.GetString("cpwd")
	if userName == "" || pwd == "" || email == "" {
		this.Redirect("/register?errmsg="+"用户名、密码、邮箱不能为空", 302)
		return
	}
	reg, _ := regexp.Compile("^[A-Za-z0-9\u4e00-\u9fa5]+@[a-zA-Z0-9_-]+(\\.[a-zA-Z0-9_-]+)+$")
	res:=reg.FindString(email)
	if res == "" {
		this.Redirect("/register?errmsg="+"邮箱格式不正确", 302)
		return
	}
	if pwd != cpwd {
		this.Redirect("/register?errmsg="+"密码不一致", 302)
		return
	}
	o:=orm.NewOrm()
	var user models.User
	user.Name=userName
	user.PassWord=pwd
	user.Email=email
	_,err:=o.Insert(&user)
	if err!=nil {
		this.Data["errmsg"] = "注册失败，请重新注册！"
		this.TplName = "register.html"
		return
	}

	emailConfig := `{"host":"smtp.sina.com","username":"gis324@sina.com","password":"19860324","port":25}`
	eml := utils.NewEMail(emailConfig)
	eml.From="gis324@sina.com"
	eml.To=[]string{"2370707273@qq.com"}
	eml.Subject="天天生鲜网站激活"
	eml.HTML=`<a href="http://192.168.248.129:8080/activate?id=`+strconv.Itoa(user.Id)+`">请点此处激活</a>`
	err = eml.Send()
	if err!=nil {
		beego.Error(err)
	}


	this.Ctx.WriteString("注册成功,请去注册邮箱激活账户！")
}

func (this *UserController) HandleActivation() {
	id,err:=this.GetInt("id")
	if err!=nil {
		this.Data["errmsg"] = "激活失败，请重新注册！"
		this.TplName = "register.html"
		return
	}

	o:=orm.NewOrm()
	var user models.User
	user.Id=id
	err = o.Read(&user)
	if err!=nil {
		this.Data["errmsg"] = "用户不存在，请先注册！"
		this.TplName = "register.html"
		return
	}
	user.Active=true
	_,err = o.Update(&user)
	if err!=nil {
		this.Data["errmsg"] = "激活失败，请重新注册！"
		this.TplName = "register.html"
		return
	}
	this.Redirect("/login",302)
}

func (this *UserController) ShowLogin()  {
	enc := this.Ctx.GetCookie("userName")
	userName,_:=base64.StdEncoding.DecodeString(enc)
	if string(userName) != "" {
		this.Data["userName"] = string(userName)
		this.Data["check"] = "checked"
	} else {
		this.Data["userName"] = ""
		this.Data["check"] = ""
	}
	this.TplName = "login.html"
}

func (this *UserController) HandleLogin()  {
	userName:=this.GetString("username")
	pwd:=this.GetString("pwd")

	if userName == "" || pwd == "" {
		this.Data["errmsg"] = "用户名和密码不能为空！"
		this.TplName = "login.html"
		return
	}

	o:=orm.NewOrm()
	var user models.User
	user.Name=userName
	err := o.Read(&user,"Name")
	if err !=nil {
		this.Data["errmsg"] = "用户名不存在，请先注册！"
		this.TplName = "login.html"
		return
	}
	if !user.Active {
		this.Data["errmsg"] = "用户未激活，请先使用邮箱激活！"
		this.TplName = "login.html"
		return
	}
	if user.PassWord != pwd {
		this.Data["errmsg"] = "密码错误！"
		this.TplName = "login.html"
		return
	}
	check:=this.GetString("check")
	if check == "on" {
		enc := base64.StdEncoding.EncodeToString([]byte(userName))
		this.Ctx.SetCookie("userName",enc,3600)
	} else {
		this.Ctx.SetCookie("userName","",-1)
	}
	this.SetSession("userName",userName)

	this.Redirect("/",302)
}

func (this *UserController) HandleLogout()  {
	this.DelSession("userName")
	this.Redirect("/login",302)
}