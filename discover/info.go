package discover

import (
	"errors"
	"regexp"

	"github.com/albenik/go-serial/v2"
)

// Fitkit Struktura obsahující informace o nalezeném FITkitu
type Fitkit struct {
	Port     string // Sériový port, na kterém se FITkit nachází
	Version  string // Verze FITkitu (většinou 1.x nebo 2.x)
	Revision string // Číslo revize FITkitu
}

// GetFitkitInfo - Vrátí informaci o FITkitu na daném portu
// Pokud nastane chyba při čtení vrátí chybu nebo vrátí chybu,
// pokud se nejedná o validní port FITkitu
func GetFitkitInfo(portName string) (*Fitkit, error) {
	fitkitRegex, _ := regexp.Compile("FITkit (.+) \\$Rev: (\\d+) \\$")

	port, err := serial.Open(
		portName,
		serial.WithBaudrate(460800),
		serial.WithDataBits(8),
		serial.WithStopBits(serial.OneStopBit),
		serial.WithParity(serial.NoParity),
		serial.WithHUPCL(true),
		serial.WithReadTimeout(500),
	)

	if err != nil {
		return nil, err
	}

	defer port.Close()

	buff := make([]byte, 64)
	received := 0

	for {
		n, err := port.Read(buff[received:])
		if err != nil {
			return nil, err
		}

		received += n

		submatch := fitkitRegex.FindSubmatch(buff)
		if len(submatch) >= 3 {
			return &Fitkit{
				Port:     portName,
				Version:  string(submatch[1]),
				Revision: string(submatch[2]),
			}, nil
		}

		if n == 0 || received >= len(buff) {
			return nil, errors.New("Unknown data")
		}
	}
}
