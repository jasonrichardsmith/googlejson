// Copyright 2015 Jason Richard Smith.
// Use of this source code is governed by a GPL-3
// license that can be found in the LICENSE file.

// Package googlejson implements json structure defined
// by the Google JSON Style Guide
//
// See the style guide for more details
// https://google-styleguide.googlecode.com/svn/trunk/jsoncstyleguide.xml

/*
Response is a structure designed to hold the contents of a json response to be
sent or received with the http package.

You can easily utilize this to marshal results from an http.Reponse

	resp, err := http.Get("http://example.com/")
	if err != nil {
		log.Fatal(err)
	}
	gresp, err := googlejson.NewFromHTTPResponse(res)

You can also easily write to a http.ResponseWriter

	func MyHandle(w http.ResponseWriter, r *http.Request) {
		gresp := googlejson.New()
		gresp.APIVersion = "1.2"
		code, err:= gresp.WriteToHTTPResponse(w)
	}

Or just to a byte slice

	gresp := googlejson.New()
	gresp.APIVersion = "1.2"
	b := gresp.Write()


All the data for the API is stored in Response.Data.Items.  These will always be stored as []json.RawMessage
which can be retrieved or set with the AddItem, CurrentItem and NextItem methods.
*/
package googlejson
