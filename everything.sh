#!/bin/sh -eu

# export from aseprite
./aseprite-export.sh

# convert performance seeds to performancde sets based on latest data from aseprite
cat ./performance-set.seeds | go run ./cmd/aseprite-deseed

# prepare Tiled Map Editor combat terrain resources for packing.
./tiled-export.sh
go run ./cmd/tiled-convert -in ./combat-terrain -out ./combat-terrain
rm ./combat-terrain/*.export

# inline png into res/images.go
./inline.sh

# pack game data into squads.data
./pack.sh
