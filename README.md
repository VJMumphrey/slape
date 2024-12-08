# SLaMO

### About

### Installation
We need to install some dependencies so that we can build and run the project. The first thing we need to install is Podman.
Follow your OS specific instructions to install podman.
After that, if you want to run the containers with a gpu, you'll need to install the cuda/rocm CDI 
https://docs.nvidia.com/datacenter/cloud-native/container-toolkit/latest/cdi-support.html

### Building
We use podamn to build the app since it is a dependancy of the project and it makes building it a lot easier. This also allows us to have reproduceable builds.

To build the app run the following command
```bash
podman build -f ./Containerfile -t slamo:latest
```

This should dump the binary files for the project into /build. 

### Cleanup
Containers are very useful for making reproduceable builds but the can take up a lot of space over time. This coupled with the fact that we also have to save storage space for models means that we need to be more cognicent of that fact. Here are some tips to remove dead resources in this project.

This command will tell you how much of your disk is currently being used by podman
```bash
podman system df
```

These commands are good for cleaning up these old resources.
```bash
podman container prune
podman image prune
podman builder prune
```

### Reference
