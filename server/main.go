package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"github.com/gorilla/mux"
)

// TODO should go into config file
var symbolSetFileArea *string // = filepath.Join(".", "symbol_files")
var staticFolder string       // = "."

func getParam(paramName string, r *http.Request) string {
	res := r.FormValue(paramName)
	if res != "" {
		return res
	}
	vars := mux.Vars(r)
	return vars[paramName]
}

// print serverMsg to server log, and return an http error with clientMsg and the specified error code (http.StatusInternalServerError, etc)
// func httpError(w http.ResponseWriter, serverMsg string, clientMsg string, errCode int) {
// 	log.Println(serverMsg)
// 	http.Error(w, clientMsg, errCode)
// }

// func readFile(fName string) ([]string, error) {
// 	bytes, err := ioutil.ReadFile(fName)
// 	if err != nil {
// 		return []string{}, err
// 	}
// 	return strings.Split(strings.TrimSpace(string(bytes)), "\n"), nil
// }

func pingHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "symbolset")
}

func versionHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, strings.Join(vInfo, "\n"))
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	html := `<h1>Symbolset</h1>`

	for _, subRouter := range subRouters {
		html = html + `<p><a href="` + subRouter.root + `"><b>` + removeInitialSlash(subRouter.root) + `</b></a>`
		html = html + " | " + subRouter.desc + "</p>\n\n"

	}
	html = html + "<p/><hr/><a href='/version'>Symbolset version info</a>"
	fmt.Fprint(w, html)
}

func isStaticPage(url string) bool {
	return url == "/" || strings.Contains(url, "externals") || strings.Contains(url, "built") || url == "/websockreg" || url == "/favicon.ico" || url == "/static/" || url == "/ipa_table.txt" || url == "/ping" || url == "/version"
}

var initialSlashRe = regexp.MustCompile("^/")

func removeInitialSlash(url string) string {
	return initialSlashRe.ReplaceAllString(url, "")
}

// TODO: Neat URL encoding...
func urlEnc(url string) string {
	return strings.Replace(strings.Replace(strings.Replace(strings.Replace(url, " ", "%20", -1), "\n", "", -1), `"`, "%22", -1), "\t", "", -1)
}

func (h urlHandler) helpHTML(root string) string {
	s := "<h2>" + h.name + "</h2> " + h.help
	if strings.Contains(h.url, "{") {
		s = s + `<p>API URL: <code>` + root + h.url + `</code></p>`
	}
	if len(h.examples) > 0 {
		//s = s + `<p>Example invocation:`
		for _, x := range h.examples {
			urlPretty := root + x
			url := root + urlEnc(x)
			s = s + `<pre><a href="` + url + `">` + urlPretty + `</a></pre>`
		}
		//s = s + "</p>"
	}
	return s
}
func isHandeledPage(url string) bool {
	for _, sub := range subRouters {
		if sub.root == url || sub.root+"/" == url {
			return true
		}
		for _, handler := range sub.handlers {
			if sub.root+handler.url == url {
				return true
			}
		}
	}
	return false
}

// UTC time with format: yyyy-MM-dd HH:mm:ss z | %Y-%m-%d %H:%M:%S %Z
var startedTimestamp = time.Now().UTC().Format("2006-01-02 15:04:05 MST")

func getVersionInfo() []string {
	res := []string{}
	var buildInfoFile = "/wikispeech/symbolset/build_info.txt"
	if _, err := os.Stat(buildInfoFile); os.IsNotExist(err) {
		var msg = fmt.Sprintf("server: build info not defined : no such file: %s\n", buildInfoFile)
		log.Print(msg)
		res = append(res, "Application name: symbolset")
		res = append(res, "Build timestamp: n/a")
		res = append(res, "Built by: user")
		tag, err := exec.Command("git", "describe", "--tags").Output()
		if err != nil {
			log.Printf("server: couldn't retrieve git release info: %v", err)
			res = append(res, "Release: unknown")
		} else {
			branch, err := exec.Command("git", "rev-parse", "--abbrev-ref", "HEAD").Output()
			if err != nil {
				log.Printf("server: couldn't retrieve git release info: %v", err)
				res = append(res, "Release: unknown")
			}
			res = append(res, strings.TrimSpace(fmt.Sprintf("Release: %s on branch %s",
				strings.TrimSpace(string(tag)),
				strings.TrimSpace(string(branch)))))
		}
	} else {

		fBytes, err := ioutil.ReadFile(buildInfoFile)
		if err != nil {
			var msg = fmt.Sprintf("server: error reading file : %v", err)
			log.Print(msg)
		}
		if _, err = os.Stat(buildInfoFile); os.IsNotExist(err) {
			var msg = fmt.Sprintf("server: error when reading content from timestamp file : %v", err)
			log.Print(msg)
		} else {
			res = strings.Split(strings.TrimSpace(string(fBytes)), "\n")
		}
	}
	res = append(res, "Started: "+startedTimestamp)
	log.Println("server: parsed version info", res)
	return res
}

