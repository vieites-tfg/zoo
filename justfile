set allow-duplicate-recipes

alias d := down
alias da := down_all
alias re := rebuild

_default:
  just -l

down:
  docker compose down

down_all:
  docker compose down -v --rmi "all"

rebuild:
  docker compose build --no-cache

init:
  yarn install

dev:
  docker compose up
