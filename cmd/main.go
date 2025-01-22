/*
Package SLaP-E is a binary that starts a pod on the local computer using as socket to podman.

Usage:

	./slape

Containerized models are spawned as needed adhering to a pipeline system.
*/
package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/StoneG24/slape/cmd/orchestration"
	"github.com/StoneG24/slape/cmd/pod"
	"github.com/containers/podman/v5/pkg/bindings/containers"
)

var (
    // holds all of the string container IDs 
    // for the containers in the app
    slapeContainers []string
)

type simple struct {
    Prompt string `json:"prompt"`
}

func simplerequest(w http.ResponseWriter, req *http.Request) {

    var simplePayload simple

    decoder := json.NewDecoder(req.Body)

    err := decoder.Decode(&simplePayload)
    if err != nil {
        fmt.Println(err)
    }

    conn, err := pod.CreateBindingConnection()
    if err != nil {
        fmt.Println(err)
        return
    }
    
    //podId, err := pod.SetupPod(conn)
    //pods.Start(conn, podId, nil)

    containerID, err := pod.DefineOllamaContainer(conn, "chatmodel")
    slapeContainers = append(slapeContainers, containerID)

    execSeshID, err := containers.ExecCreate(conn, containerID, nil)
    err = containers.ExecStartAndAttach(conn, execSeshID, nil)

    err = containers.Start(conn, containerID, nil)
    if err != nil {
        fmt.Println(err)
    }

    // TODO(frontend) prompt should come from the frontend socket or something
    prompt := "how many r's are in the word strawberry"

    err = orchestration.GetSimpleAnswer(prompt)

    err = containers.ExecRemove(conn, execSeshID, nil)
    err = containers.Stop(conn, containerID, nil)

    // TODO json the response
    w.Write([]byte("Your gay"))
}

func main() {

    http.HandleFunc("/simple", simplerequest)

    go http.ListenAndServe(":3069", nil)

}
