go-gripcontrol
================

Author: Konstantin Bokarius <kon@fanout.io>

A GRIP library for Go.

License
-------

go-gripcontrol is offered under the MIT license. See the LICENSE file.

Installation
------------

```sh
go get github.com/fanout/go-gripcontrol
```

go-gripcontrol requires jwt-go 2.2.0 and go-pubcontrol 1.0.1. To ensure that the correct version of both of these dependencies are installed use godeps:

```sh
go get github.com/tools/godep
cd $GOPATH/src/github.com/fanout/go-gripcontrol
$GOPATH/bin/godep restore
```

Usage
-----

Examples for how to publish HTTP response and HTTP stream messages to GRIP proxy endpoints via the GripPubControl class.

```go
package main

import "github.com/fanout/go-pubcontrol"
import "github.com/fanout/go-gripcontrol"
import "encoding/base64"

func main() {
    // GripPubControl can be initialized with or without an endpoint configuration.
    // Each endpoint can include optional JWT authentication info.
    // Multiple endpoints can be included in a single configuration.

    // Initialize GripPubControl with a single endpoint:
    decodedKey, err := base64.StdEncoding.DecodeString("<key>")
    if err != nil {
        panic("Failed to base64 decode the key")
    }
    pub := gripcontrol.NewGripPubControl([]map[string]interface{} {
            map[string]interface{} {
            "control_uri": "https://api.fanout.io/realm/<realm>",
            "control_iss": "<realm>", 
            "key": decodedKey}})

    // Add new endpoints by applying an endpoint configuration:
    pub.ApplyGripConfig([]map[string]interface{} {
            map[string]interface{} { "control_uri": "<myendpoint_uri_1>" },
            map[string]interface{} { "control_uri": "<myendpoint_uri_2>" }})

    // Remove all configured endpoints:
    pub.RemoveAllClients()

    // Explicitly add an endpoint as a PubControlClient instance:
    client := pubcontrol.NewPubControlClient("<myendpoint_uri>")
    // Optionally set JWT auth: client.SetAuthJwt(<claim>, "<key>")
    // Optionally set basic auth: client.SetAuthBasic("<user>", "<password>")
    pub.AddClient(client)

    // Publish across all configured endpoints:
    err = pub.PublishHttpResponse("<channel>", "Test publish!!", "", "")
    if err != nil {
        panic("Publish failed with: " + err.Error())
    }
    err = pub.PublishHttpStream("<channel>", "Test publish!!", "", "")
    if err != nil {
        panic("Publish failed with: " + err.Error())
    }
}
```

Validate the Grip-Sig request header from incoming GRIP messages. This ensures that the message was sent from a valid source and is not expired. Note that when using Fanout.io the key is the realm key, and when using Pushpin the key is configurable in Pushpin's settings.

```go
isValid := gripcontrol.ValidateSig(request.Header["Grip-Sig"][0], "<key>")
```

Long polling example via response _headers_. The client connects to a GRIP proxy over HTTP and the proxy forwards the request to the origin. The origin subscribes the client to a channel and instructs it to long poll via the response _headers_. Note that with the recent versions of Apache it's not possible to send a 304 response containing custom headers, in which case the response body should be used instead (next usage example below).

```go
package main

import "github.com/fanout/go-gripcontrol"
import "net/http"

func HandleRequest(writer http.ResponseWriter, request *http.Request) {
    // Validate the Grip-Sig header:
    if !gripcontrol.ValidateSig(request.Header["Grip-Sig"][0], "<key>") {
        http.Error(writer, "GRIP authorization failed", http.StatusUnauthorized)
        return
    }

    // Create channel header containing channel information:
    channel := gripcontrol.CreateGripChannelHeader([]*gripcontrol.Channel {
            &gripcontrol.Channel{Name: "<channel>"}})

    // Instruct the client to long poll via the response headers:
    writer.Header().Set("Grip-Hold", "response")
    writer.Header().Set("Grip-Channel", channel)
    // To optionally set a timeout value in seconds:
    // writer.Header().Set("Grip-Timeout", "<timeout_value>")
}

func main() {
    http.HandleFunc("/", HandleRequest)
    http.ListenAndServe(":80", nil)
}
```

Long polling example via response _body_. The client connects to a GRIP proxy over HTTP and the proxy forwards the request to the origin. The origin subscribes the client to a channel and instructs it to long poll via the response _body_.

```go
package main

import "github.com/fanout/go-gripcontrol"
import "net/http"
import "io"

func HandleRequest(writer http.ResponseWriter, request *http.Request) {
    // Validate the Grip-Sig header:
    if !gripcontrol.ValidateSig(request.Header["Grip-Sig"][0], "<key>") {
        http.Error(writer, "GRIP authorization failed", http.StatusUnauthorized)
        return
    }

    // Create channel list containing channel information:
    channel := []*gripcontrol.Channel {&gripcontrol.Channel{Name: "<channel>"}}

    // Create hold response body:
    body, err := gripcontrol.CreateHoldResponse(channel, nil, nil)
    // Or to optionally set a timeout value in seconds:
    // timeout := <timeout_value>
    // body, err := gripcontrol.CreateHoldResponse(channel, nil, &timeout)
    if err != nil {
        panic("Failed to create hold response: " + err.Error())
    }

    // Instruct the client to long poll via the response body:
    writer.Header().Set("Content-Type", "application/grip-instruct")
    io.WriteString(writer, body)
}

func main() {
    http.HandleFunc("/", HandleRequest)
    http.ListenAndServe(":80", nil)
}
```

