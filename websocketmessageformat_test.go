//    websocketmessageformat_test.go
//    ~~~~~~~~~
//    This module implements the WebSocketMessageFormat tests.
//    :authors: Konstantin Bokarius.
//    :copyright: (c) 2015 by Fanout, Inc.
//    :license: MIT, see LICENSE for more details.

package gripcontrol

import ("testing"
        "github.com/stretchr/testify/assert")

func TestWebSocketMessageFormatName(t *testing.T) {
    fmt := &WebSocketMessageFormat{}
    assert.Equal(t, fmt.Name(), "ws-message")
}

func TestWebSocketMessageFormatExport(t *testing.T) {
    fmt := &WebSocketMessageFormat{}
    assert.Equal(t, fmt.Export(), map[string]interface{}{
            "content": ""})
    fmt = &WebSocketMessageFormat{Content: []byte("content")}
    assert.Equal(t, fmt.Export(), map[string]interface{}{
            "content": "content"})
    fmt = &WebSocketMessageFormat{Content: []byte("content"), Binary:true}
    assert.Equal(t, fmt.Export(), map[string]interface{}{
            "content-bin": "content"})
}
