package parser

import (
	"bufio"
	"io"
	"os"
)

func ParseAddrFile(filename string) (addr []string, err error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer func() {
		if deferErr := file.Close(); deferErr != nil {
			err = deferErr
		}
	}()

	reader := bufio.NewReader(file)
	for {
		line, _, err := reader.ReadLine()
		if err == io.EOF {
			break
		}

		addr = append(addr, string(line))
	}

	return addr, nil

}