package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/mkideal/cli"
)

var (
	yearOffset = 1996
)

type serial struct {
	serialNum string
	year      int
	weeknum   int
	mfgDate   time.Time
}

// Splits out the year and week number out of a serial number
func (s *serial) splitNums() {
	if len(s.serialNum) != 11 {
		log.Fatal("Serial number must be exactly 11 characters")
	}

	parts := strings.Split(s.serialNum, "")
	s.year, _ = strconv.Atoi(parts[3] + parts[4])
	s.weeknum, _ = strconv.Atoi(parts[5] + parts[6])
	if s.weeknum > 52 {
		log.Fatal("Week number must not be higher than 52")
	}
}

// Calculates the Monday of a given year and week number
func (s *serial) firstDayOfISOWeek() {
	timezone, _ := time.LoadLocation("UTC")
	date := time.Date(s.year, 0, 0, 0, 0, 0, 0, timezone)
	isoYear, isoWeek := date.ISOWeek()

	// iterate back to Monday
	for date.Weekday() != time.Monday {
		date = date.AddDate(0, 0, -1)
		isoYear, isoWeek = date.ISOWeek()
	}

	// iterate forward to the first day of the first week
	for isoYear < s.year {
		date = date.AddDate(0, 0, 7)
		isoYear, isoWeek = date.ISOWeek()
	}

	// iterate forward to the first day of the given week
	for isoWeek < s.weeknum {
		date = date.AddDate(0, 0, 7)
		isoYear, isoWeek = date.ISOWeek()
	}

	s.mfgDate = date
}

// Parses the given serial and prints out its manufactured date
func (s *serial) parseSerial() {
	fmt.Println(s.serialNum)
	s.splitNums()
	s.year = s.year + yearOffset // Add our year offset
	s.firstDayOfISOWeek()
	fmt.Println(s.mfgDate.Format("2006-01-02"))
}

type argT struct {
	cli.Helper
	Serial   string `cli:"s,serial"   usage:"serial to parse"`
	Filename string `cli:"f,filename" usage:"filename with serials to parse"`
}

var (
	helptext = `
Parses provided cisco serial and returns manufactured date

Examples:
go run serial_to_date.go --serial FAA04459FNI
go run serial_to_date.go --filename serials.txt
`
)

func main() {
	cli.Run(new(argT), func(ctx *cli.Context) error {
		argv := ctx.Argv().(*argT)

		// Read single serial
		if argv.Serial != "" {
			s := serial{serialNum: argv.Serial}
			s.parseSerial()
		}

		// Read filename
		if argv.Filename != "" {
			file, err := os.Open(argv.Filename)
			if err != nil {
				log.Fatal(err)
			}
			defer file.Close()

			scanner := bufio.NewScanner(file)
			for scanner.Scan() {
				s := serial{serialNum: scanner.Text()}
				s.parseSerial()
				fmt.Println("")
			}
		}

		// If neither flag, print help
		if argv.Serial == "" && argv.Filename == "" {
			log.Fatal(helptext)
		}

		return nil
	})
}
