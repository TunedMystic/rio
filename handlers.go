package rio

import (
	"io/fs"
	"net/http"
)

// ------------------------------------------------------------------
//
//
// FileServer Handlers
//
//
// ------------------------------------------------------------------

// FileServerFS is an http handler which serves files from the given file system.
//
//	r.Handle("/static/", FileServerFS(fsys))
//
// .
func FileServerFS(fsys fs.FS) http.Handler {
	return http.FileServerFS(fsys)
}

// FileServer is an http handler which serves files from the file system.
// Root is the filesystem root.
// Prefix is the prefix of the request path to strip off before searching the
// filesystem for the given file.
//
//	r.Handle("/static/", FileServer("./staticfiles", "/static/"))
//
// .
func FileServer(root, prefix string) http.Handler {
	return http.StripPrefix(prefix, http.FileServer(http.Dir(root)))
}
