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
    tickers     []string
}

func (bh buttstickerHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
    fmt.Fprintf(w, bh.tickers[rand.Intn(len(bh.tickers))])
}

func main() {
    var bh buttstickerHandler

    tickerJson, err := os.Open(filepath.Join("/usr/share/tickerdata", "tickers.json"))
    if err != nil {
        log.Println(err.Error())
    }
    byteVal, _ := ioutil.ReadAll(tickerJson)
    json.Unmarshal(byteVal, &bh.tickers)
    defer tickerJson.Close()

    apiPrefix := "/api/v1"
    router := mux.NewRouter()
    apiRouter := router.PathPrefix(apiPrefix).Subrouter()

    apiRouter.Handle("/tickers/rand", bh).Methods("GET")

    http.Handle("/", apiRouter)

    log.Println("*** INFO: Available Routes:")
    apiRouter.Walk(func(route *mux.Route, router *mux.Router, ancestors []*mux.Route) error {
        t, err := route.GetPathTemplate()
        if err != nil {
            return err
        }
        log.Println("*** INFO: " + t)
        return nil
    })
    log.Fatal(http.ListenAndServe(strings.Join([]string{":", "80"}, ""), router));
}
