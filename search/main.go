package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/blevesearch/bleve"
)

type Page struct {
	Title   string
	Content string
	Path    string
}

func main() {
  // Open index
  log.Printf("Creating index\n")
	mapping := bleve.NewIndexMapping()
	index, err := bleve.New("nanoc.bleve", mapping)
	if err != nil {
		fmt.Println(err)
		return
	}

  // Fill index
	reader, err := os.Open("../output/search-index.json")
	if err != nil {
		log.Fatal(err)
	}
	dec := json.NewDecoder(reader)
	var pages []Page
	if err := dec.Decode(&pages); err != nil && err != io.EOF {
		log.Fatal(err)
	}
	for _, page := range pages {
		log.Printf("Indexing %s\n", page.Path)
		index.Index(page.Path, page)
	}

  // Create server
	http.HandleFunc("/search", func(w http.ResponseWriter, r *http.Request) {
    log.Printf("Got search request for %s\n", r.URL.Query().Get("q"))

		query := bleve.NewQueryStringQuery(r.URL.Query().Get("q"))
		search := bleve.NewSearchRequest(query)
    search.Highlight = bleve.NewHighlight()
		searchResults, err := index.Search(search)
		if err != nil {
			fmt.Println(err)
			return
		}
		json.NewEncoder(w).Encode(searchResults.Hits)
	})

  // Serve
  log.Printf("Waiting for requests\n")
	if err = http.ListenAndServe(":8094", nil); err != nil {
		log.Fatal(err)
	}
}
