package app

import "grest.dev/grest"

// Error returns a pointer to the errorUtil instance (eu).
// If eu is not initialized, it creates a new errorUtil instance and assigns it to eu.
// It ensures that only one instance of errorUtil is created and reused.
func Error() *errorUtil {
	if eu == nil {
		eu = &errorUtil{}
	}
	return eu
}

// eu is a pointer to an errorUtil instance.
// It is used to store and access the singleton instance of errorUtil.
var eu *errorUtil

// errorUtil represents an error utility.
// It embeds grest.Error, which indicates that errorUtil inherits from grest.Error.
type errorUtil struct {
	grest.Error
}
