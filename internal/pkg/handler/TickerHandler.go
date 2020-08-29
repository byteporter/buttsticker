package handler

import (
    "fmt"
    "strconv"
    "math/rand"
    "net/http"
    "strings"
    "log"
    "github.com/gorilla/mux"
    "os"
    "io/ioutil"
    "encoding/json"
)

type TickerHandler struct {
    TickerFilePath  string
}

func (th TickerHandler) GetRandomTicker(w http.ResponseWriter, req *http.Request) {
    tickers := th.readFile()
    fmt.Fprintf(w, tickers[rand.Intn(len(tickers))])
}

func (th TickerHandler) GetTickers(w http.ResponseWriter, req *http.Request) {
    tickers := th.readFile()
    fmt.Fprintf(w, strings.Join(tickers, "\n"))
}

// This should write to file instead of mutating the underlying struct
func (th *TickerHandler) PostTickers(w http.ResponseWriter, req *http.Request) {
    tickers := th.readFile()
    log.Println("*** POST Entered")

    if err := req.ParseForm(); err != nil {
        log.Println("*** POST Error: submission was all kinds of stupid")
        return
    }
    
    tickerString := req.FormValue("ticker")
    
    log.Printf("*** POST content: %s", tickerString)
    
    tickers = append(tickers, tickerString)
    fmt.Fprintf(w, strings.Join(tickers, "\n"))
}

func (th TickerHandler) GetTicker(w http.ResponseWriter, req *http.Request) {
    tickers := th.readFile()
    vars := mux.Vars(req)
    id, _ := strconv.Atoi(vars["id"])

    if id >= len(tickers) {
        w.WriteHeader(http.StatusNotFound)
        fmt.Fprintf(w, "Requested id '%d' is out of bounds", id)
    } else {
        fmt.Fprintf(w, tickers[id]) 
    }
}

func (th TickerHandler) readFile() ([]string) {
    var tickers []string

    tickerJson, err := os.Open(th.TickerFilePath)
    if err != nil {
        log.Println(err.Error())
    }
    byteVal, _ := ioutil.ReadAll(tickerJson)
    json.Unmarshal(byteVal, &tickers)
    defer tickerJson.Close()

    return tickers
}
