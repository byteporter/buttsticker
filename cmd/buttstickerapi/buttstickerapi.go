package main

import (
    "internal/pkg/handler"

    "log"
    "net/http"
    "path/filepath"
    "strings"
    "github.com/gorilla/mux"
    "github.com/gorilla/handlers"

    // "flag"
)

func main() {
    th := handler.TickerHandler{
        filepath.Join("/usr/share/tickerdata", "ticker.db"),
        filepath.Join("/usr/share/web", "achievements.html"),
        "dd6f6992-e6ee-435a-bf63-2d3b90ffd107",
    }

    apiPrefix := "/api/v1"
    router := mux.NewRouter()
    apiRouter := router.PathPrefix(apiPrefix).Subrouter()
    passwordRouter := apiRouter.PathPrefix(th.Password)

    apiRouter.HandleFunc("/tickers/rand", th.GetRandomTicker).Methods("GET")
    apiRouter.HandleFunc("/tickers", th.GetTickers).Methods("GET")
    apiRouter.HandleFunc("/tickers/{id:[0-9]+}", th.GetTicker).Methods("GET")
    passwordRouter.HandleFunc("/tickers", th.PostTickers).Methods("POST")
    apiRouter.HandleFunc("/achievements", th.GetAchievements).Methods("GET")
    passwordRouter.HandleFunc("/achievements/addform", th.GetAchievementsAddform).Methods("GET")

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
    log.Fatal(http.ListenAndServe(strings.Join([]string{":", "8080"}, ""), router));
}
