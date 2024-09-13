package errorsInWhale

import "log"

// If err parameter is not nil (so it contains an error), all programm will be aborted
func Check(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
