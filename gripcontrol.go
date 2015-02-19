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

func EncodeWebsocketEvents(events []*WebSocketEvent) string {
    out := ""
    for _, event := range events {
        if (event.Content != "") {
            out += fmt.Sprintf("%s %02x\r\n%s\r\n", event.Type, 
                    utf8.RuneCountInString(event.Content), event.Content)
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
