//    pubcontrolclient.go
//    ~~~~~~~~~
//    This module implements the PubControlClient functionality.
//    :authors: Konstantin Bokarius.
//    :copyright: (c) 2015 by Fanout, Inc.
//    :license: MIT, see LICENSE for more details.

package pubcontrol

import (
    "sync"
    "strings"
    "time"
    "net/http"
    "bytes"
    "io/ioutil"
    "strconv"
    "encoding/json"
    "github.com/dgrijalva/jwt-go"
    "encoding/base64"
)

// An internal type used to define the Publish method.
type publisher func(pcc *PubControlClient, channel string, item *Item) error

// An internal type used to define the pubCall method.
type pubCaller func(pcc *PubControlClient, uri, authHeader string,
        items []map[string]interface{}) error

// An internal type used to define the makeHttpRequest method.
type makeHttpRequester func(pcc *PubControlClient, uri, authHeader string,
        jsonContent []byte) (int, []byte, error)

// The PubControlClient struct allows consumers to publish to an endpoint of
// their choice. The consumer wraps a Format struct instance in an Item struct
// instance and passes that to the publish method. The publish method has
// an optional callback parameter that is called after the publishing is 
// complete to notify the consumer of the result.
type PubControlClient struct {
    uri string
    isWorkerRunning bool
    lock *sync.Mutex
    authBasicUser string
    authBasicPass string
    authJwtClaim map[string]interface{}
    authJwtKey []byte
    publish publisher
    pubCall pubCaller
    makeHttpRequest makeHttpRequester
    httpClient *http.Client
}

// Initialize this struct with a URL representing the publishing endpoint.
func NewPubControlClient(uri string) *PubControlClient {
    newPcc := new(PubControlClient)
    newPcc.uri = uri
    newPcc.lock = &sync.Mutex{}
    newPcc.pubCall = pubCall
    newPcc.publish = publish
    newPcc.makeHttpRequest = makeHttpRequest
    newPcc.httpClient = &http.Client{}
    return newPcc
}

// Call this method and pass a username and password to use basic
// authentication with the configured endpoint.
func (pcc *PubControlClient) SetAuthBasic(username, password string) {
    pcc.lock.Lock()
    pcc.authBasicUser = username
    pcc.authBasicPass = password
    pcc.lock.Unlock()
}

// Call this method and pass a claim and key to use JWT authentication
// with the configured endpoint.
func (pcc *PubControlClient) SetAuthJwt(claim map[string]interface{}, 
        key []byte) {
    pcc.lock.Lock()
    pcc.authJwtClaim = claim
    pcc.authJwtKey = key
    pcc.lock.Unlock()
}

// An internal method used to generate an authorization header. The
// authorization header is generated based on whether basic or JWT
// authorization information was provided via the publicly accessible
// 'set_*_auth' methods defined above.
func (pcc *PubControlClient) generateAuthHeader() (string, error) {
    if pcc.authBasicUser != "" {
        encodedCredentials := base64.StdEncoding.EncodeToString([]byte(
                strings.Join([]string{pcc.authBasicUser, ":",
                pcc.authBasicPass}, "")))
        return strings.Join([]string{"Basic ", encodedCredentials}, ""), nil
    } else if pcc.authJwtClaim != nil {
        token := jwt.New(jwt.SigningMethodHS256)
        token.Valid = true
        for k, v := range pcc.authJwtClaim {
            token.Claims[k] = v
        }
        if _, ok := pcc.authJwtClaim["exp"]; !ok {
            token.Claims["exp"] = time.Now().Add(time.Second * 3600).Unix()
        }
        tokenString, err := token.SignedString(pcc.authJwtKey)
        if err != nil {
            return "", err
        }
        return strings.Join([]string{"Bearer ", tokenString}, ""), nil
    } else {
        return "", nil
    }
}

// The publish method for publishing the specified item to the specified
// channel on the configured endpoint.
func (pcc *PubControlClient) Publish(channel string, item *Item) error {
    return pcc.publish(pcc, channel, item)
}

// An internal publish method to facilitate testing.
func publish(pcc *PubControlClient, channel string, item *Item) error {
    export, err := item.Export()
    if err != nil {
        return err
    }
    export["channel"] = channel
    uri := ""
    auth := ""
    pcc.lock.Lock()
    uri = pcc.uri
    auth, err = pcc.generateAuthHeader()
    pcc.lock.Unlock()
    if err != nil {
        return err
    }
    err = pcc.pubCall(pcc, uri, auth, [](map[string]interface{}){export})
    if err != nil {
        return err
    }
    return nil
}

// An internal method for preparing the HTTP POST request for publishing
// data to the endpoint. This method accepts the URI endpoint, authorization
// header, and a list of items to publish.
func pubCall(pcc *PubControlClient, uri, authHeader string,
        items []map[string]interface{}) error {
    uri = strings.Join([]string{uri, "/publish/"}, "")
    content := make(map[string]interface{})
    content["items"] = items
    var jsonContent []byte
    jsonContent, err := json.Marshal(content)
    if err != nil {
        return err
    }
    statusCode, body, err := pcc.makeHttpRequest(pcc, uri, authHeader,
            jsonContent)
    if err != nil {
        return err
    }
    if statusCode < 200 || statusCode >= 300 {
        return &PublishError{err: strings.Join([]string{"Failure status code: ",
                strconv.Itoa(statusCode), " with message: ",
                string(body)}, "")}
    }
    return nil
}

// An internal method used to make the HTTP request for publishing based
// on the specified URI, auth header, and JSON content. An HTTP status
// code, response body, and an error will be returned.
func makeHttpRequest(pcc *PubControlClient, uri, authHeader string,
        jsonContent []byte) (int, []byte, error) {
    var req *http.Request
    req, err := http.NewRequest("POST", uri, bytes.NewReader(jsonContent))
    if err != nil {
        return 0, nil, err
    }
    req.Header.Add("Content-Type", "application/json")
    req.Header.Add("Authorization", authHeader)
    resp, err := pcc.httpClient.Do(req)
    if err != nil {
        return 0, nil, err
    }
    defer resp.Body.Close()
    var body []byte
    body, err = ioutil.ReadAll(resp.Body)
    if err != nil {
        return 0, nil, err
    }
    return resp.StatusCode, body, nil
}

// An error struct used to represent an error encountered during publishing.
type PublishError struct {
    err string
}

// This function returns the message associated with the Publish error struct.
func (e PublishError) Error() string {
    return e.err
}
