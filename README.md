Small Language Model Prompt Engineering
=======

### About

### Installation
We need to install some dependencies so that we can build and run the project. The first thing we need to install is Podman.
Follow your OS specific instructions to install podman.

Create a folder for the models that you want to use. These should probably be *.gguf* files.

#### GPU Support
After that, if you want to run the containers with a gpu, you'll need to install the cuda/rocm CDI 
https://docs.nvidia.com/datacenter/cloud-native/container-toolkit/latest/cdi-support.html

Refer to the proper documentation for rocm.

#### Better Memory Usage
We will used crun while creating this project. It uses less memory and is faster to startup containers. This is one of the main goals for the project and so its an obvious choice.
[crun](https://github.com/containers/crun)

Download and configure it by following the correct guides for your system. Their readme should be enough. 
Once you've done that, you're going to want to make it the default runtime by editing your config file for podman. If you are using Fedora then it is already setup for you.

#### Socket
To run the app you need to turn on the podman socket. This allows the app to talk to the socket and controll its components.

For rootful uses
```bash
systemctl start podman.socket
```

For non-rootful use cases
```bash
systemctl start --user podman.socket
```

To close the socket run,
```bash
systemctl stop
```

### Building
We will use several Makefiles to build the project. These build the frontend and backend components into build/ for the backend and frontend inside frontend/build

### Cleanup
Containers are very useful for making reproduceable builds but the can take up a lot of space over time. This coupled with the fact that we also have to save storage space for models means that we need to be more cognicent of that fact. Here are some tips to remove dead resources in this project.

This command will tell you how much of your disk is currently being used by podman
```bash
podman system df
```

These commands are good for cleaning up these old resources.
```bash
podman container prune
```
```bash
podman image prune
```
```bash
podman builder prune
```

### Reference
Here are some of the research papers that we used to aid us in development.
