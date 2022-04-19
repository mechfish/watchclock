package watchclock

import (
	"context"

	log "github.com/sirupsen/logrus"
)

// Renew implements the `watchclock renew` command.
func Renew(ctx context.Context, config *Config, args []string) error {
	log.WithFields(log.Fields{"object": args[0]}).Info("Renewing object locks")
	return nil
}
