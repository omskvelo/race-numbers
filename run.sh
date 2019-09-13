#!/bin/bash -ex

REGISTERED_USERS_GOOGLE_SHEET_ID='18ZwIpjhWgrlV4lfeBPsb5XWJYn5P4W7j1Yl7wM4nXI4'


cd google-sheet-to-csv
#go run main.go -app event-table -id "$REGISTERED_USERS_GOOGLE_SHEET_ID" -sec ../client_secret.json > ../_data/participants.csv


cd ../rate-participants
#go run main.go -p ../_data/participants.csv -r ../_data/rating.csv > ../_data/participants_rated.csv


cd ../start-number-draw
#go build


cd ../render-numbers/
RENDER_TMP_DIR='../tmp'
RENDER_OUT_DIR='../_out'
rm -rf "$RENDER_TMP_DIR"
mkdir -p "$RENDER_TMP_DIR"
rm -rf "$RENDER_OUT_DIR"
mkdir -p "$RENDER_OUT_DIR"

go run main.go -bg ../_data/nomer_bg.pdf -p ../_data/participants_rated.csv -tmp "$RENDER_TMP_DIR" -out "$RENDER_OUT_DIR"
rm -rf "$RENDER_TMP_DIR"