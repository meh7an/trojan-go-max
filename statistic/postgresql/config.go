package postgresql

import (
	"github.com/p4gefau1t/trojan-go/config"
)

type PostgreSQLConfig struct {
	Enabled    bool   `json:"enabled" yaml:"enabled"`
	ServerHost string `json:"server_addr" yaml:"server-addr"`
	ServerPort int    `json:"server_port" yaml:"server-port"`
	Database   string `json:"database" yaml:"database"`
	Username   string `json:"username" yaml:"username"`
	Password   string `json:"password" yaml:"password"`
	CheckRate  int    `json:"check_rate" yaml:"check-rate"`
	SSLMode    string `json:"ssl_mode" yaml:"ssl-mode"`
}

type Config struct {
	PostgreSQL PostgreSQLConfig `json:"postgresql" yaml:"postgresql"`
}

func init() {
	config.RegisterConfigCreator(Name, func() interface{} {
		return &Config{
			PostgreSQL: PostgreSQLConfig{
				ServerPort: 5432,
				CheckRate:  30,
				SSLMode:    "disable",
			},
		}
	})
}