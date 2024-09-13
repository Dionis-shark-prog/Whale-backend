package additionalForWhale

import (
	"WhaleWebSite/pkg/errorsInWhale"
	"bufio"
	"os"
)

func GetStrings(fileName string) []string {
	var lines []string
	file, err := os.Open(fileName)
	if os.IsNotExist(err) {
		return nil
	}
	errorsInWhale.Check(err)
	defer file.Close()
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	errorsInWhale.Check(scanner.Err())
	return lines
}
