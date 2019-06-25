# api_system_iris

后端api接口开源，项目使用Go语言编写，框架采用 [iris](https://github.com/kataras/iris)，本仓库所有代码仅供学习交流

## 开源Api功能
* [Zstack授权文件](#Zstack授权文件)
    * [生成](#生成)
    * [下载](#下载)


## Zstack授权文件

主路径：{HOST}/zstack/license

### 生成

请求路径：{主路径}/generate/{授权类型}
请求方法：POST
支持的授权类型：

    Trial
	Paid
	OEM
	Free
	TrialExt
	Hybrid
 	AddOn
	HybridTrialExt

请求参数：

    User            用户名      建议为[字母，数字，中文，下划线]组成
	RequestCode     请求码      从系统管理后台获取
	Expiretime      过期日期    授权到期日期，格式为 YYYY-MM-DD
	Hostnum         主机数      授权主机数（如果授权类型为AddOn，该字段无效）
	Prodinfo        附加信息    如果授权类型为AddOn，仅支持填写为[vmware,project-management,disaster-recovery,v2v,baremetal],填写多个时使用半角符号"," 分隔。
	Cpunum          Cpu数       授权Cpu数（如果授权类型为AddOn，该字段无效）

请求返回：
* 成功：

        {
            "code": 1, //状态码
            "message": "授权文件生成完毕！", //消息
            "url": "http://host/zstack/license/download/162e64550869b305c257f436181d81bf", //本次生成的文件缓存下载地址
            "file": "zstack_license.tar.gz", //文件名
            "body":"XXXXXX..."  //base64编码后的文件内容
        }

* 失败：

        {
            "code": -1, //状态码
            "Error": "模块不存在", //错误原因
            "Tips": "请输入正确的模块名称"  //建议
        }

* 其他错误：

        {
            "code": 404,
            "ip": "127.0.0.1",
            "path": "/zstack/license/generate/xxx",
            "status": "error",
            "message": "404 Not Found"
        }


### 下载

请求路径：{主路径}/download/{MD5特征码}
请求方法：GET
MD5特征码：

    由生成服务返回，有效期24小时（服务重启或关闭则立即失效）


请求参数：无
请求返回：
* 成功：直接下载文件

* 失败：

        {
            "code": -1,
            "Error": "文件不存在",
            "Tips": "请使用生成服务生成先！"
        }