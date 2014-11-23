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
)

// [START greeting_struct]
type Greeting struct {
        Author  string
        Content string
        Date    time.Time
}

type User struct {
    Username string
    DeviceToken string
    Feeds []string
}

// [END greeting_struct]

func init() {
        http.HandleFunc("/", root)
        http.HandleFunc("/sign", sign)
        http.HandleFunc("/users", users)
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
    // res, err := http.NewRequest("http://www.feedforall.com/sample.xml")

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

        fu := User {
            Username: r.Header.Get("username"),
            DeviceToken: r.Header.Get("deviceToken"),
        }

        if u := user.Current(c); u != nil {
            fu.Username = u.String()
        }

        fmt.Fprintf(w, "Hello " + fu.Username + " device: " + fu.DeviceToken)
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