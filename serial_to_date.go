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
	"github.com/snabb/isoweek"
)

const yearOffset = 1996

type serial struct {
	serialNum     string
	year, weekNum int
	mfgDate       time.Time
}

// Splits out the year and week number out of a serial number
func (s *serial) splitNums() {
	if len(s.serialNum) != 11 {
		log.Fatal("Serial number must be exactly 11 characters")
	}

	parts := strings.Split(s.serialNum, "")
	s.year, _ = strconv.Atoi(parts[3] + parts[4])
	s.weekNum, _ = strconv.Atoi(parts[5] + parts[6])
	if s.weekNum > 52 {
		log.Fatal("Week number must not be higher than 52")
	}
}

// Parses the given serial and prints out its manufactured date
func (s *serial) parseSerial() {
	fmt.Println(s.serialNum)
	s.splitNums()
	s.year = s.year + yearOffset // Add our year offset
	s.mfgDate = isoweek.StartTime(s.year, s.weekNum, time.UTC)
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
