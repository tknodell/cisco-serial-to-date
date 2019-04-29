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
	serialNum, location string
	year, weekNum       int
	mfgDate             time.Time
}

// Parses the given serial and prints out extracted info
func (s *serial) parseSerial() {
	fmt.Println(s.serialNum)

	if len(s.serialNum) != 11 {
		log.Fatal("Serial number must be exactly 11 characters")
	}

	s.getLocation()
	s.getMfgDate()

	fmt.Println("Manufactured on: " + s.mfgDate.Format("2006-01-02"))
	fmt.Println("Manufactured in: " + s.location)
}

func (s *serial) getLocation() {
	// locationCode is first 3 characters
	locationCode := s.serialNum[0:3]

	// Location codes
	locations := map[string]string{
		"CTH": "Celestica - Thailand",
		"FAA": "Flextronics - San Jose, CA.",
		"FOC": "Foxconn - Shenzhen China",
		"JAB": "Jabil - Florida",
		"JPE": "Jabil - Malaysia",
		"JSH": "Jabil - Shanghai China",
		"PEN": "Solectron - Malaysia",
		"TAU": "Solectron - Texas",
	}

	// exists is a bool which will be true if the value exists in the map
	if value, exists := locations[locationCode]; exists {
		s.location = value
	} else {
		s.location = "Unknown"
	}
}

// Splits out the year and week number out of a serial number
// Then determines the manufactured date by adding yearOffset
func (s *serial) getMfgDate() {
	// year is numbers 4 & 5 in string
	// week is numbers 6 & 7 in string
	parts := strings.Split(s.serialNum, "")
	s.year, _ = strconv.Atoi(parts[3] + parts[4])
	s.weekNum, _ = strconv.Atoi(parts[5] + parts[6])

	if s.weekNum > 52 {
		log.Fatal("Week number must not be higher than 52")
	}

	s.year = s.year + yearOffset // Add our year offset
	s.mfgDate = isoweek.StartTime(s.year, s.weekNum, time.UTC)
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
./serial_to_date --serial FAA04459FNI
./serial_to_date --filename serials.txt
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
