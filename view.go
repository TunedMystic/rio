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

// Templates configures the default view with the given
// templateFS filesystem and the
// opts functional options.
func Templates(templatesFS fs.FS, opts ...ViewOpt) {
	defaultView = NewView(templatesFS, opts...)
}

// Render writes a template to the http.ResponseWriter.
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

// ViewOpt is a function to configure a View.
type ViewOpt func(*View)

// WithFuncMap adds the functions in the given funcMap to the View.
//
// If the View contains a funcMap function with the same name,
// then it will not be added.
func WithFuncMap(funcMap template.FuncMap) ViewOpt {
	return func(v *View) {
		for key := range funcMap {
			if _, ok := v.funcMap[key]; !ok {
				v.funcMap[key] = funcMap[key]
			}
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

// View is a collection of html templates for rendering.
type View struct {
	templates *template.Template
	funcMap   template.FuncMap
}

// NewView constructs and returns a new *View.
//
// The templateFS is a filesystem which contains html templates.
// The opts is a slice of ViewOpt funcs for optional configuration.
func NewView(templatesFS fs.FS, opts ...ViewOpt) *View {
	view, err := constructView(templatesFS, opts...)
	if err != nil {
		panic(fmt.Errorf("failed to construct view: %w", err))
	}
	return view
}

// Render writes a template to the http.ResponseWriter.
func (v *View) Render(w http.ResponseWriter, page string, status int, data any) error {
	buf := new(bytes.Buffer)

	// Write the template to the buffer first.
	if err := v.templates.ExecuteTemplate(buf, page, data); err != nil {
		return err
	}

	w.WriteHeader(status)

	// Write the contents of the buffer to the http.ResponseWriter.
	buf.WriteTo(w)

	return nil
}

// constructView constructs and returns a *View.
//
// All html templates within templatesFS are parsed and loaded.
//
// Default template functions are provided to the funcmap, but
// more can be added with the ViewOpt functions.
func constructView(templatesFS fs.FS, opts ...ViewOpt) (*View, error) {
	v := &View{
		templates: template.New(""),
		funcMap:   template.FuncMap{},
	}

	// Set the default template functions.
	v.funcMap["safe"] = DisplaySafeHTML
	v.funcMap["time"] = DisplayTime
	v.funcMap["date"] = DisplayDate
	v.funcMap["datetime"] = DisplayDateTime

	// Configure the View with with ViewOpt funcs, if any.
	for i := range opts {
		opts[i](v)
	}

	// Parse and load all templates from the given filesystem.
	//
	// Walk the templateFS filesystem, recursively.
	err := fs.WalkDir(templatesFS, ".", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		// Process all Html files.
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

	return v, err
}
