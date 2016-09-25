package diagnosticreport

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"io/ioutil"
)

// From compress/gzip/gunzip.go
const (
	gzipID1     = 0x1f
	gzipID2     = 0x8b
	gzipDeflate = 8
)

type gzippedString []byte

// MarshalJSON is overridden to decompress the content field if it is gzip
// compressed.
func (f File) MarshalJSON() ([]byte, error) {
	return json.Marshal(&struct {
		Path    string        `json:"path,omitempty"`
		Content gzippedString `json:"content,omitempty"`
	}{
		Path:    f.Path,
		Content: gzippedString(f.Content),
	})
}

func (s gzippedString) MarshalJSON() ([]byte, error) {
	if len(s) < 3 {
		return json.Marshal(string(s))
	}
	if s[0] != gzipID1 || s[1] != gzipID2 || s[2] != gzipDeflate {
		return json.Marshal(string(s))
	}
	msgReader := bytes.NewReader(s)
	gzipReader, err := gzip.NewReader(msgReader)
	if err != nil {
		return nil, err
	}
	defer gzipReader.Close()
	b, err := ioutil.ReadAll(gzipReader)
	if err != nil {
		return nil, err
	}
	return json.Marshal(string(b))
}
