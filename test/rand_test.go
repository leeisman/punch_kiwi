package test

import (
	"fmt"
	"math"
	"math/rand"
	"testing"
	"time"
)

func init() {

}

func Test_RandInt64(t *testing.T) {
	min := int64(-5000000000)
	max := int64(5000000000)
	float := math.Pow(10, -14)
	fmt.Println(float)
	rand.Seed(time.Now().UnixNano())
	randomNumber := rand.Int63n(max-min) + min
	fmt.Println(fmt.Sprintf("%v",float64(randomNumber) * float))
}
