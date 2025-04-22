alias dv := down_vol
alias l := logs
alias tb := test_backend
alias tf := test_frontend
alias ib := image_build
alias ip := image_push
alias ibp := image_build_push
alias pr := pkg_remote
alias pl := pkg_local
alias prl := pkg_remote_local

_default:
  just -l

down_vol:
  docker compose down -v

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
  docker run --rm -w /app -v $PWD:/app -e CR_PAT=$CR_PAT --entrypoint={{entrypoint}} zoo-base {{command}}

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

  docker run --rm -it --network zoo_default -v $PWD:/e2e -w /e2e cypress yarn run e2e

test_backend:
  @just _run "lerna" "run test --scope backend"

test_frontend:
  just e2e

test:
  just test_backend
  just test_frontend

lint:
  @just _run "yarn" "lint"

image_build package:
  ./image.sh build {{package}}

image_push package:
  ./image.sh push {{package}}

image_build_push package:
  ./image.sh all {{package}}

pkg_remote package:
  ./push_package.sh remote {{package}}

pkg_local package:
  ./push_package.sh local {{package}}

pkg_remote_local package:
  ./push_package.sh all {{package}}
