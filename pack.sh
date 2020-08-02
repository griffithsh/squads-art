#!/bin/sh

tar -czf squads.data \
    *.names \
    *.overworld-recipe \
    *.png \
    *.performance-set \
    *.skills \
    combat-terrain/*.terrain \
    combat-terrain/*.png \
    character-appearance/*.appearance \
    character-appearance/*.variations.png
