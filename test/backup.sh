#!/bin/bash

root="$HOME/.steam/root"
target="$(dirname $0)/root"

sources=(
  appcache/appinfo.vdf
  config/config.vdf
  config/loginusers.vdf
  steamapps/libraryfolders.vdf
  userdata/90252099/config/localconfig.vdf
  userdata/90252099/config/shortcuts.vdf
)

for val in ${sources[@]}; do
  mkdir -p $target/$(dirname $val)
  cp $root/$val $target/$val
done

tar cJf $target/../root.tar.xz $target
