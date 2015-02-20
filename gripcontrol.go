//    gripcontrol.go
//    ~~~~~~~~~
//    This module implements the GripControl features.
//    :authors: Konstantin Bokarius.
//    :copyright: (c) 2015 by Fanout, Inc.
//    :license: MIT, see LICENSE for more details.

package gripcontrol

import "fmt"
import "unicode/utf8"
import "encoding/json"
import "strings"

func CreateGripChannelHeader(channels []*Channel) string {
    var parts []string
    for _, channel := range channels {   
        s := channel.Name
        if (channel.PrevId != "") {
            s += "; prev-id=" + channel.PrevId
        }
        parts = append(parts, s)
    }
    return strings.Join(parts, ", ")
}

func DecodeWebsocketEvents(body string) ([]*WebSocketEvent, error) {
    out := make([]*WebSocketEvent, 0)
    for start := 0; start < utf8.RuneCountInString(body); {
        partialBody := body[start:]
        if (partialBody == "\r\n") {
            break
        }
        at := strings.Index(partialBody, "\r\n")
        if (at == -1) {
            return nil, &GripFormatError{err: "bad format"} 
        }
        start += at + 2
        contentStart := at + 2
        typeline := partialBody[0:at]
        at = strings.Index(typeline, " ")
        var event *WebSocketEvent
        if (at != -1) {
            etype := typeline[:at]
            var clen int
            fmt.Sscanf(typeline[at + 1:], "%x", &clen)
            content := partialBody[contentStart:contentStart + clen]           
            start += clen + 2
            event = &WebSocketEvent{Type:etype, Content:content}
        } else {
            event = &WebSocketEvent{Type:typeline}
        }
        out = append(out, event)
    }
    return out, nil
}

func EncodeWebsocketEvents(events []*WebSocketEvent) string {
    out := ""
    for _, event := range events {
        if (event.Content != "") {
            out += fmt.Sprintf("%s %x\r\n%s\r\n", event.Type, 
                    len(event.Content), event.Content)
        } else {
            out += fmt.Sprintf("%s\r\n", event.Type)
        }
    }
    return out
}

func WebsocketControlMessage(messageType string,
        args map[string]interface{}) (string, error) {
    out := make(map[string]interface{})
    if (args != nil) {
        for key, value := range args {
            out[key] = value
        }
    }
    out["type"] = messageType
    message, err := json.Marshal(out)
    return string(message), err
}

type GripFormatError struct {
    err string
}

func (e GripFormatError) Error() string {
    return e.err
}
