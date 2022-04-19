// The watchclock command updates Object Lock retention periods for
// objects in an S3 bucket.
package main

import (
	"context"
	"flag"
	"fmt"
	"os"

	"github.com/mechfish/watchclock"
)

func main() {
	conf := watchclock.Config{}

	flag.StringVar(&conf.CacheTableName, "cache-table-name", "watchclock-cache", "Name of the DynamoDB table to use for the cache.")
	flag.BoolVar(&conf.ClearCache, "clear-cache", false, "Clear and rebuild the Object Lock cache.")
	flag.BoolVar(&conf.Debug, "debug", false, "Log debug messages.")
	flag.UintVar(&conf.MinimumDays, "minimum-days", 1, "Renew all locks that will expire within this many days.")
	flag.StringVar(&conf.Region, "region", "us-east-1", "Name of the AWS region containing the bucket.")
	flag.UintVar(&conf.RenewForDays, "renew-for", 7, "Reset object lock expiration to N days from now.")
	flag.BoolVar(&conf.UpdateAllVersions, "all-versions", false, "Update locks for every version of each S3 object.")
	flag.BoolVar(&conf.SkipCache, "no-cache", false, "Do not use the Object Lock cache.")
	flag.BoolVar(&conf.CreateCache, "use-cache", false, "Create and use a cache table in DynamoDB.")

	flag.Usage = func() {
		fmt.Fprintf(flag.CommandLine.Output(), "Usage: watchclock [arguments]\n\nOptions:\n")
		flag.PrintDefaults()
	}

	flag.Parse()

	ctx := context.TODO()

	err := conf.Validate()
	if err == nil {
		err = watchclock.Renew(ctx, &conf, flag.Args())
	}
	if err != nil {
		fmt.Fprintf(flag.CommandLine.Output(), "%s\n", err.Error())
		os.Exit(1)
	}
}
