#!/bin/bash

set -e

root="$HOME/.steam/root"
target="$(dirname $0)/root"
key="age16yq3pw0wesxhj895czj7nfy4vzufupa0gex9k66ap3es3jlmtezqpnfmh8"

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

tar cJ $target | age -r $key > $target/../root.tar.xz.age
