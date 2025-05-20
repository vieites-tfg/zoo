_default:
  just -l

alias dv := down_vol
down_vol:
  docker compose down -v

alias l := logs
logs service:
  docker compose logs {{service}} -f

_build_zoo_base:
  #!/usr/bin/env bash
  if [[ "$(docker images -f reference=zoo-base | wc -l | xargs)" != "2" ]]
  then
    docker build --target base -t zoo-base .
  fi

_run entrypoint command:
  @just _build_zoo_base
  docker run --rm -w /app -v $PWD:/app --env-file .env --entrypoint={{entrypoint}} zoo-base {{command}}

init:
  @just _run "yarn" "install"

dev:
  docker compose up -d

e2e:
  #!/usr/bin/env bash
  if [[ "$(docker images -f reference=cypress | wc -l | xargs)" != "2" ]]
  then
    docker build -f Dockerfile.cypress -t cypress .
  fi

  echo "//npm.pkg.github.com/:_authToken=\${CR_PAT}" > .npmrc

  docker run --rm -it --network zoo_default \
    -e BASE_URL="http://zoo-frontend" \
    --env-file .env -v $PWD:/e2e -w /e2e cypress yarn run e2e

  rm .npmrc

alias tb := test_backend
test_backend:
  @just _run "lerna" "run test --scope @vieites-tfg/zoo-backend"

alias tf := test_frontend
test_frontend:
  just e2e

test:
  just test_backend
  just test_frontend

lint:
  @just _run "yarn" "lint"

alias ib := image_build
image_build package:
  ./image.sh build {{package}}

alias ip := image_push
image_push package:
  ./image.sh push {{package}}

alias ibp := image_build_push
image_build_push package:
  ./image.sh all {{package}}

alias pr := pkg_remote
pkg_remote package:
  ./push_package.sh remote {{package}}

alias pl := pkg_local
pkg_local package:
  ./push_package.sh local {{package}}

alias prl := pkg_remote_local
pkg_remote_local package:
  ./push_package.sh all {{package}}

cluster := "zoo-cluster"

alias ss := set_secret
set_secret:
  kubectl create secret docker-registry ghcr-secret \
    --docker-server=ghcr.io \
    --docker-username=vieites \
    --docker-password=$CR_PAT \
    -n dev

apply_ingress:
  kubectl apply -f https://raw.githubusercontent.com/kubernetes/ingress-nginx/main/deploy/static/provider/kind/deploy.yaml
  sleep 0.5
  kubectl wait --namespace ingress-nginx \
  --for=condition=ready pod \
  --selector=app.kubernetes.io/component=controller \
  --timeout=90s

create_cluster:
  #!/usr/bin/env bash
  cat <<EOF | kind create cluster --name {{cluster}} --config=-
  kind: Cluster
  apiVersion: kind.x-k8s.io/v1alpha4
  nodes:
  - role: control-plane
    kubeadmConfigPatches:
    - |
      kind: InitConfiguration
      nodeRegistration:
        kubeletExtraArgs:
          node-labels: "ingress-ready=true"
    extraPortMappings:
    - containerPort: 80
      hostPort: 8080
      protocol: TCP
  EOF

  kubectl create namespace dev
  kubectl create namespace pre
  kubectl create namespace pro

  just set_secret

alias dc := delete_cluster
delete_cluster:
  kind delete cluster -n {{cluster}}

check_hosts ns:
  #!/usr/bin/env bash
  host=$(grep "zoo-{{ns}}" /etc/hosts | wc -l | xargs)
  if [[ $host == "0" ]]
  then
    echo "127.0.0.1 zoo-{{ns}}.example.com" | sudo tee -a /etc/hosts
  fi

launch_chart ns:
  just check_hosts {{ns}}

  helm install zoo ./charts/zoo -n {{ns}}
