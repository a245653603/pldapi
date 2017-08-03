package jobs

import (
	"compress/gzip"

	"encoding/json"

	"fmt"
	"io"
	"io/ioutil"

	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/canoeist2018/pldapi/utils"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/toolbox"
)

var (
	username string = utils.Pldusername
	password string = utils.Pldpassword
	pldip    string = utils.Pldip
)

type LoginStatus struct {
	Status int
	Info   string
	Data   bool
}

func init() {
	beego.Info("当前的版本为： ", utils.Version)
	beego.Info("堡垒机 IP 地址为： ", pldip)

	env := beego.BConfig.RunMode
	beego.Info("当前的运行环境为： ", env)

	task := toolbox.NewTask("task", "0 */20 * * * *", func() error { startTask(); return nil })
	err := task.Run()
	if err != nil {
		fmt.Println(err)
	}
	toolbox.AddTask("task", task)

	task2 := toolbox.NewTask("task", "0 1 */4 * * *", func() error { loadAllDb(); return nil })
	// 当前环境如果不是开发环境，则程序开始运行时会执行数据库初始化，避免调试的时候过度拖取数据
	if env != "dev" {
		err = task2.Run()
		if err != nil {
			fmt.Println(err)
		}

	}
	toolbox.AddTask("task2", task2)
	toolbox.StartTask()

}

func startTask() {
	beego.Info("当前的会话状态是: ", utils.GetSession())
	for !isLogin() {
		var sessionIDforsuccess string
		for sessionIDforsuccess == "" {
			sessionIDforsuccess = submit()
		}
		utils.GetSession().Set("pldsession", sessionIDforsuccess)
		utils.GetSession().SetSessionID(sessionIDforsuccess)
		beego.Info("新的会话状态是: ", sessionIDforsuccess)
	}
}

func loadAllDb() {
	utils.ReloadAllDevice()
	utils.ReloadAllAcount()
	utils.ReloadAllGroup()
	beego.Info("设备、用户、分组已经成功载入！")
}

func submit() string {
	//获取OCR后的验证码
	verifyPage()
	verifyImg()
	verifyText := readVerifyFile()
	beego.Info("识别出验证码：", verifyText)

	postData := strings.NewReader("strategies=local&account=" + username + "&local=" + password + "&interface=no&cert_check=no&verify=" + verifyText + "&finger=&ajax=1&confirm_login=")

	response, err := utils.GetSession().Post(utils.GetSession().GetPldIP()+"/index.php/Public/checkLogin", postData)
	if err != nil {
		beego.Error("登陆失败！")
	}
	defer response.Body.Close()

	var reader io.Reader
	if response.Header.Get("Content-Encoding") == "gzip" {
		reader, err = gzip.NewReader(response.Body)
		if err != nil {
			return ""
		}
	} else {
		reader = response.Body
	}

	body, _ := ioutil.ReadAll(reader)
	var ls LoginStatus
	json.Unmarshal(body, &ls)
	beego.Info("堡垒机登陆信息返回值：", ls)

	if len(response.Header["Set-Cookie"]) > 2 {
		str := response.Header["Set-Cookie"][2]
		return str
	} else {
		return ""
	}
}

// 获取验证码图片，并保存在对应的文件夹中
func verifyPage() {
	response, err := utils.GetSession().Get(utils.GetSession().GetPldIP() + "/index.php/Public/verify")
	if err != nil {
		beego.Error("获取验证码失败！")
	}
	defer response.Body.Close()

	userFile := "tmp/verify.gif"
	gif, err := os.Create(userFile)
	if err != nil {
		beego.Error("创建 tmp/verify.gif 文件失败", err)
	}
	defer gif.Close()
	gifByte, _ := ioutil.ReadAll(response.Body)
	gif.Write(gifByte)
}

func verifyImg() {
	exec.Command("/bin/bash", "-c", "convert -compress none -depth 8 -alpha off -colorspace Gray tmp/verify.gif tmp/verify.tiff").Start()
	time.Sleep(2 * time.Second)
	exec.Command("/bin/bash", "-c", "convert tmp/verify.tiff -scale 100% tmp/verify2.tiff").Start()
	time.Sleep(2 * time.Second)
	exec.Command("/bin/bash", "-c", "tesseract tmp/verify2.tiff tmp/1").Start()
	time.Sleep(2 * time.Second)
}

func readVerifyFile() string {
	userFile := "tmp/1.txt"
	txt, err := os.Open(userFile)
	if err != nil {
		fmt.Print(err)
		return ""
	}
	defer txt.Close()

	content, err := ioutil.ReadAll(txt)
	str := string(content)
	str = strings.Replace(str, " ", "", -1)
	str = strings.Replace(str, "\n", "", -1)

	return str

}

func isLogin() bool {
	response, _ := utils.GetSession().Get(utils.GetSession().GetPldIP() + "/index.php/User/add/")
	defer response.Body.Close()

	if response.StatusCode == 200 {
		beego.Info("会话状态保持可用!")
		return true
	} else {
		beego.Info("会话状态已不可用!")
		return false
	}
}
