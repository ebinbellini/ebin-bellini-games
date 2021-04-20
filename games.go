package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	http.HandleFunc("/", serve)

	fmt.Println("Listening on :2987...")
	err := http.ListenAndServe(":2987", nil)
	if err != nil {
		log.Fatal(err)
	}
}

func serve(w http.ResponseWriter, r *http.Request) {
	url := r.URL.Path
	lp := filepath.Join("templates", "layout.html")
	fp := filepath.Join("static", filepath.Clean(url))

	// First try to serve from static folder
	info, err := os.Stat(fp)
	if err == nil {
		// Serve static file
		if !info.IsDir() {
			http.ServeFile(w, r, fp)
		} else {
			// If static file does not exist try templates folder
			tp := filepath.Join("templates", filepath.Clean(url), "index.html")
			_, err := os.Stat(tp)
			if err == nil {
				// Add a / to the end of the URL if there isn't on already
				if !strings.HasSuffix(url, "/") {
					http.Redirect(w, r, url+"/", http.StatusMovedPermanently)
					return
				}

				tmpl, err := template.ParseFiles(lp, tp)
				if err != nil {
					serveNotFound(w, r)
				} else {
					tmpl.ExecuteTemplate(w, "layout", nil)
				}
			} else {
				if os.IsNotExist(err) {
					// Try to serve directory contents
					http.ServeFile(w, r, fp)
				}
			}
		}
	} else {
		if os.IsNotExist(err) {
			print("Couldn't find " + fp + "\n")
			serveNotFound(w, r)
		} else {
			serveInternalError(w, r)
		}
	}
}

func serveNotFound(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotFound)
	http.ServeFile(w, r, filepath.Join("templates", "404.html"))
}

func serveInternalError(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusInternalServerError)
	http.ServeFile(w, r, filepath.Join("templates", "error.html"))
}

/*func (w http.ResponseWriter, r *http.Request) {
	url := filepath.Join("./", filepath.Clean(r.URL.Path))
	fmt.Println(url)
	if len(url) == 1 {
		url = "index.html"
	}

	http.ServeFile(w, r, url)
}*/
