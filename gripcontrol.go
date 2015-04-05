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
import "github.com/dgrijalva/jwt-go"
import "net/url"
import "encoding/base64"

// The GripControl struct provides functionality that is used in conjunction
// with GRIP proxies. This includes facilitating the creation of hold
// instructions for HTTP long-polling and HTTP streaming, parsing GRIP URIs
// into config objects, validating the GRIP-SIG header coming from GRIP
// proxies, creating GRIP channel headers, and also WebSocket-over-HTTP
// features such as encoding/decoding web socket events and generating
// control messages.

// A convenience method for creating GRIP hold response instructions for HTTP
// long-polling. This method simply passes the specified parameters to the
// create_hold method with 'response' as the hold mode.
func CreateHoldResponse(channels []*Channel, response interface{},
        timeout *int) (string, error) {
    return CreateHold("response", channels, response, timeout)
}

// A convenience method for creating GRIP hold stream instructions for HTTP
// streaming. This method simply passes the specified parameters to the
// create_hold method with 'stream' as the hold mode.
func CreateHoldStream(channels []*Channel,
        response interface{}) (string, error) {
    return CreateHold("stream", channels, response, nil)
}

// Create GRIP hold instructions for the specified mode, channels, response
// and optional timeout value. The channel parameter can be specified as
// either a string representing the channel name, a Channel instance or an
// array of Channel instances. The response parameter can be specified as
// either a string representing the response body or a Response instance.
func CreateHold(mode string, channels []*Channel, response interface{},
        timeout *int) (string, error) {
    hold := make(map[string]interface{})
    hold["mode"] = mode
    ichannels := make([]map[string]string, 0)
    for _, channel := range channels {
        ichannel := make(map[string]string)
        ichannel["name"] = channel.Name
        if channel.PrevId != "" {
            ichannel["prev-id"] = channel.PrevId
        }
        ichannels = append(ichannels, ichannel)
    }
    hold["channels"] = ichannels
    if timeout != nil {
        hold["timeout"] = timeout
    }
    iresponse := make(map[string]interface{})
    if response != nil {
        var processedResponse *Response
        switch response.(type) {
            case *Response:
                processedResponse = response.(*Response)
            case string:
                processedResponse = &Response{Body: []byte(response.(string))}
            case []byte:
                processedResponse = &Response{Body: response.([]byte)}
            default:
                return "", &GripFormatError{err: "response must be of type " + 
                        "*Response, []byte, or string"}        
        }
        if processedResponse.Code > 0 {
            iresponse["code"] = processedResponse.Code
        }
        if processedResponse.Reason != "" {
            iresponse["reason"] = processedResponse.Reason
        }
        if (processedResponse.Headers != nil &&
                len(processedResponse.Headers) > 0) {
            iresponse["headers"] = processedResponse.Headers
        }
        if processedResponse.Body != nil && len(processedResponse.Body) > 0 {
            body := string(processedResponse.Body)
            if utf8.ValidString(body) {
                iresponse["body"] = body
            } else {
                iresponse["body-bin"] =
                        base64.StdEncoding.EncodeToString(processedResponse.Body)
            }
        }
    }
    instruct := make(map[string]interface{})
    instruct["hold"] = hold
    if len(iresponse) > 0 {
        instruct["response"] = iresponse
    }
    message, err := json.Marshal(instruct)
    if err != nil {
        return "", err
    }
    return string(message), nil
}

