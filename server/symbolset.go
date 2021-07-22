package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"sort"

	//"os"
	"encoding/json"
	"path/filepath"
	"strings"

	"github.com/gorilla/mux"
	"github.com/stts-se/symbolset"
)

// JSONSymbolSet : JSON container
type JSONSymbolSet struct {
	Name    string
	Symbols []JSONSymbol
}

// JSONSymbol : JSON container
type JSONSymbol struct {
	Symbol string
	IPA    JSONIPA
	Desc   string
	Cat    string
}

// JSONIPA : JSON container
type JSONIPA struct {
	String  string
	Unicode string
}

func (rout *subRouter) addHandler(handler urlHandler) {
	rout.router.HandleFunc(handler.url, handler.handler)
	rout.handlers = append(rout.handlers, handler)
}

type subRouter struct {
	root     string
	router   *mux.Router
	handlers []urlHandler
	desc     string
}

var subRouters []*subRouter

type urlHandler struct {
	name     string
	handler  func(w http.ResponseWriter, r *http.Request)
	url      string
	help     string
	examples []string
}

func loadSymbolSetFile(fName string) (symbolset.SymbolSet, error) {
	return symbolset.LoadSymbolSet(fName)
}

var symbolsetContent = urlHandler{
	name:     "content",
	url:      "/content/{name}",
	help:     "Lists content of a named symbolset.",
	examples: []string{"/content/sv-se_ws-sampa-DEMO"},
	handler: func(w http.ResponseWriter, r *http.Request) {
		name := getParam("name", r)
		if len(strings.TrimSpace(name)) == 0 {
			msg := "symbol set should be specified by variable 'name'"
			log.Println(msg)
			http.Error(w, msg, http.StatusInternalServerError)
			return
		}
		mMut.Lock()
		symbolset0, ok := mMut.service.SymbolSets[name]
		mMut.Unlock()
		if !ok {
			msg := fmt.Sprintf("failed getting symbol set : %v", name)
			log.Println(msg)
			http.Error(w, msg, http.StatusInternalServerError)
			return
		}
		symbolset := JSONSymbolSet{Name: symbolset0.Name}
		symbolset.Symbols = make([]JSONSymbol, 0)
		for _, sym := range symbolset0.Symbols {
			symbolset.Symbols = append(symbolset.Symbols, JSONSymbol{Symbol: sym.String, IPA: JSONIPA{String: sym.IPA.String, Unicode: sym.IPA.Unicode}, Desc: sym.Desc, Cat: sym.Cat.String()})
		}

		j, err := json.Marshal(symbolset)
		if err != nil {
			msg := fmt.Sprintf("json marshalling error : %v", err)
			log.Println(msg)
			http.Error(w, msg, http.StatusInternalServerError)
			return
		}
		fmt.Fprint(w, string(j))
	},
}

func loadSymbolSets(dirName string) error {
	mMut.Lock()
	mMut.service.Clear()
	mMut.Unlock()

	symbolSets, err := symbolset.LoadSymbolSetsFromDir(dirName)
	if err != nil {
		return err
	}
	mMut.Lock()
	mMut.service.SymbolSets = symbolSets
	mMut.Unlock()

	mappersDef := filepath.Join(dirName, "mappers.txt")
	return testMappers(mappersDef)
}

func symbolsetReloadAllHandler(w http.ResponseWriter, r *http.Request) {
	err := loadSymbolSets(*symbolSetFileArea)
	if err != nil {
		msg := err.Error()
		log.Println(msg)
		http.Error(w, msg, http.StatusInternalServerError)
		return
	}
	mMut.Lock()
	j, err := json.Marshal(symbolSetNames(mMut.service.SymbolSets))
	mMut.Unlock()
	if err != nil {
		msg := fmt.Sprintf("json marshalling error : %v", err)
		log.Println(msg)
		http.Error(w, msg, http.StatusInternalServerError)
		return
	}

	fmt.Fprint(w, string(j))

}

