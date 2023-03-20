#!/bin/sh

cd packed/content
tar -czf squads.data * **/*
mv squads.data ../..
