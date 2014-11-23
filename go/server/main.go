package feedme

import (
        "net/http"
        "time"
        "fmt"
        "appengine"
        "appengine/datastore"
        "appengine/urlfetch"
        "github.com/alexjlockwood/gcm"
        "encoding/json"
        "io/ioutil"
        "strings"
)

// [START greeting_struct]
type Greeting struct {
        Author  string
        Content string
        Date    time.Time
}

type Device struct {
    DeviceToken string `json: "#deviceToken"`
    Feeds []Feed `json: "#feeds"`
}

type Feed struct {
    url string
    lastEntry string
}

type Message struct {
    By string `json: "#by"`
    Id int  `json: "#id"`
    Kids []int `json: "#kids"`
    Score int `json: "#score"`
    Time int64  `json: "#time"`
    Title string `json: "#title"`
    Type string `json: "#type"`
    Url string `json: "#url"`
  }

// [END greeting_struct]

func init() {
        http.HandleFunc("/", root)
        http.HandleFunc("/register", register)
        
        http.HandleFunc("/send", sendNotification)
        http.HandleFunc("/fetch", fetchEntry)
}

func deviceKey(c appengine.Context) *datastore.Key {
    return datastore.NewKey(c, "Device", "default_device", 0, nil)
}

func root(w http.ResponseWriter, r *http.Request) {

}



func register(w http.ResponseWriter, r *http.Request) {
    if r.Method == "POST" {
        c := appengine.NewContext(r)

        fu := Device{
            DeviceToken: r.Header.Get("deviceToken"),
            // Feeds: r.Header.
        }

        q := datastore.NewQuery("Device").Filter("DeviceToken =", fu.DeviceToken)
        qresult := q.Run(c)

        for {
            var dev Device
            _, err :=  qresult.Next(&dev)

            if err == datastore.Done {
                fmt.Fprintf(w, "No entries")

                //Store device
                key := datastore.NewIncompleteKey(c, "Device", deviceKey(c))
                _, err := datastore.Put(c, key, &fu)
        
                if err != nil {
                    http.Error(w, err.Error(), http.StatusInternalServerError)
                    return
                }
                break
            }

            if err != nil {
                fmt.Fprintf(w, "Geht nicht!")
            }

            //TODO update feeds for device

            break
        }
    } else if r.Method == "GET" {
        fmt.Fprint(w, "GETTING SHIT HERE")
    }
}




func sendNotification(w http.ResponseWriter, r *http.Request) {
    // Create the message to be sent.
    data := map[string]interface{}{"score": "MÜLL!!!", "time": "15:10"}
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

func fetchEntry(w http.ResponseWriter, r *http.Request) {
    c := appengine.NewContext(r)
    client := urlfetch.Client(c)
    resp, err := client.Get("https://hacker-news.firebaseio.com/v0/item/8863.json?print=pretty")
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    defer resp.Body.Close()
    if body, err := ioutil.ReadAll(resp.Body); err != nil {
        fmt.Fprintf(w, "Couldn't read request body: %s", err)
    } else {
        dec := json.NewDecoder(strings.NewReader(string(body)))
        var m Message
        if err := dec.Decode(&m); err != nil {
            fmt.Fprintf(w, "Couldn't decode JSON: %s", err)
        } else {
            fmt.Fprintf(w, "Value of Param1 is: %s", m.By)
        }
    }
}