package rt

var METHOD = struct {
	GET, POST, HEAD, PUT, PATCH, DELETE, OPTIONS, TRACE, CONNECT string
}{
	GET:     "GET",
	POST:    "POST",
	HEAD:    "HEAD",
	PUT:     "PUT",
	PATCH:   "PATCH",
	DELETE:  "DELETE",
	OPTIONS: "OPTIONS",
	TRACE:   "TRACE",
	CONNECT: "CONNECT",
}
