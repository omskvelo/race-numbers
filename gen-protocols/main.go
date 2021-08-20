package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

const (
	fakeLowRating = 1_000_000
)

// ---------------------------------------------------------------------------
// Utils
// ---------------------------------------------------------------------------

func dlog(format string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, format, args...)
	fmt.Fprintf(os.Stderr, "\n")
}

// ---------------------------------------------------------------------------

type User struct {
	firstName   string
	lastName    string
	name        string
	team        string
	rating      int64
	paid        bool
	startNumber int64
	surName     string
	phone       string
	category    string
	finishTime  string
}

func readCsvFile(csvFilePath string) (records [][]string, err error) {
	file, err := os.Open(csvFilePath)
	if err != nil {
		return nil, err
	}

	defer file.Close()

	csvReader := csv.NewReader(file)
	//csvReader.Comma = ';'

	records, err = csvReader.ReadAll()

	return
}

func csvRecordsToMap(records [][]string) (result []map[string]string) {

	if len(records) == 0 {
		return nil
	}

	result = make([]map[string]string, 0, len(records))

	var header []string

	for line, record := range records {
		if line == 0 {
			header = record
			continue
		}

		recordMap := make(map[string]string)

		for col, value := range record {
			if col >= len(header) {
				break
			}
			colName := header[col]
			recordMap[colName] = value
		}

		result = append(result, recordMap)
	}

	return
}

func finishedUsersFromCsvFile(csvFilePath string) (users []User, err error) {

	records, err := readCsvFile(csvFilePath)
	if err != nil {
		return nil, err
	}

	mapRecords := csvRecordsToMap(records)

	users = make([]User, 0, len(records))

	for _, record := range mapRecords {
		var user User

		user.startNumber, _ = strconv.ParseInt(record["number"], 10, 64)
		user.category = record["category"]
		user.finishTime = record["time"]

		dlog("Finished user: %v, %v, %v", user.startNumber, user.category, user.finishTime)
	}

	return
}

func participantsUsersFromCsvFile(csvFilePath string) (users []User, err error) {

	records, err := readCsvFile(csvFilePath)
	if err != nil {
		return nil, err
	}

	mapRecords := csvRecordsToMap(records)

	users = make([]User, 0, len(records))

	for _, record := range mapRecords {
		var user User

		user.paid = len(record["Оплата_"]) != 0
		user.firstName = strings.TrimSpace(record["Имя"])
		user.lastName = strings.TrimSpace(record["Фамилия"])
		user.name = fmt.Sprintf("%v %v", user.lastName, user.firstName)
		user.team = strings.TrimSpace(record["Клуб/команда"])
		user.rating = fakeLowRating

		category := strings.TrimSpace(record["Категория"])
		categoryFields := strings.Fields(category)
		if len(categoryFields) > 0 {
			user.category = categoryFields[0]
		}

		user.surName = strings.TrimSpace(record["Отчество"])
		user.phone = strings.TrimSpace(record["Телефон"])

		users = append(users, user)

		dlog("Participant: %v, %v", user.name, user.team)
	}

	return
}

func usersMap(users []User) (result map[string]User) {
	result = make(map[string]User, 0)
	for _, user := range users {
		result[user.name] = user
	}
	return
}

var (
	participantsFileName = ""
	ratingFileName       = ""
	dumpNumbers          = false
)

func main() {

	flag.StringVar(&participantsFileName, "p", "", "Participants csv file")
	flag.StringVar(&ratingFileName, "r", "", "Results csv file")

	flag.Parse()

	if len(participantsFileName) == 0 || len(ratingFileName) == 0 {
		flag.Usage()
		return
	}

	finishedUsers, err := finishedUsersFromCsvFile(ratingFileName)
	if err != nil {
		log.Fatal(err)
	}

	participants, err := participantsUsersFromCsvFile(participantsFileName)
	if err != nil {
		log.Fatal(err)
	}

	usersMap := usersMap(participants)
	resultUsers := make([]User, len(finishedUsers))
	copy(resultUsers, finishedUsers)

	for i, finishedUser := range finishedUsers {
		user, ok := usersMap[finishedUser.name]
		if ok {
			finishedUser.firstName = user.firstName
			finishedUser.lastName = user.lastName
			finishedUser.team = user.team
			finishedUser.category = user.category
		}
		resultUsers[i] = user
	}

	for _, user := range resultUsers {
		dlog("Result user: %v %v, %v, %v", user.lastName, user.firstName, user.category, user.finishTime)
	}

	writer := csv.NewWriter(os.Stdout)
	defer writer.Flush()

	writer.Write([]string{"Фамилия Имя", "Команда", "Категория", "Стартовый номер", "Время"})

	for _, resultUser := range resultUsers {

		lineArray := make([]string, 0)
		lineArray = append(lineArray, fmt.Sprintf("%v %v", resultUser.lastName, resultUser.firstName))
		lineArray = append(lineArray, resultUser.team)
		lineArray = append(lineArray, resultUser.category)
		lineArray = append(lineArray, fmt.Sprintf("%v", resultUser.startNumber))
		lineArray = append(lineArray, resultUser.finishTime)

		writer.Write(lineArray)
	}
}
