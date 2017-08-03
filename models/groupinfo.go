package models

import (
	"github.com/astaxie/beego/orm"
)

type GroupInfo struct {
	Id       int
	GroupId  int
	DeviceId int
}

func GroupInfoAdd(group *GroupInfo) (int64, error) {
	return orm.NewOrm().Insert(group)
}

func DropGroupInfo() {
	orm.NewOrm().Raw("DELETE FROM group_info").Exec()
}
