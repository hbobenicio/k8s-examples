# my-controller

A simple controller example. No kubebuilder or operator-sdk boilerplate gatekeeping.

> First you have a complex problem.
>
> Then in order to control the complexity you just hide it under the carpet with a bad abstraction,
> instead of really understanding the problem and how to solve it.
>
> Now you have two problems.

## What is the big picture?

This project demonstrates how to create a really simple custom controller called `my-controller`.

This controller watches for `v1.Namespace` resources and delegates it to
our own custom reconciler called `NamespaceReconciler`.

On changes our reconciler just ensures the presence of a really factual label to be present on all objects.

After deployment you can follow the controller logs see what is going on.

Optionally you can check the namespaces after the controller work is done to see what kind of annotation it creates,
or maybe you take a shortcut and read the sources to find out. Consider it a homework.

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

# opcional: run this to check my-controller logs
kubectl -n my-controller logs deployments/my-controller -f
```
