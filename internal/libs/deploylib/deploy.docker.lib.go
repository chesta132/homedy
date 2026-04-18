package deploylib

import (
	"context"
	"fmt"
	"homedy/config"
	"os"

	"github.com/compose-spec/compose-go/v2/types"
	"github.com/docker/compose/v5/pkg/api"
)

func DockerProjectName(id string) string {
	return fmt.Sprintf("deploy-compose-%s", id)
}

func DockerVolumeName(id string) string {
	return fmt.Sprintf("deploy-compose-%s-data", id)
}

func LoadDockerCompose(ctx context.Context, service api.Compose, id, composeContent string) (*types.Project, error) {
	// create temp bcs load project only accept path
	file, err := os.CreateTemp(config.TEMP_DIR, "deploy-compose-*")
	if err != nil {
		return nil, err
	}
	path := file.Name()
	defer os.Remove(path)

	if _, err = file.WriteString(composeContent); err != nil {
		return nil, err
	}

	// FIXME: load project include checking env file which is temp dir haven't
	// it'll return error because of it
	project, err := service.LoadProject(ctx, api.ProjectLoadOptions{
		ConfigPaths: []string{path},
		// identify by DockerProjectName(id)
		ProjectName: DockerProjectName(id),
	})
	if err != nil {
		return nil, ErrInvalidDockerCompose
	}

	return project, nil
}

func RunDockerCompose(ctx context.Context, service api.Compose, project *types.Project) error {
	return service.Up(ctx, project, api.UpOptions{})
}

type TransformDockerComposeOpt struct {
	Ports map[string][]types.ServicePortConfig // service name = ports
}

func TransformDockerCompose(id string, project *types.Project, opt TransformDockerComposeOpt) {
	var addVolume bool
	newVolumeName := DockerVolumeName(id)

	for service := range project.Services {
		temp := project.Services[service]

		// services must be always restart
		temp.Restart = "always"

		// used service for deployment
		if ports, ok := opt.Ports[service]; ok {
			temp.Ports = ports
		} else {
			// unused service for deployment omitted
			temp.Ports = []types.ServicePortConfig{}
		}

		// make sure volumes is named volume
		for i, volume := range temp.Volumes {
			if volume.Type == "bind" {
				addVolume = true
				temp.Volumes[i].Type = "volume"
				temp.Volumes[i].Source = newVolumeName
			}
		}

		project.Services[service] = temp
	}

	if addVolume {
		project.Volumes[newVolumeName] = types.VolumeConfig{}
	}
}
