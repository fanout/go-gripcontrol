//    formatter.go
//    ~~~~~~~~~
//    This module implements the Formatter interface.
//    :authors: Konstantin Bokarius.
//    :copyright: (c) 2015 by Fanout, Inc.
//    :license: MIT, see LICENSE for more details.

package pubcontrol

// The Format interface is used for all publishing formats that are
// wrapped in the Item struct. Examples of format implementations
// include JsonObjectFormat and HttpStreamFormat.
type Formatter interface {

    // The name of the format which should return a string. Examples
    // include 'json-object' and 'http-response'
    Name() string

    // The export method which should return a format-specific hash
    // containing the required format-specific data.
    Export() interface{}
}
