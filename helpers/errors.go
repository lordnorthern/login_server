package helpers

import (
	"fmt"
)

// LogError will log errors.
// This function will allow to add a more sophisticated error handling down the road.
func LogError(err error) {
	fmt.Println(err)
}
