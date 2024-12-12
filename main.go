package main

import (
	"fmt"

	"github.com/StoneG24/SLaMO/container"
	"github.com/StoneG24/SLaMO/orchestration"
)

func main() {
    err := container.Setup()
    if err != nil {
        fmt.Println(err)
        return 
    }

    // prompt should come from forked process
    prompt := "how many r's are in the word strawberry"
    err = orchestration.Startup(prompt)
    if err != nil {
        fmt.Println(err)
        return 
    }


}
