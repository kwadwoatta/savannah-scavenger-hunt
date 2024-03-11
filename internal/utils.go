package internal

import (
	"math/rand"
)

func generateUniqueLocations(existing []Location, maxX int, maxY int, length int) []Location {
	newLocations := make([]Location, 0, length)
	locationMap := make(map[Location]bool)

	// Add existing locations to the map
	for _, loc := range existing {
		locationMap[loc] = true
	}

	for i := 0; i < length; i++ {
		var newLoc Location
		for {
			newLoc.X = rand.Intn(maxX-3) + 3
			newLoc.Y = rand.Intn(maxY-3) + 3
			if !locationMap[newLoc] {
				locationMap[newLoc] = true
				break
			}
		}
		newLocations = append(newLocations, newLoc)
	}

	return newLocations
}