func symbolsetReloadOneHandler(w http.ResponseWriter, r *http.Request) {
	name := getParam("name", r)
	mMut.Lock()
	err := mMut.service.DeleteSymbolSet(name)
	mMut.Unlock()
	if err != nil {
		msg := fmt.Sprintf("couldn't delete symbolset : %v", err)
		log.Println(msg)
		http.Error(w, msg, http.StatusInternalServerError)
		return
	}

	serverPath := filepath.Join(*symbolSetFileArea, name+symbolset.SymbolSetSuffix)
	mMut.Lock()
	err = mMut.service.Load(serverPath)
	mMut.Unlock()
	if err != nil {
		msg := fmt.Sprintf("couldn't load symbolset : %v", err)
		log.Println(msg)
		http.Error(w, msg, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	msg := fmt.Sprintf("Reloaded symbol set %s", name)
	fmt.Fprint(w, msg)

}

var symbolsetReloadOne = urlHandler{
	name:     "reload",
	url:      "/reload/{name}",
	help:     "Reloads a named symbol set in the pre-defined folder.",
	examples: []string{"/reload/sv-se_nst-xsampa-DEMO"},
	handler:  symbolsetReloadOneHandler,
}

var symbolsetReloadAll = urlHandler{
	name:     "reload",
	url:      "/reload",
	help:     "Reloads all symbol set(s) in the pre-defined folder.",
	examples: []string{"/reload"},
	handler:  symbolsetReloadAllHandler,
}

var symbolsetDelete = urlHandler{
	name:     "delete",
	url:      "/delete/{name}",
	help:     "Deletes a named symbol set.",
	examples: []string{},
	handler: func(w http.ResponseWriter, r *http.Request) {
		name := getParam("name", r)
		if len(strings.TrimSpace(name)) == 0 {
			msg := "symbol set should be specified by variable 'name'"
			log.Println(msg)
			http.Error(w, msg, http.StatusInternalServerError)
			return
		}
		mMut.Lock()
		err := mMut.service.DeleteSymbolSet(name)
		mMut.Unlock()
		if err != nil {
			msg := fmt.Sprintf("couldn't delete symbolset : %v", err)
			log.Println(msg)
			http.Error(w, msg, http.StatusInternalServerError)
			return
		}

		serverPath := filepath.Join(*symbolSetFileArea, name+symbolset.SymbolSetSuffix)
		if _, err := os.Stat(serverPath); err != nil {
			if os.IsNotExist(err) {
				msg := fmt.Sprintf("couldn't locate server file for symbol set %s", name)
				log.Println(msg)
				http.Error(w, msg, http.StatusInternalServerError)
				return
			}
		}

		err = os.Remove(serverPath)
		if err != nil {
			msg := fmt.Sprintf("couldn't delete file from server : %v", err)
			log.Println(msg)
			http.Error(w, msg, http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		msg := fmt.Sprintf("Deleted symbol set %s", name)
		fmt.Fprint(w, msg)
	},
}

var symbolsetList = urlHandler{
	name:     "list",
	url:      "/list",
	help:     "Lists available symbol sets.",
	examples: []string{"/list"},
	handler: func(w http.ResponseWriter, r *http.Request) {
		mMut.Lock()
		ss := symbolSetNames(mMut.service.SymbolSets)
		mMut.Unlock()
		j, err := json.Marshal(ss)
		if err != nil {
			msg := fmt.Sprintf("failed to marshal struct : %v", err)
			log.Println(msg)
			http.Error(w, msg, http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		fmt.Fprint(w, string(j))
	},
}

/*
func symbolSetHelpHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	html := `<h1>SymbolSet</h1>

<h2>symbolset_upload</h2> Upload symbol set file
<pre><a href="/symbolset/upload">/symbolset/upload</a></pre>
		`

	fmt.Fprint(w, html)
}
*/
// var symbolsetUploadPage = urlHandler{
// 	name:     "upload (page)",
// 	url:      "/upload_page",
// 	help:     "Upload symbol set file (GUI)",
// 	examples: []string{"/upload_page"},
// 	handler: func(w http.ResponseWriter, r *http.Request) {
// 		uploadPage := filepath.Join(staticFolder, "upload_page.html")
// 		if _, err := os.Stat(uploadPage); os.IsNotExist(err) {
// 			msg := fmt.Sprintf("No such file: %s", uploadPage)
// 			httpError(w, msg, "404 page not found", http.StatusNotFound)
// 			return
// 		}
// 		http.ServeFile(w, r, uploadPage)
// 	},
// }

// var symbolsetUpload = urlHandler{
// 	name:     "upload (api)",
// 	url:      "/upload",
// 	help:     "Upload symbol set file (API)",
// 	examples: []string{},
// 	handler: func(w http.ResponseWriter, r *http.Request) {
// 		if r.Method != "POST" {
// 			http.Error(w, fmt.Sprintf("symbol set upload only accepts POST request, got %s", r.Method), http.StatusBadRequest)
// 			return
// 		}

// 		clientUUID := getParam("client_uuid", r)

// 		if strings.TrimSpace(clientUUID) == "" {
// 			msg := "doUploadSymbolSetHandler got no client uuid"
// 			log.Println(msg)
// 			http.Error(w, msg, http.StatusBadRequest)
// 			return
// 		}

// 		// (partially) lifted from https://github.com/astaxie/build-web-application-with-golang/blob/master/de/04.5.md

// 		err := r.ParseMultipartForm(32 << 20)
// 		if err != nil {
// 			log.Println(err)
// 			http.Error(w, fmt.Sprintf("doUploadSymbolSetHandler failed pars multipart form : %v", err), http.StatusInternalServerError)
// 			return
// 		}

// 		file, handler, err := r.FormFile("upload_file")
// 		if err != nil {
// 			log.Println(err)
// 			http.Error(w, fmt.Sprintf("doUploadSymbolSetHandler failed reading file : %v", err), http.StatusInternalServerError)
// 			return
// 		}
// 		defer file.Close()
// 		serverPath := filepath.Join(*symbolSetFileArea, handler.Filename)
// 		if _, err := os.Stat(serverPath); err == nil {
// 			msg := fmt.Sprintf("symbol set already exists on server in file: %s", handler.Filename)
// 			log.Println(msg)
// 			http.Error(w, msg, http.StatusInternalServerError)
// 			return
// 		}

// 		f, err := os.OpenFile(serverPath, os.O_WRONLY|os.O_CREATE, 0600)
// 		if err != nil {
// 			log.Println(err)
// 			http.Error(w, fmt.Sprintf("doUploadSymbolSetHandler failed opening local output file : %v", err), http.StatusInternalServerError)
// 			return
// 		}
// 		/* #nosec G307 */
// 		defer f.Close()
// 		_, err = io.Copy(f, file)
// 		if err != nil {
// 			msg := fmt.Sprintf("doUploadSymbolSetHandler failed copying local output file : %v", err)
// 			log.Println(msg)
// 			http.Error(w, msg, http.StatusInternalServerError)
// 			return
// 		}
// 		ss, err := loadSymbolSetFile(serverPath)
// 		if err != nil {
// 			msg := fmt.Sprintf("couldn't load symbol set file : %v", err)
// 			err = os.Remove(serverPath)
// 			if err != nil {
// 				msg = fmt.Sprintf("%v (couldn't delete file from server)", msg)
// 			} else {
// 				msg = fmt.Sprintf("%v (the uploaded file has been deleted from server)", msg)
// 			}
// 			log.Println(msg)
// 			http.Error(w, msg, http.StatusInternalServerError)
// 			return
// 		}

// 		//f.Close()

// 		mMut.Lock()
// 		mMut.service.SymbolSets[ss.Name] = ss
// 		mMut.Unlock()

// 		fmt.Fprintf(w, "%v", handler.Header)
// 	},
// }

func symbolSetNames(sss map[string]symbolset.SymbolSet) []string {
	var ssNames []string
	for ss := range sss {
		ssNames = append(ssNames, ss)
	}
	sort.Strings(ssNames)
	return ssNames
}
