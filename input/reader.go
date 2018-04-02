package input

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

var (
	reader *bufio.Reader
)

func init() {
	if reader == nil {
		reader = bufio.NewReader(os.Stdin)
	}
}

// Read reads from the command line
func Read(message string) (string, error) {
	fmt.Printf("%s < ", message)

	text, err := reader.ReadString('\n')

	if err != nil {
		return "", err
	}

	return strings.Replace(text, "\n", "", -1), nil
}
