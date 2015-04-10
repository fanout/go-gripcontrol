//    gripcontrol_test.go
//    ~~~~~~~~~
//    This module implements the GripControl tests.
//    :authors: Konstantin Bokarius.
//    :copyright: (c) 2015 by Fanout, Inc.
//    :license: MIT, see LICENSE for more details.

package gripcontrol

import ("testing"
        "time"
        "encoding/json"
        "encoding/base64"
        "github.com/dgrijalva/jwt-go"
        "github.com/stretchr/testify/assert")

func TestCreateHold(t *testing.T) {
    hold, err := CreateHold("", nil, nil, nil)
    assert.Nil(t, err)
    holdToCompare, _ := json.Marshal(map[string]interface{} { "hold":
            map[string]interface{} {
            "mode": "", "channels": getHoldChannels(nil)}})
    assert.Equal(t, hold, string(holdToCompare))
    channels := []*Channel {&Channel{Name:"test_channel1", PrevId:"prev-id"}}
    hold, err = CreateHold("mode", channels, "response", nil)
    assert.Nil(t, err)
    holdResponse, _ := getHoldResponse("response")
    holdToCompare, _ = json.Marshal(map[string]interface{} { "hold":
            map[string]interface{} {
            "mode": "mode", "channels": getHoldChannels(channels)},
            "response": holdResponse})
    assert.Equal(t, hold, string(holdToCompare))
    timeout := 1000
    hold, err = CreateHold("mode", channels, nil, &timeout)
    holdToCompare, _ = json.Marshal(map[string]interface{} { "hold":
            map[string]interface{} {
            "mode": "mode", "channels": getHoldChannels(channels),
            "timeout": 1000}})
    assert.Equal(t, hold, string(holdToCompare))
}

func TestCreateHoldStream(t *testing.T) {
    channels := []*Channel {&Channel{Name:"test_channel1", PrevId:"prev-id"}}
    hold, err := CreateHoldStream(channels, "response")
    assert.Nil(t, err)
    holdResponse, _ := getHoldResponse("response")
    holdToCompare, _ := json.Marshal(map[string]interface{} { "hold":
            map[string]interface{} {
            "mode": "stream", "channels": getHoldChannels(channels)},
            "response": holdResponse})
    assert.Equal(t, hold, string(holdToCompare))
}

func TestCreateHoldResponse(t *testing.T) {
    channels := []*Channel {&Channel{Name:"test_channel1", PrevId:"prev-id"}}
    timeout := 1000
    hold, err := CreateHoldResponse(channels, "response", &timeout)
    assert.Nil(t, err)
    holdResponse, _ := getHoldResponse("response")
    holdToCompare, _ := json.Marshal(map[string]interface{} { "hold":
            map[string]interface{} {
            "mode": "response", "channels": getHoldChannels(channels),
            "timeout": 1000},
            "response": holdResponse})
    assert.Equal(t, hold, string(holdToCompare))
}

func TestParseGripUri(t *testing.T) {
    uri := "http://api.fanout.io/realm/realm?iss=realm" +
            "&key=base64:geag121321=="
    config, err := ParseGripUri(uri)
    assert.Nil(t, err)
    assert.Equal(t, config["control_uri"], "http://api.fanout.io/realm/realm")
    assert.Equal(t, config["control_iss"], "realm")
    key, err := base64.StdEncoding.DecodeString("geag121321==")
    assert.Nil(t, err)
    assert.Equal(t, config["key"], key)
    uri = "https://api.fanout.io/realm/realm?iss=realm" +
            "&key=base64:geag121321=="
    config, err = ParseGripUri(uri)
    assert.Equal(t, config["control_uri"], "https://api.fanout.io/realm/realm")
    config, err = ParseGripUri("http://api.fanout.io/realm/realm")
    assert.Equal(t, config["control_uri"], "http://api.fanout.io/realm/realm")
    assert.False(t, doesKeyExist(config, "control_iss"))
    assert.False(t, doesKeyExist(config, "key"))
    uri = "http://api.fanout.io/realm/realm?iss=realm" +
            "&key=base64:geag121321==&param1=value1&param2=value2"
    config, err = ParseGripUri(uri)
    assert.Nil(t, err)
    assert.Equal(t, config["control_uri"], "http://api.fanout.io/realm/realm?" +
            "param1=value1&param2=value2")
    assert.Equal(t, config["control_iss"], "realm")
    assert.Equal(t, config["key"], key)
    config, err = ParseGripUri("http://api.fanout.io:8080/realm/realm/")
    assert.Nil(t, err)
    assert.Equal(t, config["control_uri"], "http://api.fanout.io:8080/realm/realm")
    uri = "http://api.fanout.io/realm/realm?iss=realm" +
            "&key=geag121321=="
    config, err = ParseGripUri(uri)
    assert.Nil(t, err)
    assert.Equal(t, config["key"], []byte("geag121321=="))
}

func doesKeyExist(obj map[string]interface{}, key string) bool {
    if _, ok := obj["control_iss"]; ok {
        return true
    }
    return false
}

func TestValidateSig(t *testing.T) {
    token := jwt.New(jwt.SigningMethodHS256)
    token.Valid = true
    token.Claims["iss"] = "realm"
    token.Claims["exp"] = time.Now().Add(time.Second * 3600).Unix()
    tokenString, _ := token.SignedString([]byte("key"))
    assert.True(t, ValidateSig(tokenString, "key"))
    token = jwt.New(jwt.SigningMethodHS256)
    token.Valid = true
    token.Claims["iss"] = "realm"
    token.Claims["exp"] = time.Now().Add(0 - time.Second * 3600).Unix()
    tokenString, _ = token.SignedString([]byte("key"))
    assert.False(t, ValidateSig(tokenString, "key")) 
    token = jwt.New(jwt.SigningMethodHS256)
    token.Valid = true
    token.Claims["iss"] = "realm"
    token.Claims["exp"] = time.Now().Add(time.Second * 3600).Unix()
    tokenString, _ = token.SignedString([]byte("key"))
    assert.False(t, ValidateSig(tokenString, "wrong_key")) 
}

