package utils

import (
	"os"
	"path/filepath"
	"runtime"
)

var (
	VERSION   string = "0.1"
	Revision  string 
	Branch    string = "-Release"
	BuildUser string = "Deny"
	BuildDate string
	GoVersion = runtime.Version()
)

func GetCurrentPath() string {
	filename,err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		filename ="./"
	}
	return filename
}

func GetAbsPath(path string) string {
	filename,err := filepath.Abs(path)
	if err != nil {
		filename = path
	}
	return filename
}

func GetBuildInfo() (info string) {
	if Revision != "" {
		Revision = "_" + Revision
	}
	info = "VERSION：" + VERSION + Branch + Revision + "\t\tGoVersion：" + GoVersion
	if BuildUser != "" {
		info += "\nBuildUser：" + BuildUser
	}
	if BuildDate != "" {
		info += "\t\tBuildTime：" + BuildDate
	}
	return
}

func FileExists(path string) bool {
	isexists,isdir := PathExists(path)
	if isdir {
		return false
	}
	return isexists
}

func PathExists(path string) (isExists,isDir bool) {
	pathinfo,err := os.Stat(path)
	if err == nil{
		return true,pathinfo.IsDir()
	}
	return false,false
}