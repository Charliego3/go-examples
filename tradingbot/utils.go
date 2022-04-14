package main

import (
	"math/rand"
	"time"
)

var ur = rand.New(rand.NewSource(time.Now().UnixNano()))

func RandomRange[T ~int | ~int32 | ~int64](min, max T) T {
	diff := max - min
	r := Random(diff)
	return T(r) + min
}

func Random[T ~int | ~int32 | ~int64](limit T) T {
	return T(ur.Int63n(int64(limit)))
	// var r T
	// switch t := limit.(type) {
	// case int:
	// 	r = T(ur.Intn(t))
	// case int32:
	// 	r = T(ur.Int31n(t))
	// case int64:
	// 	r = T(ur.Int63n(t))
	// }
	// return r
}
