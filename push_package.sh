#!/usr/bin/env bash

set -e

trap "exit 1" SIGUSR1
PID="$$"

ROOT=$(realpath $(dirname "$0"))
PACKAGES="backend frontend"
LOCAL="$ROOT/local_packages"

usage () {
	cat <<EOF
Usage: $0 <LOCATION> <PACKAGE>

LOCATION:
  remote	push remotely the indicated <PACKAGE>
  local		push locally the indicated <PACKAGE>
  all		push both remote and locally

PACKAGE:
  backend	push the 'backend' package into the <LOCATION> registry
  frontend	push the 'frontend' package into the <LOCATION> registry
  all		push all the packages into the <LOCATION> registry
EOF

exitt
}

exitt () {
	kill -SIGUSR1 "$PID"
}

build_base () {
	if [[ "$(docker images -f reference=zoo-base | wc -l | xargs)" != "2" ]]
	then
		docker build --target base -t zoo-base .
	fi
}

# 1:	a list of valid packages
to_remote () {
	for p in $1
	do
		docker run --rm -w /app -v $PWD:/app -e CR_PAT=$CR_PAT --entrypoint=yarn \
		zoo-base publish --access restricted ./packages/"$p"
	done
}

# 1:	a list of valid packages
to_local () {
	mkdir -p "$LOCAL"
	rm -rf "$LOCAL/*"

	for p in $1
	do
		cd "$ROOT/packages/$p"
		version=$(node -p "require('./package.json').version")
		yarn pack --filename "$LOCAL/zoo-${p}-${version}.tgz"
		cd "$ROOT"
	done
}

main () {
	if [[ "$#" != 2 ]]
	then
		usage
	fi

	local package=""

	case "$2" in
		backend | frontend)
			package="$2"
			;;
		all)
			package="$PACKAGES"
			;;
		*)
			usage
			;;
	esac

	case "$1" in
		remote | local | all)
			build_base
			;;
		*)
			usage
			;;
	esac

	case "$1" in
		remote)
			to_remote "$package"
			;;
		local)
			to_local "$package"
			;;
		all)
			to_remote "$package"
			to_local "$package"
			;;
	esac
}

main "$@"
