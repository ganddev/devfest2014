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
    regIDs := []string{"",""}
    msg := gcm.NewMessage(data, regIDs...)

    c := appengine.NewContext(r)
    client := urlfetch.Client(c)
    sender := &gcm.Sender{ApiKey: "XXXXXXX", Http: client}
    // Send the message and receive the response after at most two retries.
    response, err := sender.Send(msg, 2)
    if err != nil {
        fmt.Fprint(w,"Failed to send message:", err)
        return
    } else {
        fmt.Fprint(w, "Send message:" ,response)
    }
}