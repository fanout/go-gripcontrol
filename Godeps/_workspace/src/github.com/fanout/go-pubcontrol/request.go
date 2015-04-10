//    request.go
//    ~~~~~~~~~
//    This module implements the request struct.
//    :authors: Konstantin Bokarius.
//    :copyright: (c) 2015 by Fanout, Inc.
//    :license: MIT, see LICENSE for more details.

package pubcontrol

// The Request struct represents the parameters required for publishing a
// message. This includes the request type, URI, authorization header,
// exported message data, and callback function.
type request struct {
    Type string
    Uri string
    Auth string
    Export map[string]interface{}
    Callback func(result bool, err error)
}
