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
	firstName string
	lastName  string
	name      string
	team      string
	// points       int64
	rating      int64
	paid        bool
	startNumber int64
	surName     string
	phone       string
	category    string
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

		firstName := strings.TrimSpace(record["firstname"])
		lastName := strings.TrimSpace(record["lastname"])

		user.name = fmt.Sprintf("%v %v", lastName, firstName)
		//user.points, _ = strconv.ParseInt(record["pts"], 10, 64)
		rating, err := strconv.ParseInt(record["number"], 10, 64)
		if err == nil {
			user.rating = rating
		} else {
			user.rating = fakeLowRating
		}
		users = append(users, user)

		dlog("Rated user: %v, %v", user.name, user.rating)
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

		user.paid = true //len(record["Оплата_"]) != 0
		// if !user.paid {
		// 	continue
		// }
		user.firstName = strings.TrimSpace(record["Имя"])
		user.lastName = strings.TrimSpace(record["Фамилия"])
		user.name = fmt.Sprintf("%v %v", user.lastName, user.firstName)
		user.team = strings.TrimSpace(record["Клуб/команда"])
		user.rating = fakeLowRating

		category := strings.TrimSpace(record["Категория"])
		user.category = category
		// categoryFields := strings.Fields(category)
		// if len(categoryFields) > 0 {
		// user.category = categoryFields[0]
		// }

		user.surName = strings.TrimSpace(record["Отчество"])
		user.phone = strings.TrimSpace(record["Телефон"])

		users = append(users, user)

		dlog("Participant: %v, %v", user.name, user.team)
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
	dumpNumbers          = false
	startList            = false
)

func main() {

	flag.StringVar(&participantsFileName, "p", "", "Participants csv file")
	flag.StringVar(&ratingFileName, "r", "", "Rating csv file")
	flag.BoolVar(&dumpNumbers, "dump", false, "Dump numbers")
	flag.BoolVar(&startList, "startList", false, "Generate start list")

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
	allUsers := make([]User, len(participants))
	copy(allUsers, participants)

	for i, user := range participants {
		ratedUser, ok := ratedUsersMap[user.name]
		if ok {
			//user.points = ratedUser.points
			user.rating = ratedUser.rating
		}
		allUsers[i] = user
	}

	for _, user := range allUsers {
		dlog("User: %v, rating:%v", user.name, user.rating)
	}

	sortedUsers := make([]User, len(allUsers))
	copy(sortedUsers, allUsers)

	sort.Slice(sortedUsers, func(index1, index2 int) bool {
		user1 := sortedUsers[index1]
		user2 := sortedUsers[index2]
		if user2.rating > user1.rating {
			return true
		} else if user2.rating < user1.rating {
			return false
		} else {
			return strings.Compare(user1.name, user2.name) < 0
		}
	})

	for _, user := range sortedUsers {
		dlog("Sorted user: %v, rating:%v", user.name, user.rating)
	}

	allUsersMap := make(map[string]User)

	number := 1

	for i := 0; i < len(sortedUsers); i++ {

		user := sortedUsers[i]

		if user.paid {
			user.startNumber = int64(number)
			number += 1
		}

		allUsersMap[user.name] = user
		sortedUsers[i] = user
	}

	if dumpNumbers {
		for _, user := range participants {
			dlog("%v", allUsersMap[user.name].startNumber)
		}
	} else {
		writer := csv.NewWriter(os.Stdout)
		defer writer.Flush()

		if startList {
			writer.Write([]string{"Фамилия Имя", "Категория", "Номер", "Оплата"})

			for _, participantUser := range participants {

				user, ok := allUsersMap[participantUser.name]
				if !ok {
					continue
				}

				lineArray := make([]string, 0)
				lineArray = append(lineArray, user.name)
				lineArray = append(lineArray, user.category)
				startNumberString := ""
				if user.startNumber != 0 {
					startNumberString = fmt.Sprintf("%v", user.startNumber)
				}
				lineArray = append(lineArray, startNumberString)
				userPaidString := ""
				if user.paid {
					userPaidString = "+"
				}
				lineArray = append(lineArray, userPaidString)

				writer.Write(lineArray)
			}
		} else {
			writer.Write([]string{"number", "name", "team", "pts"})

			for _, user := range sortedUsers {
				lineArray := make([]string, 0)
				lineArray = append(lineArray, fmt.Sprintf("%v", user.startNumber))
				lineArray = append(lineArray, user.name)
				lineArray = append(lineArray, user.team)
				lineArray = append(lineArray, "")

				writer.Write(lineArray)
			}
		}

	}
}
