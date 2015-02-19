//    httpstreamformat.go
//    ~~~~~~~~~
//    This module implements the HttpStreamFormat struct.
//    :authors: Konstantin Bokarius.
//    :copyright: (c) 2015 by Fanout, Inc.
//    :license: MIT, see LICENSE for more details.

package gripcontrol

import "encoding/base64"
import "unicode/utf8"

type HttpStreamFormat struct {
    Content []byte
    Close bool
}

func (format HttpStreamFormat) Name() string {
    return "http-stream"
}

func (format HttpStreamFormat) Export() interface{} {
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
