package render

import (
	"bytes"
	"html/template"
	"io/fs"
	"net/http"

	"github.com/tunedmystic/rio/utils"
)

type View struct {
	templates *template.Template
}

func New(files fs.FS, pattern string) *View {
	funcs := template.FuncMap{
		"safe": func(content string) template.HTML {
			return template.HTML(content)
		},
	}

	tmpl := template.New("").Funcs(funcs)
	tmpl = template.Must(tmpl.ParseFS(files, pattern))

	return &View{
		templates: tmpl,
	}

	// view := NewView(fs, "*.html")
	// view.Render(w, "some-template", ["something", "here"])
	// view.Render404("")
	// log.Fatal()
	// logg.Info()
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
