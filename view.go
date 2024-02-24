package rio

import (
	"bytes"
	"fmt"
	"html/template"
	"io/fs"
	"net/http"
	"strings"

	"github.com/tunedmystic/rio/utils"
)

// ------------------------------------------------------------------
//
//
// Type: View
//
//
// ------------------------------------------------------------------

type View struct {
	templates *template.Template
}

func NewView(filesFS fs.FS) *View {
	funcs := template.FuncMap{
		"safe": func(content string) template.HTML {
			return template.HTML(content)
		},
	}

	tmpl := template.New("")

	// Walk the filesystem.
	err := fs.WalkDir(filesFS, ".", func(path string, d fs.DirEntry, err error) error {
		fmt.Println(path)
		if err != nil {
			return err
		}

		// Process all Html files, recursively.
		if !d.IsDir() && strings.HasSuffix(path, ".html") {
			// Read the file.
			fileBytes, err := fs.ReadFile(filesFS, path)
			if err != nil {
				return err
			}

			// Create new template.
			t := tmpl.New(path).Funcs(funcs)

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

	return &View{
		templates: tmpl,
	}
}

// Render writes a template to the http.ResponseWriter.
// .
func (v *View) Render(w http.ResponseWriter, status int, page string, data any) {
	buf := new(bytes.Buffer)

	// Write the template to the buffer first.
	// If error, then respond with a server error and return.
	if err := v.templates.ExecuteTemplate(buf, page, data); err != nil {
		utils.Http500(w, err)
		return
	}

	w.WriteHeader(status)

	// Write the contents of the buffer to the http.ResponseWriter.
	buf.WriteTo(w)
}

func (v *View) Render404(w http.ResponseWriter, data any) {
	v.Render(w, http.StatusNotFound, "404", data)
}

func (v *View) Render500(w http.ResponseWriter, data any) {
	v.Render(w, http.StatusInternalServerError, "500", data)
}
