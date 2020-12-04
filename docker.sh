#!/bin/sh

# runs everything.sh in a docker container. Make sure you have the squads repo
# as a sibling directory of this one.

docker build -t squads-art-build .
docker run -i -t -v $(pwd):/mnt:z -v $(pwd)/../squads:/squads:z squads-art-build
