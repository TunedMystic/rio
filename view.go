package rio

import (
	"bytes"
	"fmt"
	"html/template"
	"io/fs"
	"net/http"
	"strings"
)

// ------------------------------------------------------------------
//
//
// Default View
//
//
// ------------------------------------------------------------------

var defaultView = &View{templates: template.New("")}

func Templates(templatesFS fs.FS, opts ...ViewOpt) {
	defaultView = NewView(templatesFS, opts...)
}

func Render(w http.ResponseWriter, page string, status int, data any) error {
	return defaultView.Render(w, page, status, data)
}

// ------------------------------------------------------------------
//
//
// Functional Options for View
//
//
// ------------------------------------------------------------------

type ViewOpt func(*View)

func WithFuncMap(funcMap template.FuncMap) ViewOpt {
	return func(v *View) {
		for key := range funcMap {
			v.funcMap[key] = funcMap[key]
		}
	}
}

// ------------------------------------------------------------------
//
//
// Type: View
//
//
// ------------------------------------------------------------------

type View struct {
	templates *template.Template
	funcMap   template.FuncMap
}

func NewView(templatesFS fs.FS, opts ...ViewOpt) *View {
	v := &View{
		templates: template.New(""),
		funcMap:   template.FuncMap{},
	}

	// Set default functions for the func map (from utils.go).
	v.funcMap["safe"] = DisplaySafeHTML
	v.funcMap["time"] = DisplayTime
	v.funcMap["date"] = DisplayDate
	v.funcMap["datetime"] = DisplayDateTime

	// Configure with ViewOpt funcs, if any.
	for i := range opts {
		opts[i](v)
	}

	// Walk the filesystem.
	err := fs.WalkDir(templatesFS, ".", func(path string, d fs.DirEntry, err error) error {
		fmt.Println(path)
		if err != nil {
			return err
		}

		// Process all Html files, recursively.
		if !d.IsDir() && strings.HasSuffix(path, ".html") {
			// Read the file.
			fileBytes, err := fs.ReadFile(templatesFS, path)
			if err != nil {
				return err
			}

			// Create new template.
			t := v.templates.New(path).Funcs(v.funcMap)

			// Parse the template.
			if _, err := t.Parse(string(fileBytes)); err != nil {
				return err
			}
		}
		return nil
	})

	if err != nil {
		panic(err)
	}

	return v
}

// Render writes a template to the http.ResponseWriter.
// .
func (v *View) Render(w http.ResponseWriter, page string, status int, data any) error {
	buf := new(bytes.Buffer)

	// Write the template to the buffer first.
	// If error, then respond with a server error and return.
	if err := v.templates.ExecuteTemplate(buf, page, data); err != nil {
		return err
	}

	w.WriteHeader(status)

	// Write the contents of the buffer to the http.ResponseWriter.
	buf.WriteTo(w)

	return nil
}
