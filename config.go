package watchclock

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
