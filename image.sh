#!/usr/bin/env bash

set -e

trap "exit 1" SIGUSR1
PID="$$"

PACKAGES="backend frontend"

usage () {
	cat <<EOF
Usage: $0 <ACTION> <PACKAGE>

ACTION:
  build		build the indicated <PACKAGE>
  push		push the indicated <PACKAGE>
  all		do both build and push

PACKAGE:
  backend	do the <ACTION> for the image of the 'backend' package
  frontend	do the <ACTION> for the image of the 'frontend' package
  all		do the <ACTION> for all the images
EOF

exitt
}

exitt () {
	kill -SIGUSR1 "$PID"
}

# 1:	the package
get_version () {
	cd packages/"$1"
	local version=$(node -p "require('./package.json').version")
	cd - &>/dev/null

	echo "$version"
}

# 1:	a list of valid packages
push_image () {
	for p in $1
	do
		local version=$(get_version "$p")

		# Push.
		docker push ghcr.io/vieites-tfg/zoo-"$p":"$version"
		docker push ghcr.io/vieites-tfg/zoo-"$p":latest
	done
}

# 1:	a list of valid packages
build_image () {
	for p in $1
	do
		local version=$(get_version "$p")

		# Build.
		docker build --target "$p" -t ghcr.io/vieites-tfg/zoo-"$p":"$version" .
		docker tag ghcr.io/vieites-tfg/zoo-"$p":"$version" ghcr.io/vieites-tfg/zoo-"$p":latest
	done
}

main () {
	if [[ "$#" != 2 ]]
	then
		usage
	fi

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
		build)
			build_image "$package"
			;;
		push)
			push_image "$package"
			;;
		all)
			build_image "$package"
			push_image "$package"
			;;
		*)
			usage
			;;
	esac
}

main "$@"
