package controllers

import (
	"encoding/json"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/cache"
	_ "github.com/astaxie/beego/cache/redis"
	"github.com/astaxie/beego/orm"
	"houseproject/models"
	_ "syscall"
	"time"
)

type AreaController struct {
	beego.Controller
}

func (c *AreaController)RetData(resp map[string]interface{}){
	c.Data["json"]=resp
	c.ServeJSON()
}

func (c *AreaController) GetArea(){
	resp := make(map[string]interface{}) //借助map打包成JSON
	defer c.RetData(resp)             //将数据打包成JSON返回给前端
	resp["errno"]=models.RECODE_OK
	resp["errmsg"]=models.RecodeText(models.RECODE_OK)

	//先从缓存中拿数据，拿不到数据再从数据库中拿

	//1.从redis缓存中拿数据，减少数据库的访问（一般将不变的东西放入缓存）
	//创建一个全局变量、连接
	cache_conn,_:=cache.NewCache("redis",`{"key":"houseproject","conn":":6379","dbNum":"0"}`)  //redis
	areaData:=cache_conn.Get("area") //返回二进制文件
	var areajson []models.Area

	if areaData!=nil{
		//将json解析为对应的数据结构
		json.Unmarshal(areaData.([]byte),&areajson)
		resp["data"]=areajson
		beego.Info("from cache")
		return
	}

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

	resp["data"]=areas
	beego.Info("from sql")
	//将数据结构打包成JSON方便传输，存入缓存
	json_str,json_err:=json.Marshal(areas)
	if json_err!=nil{
		beego.Info("encoding error")
		return
	}
	cache_conn.Put("area",json_str,time.Second*3600)
	//cache_conn.Put("area",areas,time.Second*60)

	//3.将数据打包成JSON返回给前端
	beego.Info("query data sucess ,areas =",resp["data"],",num =",num)

}


