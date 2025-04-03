set allow-duplicate-recipes

alias d := down
alias dv := down_vol
alias dva := down_vol_all
alias re := rebuild
alias l := logs
alias tb := test_backend
alias tf := test_frontend

_default:
  just -l

down:
  docker compose down

down_vol:
  docker compose down -v

down_vol_all:
  docker compose down -v --rmi "all"
  docker rmi cypress || true

rebuild:
  docker compose build --no-cache

logs service:
  docker compose logs {{service}} -f

_run entrypoint command:
  #!/usr/bin/env bash
  if [[ "$(docker images -f reference=zoo-zoo | wc -l | xargs)" != "2" ]]
  then
    docker build -t zoo-zoo .
  fi

  docker run --rm -w /app -v $PWD:/app --entrypoint={{entrypoint}} zoo-zoo {{command}}

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
