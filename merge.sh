#!/bin/bash

mkdir -p _merged
rm _merged/*

for ((i = 1 ; i <= 99 ; i+=2 )); do
    NUM1=$(printf '%03d' $i)
    NUM2=$(printf '%03d' $((i+1)))
    FILE1="${NUM1}.pdf"
    FILE2="${NUM2}.pdf"
    echo $FILE1 $FILE2

    pdfjam --angle 180 _out/$FILE2 --outfile _tmp/$FILE2

    OUTFILE="${NUM1}-${NUM2}.pdf"

    pdfjam --nup 1x2 _out/$FILE1 _tmp/$FILE2 --outfile o.pdf
done