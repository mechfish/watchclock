package watchclock

// Renew implements the `watchclock renew` command.
func Renew(config *Config, args []string) error {
	config.Logger().Infof("Renewing object locks for %s", args[0])
	return nil
}
