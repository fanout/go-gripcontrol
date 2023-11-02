//    websocketevent.go
//    ~~~~~~~~~
//    This module implements the WebSocketEvent struct.
//    :authors: Konstantin Bokarius.
//    :copyright: (c) 2015 by Fanout, Inc.
//    :license: MIT, see LICENSE for more details.

package gripcontrol

// The WebSocketEvent struct represents WebSocket event information that is
// used with the GRIP WebSocket-over-HTTP protocol. It includes information
// about the type of event as well as an optional content field.
type WebSocketEvent struct {
	Type    string
	Content string
}
