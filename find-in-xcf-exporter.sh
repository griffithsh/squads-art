#!/bin/sh

grep -Fq "$1" export-xcfs.sh || echo "WARNING: $1 is not referenced by export-xcfs.sh"
