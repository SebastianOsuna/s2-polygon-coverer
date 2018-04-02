package main

import (
	"reflect"
	"testing"
)

func TestGeoJSONFromJSON(t *testing.T) {
	json := "{\"type\":\"Polygon\",\"coordinates\":[[[1,1], [2,2]]]}"
	parsed, err := GeoJSONFromJSON([]byte(json))

	if err != nil {
		t.Error(err)
	}

	if parsed.FeatureType != "Polygon" {
		t.Error("Expected Polygon as FeatureType, got", parsed.FeatureType)
	}

	if len(parsed.Coordinates()) != 2 {
		t.Error("Expected to see 2 coordinates, got", len(parsed.Coordinates()))
	}
}

func TestGeoJSONFromJSONValidations(t *testing.T) {
	json := "{\"type\":\"MultiPolygon\",\"coordinates\":[[[1,1], [2,2]]]}"
	parsed, err := GeoJSONFromJSON([]byte(json))

	if err == nil || parsed != nil {
		t.Error("Only Polygon type should be supported")
	}
}

func TestGeoJSONParsing(t *testing.T) {
	invalidJSON := "{someNonValidJson}"
	_, err := GeoJSONFromJSON([]byte(invalidJSON))

	if err == nil {
		t.Error("Expected JSON parsing to fail")
	}
}

func TestGeoJSONCoordinates(t *testing.T) {
	json := "{\"type\":\"Polygon\",\"coordinates\":[[[1,2], [3,4]]]}"
	parsed, _ := GeoJSONFromJSON([]byte(json))
	coords := parsed.Coordinates()

	if coords[0][0] != 1 || coords[0][1] != 2 {
		t.Error("Expected to see (1, 2) as first coordinate, got", coords[0])
	}

	if coords[1][0] != 3 || coords[1][1] != 4 {
		t.Error("Expected to see (3, 4) as second coordinate, got", coords[0])
	}
}

func TestGeoJSONFromFile(t *testing.T) {
	filePath := "./data.json"
	geo, err := GeoJSONFromFile(filePath)

	if err != nil {
		t.Error(err)
	}

	if geo.FeatureType != "Polygon" {
		t.Error("Expected Polygon as FeatureType, got", geo.FeatureType)
	}
}

func TestGeoJSONFromFileNotFound(t *testing.T) {
	filePath := "oops.json"
	_, err := GeoJSONFromFile(filePath)

	if err == nil {
		t.Error("Should fail with non-existing file")
	}
}

func TestCoverPolygon(t *testing.T) {
	geo, _ := GeoJSONFromFile("./data.json")
	expects := []string{"8ef6254", "8ef625c", "8ef62f4", "8ef62fc", "8ef63ac"}

	cells := CoverPolygon(geo, 11)

	if !reflect.DeepEqual(expects, cells) {
		t.Errorf("Got wrong cellIDs, expected %v, got %v", expects, cells)
	}
}

func TestCoverPolygonAsMap(t *testing.T) {
	geo, _ := GeoJSONFromFile("./data.json")
	expects := map[string]bool{
		"8ef6254": true,
		"8ef625c": true,
		"8ef62f4": true,
		"8ef62fc": true,
		"8ef63ac": true,
	}

	cells := CoverPolygonAsMap(geo, 11)

	if !reflect.DeepEqual(expects, cells) {
		t.Errorf("Got wrong cellIDs, expected %v, got %v", expects, cells)
	}
}

func TestGeoJsonContainsPoint(t *testing.T) {
	geo, _ := GeoJSONFromFile("./data.json")
	contained1, _ := GeoJSONContainsPoint(geo, 10.418979, -75.553494, 17)
	contained2, _ := GeoJSONContainsPoint(geo, 10.418979, -75.55394, 17)

	if !contained1 {
		t.Error("Expected to see 10.418979, -75.553494 inside the S2 polygon")
	}

	if contained2 {
		t.Error("Expected NOT to see 10.418979, -75.55394 inside the S2 polygon")
	}

	_, err := GeoJSONContainsPoint(geo, 10.418979, -75.55394, -2)

	if err == nil {
		t.Error("Expected error on non-positive level")
	}
}
