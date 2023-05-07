// Package env get env & app config, all the public field must after init()
// finished and flag.Parse().
package env

import (
	"os"
	"path/filepath"
	"time"

	"github.com/GitHub121380/golib/utils"
	"github.com/gin-gonic/gin"
)

// deploy env.
const (
	DEPLOY_ENV_TEST = "test"
	DEPLOY_ENV_TIPS = "tips"
	DEPLOY_ENV_PROD = "prod"
)

// 环境变量相关
const (
	AppIDC      = "APP_IDC"
	AppRootPath = "APP_ROOT_PATH"
)

const DefaultRootPath = "/home/homework/"

var (
	// Hostname machine hostname.
	Hostname string
	// DeployEnv deploy env where app at.
	DeployEnv string
	//Clouds to which machines belong
	CloudType string
	//current IDC Name
	IDC string
	//local IP
	LocalIP string
)

// app configuration
var (
	// AppName is global unique application name, register by service tree
	AppName string
	// 运行模式 debug|release
	RunMode string

	// app root path
	rootPath string
)

func init() {
	var err error
	if Hostname, err = os.Hostname(); err != nil || Hostname == "" {
		Hostname = os.Getenv("HOSTNAME")
	}
	// 初始化本地IP
	LocalIP = utils.GetLocalIp()
	// 初始化本地IDC名称
	IDC, CloudType = queryIDC()
	println(time.Now().String(), " : current idc:", IDC, ", cloud:", CloudType)

	if IDC == IDC_TEST {
		DeployEnv = DEPLOY_ENV_TEST
		RunMode = gin.DebugMode
	} else {
		DeployEnv = DEPLOY_ENV_PROD
		RunMode = gin.ReleaseMode
	}
}

// 虚拟机上用户通过调用SetAppName()实现应用路径的确定
func SetAppName(appName string) {
	if appName == "" {
		panic("请使用 env.SetAppName(模块名) 指定appName，一旦创建不能再修改")
	}
	AppName = appName
	setAppPath()
}

func GetAppName() string {
	return AppName
}

func setAppPath() {
	if IDC == IDC_TEST {
		// test环境支持用户通过环境变量修改rootPath，默认./
		if r := os.Getenv(AppRootPath); r != "" {
			SetRootPath(filepath.Join(r, AppName))
		} else {
			SetRootPath("./")
		}
	} else {
		// 虚拟机线上环境的日志、配置路径使用绝对路径（固定）
		SetRootPath(DefaultRootPath)
	}
	println("load conf: ", GetConfDirPath())
	println("log path: ", GetLogDirPath())
}

// SetRootPath 设置应用的根目录
func SetRootPath(r string) {
	rootPath = r
	println("SetRootPath: ", GetConfDirPath())
}

// RootPath 返回应用的根目录
func GetRootPath() string {
	if rootPath != "" {
		return rootPath
	} else {
		return DefaultRootPath
	}
}

// GetConfDirPath 返回配置文件目录绝对地址
func GetConfDirPath() string {
	if IDC == IDC_TEST {
		return filepath.Join(GetRootPath(), "conf")
	} else {
		return filepath.Join(GetRootPath(), "goapp", AppName, "conf")
	}

}

// LogRootPath 返回log目录的绝对地址
func GetLogDirPath() string {
	if IDC == IDC_TEST {
		return filepath.Join(GetRootPath(), "log")
	} else {
		return filepath.Join(GetRootPath(), "clog", "go", AppName)
	}
}