WebSocket example using golang.org/x/net/websocket. A client connects to a GRIP proxy via WebSockets and the proxy forward the request to the origin. The origin accepts the connection over a WebSocket and responds with a control message indicating that the client should be subscribed to a channel. Note that in order for the GRIP proxy to properly interpret the control messages, the origin must provide a 'grip' extension in the 'Sec-WebSocket-Extensions' header.

```go
package main

import "time"
import "net/http"
import "github.com/gorilla/websocket"
import "github.com/fanout/go-pubcontrol"
import "github.com/fanout/go-gripcontrol"

var upgrader = websocket.Upgrader{
    ReadBufferSize:  1024,
    WriteBufferSize: 1024,
    CheckOrigin: func(r *http.Request) bool { return true },
}

func GripWebSocketHandler(writer http.ResponseWriter, request *http.Request) {
    // Create the WebSocket control message:
    wsControlMessage, err := gripcontrol.WebSocketControlMessage("subscribe",
            map[string]interface{} { "channel": "<channel>" })
    if err != nil {
        panic("Unable to create control message: " + err.Error())
    }

    // Ensure that the GRIP proxy processes control messages by upgrading
    // with the Sec-WebSocket-Extensions header:
    conn, _ := upgrader.Upgrade(writer, request, http.Header {
            "Sec-WebSocket-Extensions": []string {"grip; message-prefix=\"\""}})

    // Subscribe the WebSocket to a channel:
    conn.WriteMessage(1, []byte("c:" + wsControlMessage))

    // Wait 3 seconds and publish a message to the subscribed channel:
    time.Sleep(3 * time.Second)
    pub := gripcontrol.NewGripPubControl([]map[string]interface{} {
            map[string]interface{} { "control_uri": "<myendpoint_uri>" }})
    format := &gripcontrol.WebSocketMessageFormat {
            Content: []byte("Test WebSocket Publish!!") } 
    item := pubcontrol.NewItem([]pubcontrol.Formatter{format}, "", "")
    err = pub.Publish("test_channel", item)
    if err != nil {
        panic("Publish failed with: " + err.Error())
    }
}

func main() {
    http.HandleFunc("/", GripWebSocketHandler)
    http.ListenAndServe(":80", nil)
}
```

WebSocket over HTTP example. In this case, a client connects to a GRIP proxy via WebSockets and the GRIP proxy communicates with the origin via HTTP.

```go
package main

import "github.com/fanout/go-gripcontrol"
import "github.com/fanout/go-pubcontrol"
import "io/ioutil"
import "net/http"
import "time"
import "io"

func HandleRequest(writer http.ResponseWriter, request *http.Request) {
    // Validate the Grip-Sig header:
    if !gripcontrol.ValidateSig(request.Header["Grip-Sig"][0], "<key>") {
        http.Error(writer, "GRIP authorization failed", http.StatusUnauthorized)
        return
    }

    // Set the headers required by the GRIP proxy:
    writer.Header().Set("Sec-WebSocket-Extensions", "grip; message-prefix=\"\"")
    writer.Header().Set("Content-Type", "application/websocket-events")

    body, _ := ioutil.ReadAll(request.Body)
    inEvents, err := gripcontrol.DecodeWebSocketEvents(string(body))
    if err != nil {
        panic("Failed to decode WebSocket events: " + err.Error())
    }

    if inEvents[0].Type == "OPEN" {
        // Create the WebSocket control message:
        wsControlMessage, err := gripcontrol.WebSocketControlMessage("subscribe",
                map[string]interface{} { "channel": "<channel>" })
        if err != nil {
            panic("Unable to create control message: " + err.Error())
        }

        // Open the WebSocket and subscribe it to a channel:
        outEvents := []*gripcontrol.WebSocketEvent {
                &gripcontrol.WebSocketEvent { Type: "OPEN" },
                &gripcontrol.WebSocketEvent { Type: "TEXT",
                        Content: "c:" + wsControlMessage }}
        io.WriteString(writer, gripcontrol.EncodeWebSocketEvents(outEvents))

        go func() {
            // Wait 3 seconds and publish a message to the subscribed channel:
            time.Sleep(3 * time.Second)
            pub := gripcontrol.NewGripPubControl([]map[string]interface{} {
                    map[string]interface{} { "control_uri": "<myendpoint_uri>" }})
            format := &gripcontrol.WebSocketMessageFormat {
                    Content: []byte("Test WebSocket Publish!!") } 
            item := pubcontrol.NewItem([]pubcontrol.Formatter{format}, "", "")
            err = pub.Publish("test_channel", item)
            if err != nil {
                panic("Publish failed with: " + err.Error())
            }
        }()
    }
}

func main() {
    http.HandleFunc("/", HandleRequest)
    http.ListenAndServe(":80", nil)
}
```

Parse a GRIP URI to extract the URI, ISS, and key values. The values will be returned in a hash containing 'control_uri', 'control_iss', and 'key' keys.

```go
config := gripcontrol.ParseGripUri(
    "http://api.fanout.io/realm/<myrealm>?iss=<myrealm>" +
    "&key=base64:<myrealmkey>")
```
