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

// Type is the top level json object.
// It should contain at least one data object or one error.
type Response struct {
	// Version of API being served or received.
	APIVersion string `json:"apiVersion"`
	// Context is a parameter submitted by requestor
	// as in an http.Request, this will be returned
	// to the client for context.
	Context string `json:"context"`
	// ID is a unique ID assigned to the request
	// if the API will need to reference a transaction.
	ID string `json:"id"`
	// Method represents the operation performed.
	Method string `json:"method"`
	// Params are a list of parameters submitted to the API.
	Params map[string]string `json:"params"`
	// Data holds the actual data that was returned.
	Data `json:"datai, omitempty"`
	// Errors to be returned.
	Error `json:"error, omitempty"`
}

// Shortcut to create a new Response
func New() *Response {
	r := Response{Params: make(map[string]string), Data: NewData()}
	return &r
}

// Shortcut to create a response from an http.Response
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

// Shortcut to create a copy of the response.  This
// means you can set common settings, and resuse the base Response.
func (r *Response) Copy() *Response {
	nr := Response{
		APIVersion: r.APIVersion,
		Method:     r.Method,
		Params:     r.Params,
	}
	return &nr
}

// Write the struct to byte[]
func (r *Response) Write() ([]byte, error) {
	return json.Marshal(r)
}

// Shortcut to write to an http.ResponseWriter.
func (r *Response) WriteToResponse(w http.ResponseWriter) error {
	b, err := r.Write()
	if err != nil {
		return err
	}
	w.Write(b)
	return nil
}

// Data structure holds all information specific to the data in the
// Response and the data itself
type Data struct {
	// Kind is a name of the entity being returned, such as
	// cars, orders, customers etc.
	Kind string `json:"kind"`
	// A list of field being returned.  This is a comma separated
	// string so helper methods are provided below.
	Fields string `json:"fields"`
	// Etags are an ID for the version of the data you are viewing
	// this allows to identify expired data.
	// More can be read here:
	// https://developers.google.com/gdata/docs/2.0/reference?csw=1#ResourceVersioning
	Etag string `json:"etag"`
	// The ID for this request
	ID string `json:"id"`
	// Language.
	Lang string `json:"lang"`
	// Last time record was updated
	Updated string `json:"updated"`
	// Deleted, if this was a delete request did it occur.
	Deleted bool `json:"deleted"`
	// How many items have returned this request.
	CurrentItemCount int `json:"currentItemCount"`
	// How many items could be returned in each request.
	ItemsPerPage int `json:"itemsPerPage"`
	// Where the list begins in the full list of return values.
	StartIndex int `json:"startIndex"`
	// Total number of Items matching request.
	TotalItems int `json:"totalItems"`
	// Current page.
	PageIndex int `json:"pageIndex"`
	// Total number of pages.
	TotalPages int `json:"totalPages"`
	// Link to next result or result set for paginated results.
	NextLink string `json:"nextLink"`
	// Link to previous result or result set for paginated results.
	PreviousLink string `json:"previousLink"`
	// Direct link to current data set.
	SelfLink string `json:"selfLink"`
	// Link to edit results.
	EditLink string `json:"editLink"`
	// An array of items -  this is the actual data.
	Items []json.RawMessage `json:"items"`

	// pointer to current item.
	item int
}

// Shortcut to new data object.
func NewData() *Data {
	d := Data{Items: make([]interface{})}
	return &d
}

// Add a single field to the list of fields to be returned.
func (d *Data) AddField(key string) {
	fs := d.GetFields()
	fs = append(fs, []string{key})
	d.Fields = strings.Join(fs, ",")
}

// Add a list of fields to the list of fields to be returned.
func (d *Data) AddFields(keys []string) {
	fs := d.GetFields()
	fs = append(fs, keys...)
	d.Fields = strings.Join(fs, ",")
}

// Get a list of fields.
func (d *Data) GetFields() []string {
	return strings.Split(d.Fields, ",")
}

// Add a single data item to the list of Items to be returned.
func (d *Data) AddItem(i interface{}) error {
	js, err := json.Marshall(i)
	d.Items = append(d.Items, js)
	if err != nil {
		return err
	}
	d.SetItemCount()
	return nil
}

// Set the item count to the number of Items to be returned.
func (d *Data) SetItemCount() {
	d.CurrentItemCount = len(d.Items)
}

// Get count of current items in the Items list.
func (d *Data) ItemsCount() int {
	return len(d.Items)
}

// Retrieve the item at current pointer position.
func (d *Data) CurrentItem(i interface{}) error {
	return json.UnMarshall(d.Items[d.item], i)
}

// Retrieve the next title.
func (d *Data) NextItem(i interface{}) error {
	count = d.ItemCount()
	if count == d.item+1 {
		return error.New("End of items")
	}
	d.item = d.item + 1
	return d.CurrentItem(i)
}

// Reset Item pointer to 0.
func (d *Data) ResetItems() {
	d.item = 0
}

// Error object to be returned.
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
