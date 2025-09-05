set dotenv-load

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

test_frontend:
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

test:
  just test_backend
  just test_frontend

lint:
  @just _run "yarn" "lint"

alias ib := image_build
image_build package:
  ./scripts/image.sh build {{package}}

alias ip := image_push
image_push package:
  ./scripts/image.sh push {{package}}

alias ibp := image_build_push
image_build_push package:
  ./scripts/image.sh all {{package}}

alias pr := pkg_remote
pkg_remote package:
  ./scripts/push_package.sh remote {{package}}

alias pl := pkg_local
pkg_local package:
  ./scripts/push_package.sh local {{package}}

alias prl := pkg_remote_local
pkg_remote_local package:
  ./scripts/push_package.sh all {{package}}

apply_ingress:
  kubectl apply -f https://raw.githubusercontent.com/kubernetes/ingress-nginx/main/deploy/static/provider/kind/deploy.yaml
  sleep 0.5
  kubectl wait --namespace ingress-nginx \
  --for=condition=ready pod \
  --selector=app.kubernetes.io/component=controller \
  --timeout=90s

create_cluster cluster:
  kind create cluster --name {{cluster}} --config=./cluster/kind_local.yaml

  kubectl create namespace dev
  kubectl create namespace pre
  kubectl create namespace pro

alias dc := delete_cluster
delete_cluster cluster:
  kind delete cluster -n {{cluster}}

check_hosts +ns:
  #!/usr/bin/env bash
  for n in {{ns}}
  do
    host=$(grep "zoo-$n" /etc/hosts | wc -l | xargs)
    if [[ $host == "0" ]]
    then
        echo "127.0.0.1 zoo-$n.example.com" | sudo tee -a /etc/hosts
    fi
    host=$(grep "api-zoo-$n" /etc/hosts | wc -l | xargs)
    if [[ $host == "0" ]]
    then
        echo "127.0.0.1 api-zoo-$n.example.com" | sudo tee -a /etc/hosts
    fi
  done

launch_chart ns:
  just check_hosts {{ns}}
  set -a; . ./.env; set +a
  helmfile -f ../values/helmfile.yaml.gotmpl -e {{ns}} sync

alias tc := template_chart
template_chart ns:
  set -a; . ./.env; set +a
  helmfile -f ../values/helmfile.yaml.gotmpl -e {{ns}} template

alias lc := lint_chart
lint_chart ns:
  set -a; . ./.env; set +a
  helmfile -f ../values/helmfile.yaml.gotmpl -e {{ns}} lint

alias cac := create_all_clusters
create_all_clusters:
  ./scripts/create_envs.sh

alias dac := delete_all_clusters
delete_all_clusters:
  printf "dev\npre\npro" | xargs -I% kind delete cluster --name %
