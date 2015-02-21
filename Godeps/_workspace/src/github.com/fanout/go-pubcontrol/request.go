//    request.go
//    ~~~~~~~~~
//    This module implements the request struct.
//    :authors: Konstantin Bokarius.
//    :copyright: (c) 2015 by Fanout, Inc.
//    :license: MIT, see LICENSE for more details.

package pubcontrol

type request struct {
    Type string
    Uri string
    Auth string
    Export map[string]interface{}
    Callback func(result bool, err error)
}
