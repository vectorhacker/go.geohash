package main

import "math"
import "github.com/vectorhacker/geohash"
import "fmt"
import "flag"

func haversine(θ float64) float64 {
	return (1 - math.Cos(θ)) / 2
}

type pos struct {
	φ float64 // latitude, radians
	ψ float64 // longitude, radians
}

func degPos(lat, lon float64) pos {
	return pos{lat * math.Pi / 180.0, lon * math.Pi / 180.0}
}

const rEarth = 6378.137 // km

var (
	lat       = flag.Float64("lat", 0.0, "")
	lon       = flag.Float64("lon", 0.0, "")
	precision = flag.Int("precision", 12, "")
)

func hsDist(p1, p2 pos) float64 {
	return 2 * rEarth * math.Asin(math.Sqrt(haversine(p2.φ-p1.φ)+
		math.Cos(p1.φ)*math.Cos(p2.φ)*haversine(p2.ψ-p1.ψ)))
}

func main() {
	flag.Parse()
	boxing(geohash.Encode(*lat, *lon, *precision))

}

func boxing(box *geohash.Box) {

	if box.Precision() < 1 {
		return
	}

	fmt.Println("--------------------")

	neighbors := box.Neighbors()
	fmt.Printf("Distnace from neighbors\n")

	for _, n := range neighbors {
		fmt.Printf("%f km\n", hsDist(degPos(n.Lat(), n.Lon()), degPos(box.Lat(), box.Lon())))
	}

	latDis := hsDist(degPos(box.Lat()+box.Height()/2, box.Lon()), degPos(box.Lat()-box.Height()/2, box.Lon()))
	lonDis := hsDist(degPos(box.Lat(), box.Lon()+box.Width()/2), degPos(box.Lat(), box.Lon()-box.Width()/2))
	fmt.Printf("Precision %d\t box sized\t %f km X %f km\n", box.Precision(), latDis, lonDis)
	fmt.Printf("lat = %f, lon = %f\n", box.Lat(), box.Lon())

	fmt.Println("--------------------")
	
	boxing(geohash.Encode(box.Lat(), box.Lon(), box.Precision()-1))
	
}
