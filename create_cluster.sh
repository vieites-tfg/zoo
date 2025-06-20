#!/usr/bin/env bash

set -euo pipefail

cat <<EOF | kind create cluster --name "$1" --config=-
kind: Cluster
apiVersion: kind.x-k8s.io/v1alpha4
nodes:
- role: control-plane
  extraPortMappings:
  - containerPort: 80
    hostPort: 8080
    protocol: TCP
  - containerPort: 443
    hostPort: 443
    protocol: TCP
EOF

kubectl create namespace dev
kubectl create namespace pre
kubectl create namespace pro
