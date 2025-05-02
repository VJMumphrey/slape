package pipeline

import (
	"context"
	"log"
	"os"
	"path/filepath"
	"runtime"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/mount"
	"github.com/docker/docker/client"
	"github.com/docker/go-connections/nat"
)

func CreateOllamaContainer(apiClient *client.Client, portNum string, name string, ctx context.Context, containerImage string, gpuTrue bool) (container.CreateResponse, error) {

	portSet := nat.PortSet{
		nat.Port("8000/tcp"): struct{}{}, // map 11434 TCP port
	}

	portBindings := nat.PortMap{
		nat.Port("8000/tcp"): []nat.PortBinding{
			{
				HostIP:   "0.0.0.0",
				HostPort: portNum,
			},
		},
	}

	var mountString string

	if runtime.GOOS == "windows" {
		ex, err := os.Executable()
		if err != nil {
			log.Println("idk something else")
		}

		currentPath := filepath.Dir(ex)

		log.Println(currentPath)

		mountString = currentPath + "\\models"
	}

	if runtime.GOOS == "linux" {
		mountString = os.Getenv("PWD") + "/models"
	}

	// TODO(v) add --jinja for function calling using the OpenAI API setup
	/*
		var cmds []string
		if !gpuTrue {
			cmds = []string{"-m", "/models/" + modelName, "--port", "8000", "--host", "0.0.0.0", "-fa", "--mlock", "--no-webui", "-c", strconv.Itoa(vars.ContextLength), "-cb"}
		}
	*/

	var hostconfig container.HostConfig

	// TODO(v) expand past nvidia systems.
	// ROCm will present interesting challenges. Its simpler but more setups in the config.
	switch gpuTrue {
	case true:
		hostconfig = container.HostConfig{
			Runtime:      "nvidia",
			PortBindings: portBindings,
			Binds:        []string{"ollama:/root/.ollama"},
			/*
				Mounts: []mount.Mount{{
					Type:   mount.TypeVolume,
					Source: "ollama",
					Target: "/root/.ollama",
				}},
			*/
		}
	case false:
		hostconfig = container.HostConfig{
			PortBindings: portBindings,
			Mounts: []mount.Mount{{
				Type:   mount.TypeBind,
				Source: mountString,
				Target: "/models",
			}},
		}
	}

	var createResponse container.CreateResponse
	var err error
	// create container
	createResponse, err = apiClient.ContainerCreate(ctx, &container.Config{
		ExposedPorts: portSet,
		Image:        containerImage,
		//Cmd:          cmds,
	}, &hostconfig, nil, nil, name)

	return createResponse, err
}
