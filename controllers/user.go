package controllers
import (
	"encoding/json"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"io"
	"os"
	"path/filepath"
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
func (c *UserController) UserRigster(){
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

//用户上传头像，user数据库存储头像的URL,图片文件存储在本地
//D:\Files\GO\src\houseproject\static\upload\avatar
func (c *UserController)PostAvatar(){
	resp:=make(map[string]interface{})
	defer c.RetData(resp)
	resp["errno"]=models.RECODE_OK
	resp["errmsg"]=models.RecodeText(models.RECODE_OK)

	//1获取上传的文件
	filedata,fileheader,fileerr:=c.GetFile("avatar")  //返回文件的二进制数据，文件头信息，err
	if fileerr!=nil{
		resp["errno"]=models.RECODE_REQERR
		resp["errmsg"]=models.RecodeText(models.RECODE_REQERR)
		return
	}

	//2得到文件后缀,判断是一个图片文件
	fileExt := filepath.Ext(fileheader.Filename)  //获取最后一个后缀
	if fileExt != ".jpg" && fileExt != ".png" && fileExt != ".gif" && fileExt != ".jpeg" {
		resp["errno"]=models.RECODE_DATAERR
		resp["errmsg"]=models.RecodeText(models.RECODE_DATAERR)
		return
	}

	//3拼接fileurl,将文件存在本地文件夹内
	fileDir:="D:/Files/GO/src/houseproject/static/upload/avatar/"+c.GetSession("name").(string)
	fileurl:=fileDir+"/"+fileheader.Filename
	//先创建user的文件夹
	merr := os.MkdirAll(fileDir, os.ModePerm)
	if merr != nil {
		resp["errno"]=models.RECODE_DATAERR
		resp["errmsg"]=models.RecodeText(models.RECODE_DATAERR)
		return
	}
	beego.Info(fileurl)
	desFile, ferr := os.Create(fileurl)   //创建本地的目标文件
	if ferr != nil {
		resp["errno"]=models.RECODE_DATAERR
		resp["errmsg"]=models.RecodeText(models.RECODE_DATAERR)
		return
	}
	//将浏览器客户端上传的文件拷贝到本地路径的文件里面
	_,cerr:= io.Copy(desFile, filedata)
	if cerr != nil {
		resp["errno"]=models.RECODE_DATAERR
		resp["errmsg"]=models.RecodeText(models.RECODE_DATAERR)
		return
	}

	//4从session读取user.id
	user_id:=c.GetSession("user_id")
	var user models.User
	user.Id=user_id.(int)

	//5用URL更新数据库user的avatarurl
	o:=orm.NewOrm()
	o.QueryTable("user").Filter("Id",user_id).One(&user)
	user.Avatar_url=fileurl
    _,updaterr:=o.Update(&user)
    if updaterr!=nil{
		resp["errno"]=models.RECODE_REQERR
		resp["errmsg"]=models.RecodeText(models.RECODE_REQERR)
		return
	}
	urlmap:=make(map[string]string)
	urlmap["avatar_url"]="http://192.168.43.166:8080/static/upload/avatar/"+c.GetSession("name").(string)+"/"+fileheader.Filename
	resp["data"]=urlmap
}
