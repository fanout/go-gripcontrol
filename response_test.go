//    response_test.go
//    ~~~~~~~~~
//    This module implements the Response tests.
//    :authors: Konstantin Bokarius.
//    :copyright: (c) 2015 by Fanout, Inc.
//    :license: MIT, see LICENSE for more details.

package gripcontrol

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestResponse(t *testing.T) {
	response := &Response{Code: 1, Reason: "reason",
		Headers: map[string]string{"header": "hval"},
		Body:    []byte("body")}
	assert.Equal(t, response.Code, 1)
	assert.Equal(t, response.Reason, "reason")
	assert.Equal(t, response.Headers, map[string]string{"header": "hval"})
	assert.Equal(t, response.Body, []byte("body"))
}
