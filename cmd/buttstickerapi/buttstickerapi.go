package main

import (
    "fmt"
    //    "html/template"
    "log"
    "net/http"
    "path/filepath"
    "strings"
    "strconv"
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

func (bh buttstickerHandler) GetTickers(w http.ResponseWriter, req *http.Request) {
    fmt.Fprintf(w, strings.Join(bh.tickers, "\n"))
}

// This should write to file instead of mutating the underlying struct
func (bh *buttstickerHandler) PostTickers(w http.ResponseWriter, req *http.Request) {
    log.Println("*** POST Entered")

    if err := req.ParseForm(); err != nil {
        log.Println("*** POST Error: submission was all kinds of stupid")
        return
    }
    
    tickerString := req.FormValue("ticker")
    
    log.Printf("*** POST content: %s", tickerString)
    
    bh.tickers = append(bh.tickers, tickerString)
    fmt.Fprintf(w, strings.Join(bh.tickers, "\n"))
}

func (bh buttstickerHandler) GetTicker(w http.ResponseWriter, req *http.Request) {
    vars := mux.Vars(req)
    id, _ := strconv.Atoi(vars["id"])

    if id >= len(bh.tickers) {
        w.WriteHeader(http.StatusNotFound)
        fmt.Fprintf(w, "Requested id '%d' is out of bounds", id)
    } else {
        fmt.Fprintf(w, bh.tickers[id]) 
    }
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
    apiRouter.HandleFunc("/tickers", bh.GetTickers).Methods("GET")
    apiRouter.HandleFunc("/tickers/{id:[0-9]+}", bh.GetTicker).Methods("GET")
    apiRouter.HandleFunc("/tickers", bh.PostTickers).Methods("POST")

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
