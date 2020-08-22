#!/bin/bash -ex

program_exists () {
	type "$1" &> /dev/null ;
}

main()
{
    REGISTERED_USERS_GOOGLE_SHEET_ID='18565fZBloJgaOP9YJrCwLUIYArc2HUBqVEljTlhH2TA'

    #cd google-sheet-to-csv
    #go run main.go -app event-table -id "$REGISTERED_USERS_GOOGLE_SHEET_ID" -sec ../client_secret.json > ../_data/participants.csv
    # cd ..


    # cd rate-participants
    # go run main.go -p ../_data/participants.csv -r ../_data/rating.csv > ../_data/participants_rated.csv
    # cd ..

    if ! program_exists pdftk; then
        echo "Please install pdftk"
        exit 2
    fi

    cd start-number-draw
    go build
    cd ..

    cd render-numbers
    RENDER_TMP_DIR='../_tmp'
    RENDER_OUT_DIR='../_out'
    rm -rf "$RENDER_TMP_DIR"
    mkdir -p "$RENDER_TMP_DIR"
    rm -rf "$RENDER_OUT_DIR"
    mkdir -p "$RENDER_OUT_DIR"

    go run main.go -limit 99 -bg ../_data/number_bg.pdf -p ../_data/participants_rated.csv -tmp "$RENDER_TMP_DIR" -out "$RENDER_OUT_DIR"
    #rm -rf "$RENDER_TMP_DIR"
}

main "$@"

