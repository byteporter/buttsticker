package main

import (
    "fmt"
//    "html/template"
    "log"
    "net/http"
    "path/filepath"
    "strings"
    "os"
    "encoding/json"
    "io/ioutil"
    "math/rand"

//    "github.com/gorilla/handlers"
    "github.com/gorilla/mux"
)

type buttstickerHandler struct {
    h           http.Handler
    contentDir  string
}

func (bh buttstickerHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
    tickerJson, err := os.Open(filepath.Join(bh.contentDir, "tickers.json"))
    if err != nil {
        fmt.Fprintf(w, err.Error())
    }

    byteVal, _ := ioutil.ReadAll(tickerJson)
    var tickers []string
    json.Unmarshal(byteVal, &tickers)

    fmt.Fprintf(w, tickers[rand.Intn(4)])

    defer tickerJson.Close()
}

func main() {
    router := mux.NewRouter()
    router.Handle("/", buttstickerHandler{nil, "/usr/share/tickerdata"}).Methods("GET")
    http.Handle("/", router)

    log.Fatal(http.ListenAndServe(strings.Join([]string{":", "80"}, ""), router));
}
