#!/bin/sh
version=$(git name-rev --tags --name-only $(git rev-parse HEAD))
[ $version = "undefined" ] && version=$(git rev-parse --short HEAD)
echo $version
