package routers

import (
	"github.com/astaxie/beego"
	"houseproject/controllers"
)

func init() {
	beego.Router("/", &controllers.MainController{})   //默认调用Get()函数
	//获得区域信息
    beego.Router("/api/v1.0/areas", &controllers.AreaController{},"get:GetArea")
	//
	beego.Router("/api/v1.0/houses/index", &controllers.HouseIndexController{},"get:GetHouseIndex")
	//设置session、删除session退出登录
	beego.Router("/api/v1.0/session", &controllers.SessionController{},"get:GetSessionData;delete:DeleteSession")
	//用户注册
	beego.Router("/api/v1.0/users", &controllers.UserController{},"post:UserRigster")
	//用户登录
	beego.Router("/api/v1.0/sessions", &controllers.SessionController{},"post:UserLogin")
	//上传头像
	beego.Router("/api/v1.0/user/avatar", &controllers.UserController{},"post:PostAvatar")
	//修改用户名
	beego.Router("/api/v1.0/user/name", &controllers.UserController{},"put:UpdateName")
	//获得用户信息
	beego.Router("/api/v1.0/user", &controllers.UserController{},"get:GetUserName")
	//实名认证
	beego.Router("/api/v1.0/user/auth", &controllers.UserController{},"get:GetAuth;post:UserAuth")
	//查询房源信息
	beego.Router("/api/v1.0/user/houses", &controllers.HousesController{},"get:GetHouses")
	//发布房源信息
	beego.Router("/api/v1.0/houses", &controllers.HousesController{},"post:PostHouse")
}

