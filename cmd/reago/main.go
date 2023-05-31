package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"path"
	"path/filepath"

	"github.com/dennwc/reago"
)

var (
	fComp = flag.String("comp", "./components", "components directory")
	fHTML = flag.String("html", "./example", "html page directory")
	fHost = flag.String("host", ":8080", "host to listen on")
)

func main() {
	flag.Parse()
	if err := run(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

// ExampleRecord is an example DB record. It should be defined by the user.
type ExampleRecord struct {
	ID   int
	Name string
}

// ExampleDB is an example database implementation for exposing data to the templates. It should be defined by the user.
type ExampleDB struct {
}

// Table is an example method for listing database records for templates to access.
func (d *ExampleDB) Table() []ExampleRecord {
	return []ExampleRecord{
		{ID: 1, Name: "Foo"},
		{ID: 2, Name: "Bar"},
	}
}

func run() error {
	e, err := reago.NewEngine(*fComp, ExampleDB{})
	if err != nil {
		return err
	}
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		name := path.Clean(r.URL.Path)
		ext := path.Ext(name)
		if ext != ".html" && ext != "" {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		if name == "/" {
			name = "index.html"
		} else if ext == "" {
			name += ".html"
		}
		fname := filepath.Join(*fHTML, name)
		log.Println(r.Method, r.URL.Path, "->", fname)
		err := e.RenderPage(w, fname)
		if err != nil {
			log.Println(err)
		}
	})
	log.Println("server started on http://localhost" + *fHost)
	return http.ListenAndServe(*fHost, nil)
}
