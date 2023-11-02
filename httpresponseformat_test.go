//    httpresponseformat_test.go
//    ~~~~~~~~~
//    This module implements the HttpResponseFormat tests.
//    :authors: Konstantin Bokarius.
//    :copyright: (c) 2015 by Fanout, Inc.
//    :license: MIT, see LICENSE for more details.

package gripcontrol

import (
	"encoding/base64"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestHttpResponseFormatName(t *testing.T) {
	fmt := &HttpResponseFormat{}
	assert.Equal(t, fmt.Name(), "http-response")
}

func TestHttpResponseFormatExport(t *testing.T) {
	fmt := &HttpResponseFormat{}
	assert.Equal(t, fmt.Export(), map[string]interface{}{})
	fmt = &HttpResponseFormat{Code: 1, Reason: "reason",
		Headers: map[string]string{"header": "hval"},
		Body:    []byte("body")}
	assert.Equal(t, fmt.Export(), map[string]interface{}{
		"code": 1, "reason": "reason",
		"headers": map[string]string{"header": "hval"},
		"body":    "body"})
	fmt = &HttpResponseFormat{Headers: map[string]string{},
		Body: []byte("\xbd\xb2\x3d\xbc\x20\xe2\x8c\xFF")}
	assert.Equal(t, fmt.Export(), map[string]interface{}{"body-bin": base64.StdEncoding.EncodeToString(
		[]byte("\xbd\xb2\x3d\xbc\x20\xe2\x8c\xFF"))})
}
