package utils

import (
	"fmt"
	"io/ioutil"
	"regexp"
	"strconv"

	"github.com/canoeist2018/pldapi/models"
)

func LoadAllGroup() {

	url := fmt.Sprintf("%s/index.php/Privilege/showGroupHostPriv/", s.pldip)
	resp, _ := s.Get(url)
	defer resp.Body.Close()

	doc, _ := ioutil.ReadAll(resp.Body)
	docString := string(doc)

	re, _ := regexp.Compile(`{id:"\w+",value:"dev(\d+)",pId:"dg(\d+)",name:"`)
	result := re.FindAllStringSubmatch(docString, -1)

	for _, content := range result {
		group := new(models.GroupInfo)
		group.GroupId, _ = strconv.Atoi(content[2])
		group.DeviceId, _ = strconv.Atoi(content[1])
		models.GroupInfoAdd(group)
	}

}

func ReloadAllGroup() {
	models.DropGroupInfo()
	LoadAllGroup()
}