// Parse the specified GRIP URI into a config object that can then be passed
// to the GripPubControl struct. The URI can include 'iss' and 'key' JWT
// authentication query parameters as well as any other required query string
// parameters. The JWT 'key' query parameter can be provided as-is or in base64
// encoded format.
func ParseGripUri(rawUri string) (map[string]interface{}, error) {
    uri, err := url.Parse(rawUri)
    if err != nil {
        return nil, err
    }
    params := uri.Query()
    iss := ""
    key := ""
    if _, ok := params["iss"]; ok {
        iss = params["iss"][0]
        delete(params, "iss")
    }
    if _, ok := params["key"]; ok {
        key = params["key"][0]
        delete(params, "key")
    }
    decodedKey := make([]byte, 0)
    if key != "" && key[:7] == "base64:" {
        var err error
        decodedKey, err = base64.StdEncoding.DecodeString(key[7:])
        if err != nil {
            return nil, err
        }
    }
    qs := params.Encode()
    path := uri.Path
    if path[len(path) - 1:] == "/" {
        path = path[:len(path) - 1]
    }
    controlUri := uri.Scheme + "://" + uri.Host + path
    if len(qs) > 0 {
        controlUri += "?" + qs
    }
    out := make(map[string]interface{})
    out["control_uri"] = controlUri
    if iss != "" {
        out["control_iss"] = iss
    }
    if len(decodedKey) > 0 {
        out["key"] = decodedKey
    }
    return out, nil
}

// Validate the specified JWT token and key. This method is used to validate
// the GRIP-SIG header coming from GRIP proxies such as Pushpin or Fanout.io.
// Note that the token expiration is also verified.
func ValidateSig(token, key string) bool {
    parsedToken, err := jwt.Parse(token, func(t *jwt.Token) (interface{},
            error) { return []byte(key), nil })
    if err == nil && parsedToken.Valid {
        return true;
    }
    return false
}

// Create a GRIP channel header for the specified channels. The channels
// parameter can be specified as a string representing the channel name,
// a Channel instance, or an array of Channel instances. The returned GRIP
// channel header is used when sending instructions to GRIP proxies via
// HTTP headers.
func CreateGripChannelHeader(channels []*Channel) string {
    var parts []string
    for _, channel := range channels {   
        s := channel.Name
        if channel.PrevId != "" {
            s += "; prev-id=" + channel.PrevId
        }
        parts = append(parts, s)
    }
    return strings.Join(parts, ", ")
}

// Decode the specified HTTP request body into an array of WebSocketEvent
// instances when using the WebSocket-over-HTTP protocol. A RuntimeError
// is raised if the format is invalid.
func DecodeWebSocketEvents(body string) ([]*WebSocketEvent, error) {
    out := make([]*WebSocketEvent, 0)
    for start := 0; start < utf8.RuneCountInString(body); {
        partialBody := body[start:]
        if partialBody == "\r\n" {
            break
        }
        at := strings.Index(partialBody, "\r\n")
        if at == -1 {
            return nil, &GripFormatError{err: "bad format"} 
        }
        start += at + 2
        contentStart := at + 2
        typeline := partialBody[0:at]
        at = strings.Index(typeline, " ")
        var event *WebSocketEvent
        if at != -1 {
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

// Encode the specified array of WebSocketEvent instances. The returned string
// value should then be passed to a GRIP proxy in the body of an HTTP response
// when using the WebSocket-over-HTTP protocol.
func EncodeWebSocketEvents(events []*WebSocketEvent) string {
    out := ""
    for _, event := range events {
        if event.Content != "" {
            out += fmt.Sprintf("%s %x\r\n%s\r\n", event.Type, 
                    len(event.Content), event.Content)
        } else {
            out += fmt.Sprintf("%s\r\n", event.Type)
        }
    }
    return out
}

// Generate a WebSocket control message with the specified type and optional
// arguments. WebSocket control messages are passed to GRIP proxies and
// example usage includes subscribing/unsubscribing a WebSocket connection
// to/from a channel.
func WebSocketControlMessage(messageType string,
        args map[string]interface{}) (string, error) {
    out := make(map[string]interface{})
    if args != nil {
        for key, value := range args {
            out[key] = value
        }
    }
    out["type"] = messageType
    message, err := json.Marshal(out)
    return string(message), err
}

// An error object used to represent a GRIP formatting error.
type GripFormatError struct {
    err string
}

// The function used to retrieve the message associated with
// a GripFormatError.
func (e GripFormatError) Error() string {
    return e.err
}
