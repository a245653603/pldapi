package utils

import (
	"crypto/tls"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/url"
	"regexp"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/config"
)

type SessionSaver struct {
	session string
	client  *http.Client
	conf    config.Configer
	pldip   string
}

var (
	Version     = "v0.0.3"
	Pldip       string
	PldApiIp    string
	Pldusername string
	Pldpassword string
	rootpass    string
	tomcatpass  string
	logpass     string
	s           = new(SessionSaver)
	iniconf     = "conf/config.ini"
)

func init() {
	var err error
	if s.conf, err = config.NewConfig("ini", iniconf); err != nil {
		beego.Error("打开配置文件 config.ini 错误！")
	}

	s.session = s.conf.String("pldsession")
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	s.client = &http.Client{Transport: tr, CheckRedirect: func(*http.Request, []*http.Request) error { return errors.New("close redirector") }}

	Pldip = s.conf.String("pldip")
	s.pldip = Pldip

	PldApiIp = s.conf.String("pldapiip")
	rootpass = s.conf.String("rootpass")
	tomcatpass = s.conf.String("tomcatpass")
	logpass = s.conf.String("logpass")
	Pldusername = s.conf.String("pldusername")
	Pldpassword = s.conf.String("pldpassword")
}

//session功能
func GetSession() *SessionSaver {
	return s
}

func (s *SessionSaver) GetSessionID() string {
	return s.session
}

func (s *SessionSaver) SetSessionID(session string) {
	s.session = session
}

func (s *SessionSaver) GetPldIP() string {
	return s.pldip
}

func (s *SessionSaver) Set(key string, value string) {
	s.conf.Set(key, value)
	s.conf.SaveConfigFile(iniconf)
}

func (s *SessionSaver) String(key string) string {
	return s.conf.String(key)
}

func (s *SessionSaver) Get(url string) (*http.Response, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		beego.Error("NewRequest failed with URL: %s", url)
		return nil, err
	}

	req.Header = newHTTPHeaders()
	response, err := s.client.Do(req)

	return response, err
}

func (s *SessionSaver) Post(url string, body io.Reader) (*http.Response, error) {
	req, err := http.NewRequest("POST", url, body)
	if err != nil {
		beego.Error("NewRequest failed with URL: %s", url)
		return nil, err
	}

	headers := newHTTPHeaders()
	headers.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header = headers

	response, err := s.client.Do(req)

	return response, err
}

func newHTTPHeaders() http.Header {
	headers := make(http.Header)
	headers.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8")
	headers.Set("Accept-Charset", "GBK,utf-8;q=0.7,*;q=0.3")
	headers.Set("Accept-Language", "zh-CN,zh;q=0.8")
	headers.Set("Cache-Control", "max-age=0")
	headers.Set("Connection", "keep-alive")
	headers.Set("Cookie", s.session)
	return headers
}

//查找设备ID
func FindDeviceIDbyIP(ip string) string {
	url := fmt.Sprintf("%s/index.php/Device?dev_name=%s", s.pldip, ip)
	resp, _ := s.Get(url)
	defer resp.Body.Close()
	doc, _ := goquery.NewDocumentFromResponse(resp)

	res, ok := doc.Find("table.list").Find("td.tCenter > input").Attr("value")
	if ok != true {
		return ""
	}
	return res
}

//账号相关,根据提供的设备IP地址，找到该设备对应的用户ID
//FindUserByIP("root", "10.1.1.1")
//找不到则返回空
func FindUserByDeviceIP(user, ip string) string {
	deviceID := FindDeviceIDbyIP(ip)

	if deviceID == "" {
		return ""
	}
	return FindUserByDeviceID(user, deviceID)
}

