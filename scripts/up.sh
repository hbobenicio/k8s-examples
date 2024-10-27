#!/bin/bash
set -eu -o pipefail

. ./scripts/common.sh

confirm_kubectl_context

log_info "applying manifests..."

kubectl apply --server-side -f "./k8s/my-operator/my-operator.namespace.k8s.yml"
kubectl apply --server-side -f "./k8s/my-operator/my-operator.serviceaccount.k8s.yml"
kubectl apply --server-side -f "./k8s/my-operator/my-operator.cluster-role.k8s.yml"
kubectl apply --server-side -f "./k8s/my-operator/my-operator.cluster-role-biding.k8s.yml"
kubectl apply --server-side -f "./k8s/my-operator/my-operator.deploy.k8s.yml"

log_info "setup is done."
