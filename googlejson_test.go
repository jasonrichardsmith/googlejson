// Copyright 2015 Jason Richard Smith.
// Use of this source code is governed by a GPL-3
// license that can be found in the LICENSE file.

package googlejson

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http/httptest"
	"testing"
)

var sample_json []byte

var sample_struct *Response

func init() {
	var err error
	sample_json, err = ioutil.ReadFile("googlejson.json")
	if err != nil {
		log.Fatal(err)
	}
	buf := new(bytes.Buffer)
	if err := json.Compact(buf, sample_json); err != nil {
		log.Fatal(err)
	}
	sample_json = buf.Bytes()
	sample_struct = New()
}

func TestNewFromResponse(t *testing.T) {
	fmt.Println(sample_struct)
	t.Error("Test failed")
}
