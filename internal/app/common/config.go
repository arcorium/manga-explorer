package common

import (
	"manga-explorer/internal/infrastructure/file"
	"manga-explorer/internal/util"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/spf13/viper"
)

type Config struct {
	JWTSecretKey   string `mapstructure:"JWT_SECRET_KEY"`
	JWTSigningType string `mapstructure:"JWT_SIGNING_TYPE"`
	signingMethod  jwt.SigningMethod

	RefreshTokenDuration time.Duration `mapstructure:"REFRESH_TOKEN_DURATION"`
	AccessTokenDuration  time.Duration `mapstructure:"ACCESS_TOKEN_DURATION"`

	// Mailing Related
	SMTPHost string `mapstructure:"SMTP_HOST"`
	SMTPPort uint16 `mapstructure:"SMTP_PORT"`
	SMTPUser string `mapstructure:"SMTP_USER"`
	SMTPPass string `mapstructure:"SMTP_PASS"`

	// Server Related
	Ip   string `mapstructure:"IP"`
	Port string `mapstructure:"PORT"`

	// Real Database
	DbProtocol string `mapstructure:"DB_PROTOCOL"`

	DbUser     string `mapstructure:"DB_USER"`
	DbPassword string `mapstructure:"DB_PASSWORD"`
	DbHost     string `mapstructure:"DB_HOST"`
	DbPort     uint16 `mapstructure:"DB_PORT"`
	DbName     string `mapstructure:"DB_NAME"`
	DbParam    string `mapstructure:"DB_PARAM"`
}

var conf *Config

func LoadConfig(name string, path ...string) (*Config, error) {
	viper.SetDefault("IP", "localhost")
	viper.SetDefault("PORT", 9999)
	viper.SetDefault("ACCESS_TOKEN_DURATION", time.Minute*5)
	viper.SetDefault("REFRESH_TOKEN_DURATION", time.Hour*24*30)
	viper.SetDefault("JWT_SIGNING_TYPE", "HS256")
	//
	viper.SetConfigName(name)
	viper.SetConfigType("env")
	if len(path) == 0 {
		viper.AddConfigPath(".")
	} else {
		viper.AddConfigPath(path[0])
	}
	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}
	if err := viper.Unmarshal(&conf); err != nil {
		return nil, err
	}

	conf.signingMethod = jwt.GetSigningMethod(conf.JWTSigningType)
	file.SetHostName(conf.Endpoint()) // Set default hostname
	return conf, nil
}

func (c *Config) Endpoint() string {
	return c.Ip + ":" + c.Port
}

func (c *Config) SigningMethod() jwt.SigningMethod {
	return c.signingMethod
}

func (c *Config) DatabaseDSN() string {
	password := ""
	if len(c.DbPassword) != 0 {
		password = ":" + c.DbPassword
	}
	dsn := c.DbProtocol + "://" + c.DbUser + password + "@" + c.DbHost + ":" +
		strconv.Itoa(int(c.DbPort)) + "/" + c.DbName
	if len(c.DbParam) == 0 {
		return dsn
	}
	return dsn + "?" + c.DbParam
}

// GetConfig retrieving singleton config object either created by LoadConfig or it will automatically use app.env file from current work directory
func GetConfig() *Config {
	if conf == nil {
		return util.DropError(LoadConfig("app"))
	}
	return conf
}
