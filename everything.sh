#!/bin/sh -eu

# prepare Tiled Map Editor combat terrain resources for packing.
./tiled-export.sh
go run ./cmd/tiled-convert -in ./combat-terrain -out ./combat-terrain
rm ./combat-terrain/*.export

rm -f ./character-appearance/*.appearance
rm -f ./character-appearance/*.variations.png
go run ./cmd/char-var

# inline png into res/images.go
go run ./cmd/png-inline -o ../squads/res/images.go -p res

# pack game data into squads.data
./pack.sh
