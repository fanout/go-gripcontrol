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

```Go
require 'gripcontrol'

def callback(result, message)
  if result
    puts 'Publish successful'
  else
    puts 'Publish failed with message: ' + message.to_s
  end
end

# GripPubControl can be initialized with or without an endpoint configuration.
# Each endpoint can include optional JWT authentication info.
# Multiple endpoints can be included in a single configuration.

grippub = GripPubControl.new({ 
    'control_uri' => 'https://api.fanout.io/realm/<myrealm>',
    'control_iss' => '<myrealm>',
    'key' => Base64.decode64('<myrealmkey>')})

# Add new endpoints by applying an endpoint configuration:
grippub.apply_grip_config([{'control_uri' => '<myendpoint_uri_1>'}, 
    {'control_uri' => '<myendpoint_uri_2>'}])

# Remove all configured endpoints:
grippub.remove_all_clients

# Explicitly add an endpoint as a PubControlClient instance:
pubclient = PubControlClient.new('<myendpoint_uri>')
# Optionally set JWT auth: pubclient.set_auth_jwt(<claim>, '<key>')
# Optionally set basic auth: pubclient.set_auth_basic('<user>', '<password>')
grippub.add_client(pubclient)

# Publish across all configured endpoints:
grippub.publish_http_response('<channel>', 'Test publish!')
grippub.publish_http_response_async('<channel>', 'Test async publish!',
    nil, nil, method(:callback))
grippub.publish_http_stream('<channel>', 'Test publish!')
grippub.publish_http_stream_async('<channel>', 'Test async publish!',
    nil, nil, method(:callback))

# Wait for all async publish calls to complete:
grippub.finish
```

Validate the Grip-Sig request header from incoming GRIP messages. This ensures that the message was sent from a valid source and is not expired. Note that when using Fanout.io the key is the realm key, and when using Pushpin the key is configurable in Pushpin's settings.

```Go
is_valid = GripControl.validate_sig(request['Grip-Sig'], '<key>')
```

Long polling example via response _headers_ using the WEBrick gem. The client connects to a GRIP proxy over HTTP and the proxy forwards the request to the origin. The origin subscribes the client to a channel and instructs it to long poll via the response _headers_. Note that with the recent versions of Apache it's not possible to send a 304 response containing custom headers, in which case the response body should be used instead (next usage example below).

```Go
require 'webrick'
require 'gripcontrol'

class GripHeadersResponse < WEBrick::HTTPServlet::AbstractServlet
  def do_GET(request, response)
    # Validate the Grip-Sig header:
    if !GripControl.validate_sig(request['Grip-Sig'], '<key>')
      return
    end

    # Instruct the client to long poll via the response headers:
    response.status = 200
    response['Grip-Hold'] = 'response'
    response['Grip-Channel'] = 
        GripControl.create_grip_channel_header('<channel>')
    # To optionally set a timeout value in seconds:
    # response['Grip-Timeout'] = <timeout_value>
  end
end

server = WEBrick::HTTPServer.new(:Port => 80)
server.mount "/", GripHeadersResponse
trap "INT" do server.shutdown end
server.start
```

Long polling example via response _body_ using the WEBrick gem. The client connects to a GRIP proxy over HTTP and the proxy forwards the request to the origin. The origin subscribes the client to a channel and instructs it to long poll via the response _body_.

```Go
require 'webrick'
require 'gripcontrol'

class GripBodyResponse < WEBrick::HTTPServlet::AbstractServlet
  def do_GET(request, response)
    # Validate the Grip-Sig header:
    if !GripControl.validate_sig(request['Grip-Sig'], '<key>')
      return
    end

    # Instruct the client to long poll via the response body:
    response.status = 200
    response['Content-Type'] = 'application/grip-instruct'
    response.body = GripControl.create_hold_response('<channel>')
    # Or to optionally set a timeout value in seconds:
    # response.body = GripControl.create_hold_response(
    #     '<channel>', nil, <timeout_value>)
  end
end

server = WEBrick::HTTPServer.new(:Port => 80)
server.mount "/", GripBodyResponse
trap "INT" do server.shutdown end
server.start
```

WebSocket example using the WEBrick gem and WEBrick WebSocket gem extension. A client connects to a GRIP proxy via WebSockets and the proxy forward the request to the origin. The origin accepts the connection over a WebSocket and responds with a control message indicating that the client should be subscribed to a channel. Note that in order for the GRIP proxy to properly interpret the control messages, the origin must provide a 'grip' extension in the 'Sec-WebSocket-Extensions' header. This is accomplished in the WEBrick WebSocket gem extension by adding the following line to lib/webrick/websocket/server.rb and rebuilding the gem: res['Sec-WebSocket-Extensions'] = 'grip; message-prefix=""'

```Go
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

```Go
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

```Go
config = GripControl.parse_grip_uri(
    'http://api.fanout.io/realm/<myrealm>?iss=<myrealm>' +
    '&key=base64:<myrealmkey>')
```
