#!/usr/bin/env bash

go build
./start-number-draw -name "Иван Иванов" -number '300' -o tmp/number.pdf
pdftk tmp/number.pdf background data/bg.pdf output tmp/out.pdf
