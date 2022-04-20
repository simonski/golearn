package http

import (
	"embed"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"sync"
	"text/template"

	"github.com/simonski/goutils"
)

var (
	//go:embed resources
	res embed.FS

	pages = map[string]string{
		"/":           "resources/index.html",
		"/index.html": "resources/index.html",
		"/code.js":    "resources/code.js",
		"/style.css":  "resources/style.css",
	}
)

func WebServer(wg *sync.WaitGroup) {
	defer wg.Done()
	cli := goutils.NewCLI(os.Args)
	port := cli.GetIntOrDefault("-p", 8000)
	sport := fmt.Sprintf(":%v", port)

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		page, ok := pages[r.URL.Path]

		if !ok {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		tpl, err := template.ParseFS(res, page) // page)
		if err != nil {
			log.Printf("Page %s not found in pages cache...", r.RequestURI)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		if strings.Index(page, ".html") > -1 {
			w.Header().Set("Content-Type", "text/html")
		} else if strings.Index(r.URL.Path, ".css") > -1 {
			w.Header().Set("Content-Type", "text/css")
		} else if strings.Index(r.URL.Path, ".js") > -1 {
			w.Header().Set("Content-Type", "text/javascript")
		}
		w.WriteHeader(http.StatusOK)
		data := map[string]interface{}{
			"userAgent": r.UserAgent(),
		}
		if err := tpl.Execute(w, data); err != nil {
			return
		}
	})
	http.FileServer(http.FS(res))
	fmt.Printf("Documentation started on %v\n", sport)
	err := http.ListenAndServe(sport, nil)
	if err != nil {
		panic(err)
	}
}
