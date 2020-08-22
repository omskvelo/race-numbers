#!/bin/bash

pdftk A=_out/001.pdf B=_out/002.pdf cat A1 B1 output - | pdf2ps -dLanguageLevel=3 - - | psnup -2 -f -m30 -Pa5 -pa4 | ps2pdf -dCompatibility=1.4 - A4.pdf