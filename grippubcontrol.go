//    grippubcontrol.go
//    ~~~~~~~~~
//    This module implements the GripPubControl struct and features.
//    :authors: Konstantin Bokarius.
//    :copyright: (c) 2015 by Fanout, Inc.
//    :license: MIT, see LICENSE for more details.

package gripcontrol

import "github.com/fanout/go-pubcontrol"

// The GripPubControl struct allows consumers to easily publish HTTP response
// and HTTP stream format messages to GRIP proxies. Configuring GripPubControl
// is slightly different from configuring PubControl in that the 'uri' and
// 'iss' keys in each config entry should have a 'control_' prefix.
// GripPubControl inherits from PubControl and therefore also provides all
// of the same functionality.
type GripPubControl struct {
    *pubcontrol.PubControl
}

// Initialize with or without a configuration. A configuration can be applied
// after initialization via the apply_grip_config method.
func NewGripPubControl(config []map[string]interface{}) *GripPubControl {
    gripPubControl := &GripPubControl{pubcontrol.NewPubControl(nil)}
    if config != nil && len(config) > 0 {
        gripPubControl.ApplyGripConfig(config)
    }
    return gripPubControl
}

// Apply the specified GRIP configuration to this GripPubControl instance.
// The configuration object can either be a hash or an array of hashes where
// each hash corresponds to a single PubControlClient instance. Each hash
// will be parsed and a PubControlClient will be created either using just
// a URI or a URI and JWT authentication information.
func (gpc *GripPubControl) ApplyGripConfig(config []map[string]interface{}) {
    for _, entry := range config {
        if _, ok := entry["control_uri"]; !ok {
            continue
        }
        pcc := pubcontrol.NewPubControlClient(entry["control_uri"].(string))
        if _, ok := entry["control_iss"]; ok {
            claim := make(map[string]interface{})
            claim["iss"] = entry["control_iss"]
            switch entry["key"].(type) {
                case string:
                    pcc.SetAuthJwt(claim, []byte(entry["key"].(string)))
                case []byte:
                    pcc.SetAuthJwt(claim, entry["key"].([]byte))
            }
        } else if _, ok := entry["key"]; ok {
            switch entry["key"].(type) {
                case string:
                    pcc.SetAuthBearer([]byte(entry["key"].(string)))
                case []byte:
                    pcc.SetAuthBearer(entry["key"].([]byte))
            }
        }
        gpc.AddClient(pcc)
    }
}

// Publish an HTTP response format message to all of the configured
// PubControlClients with a specified channel, message, and optional ID,
// previous ID, and callback. Note that the 'http_response' parameter can
// be provided as either an HttpResponseFormat instance or a string / byte
// array (in which case an HttpResponseFormat instance will automatically
// be created and have the 'body' field set to the specified value).
func (gpc *GripPubControl) PublishHttpResponse(channel string,
        http_response interface{}, id, prevId string) error {
    item, err := getHttpResponseItem(http_response, id, prevId)
    if err != nil {
        return err
    }
    return gpc.Publish(channel, item)
}

// Publish an HTTP stream format message to all of the configured
// PubControlClients with a specified channel, message, and optional ID,
// previous ID, and callback. Note that the 'http_stream' parameter can
// be provided as either an HttpStreamFormat instance or a string / byte
// array (in which case an HttpStreamFormat instance will automatically
// be created and have the 'content' field set to the specified value).
func (gpc *GripPubControl) PublishHttpStream(channel string,
        http_stream interface{}, id, prevId string) error {
    item, err := getHttpStreamItem(http_stream, id, prevId)
    if err != nil {
        return err
    }
    return gpc.Publish(channel, item)
}

// An internal method for returning an Item instance used for HTTP response
// publishing based on the specified parameters.
func getHttpResponseItem(http_response interface{}, id,
        prevId string) (*pubcontrol.Item, error) {
    var format *HttpResponseFormat
    switch http_response.(type) {
        case *HttpResponseFormat:
            format = http_response.(*HttpResponseFormat)
        case string:
            format = &HttpResponseFormat{Body: []byte(http_response.(string))}
        case []byte:
            format = &HttpResponseFormat{Body: http_response.([]byte)}
        default:
            return nil, &GripPublishError{err:
                "http_response parameter must be of type " +
                "*HttpResponseFormat, string, or []byte"}
    }
    return pubcontrol.NewItem([]pubcontrol.Formatter{format}, id, prevId), nil
}

// An internal method for returning an Item instance used for HTTP stream
// publishing based on the specified parameters.
func getHttpStreamItem(http_stream interface{}, id,
        prevId string) (*pubcontrol.Item, error) {
    var format *HttpStreamFormat
    switch http_stream.(type) {
        case *HttpStreamFormat:
            format = http_stream.(*HttpStreamFormat)
        case string:
            format = &HttpStreamFormat{Content: []byte(http_stream.(string))}
        case []byte:
            format = &HttpStreamFormat{Content: http_stream.([]byte)}
        default:
            return nil, &GripPublishError{err:
                "http_stream parameter must be of type " +
                "*HttpStreamFormat, string, or []byte"} 
    }
    return pubcontrol.NewItem([]pubcontrol.Formatter{format}, id, prevId), nil
}

// An error object representing an error encountered during publishing.
type GripPublishError struct {
    err string
}

// The function used to retrieve the message associated with a
// GripPublishError.
func (e GripPublishError) Error() string {
    return e.err
}
