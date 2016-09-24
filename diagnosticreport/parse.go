package diagnosticreport

import (
	"github.com/golang/protobuf/proto"
)

// Parse takes a byte array and returns a parsed protobuf class.
func Parse(data []byte) (*DiagnosticReport, error) {
	dr := &DiagnosticReport{}
	if err := proto.Unmarshal(data, dr); err != nil {
		return nil, err
	}
	return dr, nil
}
