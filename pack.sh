#!/bin/sh

shopt -s globstar

tar -czf squads.data \
    *.names \
    *.overworld-recipe \
    **/*.skills \
    combat-terrain/*.terrain \
    combat-terrain/*.png \
    character-appearance/*.appearance \
    character-appearance/*.variations.png