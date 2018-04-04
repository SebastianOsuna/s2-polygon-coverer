package coverer

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"math"
	"time"

	"github.com/golang/geo/s2"
	"github.com/sebastianosuna/s2-polygon-coverer/benchmark"
)

// GeoJSON Represents a geojson feature
type GeoJSON struct {
	FeatureType    string        `json:"type"`
	RawCoordinates [][][]float64 `json:"coordinates"`
}

// Coordinates returns the polygon coordinates from the geo json
func (g *GeoJSON) Coordinates() [][]float64 {
	return g.RawCoordinates[0]
}

// GeoJSONFromJSON returns a parsed GeoJson object from an JSON string
func GeoJSONFromJSON(jsonStr []byte) (*GeoJSON, error) {
	var data GeoJSON
	jsonerr := json.Unmarshal(jsonStr, &data)

	if jsonerr != nil {
		return nil, jsonerr
	}

	if data.FeatureType != "Polygon" {
		return nil, errors.New("GeoJsonFeatureType: only \"Polygon\" type supported")
	}

	return &data, nil
}

// GeoJSONFromFile loads a geoJson from the given file path
func GeoJSONFromFile(filePath string) (*GeoJSON, error) {
	defer benchmark.TimeTrack(time.Now(), "GeoJSONFromFile")
	raw, err := ioutil.ReadFile(filePath)

	if err != nil {
		return nil, err
	}

	return GeoJSONFromJSON(raw)
}

// CoverPolygon returns a list of CellIds (tokens) of the given level that cover the given GeoJSON
func CoverPolygon(p *GeoJSON, level int) []string {
	defer benchmark.TimeTrack(time.Now(), "coverPolygon")
	points := p.Coordinates()
	s2points := make([]s2.Point, len(points))

	for i, point := range points {
		lat := point[1]
		lng := point[0]
		s2point := s2.PointFromLatLng(s2.LatLngFromDegrees(lat, lng))
		s2points[i] = s2point
	}

	loops := make([]*s2.Loop, 1)
	myloop := s2.LoopFromPoints(s2points)
	loops[0] = myloop
	polygon := s2.PolygonFromLoops(loops)
	region := s2.Region(polygon)

	level = level - int(math.Trunc(myloop.Area()*200))
	fmt.Printf("LLLLEVEL %d, %f", level, myloop.Area())

	coverer := &s2.RegionCoverer{MaxLevel: level, MinLevel: level}
	cellunion := coverer.Covering(region)
	cellIds := make([]string, len(cellunion))

	for i, cell := range cellunion {
		cellIds[i] = cell.ToToken()
	}

	return cellIds
}

// CoverPolygonAsMap returns a map having the keys as the CellIDs covering the polygon. true is always the value
func CoverPolygonAsMap(polygon *GeoJSON, level int) map[string]bool {
	defer benchmark.TimeTrack(time.Now(), "coverPolygonAsMap")
	cells := CoverPolygon(polygon, level)

	set := make(map[string]bool, len(cells))

	for _, cell := range cells {
		set[cell] = true
	}

	return set
}

func latlngToCellID(lat, lng float64, level int) string {
	latlng := s2.LatLngFromDegrees(lat, lng)
	cell := s2.CellFromLatLng(latlng).ID()

	return cell.Parent(level).ToToken()
}

// GeoJSONContainsPoint checks if the given lat, lng are inside the given GeoJson polygon
func GeoJSONContainsPoint(geo *GeoJSON, lat, lng float64, level int) (bool, error) {
	if level <= 0 {
		return false, errors.New("S2Level: negative or zero level is not valid")
	}

	cover := CoverPolygonAsMap(geo, level)
	cellid := latlngToCellID(lat, lng, level)
	_, contained := cover[cellid]

	return contained, nil
}

func main() {}
