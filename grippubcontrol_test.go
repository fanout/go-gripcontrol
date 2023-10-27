//    grippubcontrol_test.go
//    ~~~~~~~~~
//    This module implements the GripPubControl tests.
//    :authors: Konstantin Bokarius.
//    :copyright: (c) 2015 by Fanout, Inc.
//    :license: MIT, see LICENSE for more details.

package gripcontrol

import ("testing"
        "github.com/fanout/go-pubcontrol"
        "github.com/stretchr/testify/assert")

// TODO: Use reflection to ensure that clients are configured properly.

func TestGripPubControlInitialize(t *testing.T) {
    NewGripPubControl(nil)
    NewGripPubControl([]map[string]interface{} {
            map[string]interface{} {
            "control_uri": "uri",
            "control_iss": "hello", 
            "key": []byte("key")}})
}

func TestApplyGripConfig(t *testing.T) {
    gpc := NewGripPubControl(nil)
    gpc.ApplyGripConfig([]map[string]interface{} {
            map[string]interface{} {
            "control_uri": "uri",
            "control_iss": "hello", 
            "key": "key"}})
}

func TestApplyGripConfigBearer(t *testing.T) {
    gpc := NewGripPubControl(nil)
    gpc.ApplyGripConfig([]map[string]interface{} {
            map[string]interface{} {
            "control_uri": "uri",
            "key": "key"}})
}

func TestPublishHttpResponse(t *testing.T) {
    gpc := NewGripPubControl([]map[string]interface{} {
            map[string]interface{} {
            "control_uri": "something://uri",
            "control_iss": "hello", 
            "key": "key"}})
    err := gpc.PublishHttpResponse("chan", "data", "id", "prev-id")
    assert.NotNil(t, err)
}

func TestPublishHttpStream(t *testing.T) {
    gpc := NewGripPubControl([]map[string]interface{} {
            map[string]interface{} {
            "control_uri": "something://uri",
            "control_iss": "hello", 
            "key": "key"}})
    err := gpc.PublishHttpStream("chan", "data", "id", "prev-id")
    assert.NotNil(t, err)
}

func TestGetHttpResponseItem(t *testing.T) {
    item, err := getHttpResponseItem("data", "id", "prev-id")
    assert.Nil(t, err)
    assert.Equal(t, pubcontrol.NewItem([]pubcontrol.Formatter{
            &HttpResponseFormat{Body: []byte("data")}}, "id", "prev-id"), item)
    item, err = getHttpResponseItem([]byte("data"), "id", "prev-id")
    assert.Nil(t, err)
    assert.Equal(t, pubcontrol.NewItem([]pubcontrol.Formatter{
            &HttpResponseFormat{Body: []byte("data")}}, "id", "prev-id"), item)
    fmt := &HttpResponseFormat{Code:1, Reason:"reason",
            Headers:map[string]string{"header":"hval"},
            Body:[]byte("body")}
    item, err = getHttpResponseItem(fmt, "id", "prev-id")
    assert.Nil(t, err)
    assert.Equal(t, pubcontrol.NewItem([]pubcontrol.Formatter{fmt},
            "id", "prev-id"), item)
    item, err = getHttpResponseItem(1, "id", "prev-id")
    assert.Nil(t, item)
    assert.NotNil(t, err)    
}

func TestGetHttpStreamItem(t *testing.T) {
    item, err := getHttpStreamItem("data", "id", "prev-id")
    assert.Nil(t, err)
    assert.Equal(t, pubcontrol.NewItem([]pubcontrol.Formatter{
            &HttpStreamFormat{Content: []byte("data")}}, "id", "prev-id"), item)
    item, err = getHttpStreamItem([]byte("data"), "id", "prev-id")
    assert.Nil(t, err)
    assert.Equal(t, pubcontrol.NewItem([]pubcontrol.Formatter{
            &HttpStreamFormat{Content: []byte("data")}}, "id", "prev-id"), item)
    fmt := &HttpStreamFormat{Content: []byte("content"), Close:true}
    item, err = getHttpStreamItem(fmt, "id", "prev-id")
    assert.Nil(t, err)
    assert.Equal(t, pubcontrol.NewItem([]pubcontrol.Formatter{fmt},
            "id", "prev-id"), item)
    item, err = getHttpStreamItem(1, "id", "prev-id")
    assert.Nil(t, item)
    assert.NotNil(t, err)    
}
