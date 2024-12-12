package container

import (
	"context"

	"github.com/containers/podman/v5/pkg/bindings"
    "github.com/containers/podman/v5/pkg/bindings/images"
)

// pull the ollama container
// if error is present return error else nil
func pullOllama(conn context.Context) (err error) {
    _, err = images.Pull(conn, "docker.com/r/ollama/ollama", nil)
	if err != nil {
        return err
	}

    return nil
}

// perform checks for app resources 
// add them if they are not present
func Setup() (err error) {

	conn, err := bindings.NewConnection(context.Background(), "unix:///run/podman/podman.sock")
	if err != nil {
        return err
	}

    pullOllama(conn)

    return nil
}
