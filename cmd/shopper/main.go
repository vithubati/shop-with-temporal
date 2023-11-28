package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/pkg/errors"
	"github.com/vithubati/go-core/service"
	"github.com/vithubati/shop-with-temporal/pkg/config"
	"github.com/vithubati/shop-with-temporal/services/shopper"
	"log"
	"os"
)

var (
	flags    = flag.NewFlagSet("oauth", flag.ExitOnError)
	confFile = flags.String("conf", "", "path to the platform configuration file")
)

func main() {
	err := flags.Parse(os.Args[1:])
	if err != nil {
		log.Fatal(errors.Wrapf(err, "failed to create configs from file:%v", *confFile))
	}
	configs, err := config.NewConfig(*confFile)
	if err != nil {
		log.Fatal(errors.Wrapf(err, "failed to create configs from file:%v", *confFile))
	}
	ctx := context.Background()
	fmt.Printf("configs: %+v\n ", configs)
	server, cleanup, err := shopper.NewApp(ctx, configs)
	if err != nil {
		log.Fatal(err)
	}
	defer cleanup()

	if err := service.Run(ctx, server); err != nil {
		log.Fatal(err)
	}
}
