package utils

import (
	"fmt"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/astaxie/beego"
)

// 获得某个账号ID所有的mapping用户账号
func FindUserIDbyAcctID(acctID string) []string {
	url := fmt.Sprintf("%s/index.php/Ndisk/user/accountid/%s/bakp/1", s.pldip, acctID)
	resp, _ := s.Get(url)
	defer resp.Body.Close()

	doc, _ := goquery.NewDocumentFromResponse(resp)
	userID := make([]string, 0)
	p := &userID

	doc.Find("select > option").Each(func(index int, sel *goquery.Selection) {
		account := sel.AttrOr("value", "none")
		*p = append(*p, account)
	})
	return userID
}

// 将某个账号相关的所有用户账号都设置好 mapping
func AddMappingbyAccID(accID string) {
	//输入一个账号的 ID
	url := s.pldip + "/index.php/Ndisk/setUser/"
	listAccID := FindUserIDbyAcctID(accID)
	var postData string
	if listAccID != nil {
		for _, id := range listAccID {
			postData = postData + "groupUserId%5B%5D=" + id + "&"
		}
	}
	postData = postData + "acct_id=" + accID + "&ajax=1"

	resp, err := s.Post(url, strings.NewReader(postData))
	defer resp.Body.Close()

	if err != nil {
		beego.Error("添加账号" + accID + "mapping 失败！")
	} else {
		beego.Info("添加账号" + accID + "mapping 成功！")
	}
}
