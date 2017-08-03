package controllers

import (
	"github.com/astaxie/beego"
	"github.com/canoeist2018/pldapi/utils"
)

type MainController struct {
	beego.Controller
}

func (c *MainController) Get() {
	c.Data["Website"] = "pldapi"
	c.Data["Email"] = "canoeist2018"
	c.Data["Version"] = utils.Version
	c.Data["Ip"] = utils.PldApiIp
	c.TplName = "index.tpl"
}

func (c *MainController) About() {
	c.Data["Website"] = "pldapi"
	c.Data["Email"] = "canoeist2018"
	c.Data["Version"] = utils.Version
	c.TplName = "about.tpl"
}
