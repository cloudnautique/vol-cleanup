package volumes

import (
	"strconv"

	log "github.com/Sirupsen/logrus"
	"github.com/fsouza/go-dockerclient"
)

const (
	dockerSocket             = "unix:///var/run/docker.sock"
	defaultDockerAPIVersion  = "1.19"
	fallbackDockerAPIVersion = "1.21"
)

type Volume struct {
	ID         string
	Attached   bool
	Path       string
	DockerPath string
}

type VolumesManager interface {
	GetVolumes(volumeDir string) error
	DeleteAllOrphans(noop bool) error
}

func NewVolumesManager() VolumesManager {
	client, defaultDocker := getDockerClient()
	if defaultDocker {
		return newVolumesPre19(client)
	}
	return newVolumes19(client)
}

func getDockerClient() (*docker.Client, bool) {
	client, _ := docker.NewVersionedClient(dockerSocket, defaultDockerAPIVersion)
	ver, err := client.Version()
	if err != nil {
		log.Fatalf("Could not get Docker version.")
	}

	if getFloat64(ver.Get("ApiVersion")) < getFloat64(fallbackDockerAPIVersion) {
		log.Debugf("Using API Version: %s", defaultDockerAPIVersion)
		return client, true
	}

	client, _ = docker.NewVersionedClient(dockerSocket, fallbackDockerAPIVersion)
	log.Debugf("Using API Version: %s", defaultDockerAPIVersion)
	return client, false
}

func getFloat64(val string) float64 {
	if val, err := strconv.ParseFloat(val, 64); err == nil {
		return val
	}
	log.Fatalf("Could not convert %s to float", val)
	return 0.0
}
