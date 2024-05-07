package common

import (
  "fmt"
  "github.com/caarlos0/env/v10"
  "github.com/joho/godotenv"
  "log"
  "strconv"
  "time"

  "github.com/golang-jwt/jwt"
)

type Config struct {
  JWTSecretKey   string `env:"JWT_SECRET_KEY,notEmpty"`
  JWTSigningType string `env:"JWT_SIGNING_TYPE" envDefault:"HS256"`
  signingMethod  jwt.SigningMethod

  RefreshTokenDuration      time.Duration `env:"REFRESH_TOKEN_DURATION" envDefault:"720h"`
  AccessTokenDuration       time.Duration `env:"ACCESS_TOKEN_DURATION" envDefault:"10m"`
  VerificationTokenDuration time.Duration `env:"VERIF_TOKEN_DURATION" envDefault:"24h"`

  // Mailing Related
  SMTPHost string `env:"SMTP_HOST,notEmpty"`
  SMTPPort uint16 `env:"SMTP_PORT,notEmpty"`
  SMTPUser string `env:"SMTP_USER,notEmpty"`
  SMTPPass string `env:"SMTP_PASS,notEmpty"`

  // Server Related
  Ip             string   `env:"LISTEN_IP" envDefault:"0.0.0.0"`
  Port           string   `env:"LISTEN_PORT" envDefault:"9999"`
  Dns            string   `env:"DNS"` // use it when using domain
  TrustedProxies []string `env:"TRUSTED_PROXIES" envSeparator:","`

  // Database
  DbProtocol string `env:"DB_PROTOCOL,notEmpty"`
  DbUser     string `env:"DB_USER,notEmpty"`
  DbPassword string `env:"DB_PASSWORD,notEmpty"`
  DbHost     string `env:"DB_HOST,notEmpty"`
  DbPort     uint16 `env:"DB_PORT,notEmpty"`
  DbName     string `env:"DB_NAME,notEmpty"`
  DbParam    string `env:"DB_PARAM"`
}

var conf = new(Config)

func LoadConfig(filenames ...string) (*Config, error) {
  err := godotenv.Load(filenames...)
  if err != nil {
    log.Println(".env file doesn't found, will using ENVIRONMENT VARIABLE as config")
  }

  err = env.Parse(conf)
  if err != nil {
    return nil, err
  }

  conf.signingMethod = jwt.GetSigningMethod(conf.JWTSigningType)
  return conf, nil
}

func (c *Config) Endpoint() string {
  return fmt.Sprintf("%s:%s", c.Ip, c.Port)
}

func (c *Config) DNS() string {
  if len(c.Dns) == 0 {
    return c.Endpoint()
  }
  return c.Dns
}

func (c *Config) ApiDNS(version uint) string {
  return fmt.Sprintf("%s/api/v%d", c.DNS(), version)
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

//// GetConfig retrieving singleton config object either created by LoadConfig or it will automatically use app.env file from current work directory
//func GetConfig() *Config {
//	if conf == nil {
//		return util.DropError(LoadConfig("app"))
//	}
//	return conf
//}
