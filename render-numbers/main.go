package main

import (
	"encoding/csv"
	"errors"
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

func dlog(format string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, format, args...)
	fmt.Fprintf(os.Stderr, "\n")
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

var (
	participantsFileName = ""
	bgFileName           = ""
	outDir               = ""
	tmpDir               = ""
)

func main() {

	flag.StringVar(&participantsFileName, "p", "", "Participants csv file")
	flag.StringVar(&bgFileName, "bg", "", "Background pdf file")
	flag.StringVar(&outDir, "out", "out", "Output dir")
	flag.StringVar(&tmpDir, "tmp", "tmp", "Tmp dir")

	flag.Parse()

	if len(participantsFileName) == 0 || len(bgFileName) == 0 {
		flag.Usage()
		return
	}

	users, err := participantsUsersFromCsvFile(participantsFileName)
	if err != nil {
		log.Fatal(err)
	}

	for i := 0; i < 250; i++ {

		var name string
		number := i + 1
		var numberString string
		var team string

		if i < len(users) {
			user := users[i]

			numberString = user["number"]
			name = user["name"]
			team = user["team"]

			parsedNumber, _ := strconv.ParseInt(numberString, 10, 64)
			if parsedNumber != int64(number) {
				log.Fatalln("Numbers don't match")
			}
		} else {
			numberString = fmt.Sprintf("%v", number)
		}

		tmpOutputFileName := fmt.Sprintf("%v/%03d.pdf", tmpDir, number)

		out, err := runProgram1(true, true, "../start-number-draw/start-number-draw", "--number", numberString, "--name", name, "--team", team, "-o", tmpOutputFileName)
		if err != nil {
			dlog("Error: %v\n", err)
			dlog("%v", out)
		}

		outputFileName := fmt.Sprintf("%v/%03d.pdf", outDir, number)

		out, err = runProgram1(true, true, "pdftk", tmpOutputFileName, "background", bgFileName, "output", outputFileName)
		if err != nil {
			dlog("Error: %v\n", err)
			dlog("%v", out)
		}
	}
}
