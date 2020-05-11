package controllers

import (
	"github.com/astaxie/beego"
	"houseproject/models"
)

type MainController struct {
	beego.Controller
}

func (c *MainController) Get(){
	c.TplName = "index.tpl"
	//models.InsertOrder()
	models.OrderQuery()
}
