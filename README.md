GeoHash algororithm in Go
=====

Example

```go
package main

import "github.com/vectorhacker/geohash"

func main() {
  box := geohash.Encode(46.7666, -101.4650, nil)

  boxes := box.Neighbors()

  boxesBoxes := boxes[0].Neighbors()


  // ....

  box1 := geohash.Decode("c8rf51e7n", 12)
  box1.Neighbors()

  // ...
}
```