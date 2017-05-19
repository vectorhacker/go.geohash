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

// Box represents a coordinate box
type Box struct {
	hash           string
	pres           int
	maxLat, maxLon float64
	minLat, minLon float64
}

func (b Box) String() string {
	return b.hash
}

// Lat returns the latitude of a box
func (b Box) Lat() float64 {
	return (b.minLat + b.maxLat) / 2
}

// Lon returns the longitude of a box
func (b Box) Lon() float64 {
	return (b.minLon + b.maxLon) / 2
}

// MaxLat returns the maximum latitude for a box
func (b Box) MaxLat() float64 {
	return b.maxLat
}

// MinLat returns the minimum latitude for a box
func (b Box) MinLat() float64 {
	return b.minLat
}

// MaxLon returns the maximum longitude for a box
func (b Box) MaxLon() float64 {
	return b.maxLon
}

// MinLon returns the minimum longitude for a box
func (b Box) MinLon() float64 {
	return b.minLon
}

// Height Calculates the height of a box
func (b Box) Height() float64 {
	return b.MaxLat() - b.MinLat()
}

// Width calculates the width of a box
func (b Box) Width() float64 {
	return b.MaxLon() - b.MinLon()
}

// Neighbors calculates the 8 neighbors of a box
func (b Box) Neighbors() []*Box {

	var (
		// Directly adjecent
		up    = Encode(b.Lat()+b.Height(), b.Lon(), &b.pres)
		down  = Encode(b.Lat()-b.Height(), b.Lon(), &b.pres)
		left  = Encode(b.Lat(), b.Lon()-b.Width(), &b.pres)
		right = Encode(b.Lat(), b.Lon()+b.Width(), &b.pres)

		// Corners
		upleft    = Encode(b.Lat()+b.Height(), b.Lon()-b.Width(), &b.pres)
		downleft  = Encode(b.Lat()-b.Height(), b.Lon()-b.Width(), &b.pres)
		upright   = Encode(b.Lat()+b.Height(), b.Lon()+b.Width(), &b.pres)
		downright = Encode(b.Lat()-b.Height(), b.Lon()+b.Width(), &b.pres)
	)

	return []*Box{up, down, left, right, upleft, downleft, upright, downright}
}

// Decode creates a new box from an initial hash
func Decode(hash string, pres *int) *Box {

	refine := func(i0, i1 float64, cd, mask int) (float64, float64) {
		if cd&mask == 0 {
			i1 = (i0 + i1) / 2
		} else {
			i0 = (i0 + i1) / 2
		}

		return i0, i1
	}

	p := 12
	if *pres > 0 {
		p = *pres
	}

	isEven := true
	lat0, lat1 := minlat, maxlat
	lon0, lon1 := minlon, maxlon
	latErr, lonErr := 90.0, 180.0

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
func Encode(lat, lon float64, pres *int) *Box {
	hash := ""
	precision := 12
	if *pres != 0 {
		precision = *pres
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
