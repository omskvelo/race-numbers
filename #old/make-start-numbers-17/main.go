package main

import (
	"encoding/csv"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strconv"
	"strings"

	"github.com/mitchellh/go-homedir"
	"golang.org/x/net/context"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/sheets/v4"
)

const (
	appName                 = "event-table"
	spreadsheetId           = "10VXzVu6MHksFl20dF2Aij5MfFgnGE0YweAt5swz59ps"
	key_sheetIndex          = "sheet_index"
	key_firstName           = "Имя"
	key_lastName            = "Фамилия"
	key_sheetStartNumber    = "Номер"
	key_sheetPersonalNumber = "Именной номер"
	key_hasPersonalNumber   = "has_personal_number"
	key_startNumber         = "start_number"
)

// ---------------------------------------------------------------------------
// Utils
// ---------------------------------------------------------------------------

func dlog(format string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, format, args...)
}

func interfaceToString(value interface{}) string {
	switch typedValue := value.(type) {
	case string:
		return typedValue
	default:
		return ""
	}
}

func interfaceToInt64(value interface{}) int64 {
	if value == nil {
		return 0
	}
	switch typedValue := value.(type) {
	case int64:
		return typedValue
	case int8:
		return int64(typedValue)
	case int16:
		return int64(typedValue)
	case int32:
		return int64(typedValue)
	case uint8:
		return int64(typedValue)
	case uint16:
		return int64(typedValue)
	case uint32:
		return int64(typedValue)
	case int:
		return int64(typedValue)
	case uint:
		return int64(typedValue)
	default:
		return 0
	}
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
	fmt.Printf("Go to the following link in your browser then type the "+
		"authorization code: \n%v\n", authURL)

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
	tokenCacheDir := filepath.Join(homeDir, ".credentials")
	os.MkdirAll(tokenCacheDir, 0700)
	return filepath.Join(tokenCacheDir,
		url.QueryEscape(fmt.Sprintf("%v.json", appName))), err
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
	fmt.Printf("Saving credential file to: %s\n", file)
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

	b, err := ioutil.ReadFile("client_secret.json")
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
	resp, err := srv.Spreadsheets.Values.Get(spreadsheetId, readRange).Do()
	return resp.Values, err
}

func getAllSheetValues() (values [][]interface{}, err error) {
	return sheetGet("A1:T1000")
}

func saveSheetValue(cellRange string, value string) (err error) {
	srv, err := sheetService()

	var vr sheets.ValueRange
	vr.Values = append(vr.Values, []interface{}{value})

	_, err = srv.Spreadsheets.Values.Update(spreadsheetId, cellRange, &vr).ValueInputOption("RAW").Do()
	return
}

func convertAndProcessUserInfos(userInfos [][]interface{}) (outUserInfos []map[string]interface{}) {
	if len(userInfos) == 0 {
		return
	}

	headerRow := userInfos[0]

	for rowIndex, row := range userInfos {

		if rowIndex == 0 {
			continue
		}

		userInfo := make(map[string]interface{})

		userInfo[key_sheetIndex] = rowIndex + 1

		for i, value := range row {
			if i < len(headerRow) {
				userInfo[interfaceToString(headerRow[i])] = value
			}
		}

		userInfo[key_firstName] = strings.TrimSpace(interfaceToString(userInfo[key_firstName]))
		userInfo[key_lastName] = strings.TrimSpace(interfaceToString(userInfo[key_lastName]))
		userInfo[key_startNumber] = userInfo[key_sheetStartNumber]

		if (len(interfaceToString(userInfo["Оплата"])) != 0) ||
			(len(interfaceToString(userInfo["Оплата НА"])) != 0) {
			userInfo["pay"] = "✔"
		}

		outUserInfos = append(outUserInfos, userInfo)
	}

	return
}

func sortUserInfosByStartNumber(userInfos []map[string]interface{}) (outUserInfos []map[string]interface{}) {
	sort.Slice(userInfos, func(i, j int) bool {
		return strings.Compare(interfaceToString(userInfos[i][key_startNumber]), interfaceToString(userInfos[j][key_startNumber])) < 0
	})

	return userInfos
}

