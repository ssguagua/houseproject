package controllers

import (
	_"encoding/json"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"houseproject/models"
	_"syscall"
)

type AreaController struct {
	beego.Controller
}

func (c *AreaController)RetData(resp map[string]interface{}){
	c.Data["json"]=resp
	c.ServeJSON()
}
func (c *AreaController) GetArea(){
	beego.Info("connected success")
	resp := make(map[string]interface{}) //借助map打包成JSON
	defer c.RetData(resp)             //将数据打包成JSON返回给前端
	//1.从session拿数据

	//2.从mysql数据库中拿到area数据
	areas:=[]models.Area{}
	o:=orm.NewOrm()
	num,err:=o.QueryTable("area").All(&areas)
	if err!=nil{
		resp["errno"]=models.RECODE_DBERR
		resp["errmsg"]=models.RecodeText(models.RECODE_DBERR)
		beego.Info("query fail,resp=",resp)
		return
	}
	if num==0{
		//查到0条数据
		resp["errno"]=models.RECODE_NODATA
		resp["errmsg"]=models.RecodeText(models.RECODE_NODATA)
		beego.Info("query fail,resp=",resp)
		return
	}
	resp["errno"]=models.RECODE_OK
	resp["errmsg"]=models.RecodeText(models.RECODE_OK)
	resp["data"]=areas
	beego.Info("query success,resp=",resp,"num=",num)

	//3.将数据打包成JSON返回给前端
}