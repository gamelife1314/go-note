package config

import (
	"github.com/kataras/iris"
)

var Configuration iris.Configuration

func init() {
	Configuration = iris.TOML("./config/config.tml")
}