package controllers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"strconv"
	"strings"

	"github.com/canoeist2018/pldapi/models"
	"github.com/canoeist2018/pldapi/utils"

	"github.com/astaxie/beego"
)

//返回给客户端的结果
type Result struct {
	Code int
	Info string
}

//添加使用的设备名和IP地址
type Device struct {
	Device string
	Ip     string
}

type RspondStatus struct {
	Status int
	Info   string
	Data   bool
}

type PldApiController struct {
	beego.Controller
}

func (this *PldApiController) Get() {
	var result Result
	result.Code = 1
	result.Info = `No Get method found, There only Post method to use!
	http://[ip]/pldapi`
	this.Data["json"] = result
	this.ServeJSON()
}

func (this *PldApiController) Post() {
	//获取需要添加的设备和IP地址
	var device Device
	var result Result
	json.Unmarshal(this.Ctx.Input.RequestBody, &device)

	url := "https://" + utils.Pldip + "/index.php/Device/insert/"

	ipList := strings.Split(device.Ip, ".")
	beego.Info("开始添加设备：", device)
	postData := strings.NewReader(fmt.Sprintf("dev_name=%s&4802__ip=%s&4802__ip=%s&4802__ip=%s&4802__ip=%s&dev_ip=%s&dev_type=2&dev_encoding=UTF-8&dev_disable=0&dev_prep=0&dev_priv_switch=&dev_priv_prompt=&dev_room=IDC&dev_seat=&dev_manager=&dev_remark=", device.Device, ipList[0], ipList[1], ipList[2], ipList[03], device.Ip))

	respond, err := utils.GetSession().Post(url, postData)
	if err != nil {
		beego.Error("添加设备出错：", device)
		result.Code = 1
		result.Info = fmt.Sprintf("添加设备出错： %s", device)
		return
	}

	body, _ := ioutil.ReadAll(respond.Body)
	var rs RspondStatus
	json.Unmarshal(body, &rs)
	beego.Info("堡垒机返回结果：", rs)

	result.Code = 0
	result.Info = fmt.Sprintf("添加设备成功： %s", device)

	//开始添加用户
	if deviceID := utils.FindDeviceIDbyIP(device.Ip); deviceID != "" {
		utils.AddRootUser(deviceID)
		rootID := utils.FindUserByDeviceID("root", deviceID)
		utils.AddTomcatUser(rootID, deviceID)
		utils.AddLogUser(rootID, deviceID)
	}

	this.Data["json"] = result
	this.ServeJSON()
}

func (this *PldApiController) GetDevices() {
	this.TplName = "device/list.tpl"
	this.Data["devices"] = models.GetAllDevices()
}

func (this *PldApiController) LoadAllDevice() {
	utils.ReloadAllDevice()

	result := Result{1, "设备刷新完成！"}
	this.Data["json"] = result
	this.ServeJSON()
}

func (this *PldApiController) DropDeviceInfo() {
	models.DropDeviceInfo()
	this.Data["json"] = &Result{0, "设备表已经删除！"}
	this.ServeJSON()
}

func (this *PldApiController) LoadAllAccount() {
	utils.ReloadAllAcount()

	result := Result{0, "设备刷新完成！"}
	this.Data["json"] = result
	this.ServeJSON()
}

func (this *PldApiController) GetAccounts() {
	this.TplName = "account/list.tpl"
	this.Data["accounts"] = models.GetAllAccounts()

}

func (this *PldApiController) AddAccount() {
	var device *models.DeviceInfo
	data := this.Ctx.Input.RequestBody
	ip := string(data)
	if !utils.IsIP(ip) {
		this.Data["json"] = &Result{1, "IP 地址输入无效！"}
	} else if device = models.GetDeviceId(ip); device.Id != 0 {
		utils.AddAllUser(device.Id)
		this.Data["json"] = &Result{0, "创建成功！"}
	} else {
		this.Data["json"] = &Result{1, "没有查到该 IP 地址： " + ip}
	}

	this.ServeJSON()
}

