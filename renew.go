package watchclock

import (
	log "github.com/sirupsen/logrus"
)

// Renew implements the `watchclock renew` command.
func Renew(config *Config, args []string) error {
	log.WithFields(log.Fields{"object": args[0]}).Info("Renewing object locks")
	return nil
}
