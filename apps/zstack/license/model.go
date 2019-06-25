package license

import (
	"api_system_iris/utils"
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/spf13/viper"
	"io/ioutil"
	"os/exec"
	"time"
)

const (
	Trial = iota + 1
	Paid
	OEM
	Free
	TrialExt
	Hybrid
 	AddOn
	HybridTrialExt
)

type RetErrorJson struct {
	Code 	int			`json:"code"`
	Error 	string		`json:"Error"`
	Tips 	string		`json:"Tips"`
}

type ReqCode struct {
	PrivateKey		string			`json:"privateKey"`
	LicenseRequest	string			`json:"licenseRequest"`
	LicenseReq 		LicenseRequest	`json:"-"`
}

type LicenseRequest struct{
	Thumbprint	string		`json:"thumbprint"`
	Pubkey		string		`json:"pubkey"`
}

type LicenseInfo struct {
	User 		string		`json:"user"`
	Licenseid	string		`json:"licenseid"`
	Type		int			`json:"type"`
	Hostnum		int			`json:"hostnum"`
	Issuetime 	string		`json:"issuetime"`
	Expiretime	string		`json:"expiretime"`
	Thumbprint 	string		`json:"thumbprint"`
	Prodinfo	string		`json:"prodinfo"`
	Cpunum 	 	int			`json:"Cpunum"`
}

type LicenseBody struct {
	License 	string		`json:"license"`
	Aeskey 		string		`json:"aeskey"`
	LicInfo		LicenseInfo	`json:"-"`
}

func (rt *RetErrorJson) Set(code int, error string, tips string) {
	rt.Code = code
	rt.Error = error
	rt.Tips = tips
}

func (rt *RetErrorJson) SetError(err string,tip string) {
	rt.Set(-1,err,tip)
}

func (this *ReqCode) DecodeSet(str string) bool {
	res,err := base64.StdEncoding.DecodeString(str)
	if err != nil {
		return false
	}
	json.Unmarshal(res,this)
	this.LicenseReq.DecodeSet(this.LicenseRequest)
	return true
}

func (this *LicenseRequest) DecodeSet(str string) bool {
	res,err := base64.StdEncoding.DecodeString(str)
	if err != nil {
		return false
	}
	json.Unmarshal(res,this)
	return true
}

func (this *LicenseInfo) Validate() (res RetErrorJson) {
	// TODO: 2019/6/24  参数校验
	res = RetErrorJson{
		Code:1,
		Error:"数据校验通过",
		Tips:"数据校验通过",
	}
	prodinfos := []string{
		"vmware","project-management","disaster-recovery","v2v","baremetal",
	}
	if this.Type == AddOn {
		flag := 0
		/*
			支持模块列表
				vmware project-management disaster-recovery v2v baremetal
		*/
		for _, prodinfo := range prodinfos {
			if this.Prodinfo == prodinfo {
				flag = 1
				break
			}
		}
		if flag == 0 {
			res.SetError(this.Prodinfo + "模块不存在","请输入正确的模块名称")
		}
	}
	date,err := time.Parse("2006-01-02", this.Expiretime)
	if err != nil {
		res.SetError("授权到期日期不合法","请使用 YYYY-MM-DD 格式")
	}else {
		this.Expiretime = date.Format("2006-01-02T15:04:05+08:00")
	}

	return
}

func (this *LicenseInfo) GetEncode() string{
	res,_ := json.Marshal(this)
	return base64.StdEncoding.EncodeToString(res)
}

func (this *LicenseBody) GetEncode() string {
	var buf bytes.Buffer
	this.License = this.LicInfo.GetEncode()
	res,_ := json.Marshal(this)
	breaker := new(utils.LineBreaker)
	breaker.SetWriter(&buf)
	breaker.Handle([]byte(base64.StdEncoding.EncodeToString(res)),72)
	return "\n" + string(buf.Bytes())
}

func (this *LicenseBody) SignRes() (res []byte) {
	//尚未有Go 语言的smime签名库实现，暂且调用linux 命令openssl进行签名
	prikey := utils.GetAbsPath(viper.GetString("zstack.license.privatekey"))
	pemkey := utils.GetAbsPath(viper.GetString("zstack.license.cert"))
	PWD := utils.GetCurrentPath() + "/conf/certs"
	if prikey == "" || !utils.FileExists(prikey) {
		prikey = PWD + "/zstack.key"
	}
	if pemkey == "" || !utils.FileExists(pemkey) {
		pemkey = PWD + "/zstack.pem"
	}
	data := this.GetEncode()
	openssl := "openssl smime -sign -inkey " + prikey + " -signer " + pemkey
	cmd := exec.Command("/bin/bash","-c",openssl)
	stdin,_ := cmd.StdinPipe()
	stdout,_ := cmd.StdoutPipe()
	if err := cmd.Start(); err != nil{
		fmt.Println("Execute failed when Start:" + err.Error())
		return nil
	}
	stdin.Write([]byte(data))
	stdin.Close()
	res,_ = ioutil.ReadAll(stdout)
	stdout.Close()
	return
}