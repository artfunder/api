package transport

import "errors"

// ErrorNotFound ...
var ErrorNotFound = errors.New("Object Not Found")

// ErrorInternal ...
var ErrorInternal = errors.New("Internal Error")

// ErrorBadMethod ...
var ErrorBadMethod = errors.New("Method Not Allowed")
