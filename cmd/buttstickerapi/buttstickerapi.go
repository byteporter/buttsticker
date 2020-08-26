package main

import (
    "internal/pkg/handler"

    "log"
    "net/http"
    "path/filepath"
    "strings"
    "github.com/gorilla/mux"
)


func main() {
    bh := handler.ButtstickerHandler{TickerFilePath: filepath.Join("/usr/share/tickerdata", "tickers.json")}

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
