#!/bin/sh -eu

# prepare Tiled Map Editor combat terrain resources for packing.
./tiled-export.sh
go run ./cmd/tiled-convert -in ./combat-terrain -out ./combat-terrain
go run ./cmd/tiled-convert -in ./combat-terrain -out ./packed/content/combat-terrain
rm ./combat-terrain/*.export

# cleanup previous run
rm -f ./packed/content/character-appearance/*.appearance
rm -f ./packed/content/character-appearance/*.variations.png

go run ./cmd/char-var

# pack game data into squads.data
./pack.sh

rm ./packed/content/character-appearance/*.appearance
