// Copyright 2015 Jason Richard Smith.
// Use of this source code is governed by a GPL-3
// license that can be found in the LICENSE file.

package googlejson

import (
	"bytes"
	"encoding/json"
	"fmt"
	//"github.com/davecgh/go-spew/spew"
	"io/ioutil"
	"log"
	//"reflect"
	//"net/http/httptest"
	"net/http"
	"os"
	"testing"
)

var sample_json []byte

var sample_struct *Response

var item json.RawMessage

type CarItem struct {
	Color string `json:"color"`
	Type  string `json:"type"`
}

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
	sample_struct = &Response{
		APIVersion: "0.1",
		Context:    "client_context",
		ID:         "1234ABCD",
		Method:     "cars.get",
		Params:     map[string]string{"color": "yellow"},
		Data: Data{
			Kind:             "car",
			Fields:           "color,type",
			Etag:             "08FQn8-eil7ImA9WxZbFEwo",
			ID:               "0000001",
			Lang:             "en",
			Updated:          "2010-01-07T19:58:42.949Z",
			Deleted:          false,
			CurrentItemCount: 1,
			ItemsPerPage:     1,
			StartIndex:       1,
			TotalItems:       10,
			PageIndex:        1,
			TotalPages:       10,
			Items:            make([]json.RawMessage, 0),
			SelfLink:         "https://github.com/jasonrichardsmith/google_json_style/google_json_style.json",
			EditLink:         "https://github.com/jasonrichardsmith/google_json_style/google_json_style.json?edit",
			NextLink:         "https://github.com/jasonrichardsmith/google_json_style/google_json_style.json?next",
			PreviousLink:     "https://github.com/jasonrichardsmith/google_json_style/google_json_style.json?prev",
		},
		Error: Error{
			Code:    404,
			Message: "Car Not Found",
			Errors: []ErrorItem{
				ErrorItem{
					Message:        "Car Not Found",
					ExtendedHelper: "http://url.to.more.details.example.com/",
					SendReport:     "http://report.example.com/",
				},
			},
		},
	}
	i, err := json.Marshal(CarItem{"yellow", "sedan"})
	if err != nil {
		log.Fatal(err)
	}
	sample_struct.Data.Items = append(sample_struct.Data.Items, i)

}

func TestNewFromResponse(t *testing.T) {
	f, err := os.Open("googlejson.json")
	if err != nil {
		log.Fatal(err)
	}
	r := http.Response{Body: f}
	var response *Response
	response, err = NewFromResponse(r)
	if err != nil {
		log.Fatal(err)
	}
	buf := new(bytes.Buffer)
	if err := json.Compact(buf, response.Data.Items[0]); err != nil {
		log.Fatal(err)
	}
	response.Data.Items[0] = buf.Bytes()
	if response.APIVersion != sample_struct.APIVersion {
		t.Error("Test failed")
	}
}

func TestWrite(t *testing.T) {
	b, _ := sample_struct.Write()
	if string(b) != string(sample_json) {
		t.Error("Test failed")
	}
}
