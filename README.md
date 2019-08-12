GeoHash Algororithm in Go
=====

Example

```go
package main

import "github.com/vectorhacker/go.geohash"

func main() {
  box := geohash.Encode(46.7666, -101.4650, 0) // 0 is the same as full presicion

  boxes := box.Neighbors()

  boxesBoxes := boxes[0].Neighbors()


  // ....

  box1 := geohash.Decode("c8rf51e7n", 12)
  box1.Neighbors()

  // ...
}
```
