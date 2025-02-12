Small Language Model Prompt Engineering
=======

### About

### Installation
We need to install some dependencies so that we can build and run the project. The first thing we need to install is Podman.
Follow your OS specific instructions to install podman.

Create a folder for the models that you want to use. These should probably be *.gguf* files.

#### GPU Support
After that, if you want to run the containers with a gpu, you'll need to install the nvidia continer toolkit along with the appropriate drivers if needed.
https://docs.nvidia.com/datacenter/cloud-native/container-toolkit/latest

Refer to the proper documentation for rocm.

<<<<<<< Updated upstream
#### Better Memory Usage
We will used crun while creating this project. It uses less memory and is faster to startup containers. This is one of the main goals for the project and so its an obvious choice.
[crun](https://github.com/containers/crun)

Download and configure it by following the correct guides for your system. Their readme should be enough. 
Once you've done that, you're going to want to make it the default runtime by editing your config file for podman. If you are using Fedora then it is already setup for you.

#### Socket
### Linux
=======
### Building
Currently we use make as our build system on the backend, the backend and frontend are seperate so if you want to swap out to a different frontend you can.
```bash
make -f Makefile.back
```

cleaning up the build with,
```bash
make -f Makefile.back clean
```

This makefile creates the *models* folder for all your models to be stored in. It will not clear it out. That is a manual task to done by the user.

### Socket

#### Linux

>>>>>>> Stashed changes
To run the app you need to turn on the docker socket. This allows the app to talk to the socket and controll its components.

For non-rootful use cases
```bash
sudo systemctl start docker
```

To close the socket on linux, 
```bash
sudo systemctl stop docker
```
<<<<<<< Updated upstream
### Windows
=======

#### Windows

>>>>>>> Stashed changes
For windows this process is managed by docker desktop.

### Building
We will use several Makefiles to build the project. These build the frontend and backend components into build/ for the backend and frontend inside frontend/build

### Cleanup
Containers are very useful for making reproduceable builds but the can take up a lot of space over time. This coupled with the fact that we also have to save storage space for models means that we need to be more cognicent of that fact. Here are some tips to remove dead resources in this project.

This command will tell you how much of your disk is currently being used by podman
```bash
docker system df
```

<<<<<<< Updated upstream
These commands are good for cleaning up these old resources.
=======
These commands are good for cleaning up these old docker resources.

>>>>>>> Stashed changes
```bash
docker container prune
```
```bash
docker image prune
```
```bash
docker builder prune
```

# Set env variable
```bash
env MODEL_PATH=$(pwd)/models
```

### Reference
Here are some of the research papers that we used to aid us in development.
