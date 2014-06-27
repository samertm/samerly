package server

import (
	"github.com/samertm/samerly/engine"

	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
)

var _ = fmt.Println // debugging

const indexHtml = `<!DOCTYPE html>
<html>
  <head>
    <title>samer.ly</title>
  </head>
  <body>
    <p>welcome to samer.ly, the world's most advanced url shortener!</p>
    <form action="create_url" method="post">
    <p>
      shorten this url: <input type="text" name="longurl" />
      <input type="submit" value="shorten!" />
    </p>
    </form>
  </body>
</html>`

// parseForm and checkForm pulled from hs-directory
// warning: modifies req by calling req.ParseForm()
func parseForm(req *http.Request, values ...string) (form url.Values, err error) {
	req.ParseForm()
	form = req.PostForm
	err = checkForm(form, values...)
	return
}

func checkForm(data url.Values, values ...string) error {
	for _, s := range values {
		if len(data[s]) == 0 {
			return errors.New(s + " not passed")
		}
	}
	return nil
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		io.WriteString(w, indexHtml)
	}
}

func createUrlHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		io.WriteString(w, "<!DOCTYPE html><html><body>Go <a href='/'>back</a>.</body></html>")
	} else if r.Method == "POST" {
		form, err := parseForm(r, "longurl")
		if err != nil {
			io.WriteString(w, err.Error())
			return
		}
		c := make(chan string)
		urls.AddUrl <- engine.Pair{form["longurl"][0], c}
		shortened := <-c
		io.WriteString(w,
			"<!DOCTYPE html><html><body>shortened url: <a href='"+
				"http://"+address+"/url/"+shortened+
				"'>samer.ly/url/"+shortened+"</body></html>")
	}
}

func urlHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		// R.URL.Path looks like the regex "/url/.*"
		// shave off first five chars
		shortened := r.URL.Path[5:]
		c := make(chan string)
		urls.GetUrl <- engine.Pair{shortened, c}
		if longurl, ok := <-c; ok {
			http.Redirect(w, r, longurl, 301)
		} else {
			io.WriteString(w, "<!DOCTYPE html><html><body>"+shortened+" does not exist ): I'M SO SORRY D:</body></html>")
		}
	}
}

func statsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		c := make(chan string)
		urls.GetStats <- engine.Pair{Recv: c}
		var fullstring string
		for s, ok := <-c; ok; s, ok = <-c {
			fullstring += s
		}
		fullstring = strings.Replace(fullstring, "\n", "<br>", -1)
		io.WriteString(w, "<!DOCTYPE html><html><body>" + fullstring + "</body></html>")
	}
}

var urls *engine.Urls
var address string

func init() {
	urls = engine.NewUrls()
	go urls.Run()
}

func ListenAndServe(ip string) {
	address = ip
	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/create_url", createUrlHandler)
	http.HandleFunc("/url/", urlHandler)
	http.HandleFunc("/stats", statsHandler)
	http.ListenAndServe(ip, nil)
}
