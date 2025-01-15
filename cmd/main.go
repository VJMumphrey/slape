package main

import (
	"fmt"
	"log"

	"github.com/StoneG24/SLaMO/cmd/orchestration"
	"github.com/StoneG24/SLaMO/cmd/pod"
	"github.com/containers/podman/v5/pkg/bindings/pods"
)

func main() {
    conn, err := pod.CreateBindingConnection()
    if err != nil {
        fmt.Println(err)
        return
    }

    podId, err := pod.SetupPod(conn)
    if err != nil {
        fmt.Println(err)
        return 
    }

    pods.Start(conn, podId, nil)

    // prompt should come from the frontend socket or something
    prompt := "how many r's are in the word strawberry"
    err = orchestration.Startup(prompt)
    if err != nil {
        fmt.Println(err)
        return 
    }

    // TODO: change to be on quit
    podStopReport, err := pods.Stop(conn, podId, nil)
    if err != nil {
        fmt.Println(err)
        return
    }

    log.Println(podStopReport)
}
