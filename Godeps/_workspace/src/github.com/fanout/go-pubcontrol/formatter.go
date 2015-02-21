//    formatter.go
//    ~~~~~~~~~
//    This module implements the Formatter interface.
//    :authors: Konstantin Bokarius.
//    :copyright: (c) 2015 by Fanout, Inc.
//    :license: MIT, see LICENSE for more details.

package pubcontrol

type Formatter interface {
    Name() string
    Export() interface{}
}
