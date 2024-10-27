#!/bin/bash
set -eu -o pipefail

. ./scripts/common.sh

confirm_kubectl_context

log_info "applying manifests..."

kubectl apply --server-side -f "./k8s/my-controller/my-controller.namespace.k8s.yml"
kubectl apply --server-side -f "./k8s/my-controller/my-controller.serviceaccount.k8s.yml"
kubectl apply --server-side -f "./k8s/my-controller/my-controller.cluster-role.k8s.yml"
kubectl apply --server-side -f "./k8s/my-controller/my-controller.cluster-role-biding.k8s.yml"
kubectl apply --server-side -f "./k8s/my-controller/my-controller.deploy.k8s.yml"

log_info "setup is done."