var vInfo []string

func newSubRouter(rout *mux.Router, root string, description string) *subRouter {
	var res = subRouter{
		router: rout.PathPrefix(root).Subrouter(),
		root:   root,
		desc:   description,
	}

	helpHandler := func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")

		html := "<h1>" + removeInitialSlash(res.root) + "</h1> <em>" + res.desc + "</em>"
		for _, handler := range res.handlers {
			html = html + handler.helpHTML(res.root)
		}
		fmt.Fprint(w, html)
	}

	res.router.HandleFunc("/", helpHandler)
	subRouters = append(subRouters, &res)
	return &res
}

func main() {
	port := flag.String("port", "8771", "Server `port`")
	host := flag.String("host", "127.0.0.1", "Server `host`")
	symbolSetFileArea = flag.String("ss_files", filepath.Join(".", "symbol_sets"), "location for symbol set files")

	var printUsage = func() {
		fmt.Fprintf(os.Stderr, "Usage:\n")
		flag.PrintDefaults()
	}

	flag.Usage = func() {
		printUsage()
		os.Exit(0)
	}
	flag.Parse()

	vInfo = getVersionInfo()

	err := loadSymbolSets(*symbolSetFileArea)
	if err != nil {
		log.Fatalf("failed to load symbol sets from dir %s : %v", *symbolSetFileArea, err)
	}
	log.Printf("server: loaded symbol sets from dir %s", *symbolSetFileArea)

	err = loadConverters(*symbolSetFileArea)
	if err != nil {
		log.Fatalf("failed to load converters from dir %s : %v", *symbolSetFileArea, err)
	}
	log.Printf("server: loaded converters from dir %s", *symbolSetFileArea)

	rout := mux.NewRouter().StrictSlash(true)

	rout.HandleFunc("/", indexHandler)
	rout.HandleFunc("/ping", pingHandler)
	rout.HandleFunc("/version", versionHandler)

	symbolset := newSubRouter(rout, "/symbolset", "Handle transcription symbol sets")
	symbolset.addHandler(symbolsetList)
	symbolset.addHandler(symbolsetDelete)
	symbolset.addHandler(symbolsetContent)
	symbolset.addHandler(symbolsetReloadOne)
	symbolset.addHandler(symbolsetReloadAll)
	symbolset.addHandler(symbolsetUploadPage)
	symbolset.addHandler(symbolsetUpload)

	mapper := newSubRouter(rout, "/mapper", "Map transcriptions between different symbol sets")
	mapper.addHandler(mapperList)
	mapper.addHandler(mapperMap)
	mapper.addHandler(mapperMaptable)

	converter := newSubRouter(rout, "/converter", "Convert transcriptions between languages")
	converter.addHandler(converterConvert)
	converter.addHandler(converterList)
	converter.addHandler(converterTable)

	var urls = []string{}
	errW := rout.Walk(func(route *mux.Route, router *mux.Router, ancestors []*mux.Route) error {
		url, err := route.GetPathTemplate()
		if err != nil {
			return err
		}
		if !isStaticPage(url) && !isHandeledPage(url) {
			log.Print("Unhandeled url: ", url)
			urls = append(urls, url+" (UNHANDELED)")
		} else {
			urls = append(urls, url)
		}
		return nil
	})

	if errW != nil {
		log.Printf("server failed to walk through route handlers : %v", errW)
	}

	meta := newSubRouter(rout, "/meta", "Meta API calls (list served URLs, etc)")
	meta.addHandler(metaURLsHandler(urls))
	meta.addHandler(metaExamplesHandler)

	// static
	rout.HandleFunc("/favicon.ico", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, filepath.Join(staticFolder, "favicon.ico"))
	})
	rout.HandleFunc("/ipa_table.txt", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, filepath.Join(staticFolder, "ipa_table.txt"))
	})
	rout.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir(staticFolder))))

	srv := &http.Server{
		Handler:      rout,
		Addr:         fmt.Sprintf("%s:%s", *host, *port),
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	fmt.Printf("Server started on %s\n", srv.Addr)

	log.Fatal(srv.ListenAndServe())
	fmt.Println("No fun")

}
