package main

import (
	"flag"
	"time"
	"github.com/bearbin/go-age"
	"os"
	"encoding/csv"
	"log"
	"fmt"
)

func dlog(format string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, format, args...)
	fmt.Fprintf(os.Stderr, "\n")
}

func categoryFromDateAndGender(dateString, gender string) string {
	layout := "2.1.2006"
	birthTime, _ := time.Parse(layout, dateString)
	eventTime, _ := time.Parse(layout, "16.09.2018")

	ageInt := age.AgeAt(birthTime, eventTime)

	//dlog("date = %v  gender = %v", dateString, gender)
	//dlog("birthdate = %v, age = %v", birthTime, ageInt)

	if gender == "Ж" {
		if ageInt < 40 {
			return "Ж18-39"
		} else {
			return "Ж40+"
		}
	} else {
		if ageInt < 23 {
			return "М18-22"
		} else if ageInt < 30 {
			return "М23-29"
		} else if ageInt < 40 {
			return "М30-39"
		} else if ageInt < 50 {
			return "М40-49"
		} else {
			return "М50+"
		}
	}
	return ""
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

func participantsUsersFromCsvFile(csvFilePath string) (mapRecords []map[string]string, err error) {

	records, err := readCsvFile(csvFilePath)
	if err != nil {
		return nil, err
	}

	mapRecords = csvRecordsToMap(records)

	return
}


func main() {

	participantsFileName := ""
	flag.StringVar(&participantsFileName, "p", "", "Participants csv file")

	flag.Parse()

	if len(participantsFileName) == 0 {
		flag.Usage()
		return
	}

	users, err := participantsUsersFromCsvFile(participantsFileName)
	if err != nil {
		log.Fatal(err)
	}

	for _, user := range users {
		dateString := user["Дата рождения"]
		gender := user["Пол"]
		category := categoryFromDateAndGender(dateString, gender)
		fmt.Println(category)
	}
}