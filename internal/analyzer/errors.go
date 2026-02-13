package analyzer

import "errors"

// ErrUnrecognizedFormat is returned by Dispatcher when no parser recognizes the binary format.
var ErrUnrecognizedFormat = errors.New("unrecognized binary format")
