//    websocketevent_test.go
//    ~~~~~~~~~
//    This module implements the WebSocketEvent tests.
//    :authors: Konstantin Bokarius.
//    :copyright: (c) 2015 by Fanout, Inc.
//    :license: MIT, see LICENSE for more details.

package gripcontrol

import ("testing"
        "github.com/stretchr/testify/assert")

func TestWebSocketEvent(t *testing.T) {
    we := &WebSocketEvent{Type:"type", Content:"content"}
    assert.Equal(t, we.Type, "type")
    assert.Equal(t, we.Content, "content")
}
