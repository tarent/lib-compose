package main

import (
	gorilla "github.com/gorilla/handlers"
	"github.com/tarent/lib-compose/composition"
	"net/http"
	"os"
)

var host = "127.0.0.1:8080"

func main() {
	panic(http.ListenAndServe(host, handler()))
}

func handler() http.Handler {
	mux := http.NewServeMux()
	mux.Handle("/static/", staticHandler())
	mux.Handle("/", compositionHandler())
	return gorilla.LoggingHandler(os.Stdout, mux)
}

func compositionHandler() http.Handler {
	contentFetcherFactory := func(r *http.Request) composition.FetchResultSupplier {
		defaultMetaJSON := map[string]interface{}{
			"header_text": "Hello World!",
			"request":     composition.MetadataForRequest(r),
		}

		fetcher := composition.NewContentFetcher(defaultMetaJSON)
		baseUrl := "http://" + r.Host

		fetcher.AddFetchJob(composition.NewFetchDefinition(baseUrl + "/static/layout.html").WithName("layout"))
		fetcher.AddFetchJob(composition.NewFetchDefinition(baseUrl + "/static/lorem.html").WithName("content"))
		fetcher.AddFetchJob(composition.NewFetchDefinition(baseUrl + "/static/styles.html"))

		return fetcher
	}
	return composition.NewCompositionHandler(contentFetcherFactory)
}

func staticHandler() http.Handler {
	return http.FileServer(http.Dir("./"))
}
