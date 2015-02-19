//    websocketmessageformat.go
//    ~~~~~~~~~
//    This module implements the WebSocketMessageFormat struct.
//    :authors: Konstantin Bokarius.
//    :copyright: (c) 2015 by Fanout, Inc.
//    :license: MIT, see LICENSE for more details.

package gripcontrol

type WebSocketMessageFormat struct {
    Content []byte
    Binary bool
}

func (format WebSocketMessageFormat) Name() string {
    return "ws-message"
}

func (format WebSocketMessageFormat) Export() interface{} {
    export := make(map[string]interface{})
    if (format.Binary) {
        export["content-bin"] = string(format.Content)
    } else {
        export["content"] = string(format.Content)
    }
    return export
}
