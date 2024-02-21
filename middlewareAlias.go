package httplayer

import "net/http"

// Middleware is an alias of http.HandlerFunc wrappers.
// For example the go chi compress (https://github.com/go-chi/chi/blob/1191921289e82fdc56f298a76ff254742f568ece/middleware/compress.go#L41-L44)
// is a Middleware too.
type Middleware func(h http.HandlerFunc) http.HandlerFunc
