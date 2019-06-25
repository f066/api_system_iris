package config

import "github.com/kataras/iris"

type HttpErrorModel struct {
	Code 	int			`json:"code"`
	IP		string		`json:"ip"`
	Path	string		`json:"path"`
	Status	string		`json:"status"`
	Message	string		`json:"message"`
}

func GetInfo(ctx iris.Context,status string,message string) (HttpErrorModel) {
	hem := HttpErrorModel{
		Code : ctx.GetStatusCode(),
		IP : ctx.RemoteAddr(),
		Path : ctx.Path(),
		Status : status,
		Message : message,
	}
	return hem
}
