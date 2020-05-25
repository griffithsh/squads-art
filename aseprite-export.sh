#!/bin/sh -eu

cd ./aseprite/performance-sets

for f in ./*.ase
do
    name="${f%.*}"

    $(aseprite --batch \
            --data $name.json \
            --sheet $name-sheet.png \
            --format json-array \
            --sheet-type rows \
            --ignore-empty \
            --split-layers \
            --filename-format "{tag}-{layer}-{tagframe}" \
            --all-layers \
        ./$name.ase)

    $(mv ./$name-sheet.png ../../)
done
