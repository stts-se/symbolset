package main

import (
	"context"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"log/syslog"
	"net/http"
	"os"
	"os/exec"
	"os/signal"
	"os/user"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"github.com/gorilla/mux"
)

// GLOBAL FLAGS
var symbolSetFileArea *string // = filepath.Join(".", "symbol_files")
var host, port *string

const staticFolder = "static"

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

//
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

func indexHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	html := `<h1>Symbolset</h1>`

	for _, subRouter := range subRouters {
		html = html + `<p><a href="` + subRouter.root + `"><b>` + removeInitialSlash(subRouter.root) + `</b></a>`
		html = html + " | " + subRouter.desc + "</p>\n\n"

	}
	html = html + "<p/><hr/><a href='/version'>About Symbolset</a>"
	fmt.Fprint(w, html)
}

func isStaticPage(url string) bool {
	return url == "/" || strings.Contains(url, "externals") || strings.Contains(url, "built") || url == "/websockreg" || url == "/favicon.ico" || url == "/static/" || url == "/ipa_table" || url == "/ping" || url == "/version"
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
var startedTimestamp = time.Now()              //.UTC().Format("2006-01-02 15:04:05 MST")
const timestampFmt = "2006-01-02 15:04:05 CET" // time.UnixDate // "%Y-%m-%d %H:%M:%S" // "2019-11-04 15:34 CET"

const buildInfoFile = "buildinfo.txt"

func getBuildInfo(prefix string, lines []string, defaultValue string) []string {
	for _, l := range lines {
		fs := strings.Split(l, ": ")
		if fs[0] == prefix {
			return fs
		}
	}
	return []string{prefix, defaultValue}
}

func generateAbout(w http.ResponseWriter, r *http.Request) {
	var buildInfoLines []string
	if _, err := os.Stat(buildInfoFile); os.IsNotExist(err) {
		log.Printf("no build info file, will generate about info on-the-fly")
	} else {
		bytes, err := ioutil.ReadFile(filepath.Clean(buildInfoFile))
		if err != nil {
			log.Printf("failed loading buildinfo file : %v", err)
		}
		buildInfoLines = strings.Split(strings.TrimSpace(string(bytes)), "\n")
	}

	res := [][]string{}
	res = append(res, []string{"Application name", "Symbolset"})

	// build timestamp
	res = append(res, getBuildInfo("Build timestamp", buildInfoLines, "n/a"))
	user, err := user.Current()
	if err != nil {
		log.Printf("failed reading system user name : %v", err)
	}

	// built by username
	res = append(res, getBuildInfo("Built by", buildInfoLines, user.Username))

	// git commit id and branch
	commitIDLong, err := exec.Command("git", "rev-parse", "HEAD").Output()
	var commitIDAndBranch = "unknown"
	if err != nil {
		log.Printf("couldn't retrieve git commit hash: %v", err)
	} else {
		commitID := string([]rune(string(commitIDLong)[0:7]))
		branch, err := exec.Command("git", "rev-parse", "--abbrev-ref", "HEAD").Output()
		if err != nil {
			log.Printf("couldn't retrieve git branch: %v", err)
		} else {
			commitIDAndBranch = fmt.Sprintf("%s on %s", commitID, strings.TrimSpace(string(branch)))
		}
	}
	res = append(res, getBuildInfo("Git commit", buildInfoLines, commitIDAndBranch))

	// git release tag
	releaseTag, err := exec.Command("git", "describe", "--tags").Output()
	if err != nil {
		log.Printf("couldn't retrieve git release/tag: %v", err)
		releaseTag = []byte("unknown")
	}
	res = append(res, getBuildInfo("Release", buildInfoLines, string(releaseTag)))

	res = append(res, []string{"Started", startedTimestamp.Format(timestampFmt)})
	res = append(res, []string{"Host", *host})
	//res = append(res, []string{"Port", port})
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	fmt.Fprintf(w, "<html><head><title>%s</title></head><body>", "Symbolset: About")
	fmt.Fprintf(w, "<table><tbody>")
	for _, l := range res {
		fmt.Fprintf(w, "<tr><td>%s</td><td>%s</td></tr>\n", l[0], l[1])
	}
	fmt.Fprintf(w, "</tbody></table>")
	fmt.Fprintf(w, "</body></html>")
}

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
	port = flag.String("port", "8771", "Server `port`")
	host = flag.String("host", "127.0.0.1", "Server `host`")
	logger := flag.String("logger", "stderr", "System `logger` (stderr, syslog or filename)")
	symbolSetFileArea = flag.String("ss_files", "", "`folder` with symbol set files (required)")

	var printUsage = func() {
		fmt.Fprintf(os.Stderr, "Usage:\n")
		flag.PrintDefaults()
	}

	flag.Usage = func() {
		printUsage()
		os.Exit(0)
	}
	flag.Parse()

	if *logger == "stderr" {
		// default logger
	} else if *logger == "syslog" {
		writer, err := syslog.New(syslog.LOG_INFO, "symbolset")
		if err != nil {
			log.Fatalf("Couldn't create logger: %v", err)
		}
		log.SetOutput(writer)
		log.SetFlags(0) // no timestamps etc, since syslog already prints that
	} else {
		f, err := os.OpenFile(*logger, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0600)
		if err != nil {
			log.Fatalf("Couldn't create logger: %v", err)
		}
		defer func() {
			err = f.Close()
			if err != nil {
				log.Fatalf("Couldn't close logger: %v", err)
			}
		}()
		log.SetOutput(f)
	}
	log.Println("server: created logger for " + *logger)

	if strings.TrimSpace(*symbolSetFileArea) == "" {
		fmt.Fprintf(os.Stderr, "-ss_files is required\n")
		printUsage()
		os.Exit(1)
	}
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
	rout.HandleFunc("/version", generateAbout).Methods("GET")

	symbolset := newSubRouter(rout, "/symbolset", "Handle transcription symbol sets")
	symbolset.addHandler(symbolsetList)
	symbolset.addHandler(symbolsetDelete)
	symbolset.addHandler(symbolsetContent)
	symbolset.addHandler(symbolsetReloadOne)
	symbolset.addHandler(symbolsetReloadAll)
	// symbolset.addHandler(symbolsetUploadPage)
	// symbolset.addHandler(symbolsetUpload)

	mapper := newSubRouter(rout, "/mapper", "Map transcriptions between different symbol sets")
	mapper.addHandler(mapperList)
	mapper.addHandler(mapperMap)
	mapper.addHandler(mapperMaptable)

	converter := newSubRouter(rout, "/converter", "Convert transcriptions between languages")
	converter.addHandler(converterConvert)
	converter.addHandler(converterList)
	converter.addHandler(converterTable)

	// static
	rout.HandleFunc("/ipa_table", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, filepath.Join(staticFolder, "ipa_table.txt"))
	})

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

	if _, err := os.Stat(staticFolder); os.IsNotExist(err) {
		log.Fatalf("Static folder does not exist: %s", staticFolder)
	}

	rout.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir(staticFolder))))

	// static
	rout.HandleFunc("/favicon.ico", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, filepath.Join(staticFolder, "favicon.ico"))
	})

	srv := &http.Server{
		Handler:      rout,
		Addr:         fmt.Sprintf("%s:%s", *host, *port),
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	log.Printf("Server started on %s", srv.Addr)

	stop := make(chan os.Signal, 1)

	signal.Notify(stop, os.Interrupt)
	go func() {
		if err := srv.ListenAndServe(); err != nil {
			log.Fatal(fmt.Errorf("server: couldn't start server on port %s : %v", *port, err))
		}
	}()
	log.Printf("server: server up and running using port %s", *port)

	<-stop

	// This happens after Ctrl-C
	fmt.Fprintf(os.Stderr, "\n")
	shutdown(srv)
	//log.Fatal(srv.ListenAndServe())
}

func shutdown(s *http.Server) {
	log.Println("server: shutting down...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	defer s.Shutdown(ctx)

	log.Println("server: shutdown completed")
}
