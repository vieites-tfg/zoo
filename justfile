set allow-duplicate-recipes

alias d := down
alias dv := down_vol
alias dva := down_vol_all
alias re := rebuild

_default:
  just -l

down:
  docker compose down

down_vol:
  docker compose down -v

down_vol_all:
  docker compose down -v --rmi "all"

rebuild:
  docker compose build --no-cache

init:
  yarn install

dev:
  docker compose up
