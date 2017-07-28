package smtp

import (
	"github.com/facundomedica/hectane/version"
	"github.com/hectane/go-smtpsrv"

	"time"
)

// Config stores configuration data for the SMTP server.
type Config struct {
	Addr        string `json:"addr"`
	ReadTimeout int    `json:"read_timeout"`
}

// smtpsrvConfig converts the config into one suitable for smtpsrv.
func (c *Config) smtpsrvConfig() *smtpsrv.Config {
	return &smtpsrv.Config{
		Addr:        c.Addr,
		Banner:      "Hectane " + version.Version,
		ReadTimeout: time.Duration(c.ReadTimeout) * time.Second,
	}
}
