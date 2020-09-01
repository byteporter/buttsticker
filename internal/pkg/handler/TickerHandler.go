package handler

import (
    "fmt"
    "strconv"
    "math/rand"
    "net/http"
    "log"
    "github.com/gorilla/mux"
    "html/template"
    "encoding/json"

    _ "github.com/mattn/go-sqlite3"
    "database/sql"
)

type TickerHandler struct {
    DBPath string
    AchievementsTemplatePath string
    Password string
}

func (th TickerHandler) GetRandomTicker(w http.ResponseWriter, req *http.Request) {
    tickers := th.readTickers()
    fmt.Fprint(w, tickers[rand.Intn(len(tickers))])
}

func (th TickerHandler) GetTickers(w http.ResponseWriter, req *http.Request) {
    var tickers []string
    for _, row := range th.readTickers() {
        tickers = append(tickers, row.Content)
    }

    tickersJson, err := json.Marshal(tickers)
    if err != nil {
        log.Fatal(err)
    }

    w.Header().Set("Content-Type", "application/json")
    fmt.Fprint(w, string(tickersJson))
}

func (th TickerHandler) GetTicker(w http.ResponseWriter, req *http.Request) {
    tickers := th.readTickers()
    vars := mux.Vars(req)
    id, _ := strconv.Atoi(vars["id"])

    if id >= len(tickers) {
        w.WriteHeader(http.StatusNotFound)
        fmt.Fprint(w, "Requested id '%d' is out of bounds", id)
    } else {
        fmt.Fprint(w, tickers[id]) 
    }
}

func (th TickerHandler) GetAchievements(w http.ResponseWriter, req *http.Request) {
    var tickers []string
    for _, row := range th.readTickers() {
        tickers = append(tickers, row.Content)
    }

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

func (th TickerHandler) GetAchievementsAddform(w http.ResponseWriter, req *http.Request) {
    log.Println("*** INFO: In GetAchievementsAddform")

    t, err := template.ParseFiles("/usr/share/web/add-achievements.html")
    if err != nil {
        log.Println(err)
        w.WriteHeader(http.StatusInternalServerError)
        return
    }

    tData := struct {
        FormAction string
    }{
        fmt.Sprintf("/api/v1/%s/tickers", th.Password),
    }

    t.Execute(w, tData)
}

func (th TickerHandler) PostTickers(w http.ResponseWriter, req *http.Request) {
    log.Println("*** POST Entered")

    if err := req.ParseForm(); err != nil {
        log.Println("*** ERROR:PostTickers malformed request")
        w.WriteHeader(http.StatusInternalServerError)
        return
    }
    
    tickerString := req.FormValue("ticker")
    
    log.Printf("*** POST content: %s", tickerString)
    
    db, err := sql.Open("sqlite3", th.DBPath)
    if err != nil {
        log.Fatal(err)
    }
    defer db.Close()

    _, err = db.Exec("INSERT INTO ticker(content) VALUES( $1 );", tickerString)
    if err != nil {
        log.Fatal(err)
    }

    th.GetTickers(w, req)
}

func (th TickerHandler) DeleteTicker(w http.ResponseWriter, req *http.Request) {
    log.Println("*** INFO: DeleteTicker Entered")

    if err := req.ParseForm(); err != nil {
        log.Println("*** ERROR: DeleteTicker malformed request")
        w.WriteHeader(http.StatusInternalServerError)
        return
    }

    tickerId, err := strconv.Atoi(req.FormValue("id"))
    if err != nil {
        log.Fatal(err)
    }

    log.Printf("*** INFO: Delete requested for (%T)id %d", tickerId, tickerId)

    db, err := sql.Open("sqlite3", th.DBPath)
    if err != nil {
        log.Fatal(err)
    }
    defer db.Close()

    _, err = db.Exec("UPDATE ticker SET deletetime = current_timestamp WHERE rowid = $1;", tickerId)
    if err != nil {
        log.Fatal(err)
    }

    th.GetTickers(w, req)
}

type TickerRow struct {
    RowID       int
    Content     string
    Createdate  string
}

func (th TickerHandler) readTickers() ([]TickerRow) {
    var tickers []TickerRow;

    db, err := sql.Open("sqlite3", th.DBPath)
    if err != nil {
        log.Fatal(err)
    }
    defer db.Close()
    
    rows, err := db.Query("SELECT rowid, content, createtime FROM ticker WHERE deletetime IS NULL;")
    if err != nil {
        log.Fatal(err)
    }
    defer rows.Close()
    for rows.Next() {
        var row TickerRow 

        err = rows.Scan(&row.RowID, &row.Content, &row.Createdate)
        if err != nil {
            log.Fatal(err)
        }
        tickers = append(tickers, row)
    }
    err = rows.Err()
    if err != nil {
        log.Fatal(err)
    }

    return tickers
}
