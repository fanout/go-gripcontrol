//    item.go
//    ~~~~~~~~~
//    This module implements the Item functionality.
//    :authors: Konstantin Bokarius.
//    :copyright: (c) 2015 by Fanout, Inc.
//    :license: MIT, see LICENSE for more details.

package pubcontrol

type Item struct {
    id string
    prevId string
    formats []Formatter
}

func NewItem(formats []Formatter, id, prevId string) *Item {
    newItem := new(Item)
    newItem.id = id
    newItem.prevId = prevId
    newItem.formats = formats
    return newItem
}

func (item *Item) Export() map[string]interface{} {
    out := make(map[string]interface{})
    if item.id != "" {
        out["id"] = item.id
    }
    if item.prevId != "" {
        out["prev-id"] = item.prevId
    }
    for _,format := range item.formats {
        out[format.Name()] = format.Export()
    }
    return out
}
