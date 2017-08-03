package models

import (
	"github.com/astaxie/beego/orm"
)

type DeviceInfo struct {
	Id     int
	Name   string
	Ip     string
	Dtype  string
	Status string
	Duser  int
}

func GetDeviceId(ip string) *DeviceInfo {
	var device DeviceInfo
	query := orm.NewOrm().QueryTable("device_info")
	query.Filter("ip", ip).One(&device)
	return &device
}

func DeviceInfoAdd(device *DeviceInfo) (int64, error) {
	return orm.NewOrm().Insert(device)
}

func GetAllDevices() []*DeviceInfo {
	devices := make([]*DeviceInfo, 0)
	query := orm.NewOrm().QueryTable("device_info")
	query.OrderBy("id").Limit(-1).All(&devices)
	return devices
}

func DropDeviceInfo() {
	orm.NewOrm().Raw("DELETE FROM device_info").Exec()
}
