# Small Language Model Prompt Engineering

## About

An application designed to leverage the efficiency of small language models by implementing prompt engineering, inferencing techniques, and external tooling to increase their accuracy.
This is done by incorporating them into small teams of five or less. This goal of this process is to develop a real-time team that works together to solve problems. This project is based off of the orignal project [Slape](https://github.com/StoneG24/slape). It captures some of the code and ideas from it but the some things will have to rewritten for this idea to work.

## Installation

Update again with documentation for k8s.

### GPU Support

After that, if you want to run the containers with a gpu, you'll need to install the NVIDIA Continer Toolkit along with the appropriate drivers if needed.
[NVIDIA Container Toolkit](https://docs.nvidia.com/datacenter/cloud-native/container-toolkit/latest)

Refer to the proper documentation for ROCM.

## Building
This command downloads the dependencies and builds the app for you. 
```bash
go build ./cmd/main.go
```
## Configuration
To configure the project we used a simple and unorthodox approach.
Instead of using a yaml file, we used a go file that maintains global constants and variables in the program.
These variables are referenced throughtout the program and make it easy to make changes to prompts and numerical values.

these live in the [defs file](pkg/vars/defs.go).

For prompts we also have go files that store the strings that encapsulate the prompts.

these live in [prompts](pkg/prompt/prompt.go).
**NOTE** we also a have set of security prompts for use in demos in [security promtps](pkg/prompt/secprompts.go).

This was choice was made to keep the logic simple and create a binary that could be bundled and moved to remote servers if needed.
A replacement for this would be utilizing golangs templates to better format the prompts. 
This removes the complexity of reading files and keeping SLaPE secure from outside attacks.

## Documentation
Our code uses go doc comments as a way of effectively documenting our code.

This tool is included in our tool list of the project.

To run the doc server locally run this command
```bash
godoc -index -notes="BUG|TODO|NOTE"
```
And travel to this url in your browser http://localhost:6060/pkg/github.com/VJMumphrey/slape/.

## Security
To run security checks on the repo run these commands,

```bash
go vet ./...
go tool staticcheck ./...
go tool govulncheck
```
All of these binaries should have been installed with SLaPE as part of the *tool* directive in go.mod.

## Features
We currently have several feautures to aid in the improvement of SLMs.
We have an optional Internet Search, along with an extra thinking step.

### Thinking
This is currently a prototype of a thought process meant to give the model extra time to consider characteristics and behaviors of a given problem.
To enable this, pass in a "thinking":"1" into your json request to our endpoints.

### Internet Search
This is another optional prototype. It is meant to give a model access to the internet for updated information compared to what it was trained on.
It should be noted that the model itself does not make the request. It merely generates the guery used to search the web. The rest is handled internally.

### Function Calling (WIP)

### Indexing RAG (LightRag/MiniRag) (WIP)
The code is present for guerying the database but it is untested and not integrated into the context.

## Reference

Here are some of the research papers that we used to aid us in development.
