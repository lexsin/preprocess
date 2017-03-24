package mconfig

import (
	"os"

	"github.com/larspensjo/config"
)

var Conf *config.Config

func init() {
	conf, err := config.ReadDefault("conf/app.conf")
	if err != nil {
		os.Exit(1)
	}
	Conf = conf
	return
}
