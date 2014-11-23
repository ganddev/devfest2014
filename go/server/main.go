package feedme

import (
        "net/http"
        "fmt"
        "appengine"
        "appengine/datastore"
        "appengine/urlfetch"
        "github.com/alexjlockwood/gcm"
        "encoding/json"
        "io/ioutil"
        "strings"
        "github.com/satori/go.uuid"
)

type Device struct {
    DeviceToken string `json: "#deviceToken"`
}

type PushedTopstories struct {
    TopStory string
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

type Notification struct {
    Id string
    Status string
}

type Topstories []int

func init() {
        http.HandleFunc("/", root)
        http.HandleFunc("/register", register)
        http.HandleFunc("/fetchTopstories", fetchTopstories)
        http.HandleFunc("/test", testMsg)
}

// DataStore crap

func deviceKey(c appengine.Context) *datastore.Key {
    return datastore.NewKey(c, "Device", "default_device", 0, nil)
}

func pushedTopstoriesKey(c appengine.Context) *datastore.Key {
    return datastore.NewKey(c, "PushedTopstories", "default_topstories", 0, nil)
}

func testMsg(w http.ResponseWriter, r *http.Request) {
    c := appengine.NewContext(r)
    sendPushNotifications(c, w, "test done", "http://devfest-berlin.de") 
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
        
                sendMessage(c, w, fu.DeviceToken, "Registration completed!", "http://devfest-berlin.de") 

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


func sendMessage(c appengine.Context, w http.ResponseWriter, deviceToken string, title string, url string) {
    data := map[string]interface{}{"title": title, "url" : url, "uuid" : uuid.NewV4().String()}
    regIDs := []string{deviceToken}
    msg := gcm.NewMessage(data, regIDs...)

    client := urlfetch.Client(c)
    sender := &gcm.Sender{ApiKey: "AIzaSyAmAUW2zbbqO16zIzv5IHQc9U9ZKPnl6jI", Http: client}
     
    response, err := sender.Send(msg, 2)
    if err != nil {
        fmt.Fprint(w,"Failed to send message:", err)
        return
    } else {
        datastore.Put(c, datastore.NewIncompleteKey(c, "message", nil), &data)
        fmt.Fprint(w, "Send message:" ,response)
    }
}

func fetchTopstories(w http.ResponseWriter, r *http.Request) {
    c := appengine.NewContext(r)
    client := urlfetch.Client(c)
    resp, err := client.Get("https://hacker-news.firebaseio.com/v0/topstories.json")
    
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

            sxs := fmt.Sprintf("%d", m[0])
            
            fetchItem(w, r, sxs)
        }
    }
}

func fetchItem(w http.ResponseWriter, r *http.Request, itemid string) {
    c := appengine.NewContext(r)
    client := urlfetch.Client(c)

    fu := PushedTopstories{
        TopStory: itemid,
    }
    //Check if pushed already
    q := datastore.NewQuery("PushedTopstories").Filter("TopStory =", fu.TopStory)
    qresult := q.Run(c)

    for {
            var pts PushedTopstories
            _, err := qresult.Next(&pts)

            if err == datastore.Done {
                fmt.Fprintf(w, "Not in list! ID: " + fu.TopStory)

                //Store device
                key := datastore.NewIncompleteKey(c, "PushedTopstories", pushedTopstoriesKey(c))
                _, err := datastore.Put(c, key, &fu)
        
                if err != nil {
                    http.Error(w, err.Error(), http.StatusInternalServerError)
                    return
                }
                break
            }
            
            //Item found: goodbye
            if err == nil {
                break
            }
        }

    siteURL := "https://hacker-news.firebaseio.com/v0/item/"+itemid+".json?print=pretty"
    fmt.Fprintf(w, "\nurl " + siteURL + "\n%i", itemid)

    resp, err := client.Get(siteURL)
    
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

            fmt.Fprintf(w, "Push Title: %s  URL: %s", m.Title, m.Url)
            sendPushNotifications(c, w, m.Title, m.Url)
        }
    }
}

func sendPushNotifications(c appengine.Context, w http.ResponseWriter, title string, url string) {
   q := datastore.NewQuery("Device")
    t := q.Run(c)

    for {
       var dx Device
       _, err := t.Next(&dx)

        if err != nil {
            break;
        }

        if err == nil {
            fmt.Fprintf(w, "\nSend notification to: " + dx.DeviceToken + "\n")
            sendMessage(c, w, dx.DeviceToken, title, url) 
       }
   }
}