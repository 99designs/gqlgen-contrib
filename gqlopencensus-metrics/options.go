package metrics

import (
	"os"
)

type (
	// Option for this metrics collector (none implemented at the moment)
	Option func(*config)

	config struct {
		host          string
		fieldsEnabled bool
	}
)

func defaultCollector() *Collector {
	host, _ := os.Hostname()
	return &Collector{
		config: &config{
			host:          host,
			fieldsEnabled: true,
		},
	}
}

// Host determines the host tag. By default this is the OS hostname
func Host(hostname string) Option {
	return func(c *config) {
		c.host = hostname
	}
}

// FieldsEnabled controls whether metrics at the field level are enabled (this is enabled by default)
func FieldsEnabled(enabled bool) Option {
	return func(c *config) {
		c.fieldsEnabled = enabled
	}
}
