#!/bin/sh -e

# runs everything.sh in a container. Make sure you have the squads repo
# as a sibling directory of this one.

CONTAINER_EXE=docker
if ! hash docker 2>/dev/null
then
	CONTAINER_EXE=podman
fi

$CONTAINER_EXE build -t squads-art-build .
$CONTAINER_EXE run -i -t -v $(pwd):/mnt:z -v $(pwd)/../squads:/squads:z squads-art-build
