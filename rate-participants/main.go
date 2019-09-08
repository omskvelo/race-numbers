package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"log"
	"os"
	"sort"
	"strconv"
	"strings"
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
	firstName    string
	lastName     string
	name         string
	team         string
	points       int64
	paid         bool
	start_number int64
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

func ratedUsersFromCsvFile(csvFilePath string) (users []User, err error) {

	records, err := readCsvFile(csvFilePath)
	if err != nil {
		return nil, err
	}

	mapRecords := csvRecordsToMap(records)

	users = make([]User, 0, len(records))

	for _, record := range mapRecords {
		var user User

		user.name = strings.TrimSpace(record["name"])
		user.points, _ = strconv.ParseInt(record["pts"], 10, 64)

		users = append(users, user)
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

		// paid := len(record["Оплата"]) != 0
		// if !paid {
		// 	continue
		// }
		user.firstName = strings.TrimSpace(record["Имя"])
		user.lastName = strings.TrimSpace(record["Фамилия"])
		user.team = strings.TrimSpace(record["Клуб/команда"])
		user.name = fmt.Sprintf("%v %v", user.lastName, user.firstName)

		users = append(users, user)
	}

	return
}

func ratedUsersMap(users []User) (result map[string]User) {
	result = make(map[string]User, 0)
	for _, user := range users {
		result[user.name] = user
	}
	return
}

var (
	participantsFileName = ""
	ratingFileName       = ""
)

func main() {

	flag.StringVar(&participantsFileName, "p", "", "Participants csv file")
	flag.StringVar(&ratingFileName, "r", "", "Rating csv file")

	flag.Parse()

	if len(participantsFileName) == 0 || len(ratingFileName) == 0 {
		flag.Usage()
		return
	}

	ratedUsers, err := ratedUsersFromCsvFile(ratingFileName)
	if err != nil {
		log.Fatal(err)
	}
	participants, err := participantsUsersFromCsvFile(participantsFileName)
	if err != nil {
		log.Fatal(err)
	}

	ratedUsersMap := ratedUsersMap(ratedUsers)

	for i, user := range participants {
		ratedUser, ok := ratedUsersMap[user.name]
		if ok {
			user.points = ratedUser.points
		}
		participants[i] = user
	}

	sort.Slice(participants, func(index1, index2 int) bool {
		user1 := participants[index1]
		user2 := participants[index2]
		if user2.points < user1.points {
			return true
		} else if user2.points > user1.points {
			return false
		}
		return strings.Compare(user1.name, user2.name) < 0
	})

	writer := csv.NewWriter(os.Stdout)
	defer writer.Flush()

	writer.Write([]string{"number", "name", "team", "pts"})

	for i, user := range participants {
		number := i + 1

		lineArray := make([]string, 0)
		lineArray = append(lineArray, fmt.Sprintf("%v", number))
		lineArray = append(lineArray, user.name)
		lineArray = append(lineArray, user.team)
		lineArray = append(lineArray, "")

		writer.Write(lineArray)
	}
}