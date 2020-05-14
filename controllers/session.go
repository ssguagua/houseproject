package controllers

import (
	_ "encoding/json"
	"github.com/astaxie/beego"
	_ "syscall"
	"houseproject/models"
)

type SessionController struct {
	beego.Controller
}

func (c *SessionController)RetData(resp map[string]interface{}){
	c.Data["json"]=resp
	c.ServeJSON()
}
func (c *SessionController) GetSessionData(){
	resp:=make(map[string]interface{})
	defer c.RetData(resp)
	user:=models.User{}

	resp["errno"]=models.RECODE_DATAERR
	resp["errmsg"]=models.RecodeText(models.RECODE_DATAERR)
	//先判断session，如果存在，就get session
	name:=c.GetSession("name")
	if name!=nil{
		user.Name=name.(string)
		resp["errno"]=models.RECODE_OK
		resp["errmsg"]=models.RecodeText(models.RECODE_OK)
		resp["data"]=user
	}
}
//删除session、用户退出登录
func (c *SessionController)DeleteSession(){
	resp:=make(map[string]interface{})
	defer c.RetData(resp)
	c.DelSession("name")   //没有返回值
	resp["errno"]=models.RECODE_OK
	resp["errmsg"]=models.RecodeText(models.RECODE_OK)
}

/*
增加session模块
1、需要在项目中设置启用session的标志main、config
2、先在注册中setsession,然后在首页getsession判断session是否有值
 */
