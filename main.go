package main

import (
	"math/rand"
	"time"

	"github.com/kwadwoatta/savannah-scavenger-hunt/internal"
)

func init() {
	rand.New(rand.NewSource(time.Now().UTC().UnixNano()))
}

func main() {
	internal.Execute()
}
