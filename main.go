package main

import (
	"bufio"
	"errors"
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

var locations = map[string]string{
	"CTH": "Celestica - Thailand",
	"FAA": "Flextronics - San Jose, CA.",
	"FOC": "Foxconn - Shenzhen China",
	"JAB": "Jabil - Florida",
	"JPE": "Jabil - Malaysia",
	"JSH": "Jabil - Shanghai China",
	"PEN": "Solectron - Malaysia",
	"TAU": "Solectron - Texas",
}

type serial struct {
	serialNum string
}

func printInfo(serialnum string) {
	s, err := newSerial(serialnum)
	if err != nil {
		log.Fatal(err)
	}

	mfgDate, err := s.getMfgDate()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(s.serialNum)
	fmt.Println("Manufactured on: " + mfgDate.Format("2006-01-02"))
	fmt.Println("Manufactured in: " + s.getLocation())
	fmt.Println("")
}

func newSerial(serialnum string) (*serial, error) {
	if len(serialnum) != 11 {
		return nil, errors.New("Serial number must be exactly 11 characters")
	}

	return &serial{
		serialNum: serialnum,
	}, nil
}

func (s *serial) getLocation() string {
	// locationCode is first 3 characters
	locationCode := s.serialNum[0:3]

	// exists is a bool which will be true if the value exists in the map
	if value, exists := locations[locationCode]; exists {
		return value
	}
	return "Unknown"
}

// Splits out the year and week number out of a serial number
// Then determines the manufactured date by adding yearOffset
func (s *serial) getMfgDate() (time.Time, error) {
	// year is numbers 4 & 5 in string
	// week is numbers 6 & 7 in string
	parts := strings.Split(s.serialNum, "")
	year, _ := strconv.Atoi(parts[3] + parts[4])
	weekNum, _ := strconv.Atoi(parts[5] + parts[6])

	if weekNum > 52 {
		return time.Time{}, errors.New("Week number must not be higher than 52")
	}

	year = year + yearOffset // Add our year offset
	return isoweek.StartTime(year, weekNum, time.UTC), nil
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
			printInfo(argv.Serial)
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
				printInfo(scanner.Text())
			}
		}

		// If neither flag, print help
		if argv.Serial == "" && argv.Filename == "" {
			log.Fatal(helptext)
		}

		return nil
	})
}
