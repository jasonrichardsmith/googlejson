// Copyright 2015 Jason Richard Smith.
// Use of this source code is governed by a GPL-3
// license that can be found in the LICENSE file.

// Package googlejson implements json structure defined
// by the Google JSON Style Guide
//
// See the style guide for more details
// https://google-styleguide.googlecode.com/svn/trunk/jsoncstyleguide.xml
package googlejson

import (
	"encoding/json"
	"net/http"
	"strings"
)

/*


 */
type Response struct {
	APIVersion string            `json:"apiVersion"`
	Context    string            `json:"context"`
	ID         string            `json:"id"`
	Method     string            `json:"method"`
	Params     map[string]string `json:"params"`
	Data       `json:"datai, omitempty"`
	Error      `json:"error, omitempty"`
}

func New() *Response {
	r := Response{Params: make(map[string]string), Data: NewData()}
	return &r
}

func NewFromResponse(r http.Response) (*Response, err) {
	res := New()
	defer r.Body.Close()
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return res, err
	}
	err := json.UnMarshall(body, res)
	return res, err
}

func (r *Response) Copy() *Response {
	nr := Response{
		APIVersion: r.APIVersion,
		Method:     r.Method,
		Params:     r.Params,
	}
	return &nr
}

func (r *Response) Write() ([]byte, error) {
	return json.Marshal(r)
}

func (r *Response) WriteToResponse(w http.ResponseWriter) error {
	b, err := r.Write()
	if err != nil {
		return err
	}
	w.Write(b)
	return nil
}

type Data struct {
	Title            string            `json:"title"`
	Description      string            `json:"description"`
	Kind             string            `json:"kind"`
	Fields           string            `json:"fields"`
	Etag             string            `json:"etag"`
	ID               string            `json:"id"`
	Lang             string            `json:"lang"`
	Updated          string            `json:"updated"`
	Deleted          bool              `json:"deleted"`
	CurrentItemCount int               `json:"currentItemCount"`
	ItemsPerPage     int               `json:"itemsPerPage"`
	StartIndex       int               `json:"startIndex"`
	TotalItems       int               `json:"totalItems"`
	PageIndex        int               `json:"pageIndex"`
	TotalPages       int               `json:"totalPages"`
	NextLink         string            `json:"nextLink"`
	PreviousLink     string            `json:"previousLink"`
	SelfLink         string            `json:"selfLink"`
	EditLink         string            `json:"editLink"`
	Items            []json.RawMessage `json:"items"`
	item             int
}

func NewData() *Data {
	d := Data{Items: make([]interface{})}
	return &d
}

func (d *Data) AddField(key string) {
	fs := d.GetFields()
	fs = append(fs, []string{key})
	d.Fields = strings.Join(fs, ",")
}

func (d *Data) AddFields(keys []string) {
	fs := d.GetFields()
	fs = append(fs, keys...)
	d.Fields = strings.Join(fs, ",")
}

func (d *Data) GetFields() []string {
	return strings.Split(d.Fields, ",")
}

func (d *Data) AddItem(i interface{}) error {
	js, err := json.Marshall(i)
	d.Items = append(d.Items, js)
	if err != nil {
		return err
	}
	d.SetItemCount()
	return nil
}

func (d *Data) SetItemCount() {
	d.CurrentItemCount = len(d.Items)
}

func (d *Data) ItemsCount() int {
	return len(d.Items)
}

func (d *Data) CurrentItem(i interface{}) error {
	return json.UnMarshall(d.Items[d.item], i)
}

func (d *Data) NextItem(i interface{}) error {
	count = d.ItemCount()
	if count == d.item+1 {
		return error.New("End of items")
	}
	d.item = d.item + 1
	return d.CurrentItem(i)
}

func (d *Data) ResetItems() {
	d.item = 0
}

type Error struct {
	Code    int         `json:"code"`
	Errors  []ErrorItem `json:"errors"`
	Message string      `json:"message"`
}

func NewError() *Error {
	er := Error{Errors: make([]ErrorItem)}
	return er
}

type ErrorItem struct {
	ExtendedHelper string `json:"extendedHelper"`
	Location       string `json:"location"`
	LocationType   string `json:"locationType"`
	Message        string `json:"message"`
	SendReport     string `json:"sendReport"`
}
