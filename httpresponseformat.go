//    httpresponseformat.go
//    ~~~~~~~~~
//    This module implements the HttpResponseFormat struct.
//    :authors: Konstantin Bokarius.
//    :copyright: (c) 2015 by Fanout, Inc.
//    :license: MIT, see LICENSE for more details.

package gripcontrol

import "encoding/base64"
import "unicode/utf8"

// The HttpResponseFormat struct is the format used to publish messages to
// HTTP response clients connected to a GRIP proxy.
type HttpResponseFormat struct {
    Code int
    Reason string
    Headers map[string]string
    Body []byte
}

// The name used when publishing this format.
func (format *HttpResponseFormat) Name() string {
    return "http-response"
}

// Export the message into the required format and include only the fields
// that are set. The body is exported as base64 as 'body-bin' (as opposed
// to 'body') if the value is a buffer.
func (format *HttpResponseFormat) Export() interface{} {
    export := make(map[string]interface{})
    if (format.Code > 0) {
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
