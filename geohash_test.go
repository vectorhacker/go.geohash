package geohash

import (
	"math"
	"testing"
)

func TestEncode(t *testing.T) {
	type input struct {
		Lat, Lon float64
		Pres     int
	}
	type _test struct {
		Input    input
		Expected string
	}

	tests := []_test{
		{
			Input: input{
				Lat:  18.5,
				Lon:  -67.5,
				Pres: 12,
			},
			Expected: "d7rcpzzfrczy",
		},
		{
			Input: input{
				Lat:  46.75546,
				Lon:  -101.43264,
				Pres: 12,
			},
			Expected: "c8rcgze421mw",
		},
	}

	for _, test := range tests {
		b := Encode(test.Input.Lat, test.Input.Lon, &test.Input.Pres)

		if b.String() != test.Expected {
			t.Fatalf("Expected %s got back %s", test.Expected, b)
		}
	}
}
func TestBoxDecode(t *testing.T) {
	type expected struct {
		Lat, Lon float64
	}

	type input struct {
		Hash string
		Pres int
	}

	type _test struct {
		Input     input
		Expected  expected
		Tolerance float64
	}

	tests := []_test{
		{
			Input: input{
				Hash: "c8rcgze421mw",
				Pres: 12,
			},
			Expected: expected{
				Lat: 46.755460,
				Lon: -101.432640,
			},
			Tolerance: 0.000001,
		},
		{
			Input: input{
				Hash: "c8rcgze42",
				Pres: 9,
			},
			Expected: expected{
				Lat: 46.7555,
				Lon: -101.4326,
			},
			Tolerance: 0.0001,
		},
		{
			Input: input{
				Hash: "c8rcgze",
				Pres: 9,
			},
			Expected: expected{
				Lat: 46.76,
				Lon: -101.43,
			},
			Tolerance: 0.01,
		},
	}

	for _, test := range tests {
		b := Decode(test.Input.Hash, &test.Input.Pres)
		if dif := math.Abs(b.Lat() - test.Expected.Lat); dif > test.Tolerance {
			t.Fatalf("Expected %f, %f got %f, %f", test.Expected.Lat, test.Expected.Lon, b.Lat(), b.Lon())
		}
		if dif := math.Abs(b.Lon() - test.Expected.Lon); dif > test.Tolerance {
			t.Fatalf("Expected %f, %f got %f, %f", test.Expected.Lat, test.Expected.Lon, b.Lat(), b.Lon())
		}
	}
}

func TestBoxNeighbors(t *testing.T) {
	type input struct {
		hash string
		pres int
	}

	tests := []struct {
		name  string
		input input
		want  []*Box
	}{
		{
			input: input{
				hash: "de30ds",
				pres: 6,
			},
			want: []*Box{
				&Box{
					hash: "de30dt",
					pres: 6,
				},
				&Box{
					hash: "de30de",
					pres: 6,
				},
				&Box{
					hash: "de30dk",
					pres: 6,
				},
				&Box{
					hash: "de30du",
					pres: 6,
				},
				&Box{
					hash: "de30dm",
					pres: 6,
				},
				&Box{
					hash: "de30d7",
					pres: 6,
				},
				&Box{
					hash: "de30dv",
					pres: 6,
				},
				&Box{
					hash: "de30dg",
					pres: 6,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := Decode(
				tt.input.hash,
				&tt.input.pres,
			)

			got := b.Neighbors()

			for index, want := range tt.want {
				if want.hash != got[index].hash {
					t.Fatalf("Wanted %s got %s", want, got)
				}
			}
		})
	}
}
