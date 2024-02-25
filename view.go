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

func Render(w http.ResponseWriter, page string, status int, data any) {
	defaultView.Render(w, page, status, data)
}

// ------------------------------------------------------------------
//
//
// View Functional Options
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

	// Set default functions for the func map.
	v.funcMap["safe"] = SafeHTML
	v.funcMap["time"] = TimeDisplay
	v.funcMap["date"] = DateDisplay
	v.funcMap["datetime"] = DateTimeDisplay

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
func (v *View) Render(w http.ResponseWriter, page string, status int, data any) {
	buf := new(bytes.Buffer)

	// Write the template to the buffer first.
	// If error, then respond with a server error and return.
	if err := v.templates.ExecuteTemplate(buf, page, data); err != nil {
		Http500(w, err)
		return
	}

	w.WriteHeader(status)

	// Write the contents of the buffer to the http.ResponseWriter.
	buf.WriteTo(w)
}

func (v *View) Render404(w http.ResponseWriter, data any) {
	v.Render(w, "404", http.StatusNotFound, data)
}

func (v *View) Render500(w http.ResponseWriter, data any) {
	v.Render(w, "500", http.StatusInternalServerError, data)
}
