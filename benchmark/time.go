package benchmark

import (
	"log"
	"time"
)

// TimeTrack tracks time sine the given time
// credits to: https://coderwall.com/p/cp5fya/measuring-execution-time-in-go
func TimeTrack(start time.Time, name string) {
	elapsed := time.Since(start)
	log.Printf("%s took %s", name, elapsed)
}
