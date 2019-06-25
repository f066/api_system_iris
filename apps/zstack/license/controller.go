package license

import (
	"api_system_iris/utils"
	"crypto/md5"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"github.com/google/uuid"
	"github.com/kataras/iris"
	"github.com/spf13/viper"
	"strconv"
	"strings"
	"time"
)

type ResJson struct {
	Code	int			`json:"code"`
	Message string		`json:"message"`
	Url 	string		`json:"url"`
	File	string		`json:"file"`
	Body 	string		`json:"body"`
}

type FileInfo struct {
	Name 	string		`json:"name"`
	Body 	string		`json:"body"`
}

func (this *FileInfo) SetBody(body []byte) {
	this.Body = base64.StdEncoding.EncodeToString(body)
}

func (this *FileInfo) GetBody() (res []byte) {
	res,_ = base64.StdEncoding.DecodeString(this.Body)
	return
}

func download(ctx iris.Context)  {
	md5_key := ctx.Params().Get("md5")
	if utils.Bm.IsExist(md5_key) {
		var file FileInfo
		json.Unmarshal([]byte(utils.Bm.Get(md5_key).(string)),&file)
		ctx.Header("Content-Disposition","attachment; filename="+file.Name)
		ctx.Binary(file.GetBody())
	}else{
		ctx.JSON(RetErrorJson{Code:-1, Error:"文件不存在", Tips:"请使用生成服务生成先！",})
	}
}

func generate(ctx iris.Context) {
	LicenseType := 0
	lictype := ctx.Params().Get("type")
	switch lictype {
	case "Trial":
		LicenseType = Trial
	case "Paid":
		LicenseType = Paid
	case "OEM":
		LicenseType = OEM
	case "Free":
		LicenseType = Free
	case "TrialExt":
		LicenseType = TrialExt
	case "Hybrid":
		LicenseType = Hybrid
	case "AddOn":
		LicenseType = AddOn
	case "HybridTrialExt":
		LicenseType = HybridTrialExt
	}

	User := ctx.PostValue("User")
	requestCode := ctx.PostValue("RequestCode")
	expiretime := ctx.PostValue("Expiretime")
	Hostnum,_ := strconv.Atoi(ctx.PostValue("Hostnum"))
	Prodinfo := ctx.PostValue("Prodinfo")
	Cpunum,_ := strconv.Atoi(ctx.PostValue("Cpunum"))
	if LicenseType == AddOn {
		Hostnum,Cpunum = 0,0
	}else{
		Prodinfo = "ZStack license"
	}
	Prodinfos := strings.Split(Prodinfo,",")
	//计算请求参数md5特征
	req_val := User + requestCode + expiretime + Prodinfo + lictype + ctx.PostValue("Hostnum") + ctx.PostValue("Cpunum")
	req_md5 := md5.Sum([]byte(req_val))
	req_md5Key := hex.EncodeToString(req_md5[:])

	var reqCode ReqCode
	if !reqCode.DecodeSet(requestCode) {
		ctx.JSON(iris.Map{"status":"error","message":"请求码错误！",})
		return
	}

	now := time.Now()

	var file FileInfo
	if utils.Bm.IsExist(req_md5Key) {
		//存在缓存
		json.Unmarshal([]byte(utils.Bm.Get(req_md5Key).(string)),&file)
		RetLicense(ctx,req_md5Key,file)
	}else {
		var res RetErrorJson
		var gf utils.GzFile
		tf := utils.NewTar()
		licinfo := LicenseInfo{
			User:User,
			Type:LicenseType,
			Hostnum:Hostnum,
			Issuetime:now.Format("2006-01-02T15:04:05+08:00"),
			Expiretime:expiretime,
			Thumbprint:reqCode.LicenseReq.Thumbprint,
			Cpunum:Cpunum,
		}
		fn := "zstack_license"
		for _, prodinfo := range Prodinfos {
			licinfo.Licenseid = uuid.New().String()
			licinfo.Prodinfo =  prodinfo
			licinfo.Expiretime = expiretime
			res = licinfo.Validate()
			if res.Code == 1 {
				if licinfo.Type == AddOn {
					fn = "zstack_" + prodinfo + "_license"
				}
				lb := LicenseBody{LicInfo:licinfo}
				tf.AddFile(utils.NewFile(fn,lb.SignRes()))
			}else {
				break
			}
		}
		gf.Set(tf.GetFile())
		if res.Code == 1 {
			linfo := lictype
			if licinfo.Type  == 7 {
				linfo += "_" + strings.ReplaceAll(Prodinfo,",","_")
			}
			//授权文件缓存
			file = FileInfo{Name:"Zstack_" + linfo + "_License.tar.gz",}
			file.SetBody(gf.GetFile())
			licenseContent,_ := json.Marshal(file)
			utils.Bm.Put(req_md5Key,string(licenseContent),24*time.Hour)
			RetLicense(ctx,req_md5Key,file)
		}else {
			ctx.JSON(res)
		}
	}
}

func RetLicense(ctx iris.Context,req_md5Key string,file FileInfo)  {
	host := viper.GetString("RootUrl")
	if host == "" {
		host = "http://" + ctx.Host()
	}
	url := host + "/zstack/license/download/" + req_md5Key
	ctx.JSON(ResJson{Code:1, Message:"授权文件生成完毕！", Url:url, File:file.Name, Body:file.Body,})
}