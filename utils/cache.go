package utils

import "github.com/astaxie/beego/cache"

var (
	Bm = New()
)

func New() (bm cache.Cache) {
	bm, _ = cache.NewCache("memory", `{"interval":60}`)
	return
}