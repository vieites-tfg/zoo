#!/usr/bin/env bash

set -euo pipefail

CURRENT_DIR="$(realpath "$0" | xargs dirname)"
SOPS_DIR="${CURRENT_DIR}/../sops"
ARGO_DIR="${CURRENT_DIR}/../argo"
CLUSTER_DIR="${CURRENT_DIR}/../cluster"
ENVS=("dev" "pre" "pro")
PASSWORDS=""

INGRESS_MANIFEST="https://raw.githubusercontent.com/kubernetes/ingress-nginx/main/deploy/static/provider/kind/deploy.yaml"

for ENV in "${ENVS[@]}"; do
    case "${ENV}" in
        dev)
            BANNER_TEXT="We are in DEV";;
        pre)
            BANNER_TEXT="We are in PRE";;
        pro)
            BANNER_TEXT="We are in PRO";;
    esac
    CONTEXT="kind-${ENV}"

    echo "--- Creating cluster '${ENV}' ---"
    kind create cluster --config "${CLUSTER_DIR}/kind_${ENV}.yaml" \
        --wait 5m || continue

    echo "--- Installing Ingress Controller in cluster '${ENV}' ---"
    kubectl apply -f "${INGRESS_MANIFEST}" --context "${CONTEXT}"
    kubectl wait --namespace ingress-nginx \
        --for=condition=ready pod \
        --selector=app.kubernetes.io/component=controller \
        --timeout=120s \
        --context "${CONTEXT}"

    kubectl create namespace argocd || true

    echo "--- Applying SOPS AGE Key for ArgoCD ---"
    cat "${SOPS_DIR}/age.agekey" |
        kubectl create secret generic sops-age -n argocd --context ${CONTEXT} \
        --from-file=keys.txt=/dev/stdin

    echo "--- Installing ArgoCD in cluster '${ENV}' ---"
    helm repo add argo https://argoproj.github.io/argo-helm
    helm repo update
    helm install argocd argo/argo-cd \
    -n argocd \
    -f argo/values.yaml \
    --wait \
    --version 6.11.1 \
    --kube-context "${CONTEXT}"

    kubectl wait --for=condition=Available deployment --all \
        -n argocd --context "${CONTEXT}" --timeout=5m

    echo "--- Applying ArgoCD application for '${ENV}' ---"
    kubectl apply -f "${ARGO_DIR}/argo_${ENV}.yaml" --context "${CONTEXT}"
    kubectl wait --for=condition=Ready pod -l app.kubernetes.io/name=argocd-server \
        -n argocd --context "${CONTEXT}" --timeout=5m

    echo "--- Applying ArgoCD UI banner for '${ENV}' ---"
    kubectl patch configmap argocd-cm -n argocd --context "${CONTEXT}" --type merge \
        -p "{\"data\":{\"ui.bannercontent\":\"${BANNER_TEXT}\"}}"

    echo "--- Cluster '${ENV}' setup complete ---"
    pass=$(kubectl -n argocd get secret argocd-initial-admin-secret \
        -o jsonpath="{.data.password}" | base64 -d)

    current_pass="${ENV} password: ${pass}\n"
    printf "${current_pass}"
    PASSWORDS+="${current_pass}"
done

echo "--- All environments created successfully! ---"
printf "${PASSWORDS}"
