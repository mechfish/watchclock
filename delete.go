package watchclock

// Delete implements the `watchclock delete` command.
func Delete(config *Config, args []string) error {
	config.Logger().Infof("Deleting object %s", args[0])
	return nil
}

// Undelete implements the `watchclock undelete` command.
func UnDelete(config *Config, args []string) error {
	config.Logger().Infof("Canceling deletion of object %s", args[0])
	return nil
}
