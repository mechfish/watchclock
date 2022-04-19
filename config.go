package watchclock

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
)

// A Config contains the options for a run of `watchclock`.
type Config struct {
	CacheTableName    string
	ClearCache        bool
	CreateCache       bool
	Debug             bool
	MinimumDays       uint
	Region            string
	RenewForDays      uint
	SkipCache         bool
	UpdateAllVersions bool

	session *session.Session
}

// Validate returns an error if the Config has missing or malformed values.
//
// It also sets the default values for parameters that have them.
func (c *Config) Validate() error {
	if c.CacheTableName == "" {
		c.CacheTableName = "watchclock-cache"
	}
	if c.Region == "" {
		c.Region = "us-east-1"
	}
	if c.MinimumDays == 0 {
		c.MinimumDays = 1
	}
	if c.RenewForDays == 0 {
		c.RenewForDays = 7
	}
	return nil
}

// GetSession returns an AWS session.
func (c *Config) GetSession() *session.Session {
	if c.session == nil {
		c.session = session.Must(session.NewSession(
			&aws.Config{
				Region: aws.String(c.Region),
			},
		))
	}
	return c.session
}
