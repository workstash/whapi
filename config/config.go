package config

import (
	"github.com/jinzhu/configor"
)

//LoggerFile the local where logs will be saved
var LoggerFile string

// API stores api configuration
type API struct {
	Port        string `default:"8080"`
	CorsOrigin  string `default:"*"`
	SessionPath string `required:"true"`

	GenerateQrCode bool   `default:"false"`
	QrCodeQuality  string `default:"medium"`
	QrCodeSize     uint   `default:"256"`

	EncodeBase64 bool `default:"false"`
}

// Client stores client configuration
type Client struct {
	LongName  string `default:"Whapi"`
	ShortName string `default:"Whapi"`
	Version   string `default:"2.2121.6"`
}

// Configuration implements all configurations
type Configuration struct {
	LoggerFile string `required:"true"`
	API        API    `required:"true"`
	Client     Client `required:"true"`
}

//Main unmarshal the configurations
var Main = (func() Configuration {
	var conf Configuration
	// if err := configor.Load(&conf, os.Getenv("CONFIG_PATH")); err != nil {
	if err := configor.Load(&conf, "./config/config.json"); err != nil {
		panic(err.Error())
	}
	return conf
})()
