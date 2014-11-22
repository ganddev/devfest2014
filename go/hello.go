package main

import (
    "fmt"
    "net/http"
    "appengine"
    "appengine/urlfetch"
    "github.com/alexjlockwood/gcm"
)

func init() {
    http.HandleFunc("/send", handler)
}


func handler(w http.ResponseWriter, r *http.Request) {
    // Create the message to be sent.
    data := map[string]interface{}{"score": "5x1", "time": "15:10"}
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