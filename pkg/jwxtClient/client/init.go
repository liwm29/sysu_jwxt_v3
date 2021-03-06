package client

import (
	"github.com/fatih/color"
	"github.com/liwm29/sysu_jwxt_v3/pkg/jwxtClient/global"
)

var DEFAULT_COOKIE_PATH string = "./cookie"
var DEFAULT_CAPTCHA_PATH string = "./captcha.jpg"

func printVersion() {
	color.Cyan("[ You are using SYSU-JWXT-VERSION3 ]")
	color.Cyan("[ working as a third-party package ]")
	switch global.MODE {
	case MODE_JWXT:
		color.HiYellow("MODE:MODE_JWXT  HOST:%s", global.HOST)
	case MODE_JWXT443:
		color.HiYellow("MODE:MODE_JWXT443  HOST:%s", global.HOST)
	}
}

func initClient(c *JwxtClient) {
	c.User = c.GetUserInfo()

}
