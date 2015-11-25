package volumes

import (
	log "github.com/Sirupsen/logrus"
	"github.com/fsouza/go-dockerclient"
)

type volumes19 struct {
	volumes      map[string]Volume
	dockerClient *docker.Client
}

func newVolumes19(client *docker.Client) *volumes19 {
	vol := &volumes19{
		volumes:      make(map[string]Volume),
		dockerClient: client,
	}
	var _ VolumesManager = vol
	return vol
}

func (v volumes19) GetVolumes(volumeDir string) error {
	volumes, err := v.dockerClient.ListVolumes(docker.ListVolumesOptions{})
	if err != nil {
		return err
	}

	for _, volume := range volumes {
		log.Infof("Found volume: %v", volume.Name)
		vol := &Volume{
			ID:         volume.Name,
			Path:       volume.Mountpoint,
			DockerPath: volume.Mountpoint,
			Attached:   true,
		}
		v.volumes[vol.ID] = *vol
	}

	err = v.setDetachedVolumes()
	if err != nil {
		return err
	}

	return nil
}

func (v volumes19) DeleteAllOrphans(noop bool) error {
	message := "NOOP: Deleting volume: "
	if noop == false {
		message = "Deleting volume: "
	}

	opts := docker.ListVolumesOptions{
		Filters: map[string][]string{
			"dangling": []string{"true"},
		},
	}

	danglingVolumes, err := v.dockerClient.ListVolumes(opts)
	if err != nil {
		return err
	}

	for _, volume := range danglingVolumes {
		log.Infof("%v: %v", message, volume.Name)
		if volume.Driver == "local" && noop == false {
			err = v.dockerClient.RemoveVolume(volume.Name)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func (v volumes19) DeleteVolume(id string) error {
	return v.dockerClient.RemoveVolume(id)
}

func (v volumes19) setDetachedVolumes() error {
	opts := docker.ListVolumesOptions{
		Filters: map[string][]string{
			"dangling": []string{"true"},
		},
	}

	danglingVolumes, err := v.dockerClient.ListVolumes(opts)
	if err != nil {
		return err
	}

	for _, volume := range danglingVolumes {
		vol := v.volumes[volume.Name]
		vol.Attached = false
		v.volumes[volume.Name] = vol
	}

	return nil
}
