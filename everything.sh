#!/bin/sh -eu

# export from aseprite
./aseprite-export.sh

# convert performance seeds to performancde sets based on latest data from aseprite
cat ./performance-set.seeds | go run ./cmd/aseprite-deseed

# inline png into res/images.go
./inline.sh

# pack game data into sqauds.data
./pack.sh
