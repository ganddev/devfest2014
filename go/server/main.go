package feedme

import (
        "html/template"
        "net/http"
        "time"
        "fmt"
        // "os"
    // "io"
        // "io/ioutil"

        "appengine"
        "appengine/datastore"
        "appengine/user"
        // "github.com/iand/feedparser"

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
    DeviceToken string
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
        http.HandleFunc("/sign", sign)
        http.HandleFunc("/users", users)
        
        http.HandleFunc("/send", sendNotification)
        http.HandleFunc("/fetch", fetchEntry)
}

func deviceKey(c appengine.Context) *datastore.Key {
    return datastore.NewKey(c, "Device", "default_device", 0, nil)
}

// guestbookKey returns the key used for all guestbook entries.
func guestbookKey(c appengine.Context) *datastore.Key {
        // The string "default_guestbook" here could be varied to have multiple guestbooks.
        return datastore.NewKey(c, "Guestbook", "default_guestbook", 0, nil)
}

// [START func_root]
func root(w http.ResponseWriter, r *http.Request) {

        // resp, err := http.Get("http://www.feedforall.com/sample.xml")

        // if err == nil {
        //     fmt.Fprint(w, resp.Body)
        // } else {
        //     fmt.Fprint(w, "blergh!")
        // }

    // fe := feedparser.Feed{"asdf", "asdf", "http://www.feedforall.com/sample.xml", nil}
    // res, err := http.NewRequest("GET", "http://www.feedforall.com/sample.xml", nil)

    // if err == nil {
    //     feed := feedparser.NewFeed(res.Body)
    //     fmt.Fprint(w, feed)
    // }

        // if err != nil {
        //     fmt.Fprint(w, "error: " + err)
        //     os.Exit(1)
        // } else {
        //     Feed.NewFeed(ioutil.ReadAll(response.Body))
        // }

        c := appengine.NewContext(r)
        // Ancestor queries, as shown here, are strongly consistent with the High
        // Replication Datastore. Queries that span entity groups are eventually
        // consistent. If we omitted the .Ancestor from this query there would be
        // a slight chance that Greeting that had just been written would not
        // show up in a query.
        // [START query]
        q := datastore.NewQuery("Greeting").Ancestor(guestbookKey(c)).Order("-Date").Limit(10)
        // [END query]
        // [START getall]
        greetings := make([]Greeting, 0, 10)
        if _, err := q.GetAll(c, &greetings); err != nil {
                http.Error(w, err.Error(), http.StatusInternalServerError)
                return
        }
        // [END getall]
        if err := guestbookTemplate.Execute(w, greetings); err != nil {
                http.Error(w, err.Error(), http.StatusInternalServerError)
        }
}
// [END func_root]

var guestbookTemplate = template.Must(template.New("book").Parse(`
<html>
  <head>
    <title>Go Guestbook</title>
  </head>
  <body>
    {{range .}}
      {{with .Author}}
        <p><b>{{.}}</b> wrote:</p>
      {{else}}
        <p>An anonymous person wrote:</p>
      {{end}}
      <pre>{{.Content}}</pre>
    {{end}}
    <form action="/sign" method="post">
      <div><textarea name="content" rows="3" cols="60"></textarea></div>
      <div><input type="submit" value="Sign Guestbook"></div>
    </form>
  </body>
</html>
`))

func users(w http.ResponseWriter, r *http.Request) {
    if r.Method == "POST" {
        c := appengine.NewContext(r)
        fu := Device{
            DeviceToken: r.Header.Get("deviceToken"),
        }

        // query := datastore.NewQuery("User").Filter("DeviceToken =", fu.DeviceToken).

        // var users []User
        // _, err := query.GetAll(c, &users)

        // if err != nil {
        //     fmt.Fprint(w, "BLABLA")
        // }

        key := datastore.NewIncompleteKey(c, "Device", deviceKey(c))
        _, err := datastore.Put(c, key, &fu)
        
        if err != nil {
                http.Error(w, err.Error(), http.StatusInternalServerError)
                return
        }

        fmt.Fprintf(w, "Hello " + " device: " + fu.DeviceToken)
    } else if r.Method == "GET" {
        fmt.Fprint(w, "GETTING SHIT HERE")
    }
}

// [START func_sign]
func sign(w http.ResponseWriter, r *http.Request) {
        c := appengine.NewContext(r)
        g := Greeting{
                Content: r.FormValue("content"),
                Date:    time.Now(),
        }
        if u := user.Current(c); u != nil {
                g.Author = u.String()
        }


        // We set the same parent key on every Greeting entity to ensure each Greeting
        // is in the same entity group. Queries across the single entity group
        // will be consistent. However, the write rate to a single entity group
        // should be limited to ~1/second.
        key := datastore.NewIncompleteKey(c, "Greeting", guestbookKey(c))
        _, err := datastore.Put(c, key, &g)
        if err != nil {
                http.Error(w, err.Error(), http.StatusInternalServerError)
                return
        }
        http.Redirect(w, r, "/", http.StatusFound)
}
// [END func_sign]







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