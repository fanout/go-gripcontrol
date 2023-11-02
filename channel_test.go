//    channel_test.go
//    ~~~~~~~~~
//    This module implements the Channel tests.
//    :authors: Konstantin Bokarius.
//    :copyright: (c) 2015 by Fanout, Inc.
//    :license: MIT, see LICENSE for more details.

package gripcontrol

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestChannel(t *testing.T) {
	ch := &Channel{Name: "name", PrevId: "prev-id"}
	assert.Equal(t, ch.Name, "name")
	assert.Equal(t, ch.PrevId, "prev-id")
}
