package helpers

import (
	"fmt"
)

// LogError will log errors.
// This function will allow to add a more sophisticated error handling down the road.
func LogError(err error, where ...string) {
	fmt.Println(err, where)
}
