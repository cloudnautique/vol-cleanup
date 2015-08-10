package main

import (
	"os"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/cloudnautique/go-vol/volumes"
	"github.com/codegangsta/cli"
)

func main() {
	app := cli.NewApp()
	app.Name = "vol-cleanup"
	app.Usage = "Clean up orphaned Docker volumes"

	app.Flags = []cli.Flag{
		cli.BoolFlag{
			Name:  "noop, n",
			Usage: "Run in a noop mode and log to screen",
		},
		cli.IntFlag{
			Name:  "interval, i",
			Value: 5,
			Usage: "Number of minutes between runs",
		},
	}

	app.Action = func(c *cli.Context) {
		log.Debugf("Noop %v", c.Bool("noop"))
		Start(c.Int("interval"), c.Bool("noop"))
	}

	app.Run(os.Args)
}

func Start(interval int, noop bool) error {
	for {
		log.Infof("Waking up to check for volumes...")
		vols := &volumes.Volumes{}
		err := vols.GetVolumes("/var/lib/docker/volumes")
		if err != nil {
			log.Fatalf("Error Getting volumes.", err)
			return err
		}

		err = vols.DeleteAllOrphans(noop)
		if err != nil {
			log.Fatalf("Error deleting orphaned vols.", err)
		}

		time.Sleep(time.Duration(interval) * time.Minute)
	}

	return nil
}
