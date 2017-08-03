package utils

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/astaxie/beego"
	"github.com/canoeist2018/pldapi/models"
)

func LoadAllDevice() {
	url := fmt.Sprintf("%s/index.php/Device/index/?listRows=20&p=1", s.pldip)
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
		addPage(i)
	}

}

func ReloadAllDevice() {
	models.DropDeviceInfo()
	LoadAllDevice()
}

func addPage(i int) {
	url := fmt.Sprintf("%s/index.php/Device/index/?listRows=20&p=%d", s.pldip, i)
	resp, _ := s.Get(url)
	defer resp.Body.Close()

	doc, _ := goquery.NewDocumentFromResponse(resp)

	doc.Find("table.list > tbody > tr.row ").Each(func(index int, sel *goquery.Selection) {
		device := new(models.DeviceInfo)
		id, _ := sel.Find("input").Attr("value")
		var err error
		device.Id, err = strconv.Atoi(id)
		if err != nil {
			beego.Error(id + ":获到 id 出错！无法转到数字！")
		}
		device.Name = sel.Find("td").Eq(2).Text()
		device.Ip = sel.Find("td").Eq(3).Text()
		device.Dtype = sel.Find("td").Eq(4).Text()
		device.Status = sel.Find("td").Eq(6).Text()
		duser := sel.Find("td").Eq(7).Text()
		duser = strings.Replace(duser, "(", "", -1)
		duser = strings.Replace(duser, ")", "", -1)
		device.Duser, err = strconv.Atoi(duser)
		if err != nil {
			beego.Error(id + ":获到 duser 出错！无法转到数字！")
		}

		models.DeviceInfoAdd(device)

	})
}
