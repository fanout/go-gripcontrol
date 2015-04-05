//    httpstreamformat.go
//    ~~~~~~~~~
//    This module implements the HttpStreamFormat struct.
//    :authors: Konstantin Bokarius.
//    :copyright: (c) 2015 by Fanout, Inc.
//    :license: MIT, see LICENSE for more details.

package gripcontrol

import "encoding/base64"
import "unicode/utf8"

// The HttpStreamFormat struct is the format used to publish messages to
// HTTP stream clients connected to a GRIP proxy.
type HttpStreamFormat struct {
    Content []byte
    Close bool
}

// The name used when publishing this format.
func (format *HttpStreamFormat) Name() string {
    return "http-stream"
}

// Exports the message in the required format depending on whether the
// message content is binary or not, or whether the connection should
// be closed.
func (format *HttpStreamFormat) Export() interface{} {
    export := make(map[string]interface{})
    if (format.Close) {
        export["action"] = "close"
    } else {
        if (format.Content != nil) {
            content := string(format.Content)
            if (utf8.ValidString(content)) {
                export["content"] = content
            } else {
                export["content-bin"] =
                        base64.StdEncoding.EncodeToString(format.Content)
            }
        }
    }
    return export
}
