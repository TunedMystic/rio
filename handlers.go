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

// FileServer is an http handler which serves files from the given file system.
//
//	r.Handle("/static/", FileServer(fsys))
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
//	r.Handle("/static/", FileServerDir("./staticfiles", "/static/"))
//
// .
func FileServerDir(root, prefix string) http.Handler {
	return http.StripPrefix(prefix, http.FileServer(http.Dir(root)))
}
