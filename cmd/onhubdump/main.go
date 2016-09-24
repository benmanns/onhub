package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/benmanns/onhub/diagnosticreport"
	"io"
	"io/ioutil"
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
			fmt.Fprintln(os.Stderr, "error getting report:", err)
			os.Exit(1)
		}
		defer func() {
			if err := resp.Body.Close(); err != nil {
				fmt.Fprintln(os.Stderr, "error closing gotten report:", err)
				os.Exit(1)
			}
		}()
		reader = resp.Body
	case viaFile:
		file, err := os.Open(source)
		if err != nil {
			fmt.Fprintln(os.Stderr, "error opening report:", err)
			os.Exit(1)
		}
		defer func() {
			if err := file.Close(); err != nil {
				fmt.Fprintln(os.Stderr, "error closing opened report:", err)
				os.Exit(1)
			}
		}()
		reader = file
	default:
		fmt.Fprintln(os.Stderr, "bad source mode")
		os.Exit(1)
	}

	data, err := ioutil.ReadAll(reader)
	if err != nil {
		fmt.Fprintln(os.Stderr, "error reading report:", err)
		os.Exit(1)
	}

	dr, err := diagnosticreport.Parse(data)
	if err != nil {
		fmt.Fprintln(os.Stderr, "error parsing report:", err)
		os.Exit(1)
	}

	encoder := json.NewEncoder(os.Stdout)
	if err := encoder.Encode(dr); err != nil {
		fmt.Fprintln(os.Stderr, "error encoding report:", err)
		os.Exit(1)
	}
}
