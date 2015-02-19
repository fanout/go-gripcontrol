//    httpresponseformat.go
//    ~~~~~~~~~
//    This module implements the HttpResponseFormat struct.
//    :authors: Konstantin Bokarius.
//    :copyright: (c) 2015 by Fanout, Inc.
//    :license: MIT, see LICENSE for more details.

package gripcontrol

import "encoding/base64"
import "unicode/utf8"

type HttpResponseFormat struct {
    Code string
    Reason string
    Headers map[string]string
    Body []byte
}

func (format *HttpResponseFormat) Name() string {
    return "http-response"
}

func (format *HttpResponseFormat) Export() interface{} {
    export := make(map[string]interface{})
    if (format.Code != "") {
        export["code"] = format.Code
    }
    if (format.Reason != "") {
        export["reason"] = format.Reason
    }
    if (format.Headers != nil && len(format.Headers) > 0) {
        export["headers"] = format.Headers
    }
    if (format.Body != nil) {
        body := string(format.Body)
        if (utf8.ValidString(body)) {
            export["body"] = body
        } else {
            export["body-bin"] =
                    base64.StdEncoding.EncodeToString(format.Body)
        }
    }
    return export
}
