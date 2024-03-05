package rio

import (
	"io/fs"
	"net/http"
)

// ------------------------------------------------------------------
//
//
// Basic Handlers
//
//
// ------------------------------------------------------------------

// BasicHttp is an http handler which serves a custom message
// with a status of 200 OK.
//
//	mux.Handle("/", BasicHttp("hi"))
//
// .
func BasicHttp(msg string) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, msg, http.StatusOK)
	}
	return http.HandlerFunc(fn)
}

// BasicJson is an http handler which serves a custom json message
// with a status of 200 OK.
//
//	mux.Handle("/", BasicJson("hi"))
//
// .
func BasicJson(msg string) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		writeJson(w, msg, http.StatusOK)
	}
	return http.HandlerFunc(fn)
}

// ------------------------------------------------------------------
//
//
// FileServer Handlers
//
//
// ------------------------------------------------------------------

// FileServer is an http handler which serves files from the given file system.
//
//	mux.Handle("/static/", FileServer(fsys))
//
// .
func FileServer(fsys fs.FS) http.Handler {
	return http.FileServerFS(fsys)
}

// FileServerDir is an http handler which serves files from the file system.
// Root is the filesystem root.
// Prefix is the prefix of the request path to strip off before searching the
// filesystem for the given file.
//
//	mux.Handle("/static/", FileServerDir("./staticfiles", "/static/"))
//
// .
func FileServerDir(root, prefix string) http.Handler {
	return http.StripPrefix(prefix, http.FileServer(http.Dir(root)))
}
