#!/usr/bin/env bash

rm test/number{1,2,3}.pdf
rm test/out{1,2,3}.pdf
go build
./start-number-draw -name "Иван Иванов" -number '999' -team 'ЦР' -o test/number1.pdf
./start-number-draw -name "Иван Иванов" -number '5' -o test/number2.pdf
./start-number-draw -number '5' -o test/number3.pdf
pdftk test/number1.pdf background ../_data/number_bg.pdf output test/out1.pdf
pdftk test/number2.pdf background ../_data/number_bg.pdf output test/out2.pdf
pdftk test/number3.pdf background ../_data/number_bg.pdf output test/out3.pdf
