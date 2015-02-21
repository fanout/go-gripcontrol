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
    if (!gripcontrol.ValidateSig(request.Header["Grip-Sig"][0], "<key>")) {
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
    if (!gripcontrol.ValidateSig(request.Header["Grip-Sig"][0], "<key>")) {
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
    if (err != nil) {
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

WebSocket example using the WEBrick gem and WEBrick WebSocket gem extension. A client connects to a GRIP proxy via WebSockets and the proxy forward the request to the origin. The origin accepts the connection over a WebSocket and responds with a control message indicating that the client should be subscribed to a channel. Note that in order for the GRIP proxy to properly interpret the control messages, the origin must provide a 'grip' extension in the 'Sec-WebSocket-Extensions' header. This is accomplished in the WEBrick WebSocket gem extension by adding the following line to lib/webrick/websocket/server.rb and rebuilding the gem: res['Sec-WebSocket-Extensions'] = 'grip; message-prefix=""'

```go
require 'webrick/websocket'
require 'gripcontrol'
require 'thread'

class GripWebSocket < WEBrick::Websocket::Servlet
  def socket_open(sock)
    # Subscribe the WebSocket to a channel:
    sock.puts('c:' + GripControl.websocket_control_message('subscribe',
        {'channel' => '<channel>'}))
    Thread.new { publish_message }
  end

  def publish_message
    # Wait and then publish a message to the subscribed channel:
    sleep(3)
    grippub = GripPubControl.new({'control_uri' => '<myendpoint>'})
    grippub.publish('<channel>', Item.new(
        WebSocketMessageFormat.new('Test WebSocket publish!!')))
  end
end

server = WEBrick::Websocket::HTTPServer.new(Port: 80)
server.mount "/websocket", GripWebSocket
trap "INT" do server.shutdown end
server.start
```

WebSocket over HTTP example using the WEBrick gem. In this case, a client connects to a GRIP proxy via WebSockets and the GRIP proxy communicates with the origin via HTTP.

```go
require 'webrick'
require 'gripcontrol'

class GripWebSocketOverHttpResponse < WEBrick::HTTPServlet::AbstractServlet
  def do_POST(request, response)
    # Validate the Grip-Sig header:
    if !GripControl.validate_sig(request['Grip-Sig'], '<key>')
      return
    end

    # Set the headers required by the GRIP proxy:
    response.status = 200
    response['Sec-WebSocket-Extensions'] = 'grip; message-prefix=""'
    response['Content-Type'] = 'application/websocket-events'

    in_events = GripControl.decode_websocket_events(request.body)
    if in_events[0].type == 'OPEN'
      # Open the WebSocket and subscribe it to a channel:
      out_events = []
      out_events.push(WebSocketEvent.new('OPEN'))
      out_events.push(WebSocketEvent.new('TEXT', 'c:' +
          GripControl.websocket_control_message('subscribe',
          {'channel' => '<channel>'})))
      response.body = GripControl.encode_websocket_events(out_events)
      Thread.new { publish_message }
    end
  end

  def publish_message
    # Wait and then publish a message to the subscribed channel:
    sleep(3)
    grippub = GripPubControl.new({'control_uri' => '<myendpoint>'})
    grippub.publish('<channel>', Item.new(
        WebSocketMessageFormat.new('Test WebSocket publish!!')))
  end
end

server = WEBrick::HTTPServer.new(Port: 80)
server.mount "/websocket", GripWebSocketOverHttpResponse
trap "INT" do server.shutdown end
server.start
```

Parse a GRIP URI to extract the URI, ISS, and key values. The values will be returned in a hash containing 'control_uri', 'control_iss', and 'key' keys.

```go
config = GripControl.parse_grip_uri(
    'http://api.fanout.io/realm/<myrealm>?iss=<myrealm>' +
    '&key=base64:<myrealmkey>')
```