func (this *PldApiController) ChangePassword() {
	id := this.Ctx.Input.Param(":id")
	idnum, _ := strconv.Atoi(id)
	res := utils.ChangePasswordbyAcctID(idnum)
	if res == "" {
		this.Data["json"] = &Result{1, id + " 没有找到对应的账号或不为tomcat,log账号！"}
	} else if res == "err" {
		this.Data["json"] = &Result{1, "堡垒机连接失败！"}
	} else {
		this.Data["json"] = &Result{0, id + " : " + res + " 修改成功！"}
	}
	this.ServeJSON()
}

func (this *PldApiController) ChangeAllPassword() {
	acc := this.Ctx.Input.Param(":acc")
	if acc != "tomcat" && acc != "log" {
		this.Data["json"] = &Result{1, "只能是 tomcat 或 log 账号！"}
		this.ServeJSON()
		return
	}

	accounts := models.GetAccounts(acc)
	accountsNum := len(accounts)
	if accountsNum == 0 {
		this.Data["json"] = &Result{1, "没有查询到相应的账号！"}
	} else {
		for _, account := range accounts {
			utils.ChangePasswordbyAcctID(account.Id)
		}
		this.Data["json"] = &Result{0, "共更新" + strconv.Itoa(accountsNum) + "账号！"}
	}

	this.ServeJSON()
}

func (this *PldApiController) UpdateDeviceName() {
	var result Result
	var device Device
	json.Unmarshal(this.Ctx.Input.RequestBody, &device)

	if !utils.IsIP(device.Ip) {
		result.Code = 1
		result.Info = "IP 地址输入无效"
	} else if !utils.ValidDeviceName(device.Device) {
		result.Code = 1
		result.Info = "设备名称输入无效"
	} else {
		utils.UpdateDeviceName(device.Device, device.Ip)
		result.Code = 0
		result.Info = "设备名称修改成攻"
	}

	this.Data["json"] = result
	this.ServeJSON()
}

func (this *PldApiController) GetMapping() {
	s := this.Ctx.Input.Param(":id")
	res := utils.FindUserIDbyAcctID(s)
	this.Data["json"] = res
	this.ServeJSON()
}

func (this *PldApiController) GetIpMapping() {
	s := this.Ctx.Input.Param(":id")
	this.Data["json"] = models.GetIpMapping(s)
	this.ServeJSON()
}

func (this *PldApiController) SetIpMapping() {
	s := this.Ctx.Input.Param(":id")
	if !utils.IsIP(s) {
		this.Data["json"] = &Result{1, s + " IP 地址无效！"}
		this.ServeJSON()
		return
	}

	account := models.GetIpMapping(s)

	if len(account) <= 0 {
		this.Data["json"] = &Result{1, s + " 没有可供映射的账号！"}
		this.ServeJSON()
		return
	}

	for _, v := range account {
		id := strconv.Itoa(v.Id)
		utils.AddMappingbyAccID(id)
	}
	this.Data["json"] = &Result{0, "已刷新" + s + " 映射的共" + strconv.Itoa(len(account)) + "账号！"}
	this.ServeJSON()
}

func (this *PldApiController) SetMapping() {
	s := string(this.Ctx.Input.RequestBody)
	utils.AddMappingbyAccID(s)
	this.Data["json"] = s
	this.ServeJSON()
}

func (this *PldApiController) SetAllMapping() {
	account := models.GetAllMappingAccount()
	if len(account) <= 0 {
		this.Data["json"] = &Result{1, "没有可供映射的账号！"}
		this.ServeJSON()
		return
	}

	for _, v := range account {
		id := strconv.Itoa(v.Id)
		utils.AddMappingbyAccID(id)
	}
	this.Data["json"] = &Result{0, "已刷新映射,共" + strconv.Itoa(len(account)) + "账号！"}
	this.ServeJSON()

}

func (this *PldApiController) LoadAllGroup() {
	utils.ReloadAllGroup()
	this.Data["json"] = &Result{0, "账号表刷新完成！"}
	this.ServeJSON()
}

func (this *PldApiController) LoadAll() {
	utils.ReloadAllDevice()
	utils.ReloadAllAcount()
	utils.ReloadAllGroup()
	this.Data["json"] = &Result{0, "设备、用户、分组已经成功载入！"}
	this.ServeJSON()
}
