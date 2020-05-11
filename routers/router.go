package routers

import (
	"github.com/astaxie/beego"
	"houseproject/controllers"
)

func init() {
    beego.Router("/", &controllers.MainController{},"get:Get")
}

