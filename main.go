package main

import (
	"flag"
	"fmt"
	"net/http"
	"net/rpc"
)

var (
	welcome    = flag.String("welcome", "welcome to goto", "欢迎语")
	listenAddr = flag.String("http", "8080", "http listen address")
	dataFile   = flag.String("file", "store.gob", "data store file name")
	hostname   = flag.String("host", "localhost", "http host name")
	rpcEnabled = flag.Bool("rpc", false, "enbale RPC server")
)

var store *URLStore

func main() {

	fmt.Println(*welcome)
	flag.Parse()
	store = NewURLStore(*dataFile)

	if *rpcEnabled {
		rpc.RegisterName("Store", store)
		rpc.HandleHTTP()
	}

	http.HandleFunc("/", Redirect)
	http.HandleFunc("/add", Add)
	http.ListenAndServe(":"+*listenAddr, nil)
}

func Redirect(w http.ResponseWriter, r *http.Request) {
	key := r.URL.Path[1:]
	url := store.Get(key)
	if url == "" {
		http.NotFound(w, r)
	}
	http.Redirect(w, r, url, http.StatusFound)
}

func Add(w http.ResponseWriter, r *http.Request) {
	url := r.FormValue("url")
	if url == "" {
		w.Header().Set("Content-Type", "text/html")
		fmt.Fprint(w, AddForm)
		return
	}
	var key string
	if err := store.Put(&url, &key); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	fmt.Fprintf(w, "http://%s/%s", *hostname+":"+*listenAddr, key)
}

const AddForm = `
<form method="POST" action="/add">
URL: <input type="text" name="url">
<input type="submit" value="Add">
</form>
`
