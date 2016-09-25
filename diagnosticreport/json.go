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

func (dr DiagnosticReport) MarshalJSON() ([]byte, error) {
	return json.Marshal(&struct {
		Version          string           `json:"version,omitempty"`
		Files            []*File          `json:"files,omitempty"`
		StormVersion     string           `json:"stormVersion,omitempty"`
		WhirlwindVersion string           `json:"whirlwindVersion,omitempty"`
		NetworkConfig    string           `json:"networkConfig,omitempty"`
		FileLengths      []*FileLength    `json:"fileLengths,omitempty"`
		WanInfo          string           `json:"wanInfo,omitempty"`
		CommandOutputs   []*CommandOutput `json:"commandOutputs,omitempty"`
		InfoJSON         json.RawMessage  `json:"infoJSON,omitempty"`
		Unknown1         int32            `json:"unknown1,omitempty"`
		UnknownPairs     []*UnknownPair   `json:"unknownPairs,omitempty"`
		UnixTime         int32            `json:"unixTime,omitempty"`
	}{
		Version:          dr.Version,
		Files:            dr.Files,
		StormVersion:     dr.StormVersion,
		WhirlwindVersion: dr.WhirlwindVersion,
		NetworkConfig:    dr.NetworkConfig,
		FileLengths:      dr.FileLengths,
		WanInfo:          dr.WanInfo,
		CommandOutputs:   dr.CommandOutputs,
		InfoJSON:         json.RawMessage(dr.InfoJSON),
		Unknown1:         dr.Unknown1,
		UnknownPairs:     dr.UnknownPairs,
		UnixTime:         dr.UnixTime,
	})
}
