package controllers
import (
	"encoding/json"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	_ "syscall"
	"houseproject/models"
)

type UserController struct {
	beego.Controller
}

func (c *UserController)RetData(resp map[string]interface{}){
	c.Data["json"]=resp
	c.ServeJSON()
}
func (c *UserController) PostUser(){
	//获取前端传来的JSON数据
	resp:=make(map[string]interface{})
	defer c.RetData(resp)
	json.Unmarshal(c.Ctx.Input.RequestBody,&resp)   //RequestBody获取到请求中的数据
													//Unmarshal方法将数据转换为json类型并保存到resp里
	//把数据插入到数据库
	o:=orm.NewOrm()
	user:=models.User{}
	user.Password_hash=(resp["password"]).(string)
	user.Name=(resp["mobile"]).(string)
	user.Mobile=(resp["mobile"]).(string)
	id,err:=o.Insert(&user)
	if err!=nil{
		resp["errno"]=models.RECODE_NODATA
		resp["errmsg"]=models.RecodeText(models.RECODE_NODATA)
		return
	}
	beego.Info("register success,id=",id)
	resp["errno"]=models.RECODE_OK
	resp["errmsg"]=models.RecodeText(models.RECODE_OK)

	//在注册结束后，set session
	//登录时就可以get session
	c.SetSession("name",user.Name)
}
//用户注册的步骤
/*
1、设置路由
2、添加user.go,写POST user的代码
注意两点：
配置文件中一定要设置copyrequestbody=true
json.Unmarshal(c.Ctx.Input.RequestBody,&resp)
 */