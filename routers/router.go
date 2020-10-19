package routers

import (
	"fyoukuapi/controllers"
	"github.com/astaxie/beego"
)

func init() {
	beego.Router("/", &controllers.MainController{})

	beego.Include(&controllers.UserController{})
	beego.Include(&controllers.VideoController{})
}
