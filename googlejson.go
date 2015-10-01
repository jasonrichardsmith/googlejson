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
	"errors"
	"io/ioutil"
	"net/http"
	"strings"
)

// Type is the top level json object.
// It should contain at least one data object or one error.
type Response struct {
	// Version of API being served or received.
	APIVersion string `json:"apiVersion,omitempty"`

	// Context is a parameter submitted by requestor
	// as in an http.Request, this will be returned
	// to the client for context.
	Context string `json:"context,omitempty"`

	// ID is a unique ID assigned to the request
	// if the API will need to reference a transaction.
	ID string `json:"id,omitempty"`

	// Method represents the operation performed.
	Method string `json:"method,omitempty"`

	// Params are a list of parameters submitted to the API.
	Params map[string]string `json:"params,omitempty"`

	// Data holds the actual data that was returned.
	Data `json:"data,omitempty"`

	// Errors to be returned.
	Error `json:"error,omitempty"`
}

// Shortcut to create a new Response
func New() *Response {
	r := Response{Params: make(map[string]string), Data: *NewData()}
	return &r
}

// Shortcut to create a response from an http.Response
func NewFromHTTPResponse(r http.Response) (*Response, error) {
	res := New()
	defer r.Body.Close()
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return res, err
	}
	err = json.Unmarshal(body, res)
	return res, err
}

// Shortcut to create a copy of the response.  This
// means you can set common settings, and re-use the base Response.
func (r *Response) Copy() *Response {
	nr := Response{
		APIVersion: r.APIVersion,
		Method:     r.Method,
		Params:     r.Params,
	}
	return &nr
}

// Write the struct to byte[].
// SetItemCount will be called prior to writing.
func (r *Response) Write() ([]byte, error) {
	r.Data.SetItemCount()
	return json.Marshal(r)
}

// Shortcut to write to an http.ResponseWriter.
func (r *Response) WriteToHTTPResponse(w http.ResponseWriter) error {
	b, err := r.Write()
	if err != nil {
		return err
	}
	return w.Write(b)
}

// Data structure holds all information specific to the data in the
// Response and the data itself
type Data struct {
	// Kind is a name of the entity being returned, such as
	// cars, orders, customers etc.
	Kind string `json:"kind,omitempty"`

	// A list of field being returned.  This is a comma separated
	// string so helper methods are provided below.
	Fields string `json:"fields,omitempty"`

	// Etags are an ID for the version of the data you are viewing
	// this allows to identify expired data.
	// More can be read here:
	// https://developers.google.com/gdata/docs/2.0/reference?csw=1#ResourceVersioning
	Etag string `json:"etag,omitempty"`

	// The ID for this request
	ID string `json:"id,omitempty"`

	// Language.
	Lang string `json:"lang,omitempty"`

	// Last time record was updated
	Updated string `json:"updated,omitempty"`

	// Deleted, if this was a delete request did it occur.
	Deleted bool `json:"deleted"`

	// How many items have returned this request.
	CurrentItemCount int `json:"currentItemCount,omitempty"`

	// How many items could be returned in each request.
	ItemsPerPage int `json:"itemsPerPage,omitempty"`

	// Where the list begins in the full list of return values.
	StartIndex int `json:"startIndex,omitempty"`

	// Total number of Items matching request.
	TotalItems int `json:"totalItems,omitempty"`

	// Current page.
	PageIndex int `json:"pageIndex,omitempty"`

	// Total number of pages.
	TotalPages int `json:"totalPages,omitempty"`

	// Direct link to current data set.
	SelfLink string `json:"selfLink,omitempty"`

	// Link to edit results.
	EditLink string `json:"editLink,omitempty"`

	// Link to next result or result set for paginated results.
	NextLink string `json:"nextLink,omitempty"`

	// Link to previous result or result set for paginated results.
	PreviousLink string `json:"previousLink,omitempty"`

	// An array of items -  this is the actual data.
	Items []json.RawMessage `json:"items,omitempty"`

	// pointer to current item.
	item int
}

// Shortcut to new data object.
func NewData() *Data {
	d := Data{Items: make([]json.RawMessage, 0)}
	return &d
}

// Add a single field to the list of fields to be returned.
func (d *Data) AddField(keys ...string) {
	fs := d.GetFields()
	for _, key := range keys {
		fs = append(fs, key)
	}
	d.Fields = strings.Join(fs, ",")
}

// Get a list of fields.
func (d *Data) GetFields() []string {
	if d.Fields == "" {
		return make([]string, 0)
	}
	return strings.Split(d.Fields, ",")
}

// Add a single data item to the list of Items to be returned.
func (d *Data) AddItem(i interface{}) error {
	js, err := json.Marshal(i)
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
	return json.Unmarshal(d.Items[d.item], i)
}

// Retrieve the next item.
func (d *Data) NextItem(i interface{}) error {
	count := d.ItemsCount()
	if count == d.item+1 {
		return errors.New("End of items")
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
	// Integer code representing an error code
	Code int `json:"code,omitempty"`

	// Message for error.
	Message string `json:"message,omitempty"`

	// An array of error message item data.
	Errors []ErrorItem `json:"errors,omitempty"`
}

// Shortcut to create Error.
func NewError() *Error {
	er := Error{Errors: make([]ErrorItem, 0)}
	return &er
}

// Details relating to the error being returned.
// For more details see
// https://google-styleguide.googlecode.com/svn/trunk/jsoncstyleguide.xml
type ErrorItem struct {
	Message        string `json:"message,omitempty"`
	Location       string `json:"location,omitempty"`
	LocationType   string `json:"locationType,omitempty"`
	ExtendedHelper string `json:"extendedHelper,omitempty"`
	Domains        string `json:"domain,omitempty"`
	Reason         string `json:"reason,omitempty"`
	SendReport     string `json:"sendReport,omitempty"`
}
