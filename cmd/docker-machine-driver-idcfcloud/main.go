package main

import (
	"os"
	"path"

	"github.com/dhayakawa-idcf/docker-machine-driver-idcfcloud/pkg/driver"
	"github.com/dhayakawa-idcf/docker-machine-driver-idcfcloud/pkg/version"
	"github.com/docker/machine/libmachine/drivers/plugin"
	"github.com/urfave/cli"
)

func main() {
	app := cli.NewApp()
	app.Name = path.Base(os.Args[0])
	app.Usage = "Docker Machine plugin binary for IDCF Cloud. Please use it through the main 'docker-machine' binary."
	app.Version = version.Version
	app.Action = func(c *cli.Context) {
		plugin.RegisterDriver(driver.NewDriver())
	}
	app.Run(os.Args)
}
