// Package geohash is a simple implementation of the Public Domain Geohash algorithm.
// It represents the hashes using boxes which contain the resulting hashes and can reverse the box,
// using the geohash and the precision to recreate it.
package geohash

import (
	"strings"
)

const (
	maxlat = 90.0
	maxlon = 180.0
	minlat = -90.0
	minlon = -180.0
	base32 = "0123456789bcdefghjkmnpqrstuvwxyz"
)

var bits = []int{16, 8, 4, 2, 1}

// Box represents a coordinate box. It knows it's width and height, and can figure out it's neighbors.
// It implements the Stringer interface, printing out the geohash value.
type Box struct {
	hash           string
	pres           int
	maxLat, maxLon float64
	minLat, minLon float64
}

// Decode creates a new box from an initial hash
func Decode(hash string) *Box {

	refine := func(minInterval, maxInterval float64, cd, mask int) (float64, float64) {
		if cd&mask == 0 {
			maxInterval = (minInterval + maxInterval) / 2
		} else {
			minInterval = (minInterval + maxInterval) / 2
		}

		return minInterval, maxInterval
	}

	p := len(hash)

	isEven := true
	lat0, lat1 := minlat, maxlat
	lon0, lon1 := minlon, maxlon
	latErr, lonErr := maxlat, maxlon

	for i := 0; i < len(hash); i++ {
		c := hash[i]

		cd := strings.Index(string(base32), string(c))

		for j := 0; j < 5; j++ {
			mask := bits[j]

			if isEven {
				lonErr /= 2
				lon0, lon1 = refine(lon0, lon1, cd, mask)
			} else {
				latErr /= 2
				lat0, lat1 = refine(lat0, lat1, cd, mask)
			}
			isEven = !isEven
		}
	}

	return &Box{hash: hash, pres: p, minLat: lat0, maxLat: lat1, minLon: lon0, maxLon: lon1}
}

// Encode takes coordinates and makes a geohash box
func Encode(lat, lon float64, pres int) *Box {
	hash := ""
	precision := 12
	if pres > 0 {
		precision = pres
	}

	isEven := true
	bit := 0
	ch := 0
	lat0, lat1 := minlat, maxlat
	lon0, lon1 := minlon, maxlon

	for len(hash) < precision {
		if isEven {
			mid := (lon0 + lon1) / 2.0

			if lon > mid {
				ch |= bits[bit]
				lon0 = mid
			} else {
				lon1 = mid
			}
		} else {
			mid := (lat0 + lat1) / 2.0

			if lat > mid {
				ch |= bits[bit]
				lat0 = mid
			} else {
				lat1 = mid
			}
		}

		isEven = !isEven

		if bit < 4 {
			bit++
		} else {
			hash += string(base32[ch])
			bit = 0
			ch = 0
		}
	}

	return &Box{hash: hash, maxLat: lat1, minLat: lat0, maxLon: lon1, minLon: lon0}
}

// Lat returns the latitude of a box
func (b Box) Lat() float64 {
	return (b.minLat + b.maxLat) / 2
}

// Lon returns the longitude of a box
func (b Box) Lon() float64 {
	return (b.minLon + b.maxLon) / 2
}

// Height returns the height of the box
func (b Box) Height() float64 {
	return b.maxLat - b.minLat
}

// Width returns the width of the box
func (b Box) Width() float64 {
	return b.maxLon - b.minLon
}

// Precision returns the precision of a box
func (b Box) Precision() int {
	return b.pres
}

// Neighbors calculates the 8 neighboring boxes of a box
func (b Box) Neighbors() []*Box {

	var (
		// Directly adjecent
		up    = Encode(b.Lat()+b.Height(), b.Lon(), b.Precision())
		down  = Encode(b.Lat()-b.Height(), b.Lon(), b.Precision())
		left  = Encode(b.Lat(), b.Lon()-b.Width(), b.Precision())
		right = Encode(b.Lat(), b.Lon()+b.Width(), b.Precision())

		// Corners
		upleft    = Encode(b.Lat()+b.Height(), b.Lon()-b.Width(), b.Precision())
		downleft  = Encode(b.Lat()-b.Height(), b.Lon()-b.Width(), b.Precision())
		upright   = Encode(b.Lat()+b.Height(), b.Lon()+b.Width(), b.Precision())
		downright = Encode(b.Lat()-b.Height(), b.Lon()+b.Width(), b.Precision())
	)

	return []*Box{up, down, left, right, upleft, downleft, upright, downright}
}

// String returns the geohash of the box
func (b Box) String() string {
	return b.hash
}


// Hash returns the geohash
func (b Box) Hash() string {
	return b.hash
}