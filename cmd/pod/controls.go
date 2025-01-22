package pod

import (
	"context"

	"github.com/containers/podman/v5/pkg/bindings"
	"github.com/containers/podman/v5/pkg/bindings/containers"
	"github.com/containers/podman/v5/pkg/specgen"
)

// CreateBindingConnection creates a rootless client for podman.
// takes in a context attribute and returns and error.
// for rootful use '/run/podman/podman.sock'
func CreateBindingConnection() (context.Context, error) {
	conn, err := bindings.NewConnection(
		context.Background(),
		"unix:///run/user/1000/podman/podman.sock",
	)
	if err != nil {
		return nil, err
	}

	return conn, nil
}

// SetupPod creates the "slap-e" pod based on a spec written in the code
/*
func SetupPod(conn context.Context) (string, error) {
	// TODO(v) This should be converted to specfile or a kubernetes config
	// for better communication.

	exist, err := pods.Exists(conn, "slape", nil)
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
	specGenerator.PodBasicConfig.Name = "slape"
	specGenerator.PodBasicConfig.Devices = devices
	specGenerator.RestartPolicy = "always"

	podSpec := types.PodSpec{
		PodSpecGen: *specGenerator,
	}

	// TODO(v) we should have a predefined spec file for clarity, ease of deployment
	podCreateReport, err := pods.CreatePodFromSpec(conn, &podSpec)
	if err != nil {
		return "", err
	}

	return podCreateReport.Id, nil
}
*/

// SetupOllamaContainer creates the ollama container for chatting with.
// this may also be usefull for creating an embedding model
// or a small model for sepeculative decoding
func DefineOllamaContainer(conn context.Context, name string) (string, error) {

	specGenerator := specgen.NewSpecGenerator("", false)
	specGenerator.Name = name
	specGenerator.RestartPolicy = "always"
	//specGenerator.Pod = "slamo"
	specGenerator.Terminal = &[]bool{true}[0]
	//specGenerator.PortsMapping

	// TODO(v) look into using ExecCreate for running the model in
	containerReport, err := containers.CreateWithSpec(conn, specGenerator, nil)
	if err != nil {
		return "", err
	}

	return containerReport.ID, nil
}
