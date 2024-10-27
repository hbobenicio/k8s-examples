# my-controller

A simple controller example. No kubebuilder or operator-sdk boilerplate gatekeeping.

## Local minikube/docker step-by-step setup

```bash
minikube start

# so we can build our docker images so minikube can see them
eval $(minikube docker-env)

# let's download our dependencies upfront so we have consistently fast docker builds
go mod vendor

# little helper to build our images
./script/docker-build-all.sh

# you'll be prompted for confirming your kubectl context which would be used as a target
./script/up.sh

```
