#!/bin/sh

set -e

cd $(dirname $0)

age -d -i $1 < root.tar.xz.age > root.tar.xz
tar xf root.tar.xz -C ..