func TestCreateGripChannelHeader(t *testing.T) {
    header := CreateGripChannelHeader(
            []*Channel {&Channel{Name:"channel"}})
    assert.Equal(t, header, "channel")
    header = CreateGripChannelHeader(
            []*Channel {&Channel{Name:"channel", PrevId:"prev-id"}})
    assert.Equal(t, header, "channel; prev-id=prev-id")
    header = CreateGripChannelHeader(
        []*Channel {&Channel{Name:"channel1", PrevId:"prev-id1"},
        &Channel{Name:"channel2", PrevId:"prev-id2"}})
    assert.Equal(t, header,
        "channel1; prev-id=prev-id1, channel2; prev-id=prev-id2")
}

func TestDecodeWebSocketEvents(t *testing.T) {
    events, err := DecodeWebSocketEvents("OPEN\r\nTEXT 5\r\nHello" + 
        "\r\nTEXT 0\r\n\r\nCLOSE\r\nTEXT\r\nCLOSE\r\n")
    assert.Nil(t, err)
    assert.Equal(t, len(events), 6)
    assert.Equal(t, events[0].Type, "OPEN")
    assert.Equal(t, events[0].Content, "")
    assert.Equal(t, events[1].Type, "TEXT")
    assert.Equal(t, events[1].Content, "Hello")
    assert.Equal(t, events[2].Type, "TEXT")
    assert.Equal(t, events[2].Content, "")
    assert.Equal(t, events[3].Type, "CLOSE")
    assert.Equal(t, events[3].Content, "")
    assert.Equal(t, events[4].Type, "TEXT")
    assert.Equal(t, events[4].Content, "")
    assert.Equal(t, events[5].Type, "CLOSE")
    assert.Equal(t, events[5].Content, "")
    events, err = DecodeWebSocketEvents("OPEN\r\n")
    assert.Nil(t, err)
    assert.Equal(t, len(events), 1)
    assert.Equal(t, events[0].Type, "OPEN")
    assert.Equal(t, events[0].Content, "")
    events, err = DecodeWebSocketEvents("TEXT 5\r\nHello\r\n")
    assert.Nil(t, err)
    assert.Equal(t, len(events), 1)
    assert.Equal(t, events[0].Type, "TEXT")
    assert.Equal(t, events[0].Content, "Hello")
    events, err = DecodeWebSocketEvents("TEXT 5")
    assert.Nil(t, events)
    assert.NotNil(t, err)
    events, err = DecodeWebSocketEvents("OPEN\r\nTEXT")
    assert.Nil(t, events)
    assert.NotNil(t, err)
}

func TestEncodeWebSocketEvents(t *testing.T) {
    events := EncodeWebSocketEvents([]*WebSocketEvent{
            &WebSocketEvent{"TEXT", "Hello"},
            &WebSocketEvent{"TEXT", ""}})
    assert.Equal(t, events, "TEXT 5\r\nHello\r\nTEXT\r\n")
}

func TestWebSocketControlMessage(t *testing.T) {
    message, err := WebSocketControlMessage("type", nil)
    assert.Nil(t, err)
    assert.Equal(t, message, "{\"type\":\"type\"}")
    message, err = WebSocketControlMessage("type", map[string]interface{} {
            "arg1": "val1", "arg2": "val2"})
    assert.Nil(t, err)
    assert.Equal(t, message, 
            "{\"arg1\":\"val1\",\"arg2\":\"val2\",\"type\":\"type\"}")
}

func TestGetHoldChannels(t *testing.T) {
    channels := getHoldChannels([]*Channel {
            &Channel{Name:"channel1", PrevId:""},
            &Channel{Name:"channel2", PrevId:"prev-id"}})
    assert.Equal(t, channels[0]["name"], "channel1")
    isPrevIdPresent := true
    if _, ok := channels[0]["prev-id"]; !ok {
        isPrevIdPresent = false
    }
    assert.False(t, isPrevIdPresent)
    assert.Equal(t, channels[1]["name"], "channel2")
    assert.Equal(t, channels[1]["prev-id"], "prev-id")
}

func TestGetHoldResponse(t *testing.T) {
    response, err := getHoldResponse("response")
    assert.Nil(t, err)
    assert.Equal(t, response, map[string]interface{} {"body": "response"})
    response, err = getHoldResponse([]byte("response"))
    assert.Nil(t, err)
    assert.Equal(t, response, map[string]interface{} {"body": "response"})
    response, err = getHoldResponse([]byte("\xbd\xb2\x3d\xbc\x20\xe2\x8c\xFF"))
    assert.Nil(t, err)
    assert.Equal(t, response, map[string]interface{} {"body-bin":
            base64.StdEncoding.EncodeToString(
            []byte("\xbd\xb2\x3d\xbc\x20\xe2\x8c\xFF"))})
    response, err = getHoldResponse(1)
    assert.NotNil(t, err)
    assert.Nil(t, response)
    resp := &Response{1, "reason",
            map[string]string {"head": "hval"}, []byte("response")}
    response, err = getHoldResponse(resp)
    assert.Nil(t, err)
    assert.Equal(t, response, map[string]interface{} {
            "body": "response", "headers": map[string]string {"head": "hval"},
            "reason": "reason", "code": 1})
}
