//    response.go
//    ~~~~~~~~~
//    This module implements the Response struct.
//    :authors: Konstantin Bokarius.
//    :copyright: (c) 2015 by Fanout, Inc.
//    :license: MIT, see LICENSE for more details.

package gripcontrol

// The Response struct is used to represent a set of HTTP response data.
// Populated instances of this struct are serialized to JSON and passed
// to the GRIP proxy in the body. The GRIP proxy then parses the message
// and deserialized the JSON into an HTTP response that is passed back 
// to the client.
type Response struct {
    Code int
    Reason string
    Headers map[string]string
    Body []byte
}
