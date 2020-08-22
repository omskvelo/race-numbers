#!/bin/bash

program_exists () {
	type "$1" &> /dev/null ;
}

main()
{
    if ! program_exists pdfjam; then
        echo "Please install pdfjam"
        exit 2
    fi

    mkdir -p _merged
    mkdir -p _tmp
    rm _merged/*
    rm _tmp/*

    for ((i = 1 ; i <= 100 ; i+=2 )); do
        NUM1=$(printf '%03d' $i)
        NUM2=$(printf '%03d' $((i+1)))
        FILE1="${NUM1}.pdf"
        FILE2="${NUM2}.pdf"
        echo $FILE1 $FILE2

        pdfjam --angle 180 --fitpaper true --rotateoversize true _out/$FILE2 --outfile _tmp/$FILE2

        OUTFILE="${NUM1}-${NUM2}.pdf"

        pdfjam --nup 1x2 _out/$FILE1 _tmp/$FILE2 --outfile _merged/$OUTFILE
    done
}

main "$@"
