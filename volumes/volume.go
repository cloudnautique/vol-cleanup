package volumes

import (
	"io/ioutil"
	"os"
	"path"

	log "github.com/Sirupsen/logrus"
	"github.com/fsouza/go-dockerclient"
)

const (
	dockerSocket     = "unix:///var/run/docker.sock"
	dockerAPIVersion = "1.19"
)

type Volume struct {
	ID         string
	Attached   bool
	Path       string
	DockerPath string
}

type Volumes map[string]Volume

func getDockerClient() *docker.Client {
	client, _ := docker.NewVersionedClient(dockerSocket, dockerAPIVersion)
	return client
}

func (v Volumes) GetVolumes(volumeDir string) error {
	client := getDockerClient()
	info, err := client.Info()
	if err != nil {
		return err
	}
	dockerPfx := info.Get("DockerRootDir")

	// Get all VFS Docker volumes from Disk.
	files, err := ioutil.ReadDir(volumeDir)
	if err != nil {
		return err
	}

	for _, f := range files {
		log.Infof("Found volume: %v", f.Name())
		filePath := path.Join(volumeDir, f.Name())
		dockerPath := path.Join(dockerPfx, "volumes", f.Name(), "_data")

		volume := &Volume{
			ID:         f.Name(),
			Path:       filePath,
			DockerPath: dockerPath,
		}

		v[volume.DockerPath] = *volume
		log.Debugf("Volume path: %v", volume.Path)
	}

	err = v.setAttachedVolumes()
	if err != nil {
		return err
	}

	return nil
}

func (v Volumes) DeleteAllOrphans(noop bool) error {
	message := "NOOP: Deleting volume: "
	if noop == false {
		message = "Deleting volume: "
	}

	for key, volume := range v {
		if volume.Attached == false {
			log.Infof("%v: %v", message, key)
			if noop == false {
				err := rmVolume(volume.Path)
				if err != nil {
					log.Errorf("%v", err)
				}
			}
		}
	}
	return nil
}

func (v Volumes) DeleteVolume(id string) error {
	for _, volume := range v {
		if volume.ID == id && volume.Attached == false {
			err := rmVolume(volume.Path)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func rmVolume(volPath string) error {
	err := os.RemoveAll(volPath)
	return err
}

func (v Volumes) setAttachedVolumes() error {
	client := getDockerClient()

	existingContainers, err := client.ListContainers(
		docker.ListContainersOptions{
			All: true,
		})

	if err != nil {
		return err
	}

	// loop over existing containers
	for _, container := range existingContainers {
		containerInfo, _ := client.InspectContainer(container.ID)
		for _, val := range containerInfo.Volumes {
			if _, exists := v[val]; exists {
				volume := v[val]
				volume.Attached = true
				v[val] = volume
			}
		}
	}

	return nil
}
