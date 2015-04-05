//    channel.go
//    ~~~~~~~~~
//    This module implements the Channel struct.
//    :authors: Konstantin Bokarius.
//    :copyright: (c) 2015 by Fanout, Inc.
//    :license: MIT, see LICENSE for more details.

package gripcontrol

// The Channel class is used to represent a channel in for a GRIP proxy and
// tracks the previous ID of the last message.
type Channel struct {
    Name string
    PrevId string
}
