package main

import (
    "fmt"
    "net/http"
    "appengine"
    "appengine/urlfetch"
    "github.com/alexjlockwood/gcm"
    "encoding/json"
    "io/ioutil"
    "strings"
)

type Topstories struct {
    Stories []int
}

type Message struct {
    By string `json: "#by"`
    Id int  `json: "#id"`
    Kids []int `json: "#kids"`
    Score int `json: "#score"`
    Time int64  `json: "#time"`
    Title string `json: "#title"`
    Popel string `json: "#type"`
    Url string `json: "#url"`
  }

func init() {
    http.HandleFunc("/send", handler)
    http.HandleFunc("/fetch", handler2)
}


func handler(w http.ResponseWriter, r *http.Request) {
    // Create the message to be sent.
    data := map[string]interface{}{"score": "7x1", "time": "15:10"}
    regIDs := []string{"APA91bHspbmRxp8nw846IuKSvRUg27ElSe2W1BLyPChUfyvhkyz7aun6YaV8z-reviUY4oqMQ1PNMP4_jt-n7QdSYt9e3utfHm-tByW3ticXEbHYZhHqREE4daxKuxJz2_XZgY3XwUvJCpU4g1uUqckko53-jZQfTCSAIyYWSRwFgwzU_wH0tMA"}
    msg := gcm.NewMessage(data, regIDs...)

    c := appengine.NewContext(r)
    client := urlfetch.Client(c)
    sender := &gcm.Sender{ApiKey: "AIzaSyAmAUW2zbbqO16zIzv5IHQc9U9ZKPnl6jI", Http: client}
    // Send the message and receive the response after at most two retries.
    response, err := sender.Send(msg, 2)
    if err != nil {
        fmt.Fprint(w,"Failed to send message:", err)
        return
    } else {
        fmt.Fprint(w, "Send message:" ,response)
    }
}

func handler2(w http.ResponseWriter, r *http.Request) {
    c := appengine.NewContext(r)
    client := urlfetch.Client(c)
    resp, err := client.Get("https://hacker-news.firebaseio.com/v0/topstories")
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    defer resp.Body.Close()
    if body, err := ioutil.ReadAll(resp.Body); err != nil {
        fmt.Fprintf(w, "Couldn't read request body: %s", err)
    } else {
        dec := json.NewDecoder(strings.NewReader(string(body)))
        var m Topstories
        if err := dec.Decode(&m); err != nil {
            fmt.Fprintf(w, "Couldn't decode JSON: %s", err)
        } else {
            fmt.Fprintf(w, "Value of Param1 is: %s", m.Stories)
        }
    }
}
