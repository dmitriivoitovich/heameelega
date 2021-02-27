package config

import (
	"encoding/base64"
	"strings"

	"github.com/spf13/viper"
)

const (
	configName = "config"
	configType = "yaml"
	envPrefix  = "VIPER"
)

type DBConf struct {
	Host     string
	User     string
	Password string
	DBName   string
	Port     uint32
}

type TLSConf struct {
	Enabled bool
	Cert    []byte
	Key     []byte
}

func AppHost() string {
	return getStr("app.host")
}

func AppTLS() TLSConf {
	cert, _ := base64.StdEncoding.DecodeString(getStr("app.tls.cert"))
	key, _ := base64.StdEncoding.DecodeString(getStr("app.tls.key"))

	return TLSConf{
		Enabled: getBool("app.tls.enabled"),
		Cert:    cert,
		Key:     key,
	}
}

func DBConfig() DBConf {
	return DBConf{
		Host:     getStr("db.host"),
		Port:     getUInt32("db.port"),
		User:     getStr("db.user"),
		Password: getStr("db.password"),
		DBName:   getStr("db.name"),
	}
}

func Read() error {
	viper.SetConfigName(configName)
	viper.SetConfigType(configType)
	viper.AddConfigPath(".")
	viper.SetEnvPrefix(envPrefix)
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AutomaticEnv()
	viper.AllowEmptyEnv(true)

	return viper.ReadInConfig()
}

func getStr(key string) string {
	return viper.GetString(key)
}

func getBool(key string) bool {
	return viper.GetBool(key)
}

func getUInt32(key string) uint32 {
	return viper.GetUint32(key)
}
