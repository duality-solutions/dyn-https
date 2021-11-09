package util

import (
	"math/rand"
	"time"
)

func RandomUIntRange(umin, umax uint) uint {
	rand.Seed(time.Now().UnixNano())
	min := int(umin)
	max := int(umax)
	return uint(rand.Intn(max-min+1) + min)
}
