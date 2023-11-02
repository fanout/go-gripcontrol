//    httpstreamformat_test.go
//    ~~~~~~~~~
//    This module implements the HttpStreamFormat tests.
//    :authors: Konstantin Bokarius.
//    :copyright: (c) 2015 by Fanout, Inc.
//    :license: MIT, see LICENSE for more details.

package gripcontrol

import (
	"encoding/base64"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestHttpStreamFormatName(t *testing.T) {
	fmt := &HttpStreamFormat{}
	assert.Equal(t, fmt.Name(), "http-stream")
}

func TestHttpStreamFormatExport(t *testing.T) {
	fmt := &HttpStreamFormat{}
	assert.Equal(t, fmt.Export(), map[string]interface{}{})
	fmt = &HttpStreamFormat{Content: []byte("content")}
	assert.Equal(t, fmt.Export(), map[string]interface{}{
		"content": "content"})
	fmt = &HttpStreamFormat{Content: []byte("content"), Close: true}
	assert.Equal(t, fmt.Export(), map[string]interface{}{
		"action": "close"})
	fmt = &HttpStreamFormat{Close: true}
	assert.Equal(t, fmt.Export(), map[string]interface{}{
		"action": "close"})
	fmt = &HttpStreamFormat{
		Content: []byte("\xbd\xb2\x3d\xbc\x20\xe2\x8c\xFF")}
	assert.Equal(t, fmt.Export(), map[string]interface{}{"content-bin": base64.StdEncoding.EncodeToString(
		[]byte("\xbd\xb2\x3d\xbc\x20\xe2\x8c\xFF"))})
}
