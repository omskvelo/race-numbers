#!/usr/bin/env bash

mkdir -p out
rm out/*

go build
./start-number-draw -name "Иван Иванов" -number '999' -team 'ЦР' -o out/number1.pdf
./start-number-draw -name "Иван Иванов" -number '5' -o out/number2.pdf
./start-number-draw -number '555' -o out/number3.pdf
pdftk out/number1.pdf background ../_data/number_bg.pdf output out/out1.pdf
pdftk out/number2.pdf background ../_data/number_bg.pdf output out/out2.pdf
pdftk out/number3.pdf background ../_data/number_bg.pdf output out/out3.pdf
