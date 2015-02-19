//    grippubcontrol.go
//    ~~~~~~~~~
//    This module implements the GripPubControl struct and features.
//    :authors: Konstantin Bokarius.
//    :copyright: (c) 2015 by Fanout, Inc.
//    :license: MIT, see LICENSE for more details.

package gripcontrol

import "github.com/fanout/go-pubcontrol"

type GripPubControl struct {
    *pubcontrol.PubControl
}

func NewGripPubControl(config []map[string]interface{}) *GripPubControl {
    gripPubControl := &GripPubControl{pubcontrol.NewPubControl(nil)}
    if config != nil && len(config) > 0 {
        gripPubControl.ApplyGripConfig(config)
    }
    return gripPubControl
}

func (gpc *GripPubControl) ApplyGripConfig(config []map[string]interface{}) {
    for _, entry := range config {
        if _, ok := entry["control_uri"]; !ok {
            continue
        }
        pcc := pubcontrol.NewPubControlClient(entry["control_uri"].(string))
        if _, ok := entry["control_iss"]; ok {
            claim := make(map[string]interface{})
            claim["iss"] = entry["control_iss"]
            pcc.SetAuthJwt(claim, entry["key"].([]byte))
        }
        gpc.AddClient(pcc)
    }
}

func (gpc *GripPubControl) PublishHttpResponse(channel string, http_response interface{},
        id, prevId string) error {
    var format *HttpResponseFormat
    switch http_response.(type) {
        case *HttpResponseFormat:
            format = http_response.(*HttpResponseFormat)
        case string:
            format = &HttpResponseFormat{Body: []byte(http_response.(string))}
        case []byte:
            format = &HttpResponseFormat{Body: http_response.([]byte)}
        default:
            return &GripPublishError{err:
                "http_response parameter must be of type " +
                "*HttpResponseFormat, string, or []byte"}
    }
    item := pubcontrol.NewItem([]pubcontrol.Formatter{format}, id, prevId)
    return gpc.Publish(channel, item)
}

func (gpc *GripPubControl) PublishHttpStream(channel string, http_stream interface{},
        id, prevId string) error {
    var format *HttpStreamFormat
    switch http_stream.(type) {
        case *HttpStreamFormat:
            format = http_stream.(*HttpStreamFormat)
        case string:
            format = &HttpStreamFormat{Content: []byte(http_stream.(string))}
        case []byte:
            format = &HttpStreamFormat{Content: http_stream.([]byte)}
        default:
            return &GripPublishError{err:
                "http_stream parameter must be of type " +
                "*HttpStreamFormat, string, or []byte"} 
    }
    item := pubcontrol.NewItem([]pubcontrol.Formatter{format}, id, prevId)
    return gpc.Publish(channel, item)
}

type GripPublishError struct {
    err string
}

func (e GripPublishError) Error() string {
    return e.err
}
