package main

import (
	"os"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/cloudnautique/vol-cleanup/volumes"
	"github.com/codegangsta/cli"
)

const (
	dockerVolumeDirectory = "/var/lib/docker/volumes"
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
	if _, err := os.Stat(dockerVolumeDirectory); os.IsNotExist(err) {
		log.Fatalf("Could not open Volume directory: %s", dockerVolumeDirectory)
	}

	for {
		log.Infof("Waking up to check for volumes...")
		vols := volumes.NewVolumesManager()
		err := vols.GetVolumes(dockerVolumeDirectory)
		if err != nil {
			log.Fatalf("Could not get volumes: %s", err)
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
