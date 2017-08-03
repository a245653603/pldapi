package routers

import (
	"github.com/astaxie/beego"
	"github.com/canoeist2018/pldapi/controllers"
)

func init() {

	beego.Router("/pldapi", &controllers.MainController{})
	beego.Router("/pldapi/about", &controllers.MainController{}, "get:About")

	beego.Router("/v1/device", &controllers.PldApiController{}, "get:Get;put:UpdateDeviceName;post:Post;delete:DropDeviceInfo")
	beego.Router("/v1/device/list", &controllers.PldApiController{}, "get:GetDevices")
	beego.Router("/v1/device/load", &controllers.PldApiController{}, "get:LoadAllDevice")
	beego.Router("/v1/account/load", &controllers.PldApiController{}, "get:LoadAllAccount")
	beego.Router("/v1/account/list", &controllers.PldApiController{}, "get:GetAccounts")
	beego.Router("/v1/account/", &controllers.PldApiController{}, "post:AddAccount")
	beego.Router("/v1/account/changepassword/:id([0-9]+)", &controllers.PldApiController{}, "put:ChangePassword")
	beego.Router("/v1/account/changeallpassword/:acc([a-z]{3,7})", &controllers.PldApiController{}, "put:ChangeAllPassword")
	beego.Router("/v1/mapping/list/:id([0-9]+)", &controllers.PldApiController{}, "get:GetMapping") // http://localhost:21080/v1/mapping/list/3550
	beego.Router("/v1/mapping/ip/:id([0-9.]+)", &controllers.PldApiController{}, "get:GetIpMapping;put:SetIpMapping")
	beego.Router("/v1/mapping/update", &controllers.PldApiController{}, "put:SetMapping")
	beego.Router("/v1/mapping/update/all", &controllers.PldApiController{}, "put:SetAllMapping")
	beego.Router("/v1/group/load", &controllers.PldApiController{}, "get:LoadAllGroup")
	beego.Router("/v1/loadall", &controllers.PldApiController{}, "get:LoadAll")
}
