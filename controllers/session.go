package controllers

import (
	"encoding/json"
	_ "encoding/json"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"houseproject/models"
	_ "syscall"
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
func (c *SessionController) UserLogin(){
	resp:=make(map[string]interface{})
	defer c.RetData(resp)     //返回resp给前端

	//1、从前端得到用户信息
	json.Unmarshal(c.Ctx.Input.RequestBody,&resp)
	//2、判断是否合法
	if resp["mobile"]==nil || resp["password"]==nil{
		resp["errno"]=models.RECODE_DATAERR
		resp["errmsg"]=models.RecodeText(models.RECODE_DATAERR)
		return
	}
	//3、与数据库匹配、判断密码是否正确
	o:=orm.NewOrm()
	user:=models.User{Name:resp["mobile"].(string)}
	_,err:=o.QueryTable("user").Filter("name",user.Name).All(&user)
	//err:=o.QueryTable("user").Filter("name",user.Name).One(&user)
	if err!=nil{
		resp["errno"]=models.RECODE_DATAERR
		resp["errmsg"]=models.RecodeText(models.RECODE_DATAERR)
		return
	}
	if user.Password_hash!=resp["password"]{
		resp["errno"]=models.RECODE_DATAERR
		resp["errmsg"]=models.RecodeText(models.RECODE_DATAERR)
		return
	}
	//4、添加session
	c.SetSession("name",user.Name)
	c.SetSession("user_id",user.Id)
	c.SetSession("mobile",user.Mobile)
	//5、返回JSON数据给前端
	resp["errno"]=models.RECODE_OK
	resp["errmsg"]=models.RecodeText(models.RECODE_OK)
}
