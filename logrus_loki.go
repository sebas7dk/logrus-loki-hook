package hook

import (
	"fmt"
	"time"

	"github.com/afiskon/promtail-client/promtail"
	"github.com/sirupsen/logrus"
)

var supportedLevels = []logrus.Level{logrus.TraceLevel, logrus.DebugLevel, logrus.InfoLevel, logrus.WarnLevel, logrus.ErrorLevel}

// Config defines configuration for hook for Loki
type Config struct {
	URL                string
	Labels             string
	BatchWait          time.Duration
	BatchEntriesNumber int
}

func (c *Config) setDefault() {
	if c.URL == "" {
		c.URL = "http://localhost:3100/api/prom/push"
	}
	if c.Labels == "" {
		c.Labels = "{source=\"" + "test" + "\",job=\"" + "job" + "\"}"
	}
	if c.BatchWait == time.Second {
		c.BatchWait = 5 * time.Second
	}
	if c.BatchEntriesNumber == 0 {
		c.BatchEntriesNumber = 10000
	}

}

type Hook struct {
	client promtail.Client
}

// NewHook creates a new hook for Loki
func NewHook(c *Config) (*Hook, error) {
	if c == nil {
		c = &Config{}
	}
	c.setDefault()
	conf := promtail.ClientConfig{
		PushURL:            c.URL,
		Labels:             c.Labels,
		BatchWait:          c.BatchWait,
		BatchEntriesNumber: c.BatchEntriesNumber,
		SendLevel:          promtail.DEBUG,
		PrintLevel:         promtail.DISABLE,
	}
	loki, err := promtail.NewClientJson(conf)
	if err != nil {
		return nil, fmt.Errorf("unable to init promtail client: %v", err)
	}
	return &Hook{
		client: loki,
	}, nil
}

// Fire implements interface for logrus
func (hook *Hook) Fire(entry *logrus.Entry) error {

	switch entry.Level {
	case logrus.DebugLevel:
		hook.client.Debugf(entry.Message)
	case logrus.InfoLevel:
		hook.client.Infof(entry.Message)
	case logrus.WarnLevel:
		hook.client.Warnf(entry.Message)
	case logrus.ErrorLevel:
		hook.client.Errorf(entry.Message)
	case logrus.TraceLevel:
		hook.client.Debugf(entry.Message)
	default:
		return fmt.Errorf("unknown log level")
	}

	return nil
}

// Levels retruns supported levels
func (hook *Hook) Levels() []logrus.Level {
	return supportedLevels
}
