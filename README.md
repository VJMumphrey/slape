# Small Language Model Prompt Engineering

### About

An application designed to leverage the efficiency of small language models by implementing prompt engineering, inferencing techniques, and external tooling to increase their accuracy.

We do this by creating what we call pipelines. Pipelines are a way of orchestrating models in a linear fashion to create better answers using multiple models versus one large model.
This is a test to prove that small models (SLMs) can punch above their size and get the job done. Our project is useful in scenarios where there is little memory to work with.

We currently have four pipelines.
1. Simple
2. Chain of Models
3. Debate
4. Embedding

All pipelines, except Embedding, have access to the tools and functionality.

##### Simple Pipeline
This pipeline is meant to be used to when only a single language model is desired.

##### Chain of Models
This pipeline is meant to be used when you want to emulate a Chain of Thought process. This first model in the pipeline generates initial thoughts and answers the question.
The final model generates the answer that is returned to the client. The intermediate models will operate on the previous answer, either affirming it or correcting it.

##### Debate of Models
This pipeline is meant to be used when you want models to debate on a topic. It has been found that debate helps models to generate better answers.
This is especially true when there is expert level knowledge present in the debate which usually challenging to attain when you can't run LLMs. With this pipeline, the SLMs can each
act as a expert.

### Installation

1. We need to install some dependencies so that we can build and run the project. The first thing we need to install is Docker.
Follow your OS specific instructions to install Docker.

2. This project also uses other projects as submodules. To get these projects run the following command to clone the repo and get those dependencies.
```bash
git clone --recursive https://github.com/StoneG24/slape.git
```

3. Create a folder for the models that you want to use. These should probably be _.gguf_ files.

#### GPU Support

After that, if you want to run the containers with a gpu, you'll need to install the nvidia continer toolkit along with the appropriate drivers if needed.
https://docs.nvidia.com/datacenter/cloud-native/container-toolkit/latest

Refer to the proper documentation for rocm.

### Building
This command downloads everything and builds the app for you. A models folder is still needed in the folder that SLaPE is ran in.
```bash
go install github.com/StoneG24/slape@latest
```
### Socket

#### Linux

To run the app you need to turn on the docker socket. This allows the app to talk to the socket and controll its components.

For non-rootful use cases

```bash
sudo systemctl start docker
```

To close the socket on linux,

```bash
sudo systemctl stop docker
```

#### Windows

For windows this process is managed by docker desktop.

### Cleanup

Containers are very useful for making reproduceable builds but the can take up a lot of space over time. This coupled with the fact that we also have to save storage space for models means that we need to be more cognicent of that fact. Here are some tips to remove dead resources in this project.

This command will tell you how much of your disk is currently being used by docker

```bash
docker system df
```

These commands are good for cleaning up these old docker resources. SLaPE should clean these up but, currently, if errors occur it won't.

```bash
docker container prune
```

```bash
docker image prune
```

```bash
docker builder prune
```

### Documentation
Our code uses go doc comments as a way of effectively documenting our code.

This tool is included in our tool list of the project.

To run the doc server locally run this command
```bash
godoc -index -notes="BUG|TODO|NOTE"
```

And travel to this url in your browser http://localhost:6060/pkg/github.com/StoneG24/slape/.

### Indexing RAG (LightRag/MiniRag)

### Reference

Here are some of the research papers that we used to aid us in development.
