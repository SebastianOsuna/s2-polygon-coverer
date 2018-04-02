package main

import (
	"log"
	"s2-polygon-coverer/input"
	"strconv"
	"strings"
)

func main() {
	polygon, err := GeoJSONFromFile("./data.json")

	if err != nil {
		log.Fatal(err)
	}

	for true {
		latlng, err := input.Read("\nLat, Lng, Level:")

		if err == nil {
			latslngs := strings.Split(latlng, ",")
			lat, _ := strconv.ParseFloat(strings.TrimSpace(latslngs[0]), 64)
			lng, _ := strconv.ParseFloat(strings.TrimSpace(latslngs[1]), 64)
			level, _ := strconv.Atoi(strings.TrimSpace(latslngs[2]))

			contained, err := GeoJSONContainsPoint(polygon, lat, lng, level)

			if err != nil {
				log.Fatal(err)
			}

			if contained {
				println("IT IS INSIDE!")
			} else {
				println("NOPE")
			}
		} else {
			println("try again...")
		}

	}
}
