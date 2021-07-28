package watchclock

import (
	"flag"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"go.uber.org/zap"
)

// A Config contains the options for a run of `watchclock`.
type Config struct {
	CacheTableName    string
	ClearCache        bool
	Debug             bool
	MinimumDays       uint
	Region            string
	RenewForDays      uint
	SkipCache         bool
	UpdateAllVersions bool

	session *session.Session
	logger  *zap.SugaredLogger
}

// A Logger provides functions for structured logging.
//
// This is a wrapper around a zap logger; see the go.uber.org/zap docs
type Logger interface {
	Infow(string, ...interface{})
	Infof(string, ...interface{})
	Debugw(string, ...interface{})
	Debugf(string, ...interface{})
}

// Logger returns a Logger for the application.
func (c *Config) Logger() Logger {
	if c.logger == nil {
		var err error
		var logger *zap.Logger
		if c.Debug {
			logger, err = zap.NewDevelopment()
		} else {
			logger, err = zap.NewProduction()
		}
		if err != nil {
			panic(err)
		}
		c.logger = logger.Sugar()
	}
	return c.logger
}

// Cleanup cleans up the Config by e.g. flushing the logs.
func (c *Config) Cleanup() {
	// ignoring logger.Sync() errors because attempting to sync stdout/stderr always returns EINVAL;
	// see e.g. https://github.com/uber-go/zap/issues/772
	c.logger.Sync()
}

// DeclareFlags adds the given command's CLI flags to the given flagSet.
func (c *Config) DeclareFlags(commandName string, flagSet *flag.FlagSet) {
	flagSet.StringVar(&c.CacheTableName, "cache-table", "watchclock-cache", "Name of the DynamoDB table to use for the cache.")
	flagSet.BoolVar(&c.ClearCache, "clear-cache", false, "Clear and rebuild the Object Lock cache.")
	flagSet.BoolVar(&c.Debug, "debug", false, "Log debug messages.")
	flagSet.UintVar(&c.MinimumDays, "minimum-days", 1, "Renew all locks that will expire within this many days.")
	flagSet.StringVar(&c.Region, "region", "us-east-1", "Name of the AWS region containing the bucket.")
	flagSet.UintVar(&c.RenewForDays, "renew-for", 7, "Reset object lock expiration to N days from now.")
	flagSet.BoolVar(&c.UpdateAllVersions, "all-versions", false, "Update locks for every version of each S3 object.")
	flagSet.BoolVar(&c.SkipCache, "no-cache", false, "Do not use the Object Lock cache.")
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
