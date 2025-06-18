package pkg

import (
	"krathub/internal/conf"
)

var AppConf *conf.App

func SetAppConf(bootstrapApp *conf.App) {
	AppConf = bootstrapApp
}
