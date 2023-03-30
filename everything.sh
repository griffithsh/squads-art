#!/bin/sh -eu

# Export all the XCFs as PNGs
./export-xcfs.sh

# Prepare Tiled Map Editor combat terrain resources for packing, moving them into packed/content.
./tiled-export.sh
go run ./cmd/tiled-convert -in ./combat-terrain -out ./packed/content/combat-terrain
find combat-terrain -name '*.png' -exec cp "{}" packed/content/"{}" \;
rm ./combat-terrain/*.export

# Cleanup previous run
find ./packed/content/character-appearance/ -type f | egrep '(\.appearance|\.variations.png)$' | xargs rm -f

go run ./cmd/char-var

# Pack game data into squads.data
./pack.sh

find ./packed/content/character-appearance/ -type f | egrep '(\.appearance|\.variations.png)$' | xargs rm -f

find . -type f -iname '*.xcf' | sed "s/^.\///g" | xargs -I{} ./find-in-xcf-exporter.sh {}
