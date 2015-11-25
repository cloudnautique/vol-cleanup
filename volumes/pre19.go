package volumes

import (
	"io/ioutil"
	"os"
	"path"

	log "github.com/Sirupsen/logrus"
	"github.com/fsouza/go-dockerclient"
)

type volumesPre19 struct {
	volumes      map[string]Volume
	dockerClient *docker.Client
}

func newVolumesPre19(client *docker.Client) *volumesPre19 {
	vol := &volumesPre19{
		volumes:      make(map[string]Volume),
		dockerClient: client,
	}
	var _ VolumesManager = vol
	return vol
}

func (v *volumesPre19) GetVolumes(volumeDir string) error {
	info, err := v.dockerClient.Info()
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

		v.volumes[volume.DockerPath] = *volume
		log.Debugf("Volume path: %v", volume.Path)
	}

	err = v.setAttachedVolumes()
	if err != nil {
		return err
	}

	return nil
}

func (v volumesPre19) DeleteAllOrphans(noop bool) error {
	message := "NOOP: Deleting volume: "
	if noop == false {
		message = "Deleting volume: "
	}

	for key, volume := range v.volumes {
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

func (v volumesPre19) DeleteVolume(id string) error {
	for _, volume := range v.volumes {
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
	return os.RemoveAll(volPath)
}

func (v volumesPre19) setAttachedVolumes() error {
	existingContainers, err := v.dockerClient.ListContainers(
		docker.ListContainersOptions{
			All: true,
		})

	if err != nil {
		return err
	}

	// loop over existing containers
	for _, container := range existingContainers {
		containerInfo, _ := v.dockerClient.InspectContainer(container.ID)
		for _, val := range containerInfo.Volumes {
			if _, exists := v.volumes[val]; exists {
				volume := v.volumes[val]
				volume.Attached = true
				v.volumes[val] = volume
			}
		}
	}

	return nil
}
