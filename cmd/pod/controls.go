package pod

import (
	"context"
	"fmt"

	"github.com/containers/podman/v5/pkg/bindings"
	"github.com/containers/podman/v5/pkg/bindings/containers"
	"github.com/containers/podman/v5/pkg/bindings/pods"
	"github.com/containers/podman/v5/pkg/domain/entities/types"
	"github.com/containers/podman/v5/pkg/specgen"
)

// rootless socket for podman
// for rootful use '/run/podman/podman.sock'
func CreateBindingConnection () (context.Context, error) {
	conn, err := bindings.NewConnection(
        context.Background(), 
        "unix:///run/user/1000/podman/podman.sock",
        )
	if err != nil {
        return nil, err
	}

    return conn, nil
}

// creates the "slamo" pod based on a spec written in the code
// TODO: This should be converted to specfile or a kubernetes config
// for better communication.
func SetupPod(conn context.Context) (string, error) {

    exist, err := pods.Exists(conn, "slamo", nil) 
    if err != nil {
        fmt.Println("Pod slamo already exists")
        return "", err
    }

    if exist {
        pods.Start(conn, "slamo", nil)
        return "", nil
    }

    devices := []string{
        "nvidia.com/gpu=all",
    }


    specGenerator := specgen.NewPodSpecGenerator()
    specGenerator.PodBasicConfig.Name = "slamo"
    specGenerator.PodBasicConfig.Devices = devices
    specGenerator.RestartPolicy = "always"

    podSpec := types.PodSpec {
        PodSpecGen: *specGenerator,
    }


    // NOTE: we should have predefined specs to use for clarity, ease of deployment
    podCreateReport, err := pods.CreatePodFromSpec(conn, &podSpec )
    if err != nil {
        return "", err
    }

    return podCreateReport.Id, nil
}

// create the ollama container for chatting with
// this may also be usefull for creating an embedding model 
// or a small model for sepeculative decoding
// NOTE: This must be ran after the "slamo" pod is created
func SetupOllamaContainer(conn context.Context) (string, error) {

    specGenerator := specgen.SpecGenerator{}
    specGenerator.Name = "chatModel"
    specGenerator.RestartPolicy = "always"
    specGenerator.Pod = "slamo"
    // specGenerator.Terminal = true
    // specGenerator.PortMappings 

    // TODO: look into using ExecCreate for running the modelin
    containerReport, err := containers.CreateWithSpec(conn, &specGenerator, nil)
    if err != nil {
        return "", err
    }

    return containerReport.ID, nil
}
