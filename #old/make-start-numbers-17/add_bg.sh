#!/bin/bash -ex

cd "$(dirname "${BASH_SOURCE[0]}")"

mkdir -p out

cd out_num

for f in *.pdf; do
	pdftk $f background ../bg.pdf output ../out/$f
done