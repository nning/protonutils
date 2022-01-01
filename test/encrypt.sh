#!/bin/bash

set -e

target="$(dirname $0)/root"
key="age16yq3pw0wesxhj895czj7nfy4vzufupa0gex9k66ap3es3jlmtezqpnfmh8"

tar cJ $target | age -r $key > $target/../root.tar.xz.age
