// @APIVersion 1.0.0
// @Title beego Test API
// @Description beego has a very cool tools to autogenerate documents for your API
// @Contact astaxie@gmail.com
// @TermsOfServiceUrl http://beego.me/
// @License Apache 2.0
// @LicenseUrl http://www.apache.org/licenses/LICENSE-2.0.html
package routers

import (
	"RegistrationCenter/controllers"

	"github.com/astaxie/beego"
)

func init() {
	ns := beego.NewNamespace("/api",
		beego.NSRouter("/discover/?:id", &controllers.DiscoverController{}, "get:DiscoverById"),
		beego.NSRouter("/discover/services", &controllers.DiscoverController{}, "get:DiscoverServices"),

		beego.NSRouter("/register/config", &controllers.RegisterController{}, "get:RegisterConfig"),
		beego.NSRouter("/register/service", &controllers.RegisterController{}, "post:RegisterService"),

		beego.NSRouter("/deregister/?:id", &controllers.RegisterController{}, "get:DeregisterService"),

		beego.NSRouter("/health/check", &controllers.HealthCheckController{}, "get:HealthCheck"),
	)
	beego.AddNamespace(ns)
}
