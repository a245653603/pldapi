package main

import (
	_ "github.com/canoeist2018/pldapi/jobs"
	_ "github.com/canoeist2018/pldapi/routers"
	_ "github.com/canoeist2018/pldapi/utils"

	"github.com/astaxie/beego"
)

func main() {
	beego.Run()
}
