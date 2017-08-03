package models

import (
	"github.com/astaxie/beego/orm"
	_ "github.com/go-sql-driver/mysql"
)

func init() {
	databaseUser := "pldapi"
	databasePass := "pldapi"
	orm.RegisterDataBase("default", "mysql", databaseUser+":"+databasePass+"@/pldapi?charset=utf8")
	orm.RegisterModel(new(DeviceInfo), new(AccountInfo), new(GroupInfo))
}
