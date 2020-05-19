package controllers

import (
	"encoding/json"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"strconv"
	_ "syscall"
	"houseproject/models"
)

type HousesController struct {
	beego.Controller
}

func (c *HousesController)RetData(resp map[string]interface{}){
	c.Data["json"]=resp
	c.ServeJSON()
}
//获取房屋信息
func (c *HousesController) GetHouses(){
	resp := make(map[string]interface{})
	defer c.RetData(resp)
	resp["errno"] = models.RECODE_OK
	resp["errmsg"] = models.RecodeText(models.RECODE_OK)

	//一对多查询
	var houses []models.House
	user_id:=c.GetSession("user_id").(int)  //********

	//从数据库读取房屋信息
	o:=orm.NewOrm()
	num,rerr:=o.QueryTable("house").Filter("user__id",user_id).All(&houses)
	//如果数据库中有带_的字段，查询时用两个__
	if rerr!=nil{
		resp["errno"]=models.RECODE_DBERR
		resp["errmsg"]=models.RecodeText(models.RECODE_DBERR)
		return
	}
	if num==0{
		resp["errno"]=models.RECODE_NODATA
		resp["errmsg"]=models.RecodeText(models.RECODE_NODATA)
		return
	}
	maphouse:=make(map[string][]models.House)
	maphouse["houses"]=houses
	resp["data"]=maphouse

}
//发布房源信息
func (c *HousesController) PostHouse(){
	resp := make(map[string]interface{})
	defer c.RetData(resp)
	resp["errno"] = models.RECODE_OK
	resp["errmsg"] = models.RecodeText(models.RECODE_OK)

	//1获取前端数据
	respdata:=make(map[string]interface{})
	json.Unmarshal(c.Ctx.Input.RequestBody,&respdata)

	//2判断前端数据的合法性

	//3将房源信息插入house表中
	house:=models.House{}
	house.Title=respdata["title"].(string)
	data:=respdata["price"].(string)
	house.Price,_=strconv.Atoi(data)     //数据库中price是整型
	house.Address=respdata["address"].(string)
	dataroom:=respdata["room_count"].(string)
	house.Room_count,_=strconv.Atoi(dataroom)
	dataacr:=respdata["acreage"].(string)
	house.Acreage,_=strconv.Atoi(dataacr)
	house.Unit=respdata["unit"].(string)
	dataca:=respdata["capacity"].(string)
	house.Capacity,_=strconv.Atoi(dataca)
	house.Beds=respdata["beds"].(string)
	datamin:=respdata["min_days"].(string)
	house.Min_days,_=strconv.Atoi(datamin)
	datamax:=respdata["max_days"].(string)
	house.Max_days,_=strconv.Atoi(datamax)
	datad:=respdata["deposit"].(string)
	house.Deposit,_=strconv.Atoi(datad)

	dataarea:=respdata["area_id"].(string)
	areaid,_:=strconv.Atoi(dataarea)
	area:=models.Area{Id:areaid}  //关联查询,一对一
	house.Area=&area

	user:=models.User{Id:c.GetSession("user_id").(int)}
	house.User=&user

	facilities:=[]models.Facility{}   //关联查询,多对多
	for _,f := range respdata["facility"].([]interface{}){
		fid,_:=strconv.Atoi(f.(string))
		fac:=models.Facility{Id:fid}
		facilities=append(facilities,fac)
	}
	//house.Facilities=&facilities

	o:=orm.NewOrm()
	houseid,err:=o.Insert(&house)         //为了获得house在库中的id主键
	if err!=nil{
		resp["errno"]=models.RECODE_DBERR
		resp["errmsg"]=models.RecodeText(models.RECODE_DBERR)
		beego.Info("111111111111")
		beego.Info(err)
		return
	}

	//3将对应的房源信息插入facility和house的多对多关系到表中
	house.Id=int(houseid)
	m2m:=o.QueryM2M(&house,"Facilities")  //第一个参数是对象，对象的主键必须存在,第二个参数是对象需要修改的M2M字段
	num,errm2m:=m2m.Add(facilities)
	if errm2m!=nil||num==0{
		resp["errno"]=models.RECODE_DBERR
		resp["errmsg"]=models.RecodeText(models.RECODE_DBERR)
		beego.Info("222222222222")
		return
	}

	data1:=make(map[string]interface{})
	data1["house_id"]=strconv.Itoa(house.Id)
	resp["data"]=data1
}
