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
        // "strconv"
)

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

type Topstories []int

func init() {
        http.HandleFunc("/", root)
        http.HandleFunc("/register", register)
        
        // http.HandleFunc("/send", sendNotification)
        http.HandleFunc("/fetch", fetchEntry)
        http.HandleFunc("/fetchTopstories", fetchTopstories)
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


func sendMessage(c appengine.Context, w http.ResponseWriter, deviceToken string, title string, url string) {
    data := map[string]interface{}{"title": title, "url" : url}
    regIDs := []string{deviceToken}
    msg := gcm.NewMessage(data, regIDs...)

    client := urlfetch.Client(c)
    sender := &gcm.Sender{ApiKey: "AIzaSyAmAUW2zbbqO16zIzv5IHQc9U9ZKPnl6jI", Http: client}
     
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
    sendMessage(c, w, "APA91bHUV33aJbJlAZeeRALozg5vegfgSdUMpE1gyOgDeYTdfoHytSAM_4dLpDkKLbcmUfgh6XysK3cP4MpcH5jy5AbPeyNY32T1uOzmlIVtO2E0W-a3KJGgZXMNUCgH2U1lody9JQjb5usbPtfRrj1VI1Us_YkHICCRfbfeI0y7PF0r40SmFD0", "Geilo", "http://istdasgeil.de")

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
            _, err := datastore.Put(c, datastore.NewIncompleteKey(c, "item", nil), &m)
            if err != nil {
                http.Error(w, err.Error(), http.StatusInternalServerError)
                return
            }
        }
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

func fetchItem(w http.ResponseWriter, r *http.Request, itemid string){
    c := appengine.NewContext(r)
    client := urlfetch.Client(c)

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
            fmt.Fprintf(w, "Value of Param1 is: %s", m.By)
            _, err := datastore.Put(c, datastore.NewIncompleteKey(c, "item", nil), &m)
            if err != nil {
                http.Error(w, err.Error(), http.StatusInternalServerError)
                return
            }
        }

    }
}