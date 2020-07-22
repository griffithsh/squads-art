#!/bin/sh -e

# Make sure it doesnt break when there are no .tsx/.tmx files
shopt -s nullglob

cd combat-terrain/

# Make sure there isn't any data from old runs.
rm -f ./*.tmx.export
rm -f ./*.tsx.export
rm -f ./*.terrain

for f in *.tmx
do
    tiled --export-map json $f ./$f.export
done

for f in *.tsx
do
    tiled --export-tileset json $f ./$f.export
done
