package routers

import (
	"github.com/HelloWorldZQ/quintinblog/controllers"
	"github.com/astaxie/beego"
)

func init() {
    beego.Router("/", &controllers.MainController{})
}
