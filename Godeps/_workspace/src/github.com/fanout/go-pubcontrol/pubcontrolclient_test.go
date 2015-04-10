//    pubcontrolclient_test.go
//    ~~~~~~~~~
//    This module implements the PubControlClient tests.
//    :authors: Konstantin Bokarius.
//    :copyright: (c) 2015 by Fanout, Inc.
//    :license: MIT, see LICENSE for more details.

package pubcontrol

import ("testing"
        "strings"
        "fmt"
        "bytes"
        "encoding/base64"
        "net/url"
        "net/http"
        "net/http/httptest"
        "encoding/json"
        "github.com/stretchr/testify/assert")

func TestPccInitialize(t *testing.T) {
    pcc := NewPubControlClient("uri")
    assert.Equal(t, pcc.uri, "uri")
    assert.NotNil(t, pcc.lock)
}

func TestPccSetAuthBasic(t *testing.T) {
    pcc := NewPubControlClient("uri")
    pcc.SetAuthBasic("user", "pass")
    assert.Equal(t, pcc.authBasicUser, "user")
    assert.Equal(t, pcc.authBasicPass, "pass")
}

func TestPccSetAuthJwt(t *testing.T) {
    pcc := NewPubControlClient("uri")
    pcc.SetAuthJwt(map[string]interface{}{"iss": "iss"}, []byte("key=="))
    assert.Equal(t, pcc.authJwtClaim, map[string]interface{}{"iss": "iss"})
    assert.Equal(t, pcc.authJwtKey, []byte("key=="))
}

func TestPccGenerateAuthHeaderBasic(t *testing.T) {
    pcc := NewPubControlClient("uri")
    pcc.SetAuthBasic("user", "pass")
    authHeader, err := pcc.generateAuthHeader()
    assert.Nil(t, err)
    assert.Equal(t, authHeader, strings.Join([]string{"Basic ",
             base64.StdEncoding.EncodeToString([]byte("user:pass"))}, ""))
}

func TestPccGenerateAuthHeaderJwt(t *testing.T) {
    pcc := NewPubControlClient("uri")
    pcc.SetAuthJwt(map[string]interface{}{"iss": "iss", "exp": 1428374723},
            []byte("key=="))
    authHeader, err := pcc.generateAuthHeader()
    assert.Nil(t, err)
    assert.Equal(t, authHeader, "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9" +
            ".eyJleHAiOjE0MjgzNzQ3MjMsImlzcyI6ImlzcyJ9.33naU1OzEkqe" +
            "7plBOXsbfxwzhFGo9c3ggTQpygQAKRw")
}

var pubCallResults []interface{} = nil

func pubCallTestMethod(pcc *PubControlClient, uri, authHeader string,
        items []map[string]interface{}) error {
    pubCallResults = append(pubCallResults, uri, authHeader, items)
    return nil
}

func pubCallTestMethodFailure(pcc *PubControlClient, uri, authHeader string,
        items []map[string]interface{}) error {
    return &PublishError{err: "error"}
}

func TestPccPublish(t *testing.T) {
    pubCallResults = nil
    formats := make([]Formatter, 0)
    formats = append(formats, fmt1a)
    item := NewItem(formats, "id", "prev-id")
    pcc := NewPubControlClient("uri")
    pcc.pubCall = pubCallTestMethod
    pcc.SetAuthBasic("user", "pass")
    err := pcc.Publish("chan", item)
    assert.Nil(t, err)
    assert.Equal(t, pubCallResults[0], "uri")
    assert.Equal(t, pubCallResults[1], strings.Join([]string{"Basic ",
             base64.StdEncoding.EncodeToString([]byte("user:pass"))}, ""))
    export, err := item.Export()
    export["channel"] = "chan"
    assert.Equal(t, pubCallResults[2], [](map[string]interface{}){export})
}

func TestPccPublishNoAuth(t *testing.T) {
    pubCallResults = nil
    formats := make([]Formatter, 0)
    formats = append(formats, fmt1a)
    item := NewItem(formats, "id", "prev-id")
    pcc := NewPubControlClient("uri")
    pcc.pubCall = pubCallTestMethod
    err := pcc.Publish("chan", item)
    assert.Nil(t, err)
    assert.Equal(t, pubCallResults[0], "uri")
    assert.Equal(t, pubCallResults[1], "")
    export, err := item.Export()
    export["channel"] = "chan"
    assert.Equal(t, pubCallResults[2], [](map[string]interface{}){export})
}

