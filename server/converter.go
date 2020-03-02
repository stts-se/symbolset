package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"sync"

	"github.com/stts-se/symbolset/converter"
)

var cMut = struct {
	sync.RWMutex
	service map[string]converter.Converter
}{
	service: make(map[string]converter.Converter),
}

// JSONConverted : JSON container
type JSONConverted struct {
	Converter string
	Input     string
	Result    string
}

var converterConvert = urlHandler{
	name: "convert",
	url:  "/convert/{converter}/{trans}",
	help: "Maps a transcription using a specified converter.",
	examples: []string{"/convert/enusampa_svsampa-DEMO/%22 D i s",
		"/convert/enusampa_svsampa-DEMO/%22 D EI . z i"},
	handler: func(w http.ResponseWriter, r *http.Request) {
		convName := getParam("converter", r)
		trans := trimTrans(getParam("trans", r))
		if len(strings.TrimSpace(convName)) == 0 {
			msg := "converter name should be specified by variable 'converter'"
			log.Println(msg)
			http.Error(w, msg, http.StatusInternalServerError)
			return
		}
		if len(strings.TrimSpace(trans)) == 0 {
			msg := "input trans should be specified by variable 'trans'"
			log.Println(msg)
			http.Error(w, msg, http.StatusInternalServerError)
			return
		}
		cMut.RLock()
		defer cMut.RUnlock()
		conv, ok := cMut.service[convName]
		if !ok {
			msg := fmt.Sprintf("no converter named : %s", convName)
			log.Println(msg)
			http.Error(w, msg, http.StatusInternalServerError)
			return
		}
		result0, err := conv.Convert(trans)
		if err != nil {
			msg := fmt.Sprintf("failed converting transcription : %v", err)
			log.Println(msg)
			http.Error(w, msg, http.StatusInternalServerError)
			return
		}
		result := JSONConverted{Input: trans, Result: result0, Converter: convName}
		j, err := json.Marshal(result)
		if err != nil {
			msg := fmt.Sprintf("json marshalling error : %v", err)
			log.Println(msg)
			http.Error(w, msg, http.StatusInternalServerError)
			return
		}
		fmt.Fprint(w, string(j))
	},
}

// JSONConverter : JSON container
type JSONConverter struct {
	Name  string
	From  string
	To    string
	Rules []JSONCRule
}

// JSONCRule : JSON container
type JSONCRule struct {
	Type string
	From string
	To   string
}

var converterTable = urlHandler{
	name:     "table",
	url:      "/table/{converter}",
	help:     "Lists map table for a specified converter.",
	examples: []string{"/table/enusampa_svsampa-DEMO"},
	handler: func(w http.ResponseWriter, r *http.Request) {
		convName := getParam("converter", r)
		if len(strings.TrimSpace(convName)) == 0 {
			msg := "converter name should be specified by variable 'converter'"
			log.Println(msg)
			http.Error(w, msg, http.StatusInternalServerError)
			return
		}
		cMut.RLock()
		defer cMut.RUnlock()
		conv, ok := cMut.service[convName]
		if !ok {
			msg := fmt.Sprintf("no converter named : %s", convName)
			log.Println(msg)
			http.Error(w, msg, http.StatusInternalServerError)
			return
		}
		jConv := JSONConverter{Name: conv.Name, From: conv.From.Name, To: conv.To.Name}
		jConv.Rules = make([]JSONCRule, 0)
		for _, rule := range conv.Rules {
			jConv.Rules = append(jConv.Rules, JSONCRule{Type: rule.Type(), From: rule.FromString(), To: rule.ToString()})
		}

		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		j, err := json.Marshal(jConv)
		if err != nil {
			msg := fmt.Sprintf("json marshalling error : %v", err)
			log.Println(msg)
			http.Error(w, msg, http.StatusInternalServerError)
			return
		}
		fmt.Fprint(w, string(j))
	},
}
var converterList = urlHandler{
	name:     "list",
	url:      "/list",
	help:     "Lists available converters.",
	examples: []string{"/list"},
	handler: func(w http.ResponseWriter, r *http.Request) {
		cMut.RLock()
		defer cMut.RUnlock()
		cs := []string{}
		for key := range cMut.service {
			cs = append(cs, key)
		}
		j, err := json.Marshal(cs)
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

func loadConverters(dirName string) error {
	cMut.Lock()
	mMut.RLock()

	cMut.service = make(map[string]converter.Converter)

	convs, testRes, err := converter.LoadFromDir(mMut.service.SymbolSets, dirName)
	if err != nil {
		return err
	}
	allOK := true
	for cName, tr := range testRes {
		if !tr.OK {
			allOK = false
			log.Printf("INIT TESTS FAILED FOR %s: %v", cName, tr)
		}
		log.Println("server: loaded converter", cName)
	}
	for _, conv := range convs {
		cMut.service[conv.Name] = conv
	}

	defer cMut.Unlock()
	defer mMut.RUnlock()

	if !allOK {
		return fmt.Errorf("FAIL")
	}
	return nil
}
