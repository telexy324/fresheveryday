package utils

import "github.com/astaxie/beego"

func GetUserName(this *beego.Controller) string {
	//userName:=this.GetSession("userName").(string)
	//this.Data["userName"] = userName
	//return userName

	userNameRaw:=this.GetSession("userName")
	if userNameRaw !=nil {
		userName := userNameRaw.(string)
		this.Data["userName"] = userName
		return userName
	} else {
		this.Data["userName"] = ""
		return ""
	}
}
