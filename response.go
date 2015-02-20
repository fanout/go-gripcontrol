//    response.go
//    ~~~~~~~~~
//    This module implements the Response struct.
//    :authors: Konstantin Bokarius.
//    :copyright: (c) 2015 by Fanout, Inc.
//    :license: MIT, see LICENSE for more details.

package gripcontrol

type Response struct {
    Code int
    Reason string
    Headers map[string]string
    Body []byte
}
