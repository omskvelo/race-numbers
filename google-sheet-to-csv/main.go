package main

import (
	"encoding/csv"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"path/filepath"

	"github.com/mitchellh/go-homedir"
	"golang.org/x/net/context"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/sheets/v4"
)

var (
	appName string
	sheetId string
	secretFileName string
)

// ---------------------------------------------------------------------------
// Utils
// ---------------------------------------------------------------------------

func dlog(format string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, format, args...)
	fmt.Fprintf(os.Stderr, "\n")
}

func interfaceToString(value interface{}) string {
	switch typedValue := value.(type) {
	case string:
		return typedValue
	default:
		return ""
	}
}


func checkError(message string, err error) {
	if err != nil {
		log.Fatal(message, err)
	}
}

func interfaceArrayToStringsArray(args []interface{}) []string {

	result := make([]string, 0, len(args))

	for _, arg := range args {
		s := interfaceToString(arg)
		result = append(result, s)
	}

	return result
}

func interfaceArrayArrayToStringsArrayArray(args [][]interface{}) [][]string {
	result := make([][]string, 0, len(args))

	for _, arg := range args {
		converted := interfaceArrayToStringsArray(arg)
		result = append(result, converted)
	}

	return result
}


// ---------------------------------------------------------------------------
// Sheets api
// ---------------------------------------------------------------------------

func sheetsGetClient(ctx context.Context, config *oauth2.Config) *http.Client {
	cacheFile, err := sheetsTokenCacheFile()
	if err != nil {
		log.Fatalf("Unable to get path to cached credential file. %v", err)
	}
	tok, err := sheetsTokenFromFile(cacheFile)
	if err != nil {
		tok = sheetsGetTokenFromWeb(config)
		sheetsSaveToken(cacheFile, tok)
	}
	return config.Client(ctx, tok)
}

func sheetsGetTokenFromWeb(config *oauth2.Config) *oauth2.Token {
	authURL := config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
	dlog("Go to the following link in your browser then type the "+
		"authorization code: \n%v", authURL)

	var code string
	if _, err := fmt.Scan(&code); err != nil {
		log.Fatalf("Unable to read authorization code %v", err)
	}

	tok, err := config.Exchange(oauth2.NoContext, code)
	if err != nil {
		log.Fatalf("Unable to retrieve token from web %v", err)
	}
	return tok
}

func sheetsTokenCacheFile() (string, error) {
	homeDir, err := homedir.Dir()
	if err != nil {
		return "", err
	}
	tokenCacheDir := filepath.Join(homeDir, ".google-api-credentials")
	os.MkdirAll(tokenCacheDir, 0700)
	filename := fmt.Sprintf("%v-%v.json", appName, sheetId)
	return filepath.Join(tokenCacheDir,
		url.QueryEscape(filename)), err
}

func sheetsTokenFromFile(file string) (*oauth2.Token, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	t := &oauth2.Token{}
	err = json.NewDecoder(f).Decode(t)
	defer f.Close()
	return t, err
}

func sheetsSaveToken(file string, token *oauth2.Token) {
	dlog("Saving credential file to: %s", file)
	f, err := os.OpenFile(file, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		log.Fatalf("Unable to cache oauth token: %v", err)
	}
	defer f.Close()
	json.NewEncoder(f).Encode(token)
}

// ---------------------------------------------------------------------------
// Main program
// ---------------------------------------------------------------------------

func sheetService() (service *sheets.Service, err error) {
	ctx := context.Background()

	b, err := ioutil.ReadFile(secretFileName)
	if err != nil {
		log.Fatalf("Unable to read client secret file: %v", err)
	}

	config, err := google.ConfigFromJSON(b, "https://www.googleapis.com/auth/spreadsheets")
	if err != nil {
		log.Fatalf("Unable to parse client secret file to config: %v", err)
	}
	client := sheetsGetClient(ctx, config)

	service, err = sheets.New(client)
	return
}

func sheetGet(readRange string) (values [][]interface{}, err error) {
	srv, err := sheetService()
	resp, err := srv.Spreadsheets.Values.Get(sheetId, readRange).Do()
	return resp.Values, err
}

func getAllSheetValues() (values [][]interface{}, err error) {
	return sheetGet("A1:T1000")
}

func saveSheetValue(spreadsheetId string, cellRange string, value string) (err error) {
	srv, err := sheetService()

	var vr sheets.ValueRange
	vr.Values = append(vr.Values, []interface{}{value})

	_, err = srv.Spreadsheets.Values.Update(spreadsheetId, cellRange, &vr).ValueInputOption("RAW").Do()
	return
}

func preprocessValuesForCsv(values [][]string) (resultValues [][]string) {

	maxLen := 0
	for _, value := range values {
		valueLen := len(value)
		if valueLen > maxLen {
			maxLen = valueLen
		}
	}

	resultValues = make([][]string, 0, len(values))
	for _, value := range values {
		resultValue := value
		for len(resultValue) < maxLen {
			resultValue = append(resultValue, "")
		}
		resultValues = append(resultValues, resultValue)
	}

	return
}

func main() {

	flag.StringVar(&sheetId,"id", "", "Sheet id")
	flag.StringVar(&appName,"app", "", "App name")
	flag.StringVar(&secretFileName,"sec", "", "Client secret file name")

	flag.Parse()

	if len(sheetId) == 0 || len(appName) == 0 || len(secretFileName) == 0 {
		flag.Usage()
		return
	}

	rawValues, err := getAllSheetValues()
	values := interfaceArrayArrayToStringsArrayArray(rawValues)
	values = preprocessValuesForCsv(values)

	if err != nil {
		log.Fatal(err)
	}

	writer := csv.NewWriter(os.Stdout)
	defer writer.Flush()

	for _, value := range values {
		err := writer.Write(value)
		checkError("Cannot write to file", err)
	}
}