func TestPccPublishErrorItem(t *testing.T) {
    pubCallResults = nil
    formats := make([]Formatter, 0)
    formats = append(formats, fmt1a)
    formats = append(formats, fmt1b)
    item := NewItem(formats, "", "")
    pcc := NewPubControlClient("uri")
    pcc.pubCall = pubCallTestMethod
    err := pcc.Publish("chan", item)
    assert.NotNil(t, err)
}

func TestPccPublishErrorPubCall(t *testing.T) {
    pubCallResults = nil
    formats := make([]Formatter, 0)
    formats = append(formats, fmt1a)
    item := NewItem(formats, "", "")
    pcc := NewPubControlClient("uri")
    pcc.pubCall = pubCallTestMethodFailure
    err := pcc.Publish("chan", item)
    assert.NotNil(t, err)
}

var makeHttpRequestResults []interface{} = nil
func makeHttpRequestTestMethod(pcc *PubControlClient, uri, authHeader string,
        jsonContent []byte) (int, []byte, error) {
    makeHttpRequestResults = append(makeHttpRequestResults, uri, authHeader,
            jsonContent)
    return 200, nil, nil
}
func makeHttpRequestTestMethodFailure(pcc *PubControlClient, uri,
        authHeader string, jsonContent []byte) (int, []byte, error) {
    return 300, []byte("body"), &PublishError{err: "message"}
}

func TestPccPubCall(t *testing.T) {
    makeHttpRequestResults = nil
    items := make([]map[string]interface{}, 0)
    items = append(items, map[string]interface{}{"item": "value"})
    pcc := NewPubControlClient("uri")
    pcc.makeHttpRequest = makeHttpRequestTestMethod
    err := pcc.pubCall(pcc, "http://uri.com", "auth header", items)
    assert.Nil(t, err)
    assert.Equal(t, makeHttpRequestResults[0], "http://uri.com/publish/")
    assert.Equal(t, makeHttpRequestResults[1], "auth header")
    content := make(map[string]interface{})
    content["items"] = items
    jsonContent, _ := json.Marshal(content)
    assert.Equal(t, makeHttpRequestResults[2], jsonContent)
}

func TestPccPubCallError(t *testing.T) {
    pcc := NewPubControlClient("uri")
    pcc.makeHttpRequest = makeHttpRequestTestMethodFailure
    err := pcc.pubCall(pcc, "http://uri.com", "", nil)
    assert.NotNil(t, err)
}

func TestPccMakeHttpRequest(t *testing.T) {
    pcc := NewPubControlClient("uri")    
    server := httptest.NewServer(http.HandlerFunc(func(
            writer http.ResponseWriter, request *http.Request) {
        assert.Equal(t, request.Method, "POST")
        assert.Equal(t, request.RequestURI, "http://uri.com/")
        buf := new(bytes.Buffer)
        buf.ReadFrom(request.Body)
        assert.Equal(t, buf.String(), "content")
        assert.Equal(t, request.Header.Get("Content-Type"),
                "application/json")
        assert.Equal(t, request.Header.Get("Authorization"),
                "auth header")
        writer.WriteHeader(200)
        fmt.Fprintln(writer, "body")
    }))
    defer server.Close()
    transport := &http.Transport{
        Proxy: func(req *http.Request) (*url.URL, error) {
            return url.Parse(server.URL)
        },
    }
    pcc.httpClient = &http.Client{Transport: transport}
    code, body, err := pcc.makeHttpRequest(pcc, "http://uri.com", "auth header",
            []byte("content"))
    assert.Equal(t, code, 200)
    assert.Equal(t, string(body), "body\n")
    assert.Nil(t, err)
}

func TestPccMakeHttpRequestError(t *testing.T) {
    pcc := NewPubControlClient("uri")    
    code, body, err := pcc.makeHttpRequest(pcc, "xxx://uri.com", "auth header",
            []byte("content"))
    assert.Equal(t, code, 0)
    assert.Equal(t, body, []byte(nil))
    assert.NotNil(t, err)
}