//
func FindUserByDeviceID(user, deviceID string) string {
	url := fmt.Sprintf("%s/index.php/DevAcct/index/devid/%s", s.pldip, deviceID)
	resp, _ := s.Get(url)
	defer resp.Body.Close()

	doc, _ := goquery.NewDocumentFromResponse(resp)
	users := make(map[string]string)

	doc.Find("table.list > tbody > tr").Each(func(index int, sel *goquery.Selection) {
		account, _ := sel.Find("td.tCenter >input").Attr("value")
		username := sel.Find("td > span.acct").Text()
		users[username] = account
	})
	return users[user]
}

func AddRootUser(deviceID string) {
	urladdress := s.pldip + "/index.php/DevAcct/insert/"
	//postData := strings.NewReader("acct_name=root&acct_epw=" + rootpass + "&repassword=" + rootpass + "&dev_id=" + deviceID + "&app_id=1&app_name=SSH&app_ico=ico%2Fapp_SecureCRT.png&app_type=0&app_port=22&app_login_prompt=&app_pwd_prompt=assword%3A&app_parm=%2Fhomeunlock&acct_rdp_flag=0&acct_type=0&acct_disable=0&acct_prompt=%24&acct_priv=1&allow_join=0&acct_sync=0&acct_priv_login=0&script_pwd=0&acct_joint_ausr=")
	v := url.Values{}
	v.Set("acct_name", "root")
	v.Set("acct_epw", rootpass)
	v.Set("repassword", rootpass)
	v.Set("dev_id", deviceID)
	v.Set("app_id", "1")
	v.Set("app_name", "SSH")
	v.Set("app_ico", "ico/app_SecureCRT.png")
	v.Set("app_type", "0")
	v.Set("app_port", "22")
	v.Set("app_login_prompt", "")
	v.Set("app_pwd_prompt", "assword:")
	v.Set("app_parm", "/homeunlock")
	v.Set("acct_rdp_flag", "0")
	v.Set("acct_type", "0")
	v.Set("acct_disable", "0")
	v.Set("acct_prompt", "$")
	v.Set("acct_priv", "1")
	v.Set("allow_join", "0")
	v.Set("acct_sync", "0")
	v.Set("acct_priv_login", "0")
	v.Set("script_pwd", "0")
	v.Set("acct_joint_ausr", "")
	postBody := v.Encode()

	resp, err := s.Post(urladdress, strings.NewReader(postBody))
	defer resp.Body.Close()

	if err != nil {
		beego.Error("添加 root 用户失败！")
	}
	beego.Info("添加 root 用户成功！")
}

func AddTomcatUser(rootID, deviceID string) {
	urladdress := s.pldip + "/index.php/DevAcct/insert/"
	//postData := strings.NewReader("acct_name=tomcat&acct_epw=" + tomcatpass + "&repassword=" + tomcatpass + "&dev_id=" + deviceID + "&app_id=1&app_name=SSH&app_ico=ico%2Fapp_SecureCRT.png&app_type=0&app_port=22&app_login_prompt=&app_pwd_prompt=assword%3A&app_parm=%2Fhomeunlock&acct_rdp_flag=0&acct_type=0&acct_disable=0&acct_prompt=%24&acct_priv=0&allow_join=0&acct_sync=1&acct_priv_login=" + rootID + "&entrust=1&script_pwd=14&acct_joint_ausr=")
	v := url.Values{}
	v.Set("acct_name", "tomcat")
	v.Set("acct_epw", tomcatpass)
	v.Set("repassword", tomcatpass)
	v.Set("dev_id", deviceID)
	v.Set("app_id", "1")
	v.Set("app_name", "SSH")
	v.Set("app_ico", "ico/app_SecureCRT.png")
	v.Set("app_type", "0")
	v.Set("app_port", "22")
	v.Set("app_login_prompt", "")
	v.Set("app_pwd_prompt", "assword:")
	v.Set("app_parm", "/homeunlock")
	v.Set("acct_rdp_flag", "0")
	v.Set("acct_type", "0")
	v.Set("acct_disable", "0")
	v.Set("acct_prompt", "$")
	v.Set("acct_priv", "0")
	v.Set("allow_join", "0")
	v.Set("acct_sync", "1")
	v.Set("acct_priv_login", rootID)
	v.Set("entrust", "1")
	v.Set("script_pwd", "14")
	v.Set("acct_joint_ausr", "")
	postBody := v.Encode()

	resp, err := s.Post(urladdress, strings.NewReader(postBody))
	defer resp.Body.Close()

	if err != nil {
		beego.Error("添加 tomcat 用户失败！")
	}
	beego.Info("添加 tomcat 用户成功！")
}

