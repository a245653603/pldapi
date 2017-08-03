package models

import (
	"github.com/astaxie/beego/orm"
)

type AccountInfo struct {
	Id       int
	Name     string
	DeviceId int
	Ip       string
	Status   string
	Priv     string
	Sync     string
	Mapping  string
}

func AccountInfoAdd(account *AccountInfo) (int64, error) {
	return orm.NewOrm().Insert(account)
}

func GetAllAccounts() []*AccountInfo {
	accounts := make([]*AccountInfo, 0)
	query := orm.NewOrm().QueryTable("account_info")
	query.OrderBy("id").Limit(-1).All(&accounts)
	return accounts
}

func GetAccounts(acc string) []*AccountInfo {
	accounts := make([]*AccountInfo, 0)
	query := orm.NewOrm().QueryTable("account_info")
	query.Filter("name", acc).Limit(-1).All(&accounts)
	return accounts
}

func GetIpMapping(ip string) []*AccountInfo {
	accounts := make([]*AccountInfo, 0)
	query := orm.NewOrm().QueryTable("account_info")
	query.OrderBy("id").Filter("ip", ip).Filter("mapping", "映射").All(&accounts)
	return accounts
}

func GetAllMappingAccount() []*AccountInfo {
	accounts := make([]*AccountInfo, 0)
	query := orm.NewOrm().QueryTable("account_info")
	query.OrderBy("id").Filter("mapping", "映射").All(&accounts)
	return accounts
}

func GetAccountById(id int) *AccountInfo {
	var account AccountInfo
	query := orm.NewOrm().QueryTable("account_info")
	query.Filter("id", id).One(&account)
	return &account
}

func DropAccountInfo() {
	orm.NewOrm().Raw("DELETE FROM account_info").Exec()
}
