package main

import (
	"encoding/json"
	"flag"
	"github.com/benmanns/onhub/diagnosticreport"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
)

type loadMode int

const (
	viaHTTP loadMode = iota
	viaFile
)

func main() {
	flag.Parse()
	source := flag.Arg(0)
	if source == "" {
		source = "http://192.168.86.1/api/v1/diagnostic-report"
	}
	var mode loadMode
	if uri, err := url.ParseRequestURI(source); err != nil || uri.Scheme == "" {
		mode = viaFile
	} else {
		mode = viaHTTP
	}
	var reader io.Reader
	switch mode {
	case viaHTTP:
		resp, err := http.Get(source)
		if err != nil {
			log.Fatalln("error getting report:", err)
		}
		defer func() {
			if err := resp.Body.Close(); err != nil {
				log.Fatalln("error closing gotten report:", err)
			}
		}()
		reader = resp.Body
	case viaFile:
		file, err := os.Open(source)
		if err != nil {
			log.Fatalln("error opening report:", err)
		}
		defer func() {
			if err := file.Close(); err != nil {
				log.Fatalln("error closing opened report:", err)
			}
		}()
		reader = file
	default:
		log.Fatalln("bad source mode")
	}

	data, err := ioutil.ReadAll(reader)
	if err != nil {
		log.Fatalln("error reading report:", err)
	}

	dr, err := diagnosticreport.Parse(data)
	if err != nil {
		log.Fatalln("error parsing report:", err)
	}

	encoder := json.NewEncoder(os.Stdout)
	if err := encoder.Encode(dr); err != nil {
		log.Fatalln("error encoding report:", err)
	}
}
