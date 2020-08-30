package handler

import (
    "fmt"
    "strconv"
    "math/rand"
    "net/http"
    "strings"
    "log"
    "github.com/gorilla/mux"
    "html/template"

    _ "github.com/mattn/go-sqlite3"
    "database/sql"
)

type TickerHandler struct {
    DBPath string
    AchievementsTemplatePath string
}

func (th TickerHandler) GetRandomTicker(w http.ResponseWriter, req *http.Request) {
    tickers := th.readTickers()
    fmt.Fprintf(w, tickers[rand.Intn(len(tickers))])
}

func (th TickerHandler) GetTickers(w http.ResponseWriter, req *http.Request) {
    tickers := th.readFile()
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

func (th TickerHandler) GetAchievements(w http.ResponseWriter, req *http.Request) {
    tickers := th.readTickers()

    log.Println("*** INFO: In GetAchievements")

    t, err := template.ParseFiles(th.AchievementsTemplatePath)
    if err != nil {
        log.Println(err)
        w.WriteHeader(http.StatusInternalServerError)
        return
    }

    tData := struct {
        Tickers []string
    }{
        tickers,
    }

    t.Execute(w, tData)
}

// This should write to file instead of mutating the underlying struct
func (th TickerHandler) PostTickers(w http.ResponseWriter, req *http.Request) {
    log.Println("*** POST Entered")

    if err := req.ParseForm(); err != nil {
        log.Println("*** POST Error: submission was all kinds of stupid")
        return
    }
    
    tickerString := req.FormValue("ticker")
    
    log.Printf("*** POST content: %s", tickerString)
    
    db, err := sql.Open("sqlite3", th.DBPath)
    if err != nil {
        log.Fatal(err)
    }
    defer db.Close()

    _, err = db.Exec("INSERT INTO ticker(content) VALUES( $1 )", tickerString)
    if err != nil {
        log.Fatal(err)
    }

    th.GetTickers(w, req)
}

func (th TickerHandler) readTickers() ([]string) {
    var tickers []string;

    db, err := sql.Open("sqlite3", th.DBPath)
    if err != nil {
        log.Fatal(err)
    }
    defer db.Close()
    
    rows, err := db.Query("SELECT content FROM ticker;")
    if err != nil {
        log.Fatal(err)
    }
    defer rows.Close()
    for rows.Next() {
        var content string
        err = rows.Scan(&content)
        if err != nil {
            log.Fatal(err)
        }
        tickers = append(tickers, content)
    }
    err = rows.Err()
    if err != nil {
        log.Fatal(err)
    }

    return tickers
}

func (th TickerHandler) readFile() ([]string) {
    return th.readTickers()
}