func generateTableNumbers() {
	dlog("Getting all sheet values\n")
	values, err := getAllSheetValues()
	if err != nil {
		log.Fatalf("%v\n", err)
	}

	dlog("Converting and processing user infos\n")
	userInfos := convertAndProcessUserInfos(values)

	dlog("Sorting user infos\n")
	userInfos = sortUserInfosByStartNumber(userInfos)

	dlog("User infos: %v", userInfos)

	counter := 0

	for _, userInfo := range userInfos {
		startNumberString := interfaceToString(userInfo[key_startNumber])
		if len(startNumberString) == 0 {
			continue
		}
		startNumberInt, _ := strconv.ParseInt(startNumberString, 10, 64)
		name := fmt.Sprintf("%v %v", userInfo[key_lastName], userInfo[key_firstName])
		team := interfaceToString(userInfo["Клуб/команда"])
		outputFileName := fmt.Sprintf("out/%03d.pdf", startNumberInt)

		out, err := runProgram1(true, true, "/Src/omskvelo.ru/start-number-draw/start-number-draw", "--number", startNumberString, "--name", name, "--team", team, "-o", outputFileName)
		if err != nil {
			dlog("Error: %v\n", err)
		}
		fmt.Printf("%s\n", out)

		counter++
		//if counter > 10 {
		//	break
		//}
	}
}

type User struct {
	name         string
	points       int64
	start_number int64
}

func usersFromCsvFile(csvFilePath string) (userMap map[string]*User, users []*User, err error) {

	file, err := os.Open(csvFilePath)
	if err != nil {
		return nil, nil, err
	}

	csvReader := csv.NewReader(file)
	//csvReader.Comma = ';'

	records, err := csvReader.ReadAll()
	if err != nil {
		return nil, nil, err
	}

	userMap = make(map[string]*User)

	for i, record := range records {
		if len(record) < 2 {
			dlog("Line %v has less 2 columns\n", i)
			continue
		}

		name := record[0]
		name = strings.TrimSpace(name)

		if len(name) == 0 {
			continue
		}

		pointsString := record[1]

		user, ok := userMap[name]
		if !ok {
			user = new(User)
			user.name = name
			userMap[name] = user
		}

		points, _ := strconv.ParseInt(pointsString, 10, 64)

		user.points += points

		//fmt.Printf("%v: %v\n", user.name, user.points)
	}

	users = make([]*User, len(userMap))

	i := 0
	for _, user := range userMap {
		users[i] = user
		i++
	}

	sort.Slice(users, func(i, j int) bool {
		points1 := users[i].points
		points2 := users[j].points
		if points2 < points1 {
			return true
		} else if points2 > points1 {
			return false
		}
		return strings.Compare(users[i].name, users[j].name) < 0
	})

	//for i, user := range users {
	//	user.start_number = int64(i) + 1
	//}

	return
}

func generateRatingNumbers(filePath string) {
	_, users, err := usersFromCsvFile(filePath)
	if err != nil {
		dlog("%v\n", err)
		return
	}

	for _, user := range users {
		startNumberString := fmt.Sprintf("%v", user.start_number)
		outputFileName := fmt.Sprintf("out/%03d.pdf", user.start_number)

		out, err := runProgram1(true, true, "/Src/omskvelo.ru/start-number-draw/start-number-draw", "--number", startNumberString, "--name", user.name, "-o", outputFileName)
		if err != nil {
			dlog("Error: %v\n", err)
		}
		fmt.Printf("%s\n", out)
	}
}

func main() {

	ratedPtr := flag.Bool("r", false, "Generate rated numbers")
	sourceCsvFilePtr := flag.String("i", "", "Input csv file")
	flag.Parse()

	if *ratedPtr {
		generateRatingNumbers(*sourceCsvFilePtr)
	} else {
		generateTableNumbers()
	}

}

func runProgram1(run, print bool, program string, args ...string) (string, error) {
	outStrings, err := runProgram(run, print, program, args...)
	if err != nil {
		return "", err
	}
	if len(outStrings) == 0 {
		return "", nil
	}

	return outStrings[0], nil
}

func runProgram(run, print bool, program string, args ...string) ([]string, error) {

	if print {
		dlog("Running %v", cmdString(program, args))
	}

	cmd := exec.Command(program, args...)

	outString := ""

	out, err := cmd.Output()
	if len(out) != 0 {
		outString = string(out)
	}
	if err != nil {
		outInsert := ""
		if len(outString) != 0 {
			outInsert = fmt.Sprintf(". %v", outString)
		}
		return nil, errors.New(fmt.Sprintf("Error running \"%v\": %v.%v", cmdString(program, args), err, outInsert))
	}

	if len(outString) == 0 {
		return nil, nil
	}

	outStrings := strings.Split(outString, "\n")
	return outStrings, nil
}

func cmdString(program string, args []string) string {
	comps := make([]string, 0)

	comps = append(comps, program)

	for _, arg := range args {
		if strings.Contains(arg, " ") {
			comps = append(comps, fmt.Sprintf("'%v'", arg))
		} else {
			comps = append(comps, arg)
		}
	}

	result := strings.Join(comps, " ")
	return result
}
