package err

import "log"

// Must is a generic function used to crash the program on
// errors that may indicate a bug or fatal error requiring
// correction.
//
// Type Parameters:
//   - T: any type
//
// Parameters:
//   - X: T - the value to be returned if no error is encountered
//   - err: error - the error to be checked
//
// Returns:
//   - T: the value X if err is nil; otherwise, logs the error and exits
func Must[T any](X T, err error) T {
	if err != nil {
		log.Fatalln(err)
	}
	return X
}
