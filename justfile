set allow-duplicate-recipes

alias d := down
alias dv := down_vol
alias dva := down_vol_all
alias re := rebuild
alias l := logs
alias tb := test_backend

_default:
  just -l

update_yarn:
  rm yarn.lock || true
  rm -rf node_modules || true
  rm -rf packages/**/node_modules || true
  rm packages/frontend/yarn.lock || true
  rm packages/backend/yarn.lock || true
  yarn install
  cp yarn.lock packages/backend
  cp yarn.lock packages/frontend
  rm yarn.lock
  rm -rf node_modules

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

init:
  docker compose run --rm --no-deps frontend yarn install
  docker compose run --rm --no-deps backend yarn install

dev:
  docker compose up -d


[working-directory: './packages/frontend']
e2e:
  #!/usr/bin/env bash
  if [[ "$(docker images -f reference=cypress | wc -l | xargs)" != "2" ]]
  then
    docker build -f Dockerfile.cypress -t cypress .
  fi

  docker run --rm --network zoo_default -it -v $PWD:/e2e -w /e2e --entrypoint=npx cypress cypress run

[working-directory: 'packages/backend']
test_backend:
  docker run --rm -v "$(pwd)":/app -w /app node:18-slim sh -c "yarn install && yarn test"
