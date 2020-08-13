package httpapi

import (
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"

	"github.com/yeqown/log"
)

// here should be all unexported functions and type

var (
	tpl404       *template.Template
	tplReport    *template.Template
	tplHome      *template.Template
	tplHighscore *template.Template
	tplAbout     *template.Template
)

func init() {
	loadTpls()
}

func loadTpls() {
	tpl404 = template.Must(
		template.New("404.html").Delims("[[", "]]").
			ParseFiles("tpl/404.html", "tpl/footer.html", "tpl/header.html"))

	tplReport = template.Must(
		template.New("report.html").Delims("[[", "]]").
			ParseFiles("tpl/report.html", "tpl/footer.html", "tpl/header.html"))

	tplHome = template.Must(
		template.New("home.html").Delims("[[", "]]").
			ParseFiles("tpl/home.html", "tpl/footer.html", "tpl/header.html"))

	fns := template.FuncMap{"add": add, "formatScore": formatScore}
	tplHighscore = template.Must(
		template.New("high_scores.html").Delims("[[", "]]").Funcs(fns).
			ParseFiles("tpl/high_scores.html", "tpl/footer.html", "tpl/header.html"))

	tplAbout = template.Must(
		template.New("about.html").Delims("[[", "]]").
			ParseFiles("tpl/about.html", "tpl/footer.html", "tpl/header.html"))
}

func Error(w http.ResponseWriter, statusCode int, err error) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.WriteHeader(statusCode)
	_, _ = fmt.Fprintln(w, err.Error())
}

// JSON write json format message to client
func JSON(w http.ResponseWriter, statusCode int, v interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	d, err := json.Marshal(v)
	if err != nil {
		log.Errorf("JSON failed to marshal, err=%v", err)
		return
	}

	// log.Debugf("JSON with s=%s", string(d))
	_, _ = fmt.Fprintln(w, string(d))
}

// renderHTML render html file to client
func renderHTML(
	w http.ResponseWriter,
	statusCode int,
	tpl *template.Template,
	data interface{},
) {
	// log.Debugf("renderHTML render data=%+v", data)
	w.WriteHeader(statusCode)
	if err := tpl.Execute(w, data); err != nil {
		log.Errorf("renderHTML failed to execute, err=%v", err)
	}
}

func add(x, y int) int {
	return x + y
}

func formatScore(x float64) string {
	return fmt.Sprintf("%.2f", x)
}