func AddLogUser(rootID, deviceID string) {
	urladdress := s.pldip + "/index.php/DevAcct/insert/"
	//postData := strings.NewReader("acct_name=log&acct_epw=" + logpass + "&repassword=" + logpass + "&dev_id=" + deviceID + "&app_id=1&app_name=SSH&app_ico=ico%2Fapp_SecureCRT.png&app_type=0&app_port=22&app_login_prompt=&app_pwd_prompt=assword%3A&app_parm=%2Fhomeunlock&acct_rdp_flag=0&acct_type=0&acct_disable=0&acct_prompt=%24&acct_priv=0&allow_join=0&acct_sync=1&acct_priv_login=" + rootID + "&entrust=1&script_pwd=14&acct_joint_ausr=")
	v := url.Values{}
	v.Set("acct_name", "log")
	v.Set("acct_epw", logpass)
	v.Set("repassword", logpass)
	v.Set("dev_id", deviceID)
	v.Set("app_id", "1")
	v.Set("app_name", "SSH")
	v.Set("app_ico", "ico/app_SecureCRT.png")
	v.Set("app_type", "0")
	v.Set("app_port", "22")
	v.Set("app_login_prompt", "")
	v.Set("app_pwd_prompt", "assword:")
	v.Set("app_parm", "/homeunlock")
	v.Set("acct_rdp_flag", "0")
	v.Set("acct_type", "0")
	v.Set("acct_disable", "0")
	v.Set("acct_prompt", "$")
	v.Set("acct_priv", "0")
	v.Set("allow_join", "0")
	v.Set("acct_sync", "1")
	v.Set("acct_priv_login", rootID)
	v.Set("entrust", "1")
	v.Set("script_pwd", "14")
	v.Set("acct_joint_ausr", "")
	postBody := v.Encode()

	resp, err := s.Post(urladdress, strings.NewReader(postBody))
	defer resp.Body.Close()

	if err != nil {
		beego.Error("添加 log 用户失败！")
	}
	beego.Info("添加 log 用户成功！")
}

func AddAllUser(deviceId int) {
	id := strconv.Itoa(deviceId)
	rootId := FindUserByDeviceID("root", id)
	if rootId == "" {
		AddRootUser(id)
		rootId = FindUserByDeviceID("root", id)
	}
	if FindUserByDeviceID("tomcat", id) == "" {
		AddTomcatUser(rootId, id)
	}
	if FindUserByDeviceID("log", id) == "" {
		AddLogUser(rootId, id)
	}
}

//更新设备名称
func UpdateDeviceName(devicename, ip string) bool {
	deviceID := FindDeviceIDbyIP(ip)
	if deviceID == "" {
		return false
	}
	url := s.pldip + "/index.php/Device/update/"
	postData := strings.NewReader("dev_name=" + devicename + "&dev_ip=" + ip + "&dev_type=2&dev_encoding=UTF-8&dev_disable=0&dev_priv_switch=&dev_priv_prompt=&dev_room=IDC&dev_seat=&dev_manager=&dev_remark=&id=" + deviceID)
	resp, err := s.Post(url, postData)
	defer resp.Body.Close()

	if err != nil {
		beego.Error("更新设备名称失败！")
		return false
	}
	return true
}

//其它功能
func IsIP(ip string) bool {
	return net.ParseIP(ip) != nil
}

func ValidDeviceName(device string) bool {
	result, _ := regexp.MatchString(`^[\w.-]{3,30}$`, device)
	return result
}

func strip(s string) string {
	return strings.TrimSpace(s)
}
