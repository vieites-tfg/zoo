alias dv := down_vol
alias l := logs
alias tb := test_backend
alias tf := test_frontend
alias b := image_build
alias p := image_push
alias bp := build_and_push
alias pp := push_package

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

build_and_push package version:
  just build_image {{package}} {{version}}
  just push_image {{package}} {{version}}

push_package package:
  @just _build_zoo_base
  @just _run "yarn" "publish --access restricted ./packages/{{package}}"

push_packages:
  @just push_package frontend
  @just push_package backend
