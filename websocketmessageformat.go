//    websocketmessageformat.go
//    ~~~~~~~~~
//    This module implements the WebSocketMessageFormat struct.
//    :authors: Konstantin Bokarius.
//    :copyright: (c) 2015 by Fanout, Inc.
//    :license: MIT, see LICENSE for more details.

package gripcontrol

// The WebSocketMessageFormat struct is the format used to publish data to
// WebSocket clients connected to GRIP proxies.
type WebSocketMessageFormat struct {
	Content []byte
	Binary  bool
}

// The name used when publishing this format.
func (format *WebSocketMessageFormat) Name() string {
	return "ws-message"
}

// Exports the message in the required format depending on whether the
// binary field is set to true or false.
func (format *WebSocketMessageFormat) Export() interface{} {
	export := make(map[string]interface{})
	if format.Binary {
		export["content-bin"] = string(format.Content)
	} else {
		export["content"] = string(format.Content)
	}
	return export
}
