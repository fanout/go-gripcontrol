//    item_test.go
//    ~~~~~~~~~
//    This module implements the Item tests.
//    :authors: Konstantin Bokarius.
//    :copyright: (c) 2015 by Fanout, Inc.
//    :license: MIT, see LICENSE for more details.

package pubcontrol

import ("testing"
        "github.com/stretchr/testify/assert")

type FormatTestStruct1 struct {
    value string
}
func (format *FormatTestStruct1) Name() string {
    return "test-format"
}
func (format *FormatTestStruct1) Export() interface{} {
    return format.value
}

type FormatTestStruct2 struct {
    value string
}
func (format *FormatTestStruct2) Name() string {
    return "test-format2"
}
func (format *FormatTestStruct2) Export() interface{} {
    return format.value
}

var fmt1a Formatter = &FormatTestStruct1{value:"value1a"}
var fmt1b Formatter = &FormatTestStruct1{value:"value1b"}
var fmt2 Formatter = &FormatTestStruct2{value:"value2"}

func TestItemInitialize(t *testing.T) {
    formats := make([]Formatter, 0);
    formats = append(formats, fmt1a);
    item := NewItem(formats, "id", "prev-id");
    assert.Equal(t, item.id, "id");
    assert.Equal(t, item.prevId, "prev-id");
    assert.Equal(t, item.formats[0], fmt1a);
}

func TestItemExport(t *testing.T) {
    formats := make([]Formatter, 0);
    formats = append(formats, fmt1a);
    formats = append(formats, fmt2);
    item := NewItem(formats, "id", "prev-id");
    export, err := item.Export();
    assert.Nil(t, err);
    assert.Equal(t, export["id"], "id");
    assert.Equal(t, export["prev-id"], "prev-id");
    assert.Equal(t, export["test-format"], "value1a");
    assert.Equal(t, export["test-format2"], "value2");
    formats = make([]Formatter, 0);
    formats = append(formats, fmt1b);
    item = NewItem(formats, "", "");
    export, err = item.Export();
    assert.Nil(t, err);
    assert.Equal(t, export["test-format"], "value1b");
    if _, ok := export["id"]; ok {
        t.Log("id not set")
        t.FailNow()
    }
    if _, ok := export["prev-id"]; ok {
        t.Log("prev-id not set")
        t.FailNow()
    }
}

func TestItemExportError(t *testing.T) {
    formats := make([]Formatter, 0);
    formats = append(formats, fmt1a);
    formats = append(formats, fmt1b);
    item := NewItem(formats, "id", "prev-id");
    export, err := item.Export();
    assert.Nil(t, export);
    assert.NotNil(t, err);
}
