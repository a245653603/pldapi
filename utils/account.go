package utils

import (
	"fmt"
	"net/url"
	"regexp"
	"strconv"
	"strings"

	"github.com/canoeist2018/pldapi/models"

	"github.com/PuerkitoBio/goquery"
	"github.com/astaxie/beego"
)

func LoadAllAccount() {
	url := fmt.Sprintf("%s/index.php/DevAcct/index/?listRows=20&p=1", s.pldip)
	resp, _ := s.Get(url)
	defer resp.Body.Close()

	doc, _ := goquery.NewDocumentFromResponse(resp)

	page := doc.Find("div.fRig > a").Last().Text()
	page = strings.Replace(page, ".", "", -1)
	page = strings.TrimSpace(page)
	pageNum, err := strconv.Atoi(page)
	if err != nil {
		beego.Error("PageNum 获取错误：", pageNum)
	}

	for i := 1; i <= pageNum; i++ {
		addAccount(i)
	}

}

func ReloadAllAcount() {
	models.DropAccountInfo()
	LoadAllAccount()
}

func addAccount(i int) {
	url := fmt.Sprintf("%s/index.php/DevAcct/index/?listRows=20&p=%d", s.pldip, i)
	resp, _ := s.Get(url)
	defer resp.Body.Close()

	doc, _ := goquery.NewDocumentFromResponse(resp)

	doc.Find("table.list > tbody > tr.row ").Each(func(index int, sel *goquery.Selection) {
		account := new(models.AccountInfo)
		id, _ := sel.Find("input").Attr("value")
		var err error
		account.Id, err = strconv.Atoi(id)
		if err != nil {
			beego.Error("获到 id 出错！无法转到数字！", account)
		}

		account.Name = sel.Find("td").Eq(2).Text()

		device, _ := sel.Find("td").Eq(4).Find("a").Attr("onclick")
		re, _ := regexp.Compile(`\d+`)
		device = re.FindString(device)
		account.DeviceId, err = strconv.Atoi(device)
		if err != nil {
			beego.Error("获取设备 id 出错！无法转换到数字！", device)
		}

		ip := sel.Find("td").Eq(4).Text()
		re, _ = regexp.Compile(`<([\d.]+)>`)
		ip = re.FindStringSubmatch(ip)[1] // FindString(ip)
		account.Ip = ip

		account.Status = sel.Find("td").Eq(5).Text()
		account.Priv = sel.Find("td").Eq(6).Text()
		account.Sync = sel.Find("td").Eq(7).Text()
		account.Mapping, _ = sel.Find("td > a").Last().Attr("title")

		_, err = models.AccountInfoAdd(account)
		if err != nil {
			fmt.Print(err)
		}
	})

}

func ChangePasswordbyAcctID(id int) string {
	urladdress := fmt.Sprintf("%s/index.php/DevAcct/resetPasswd/", s.pldip)
	idstring := strconv.Itoa(id)

	account := models.GetAccountById(id)
	if account.Name == "tomcat" {
		v := url.Values{}
		v.Set("acct_epw", tomcatpass)
		v.Set("repassword", tomcatpass)
		v.Set("id", idstring)
		postBody := v.Encode()
		resp, err := s.Post(urladdress, strings.NewReader(postBody))
		beego.Info(postBody)

		defer resp.Body.Close()
		if err != nil {
			return "err"
		}
		return "tomcat"
	} else if account.Name == "log" {
		v := url.Values{}
		v.Set("acct_epw", logpass)
		v.Set("repassword", logpass)
		v.Set("id", idstring)
		postBody := v.Encode()
		resp, err := s.Post(urladdress, strings.NewReader(postBody))
		defer resp.Body.Close()
		if err != nil {
			return "err"
		}
		return "log"

	}

	return ""
}
