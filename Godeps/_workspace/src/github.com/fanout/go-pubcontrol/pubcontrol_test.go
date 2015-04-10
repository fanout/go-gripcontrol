//    pubcontrol_test.go
//    ~~~~~~~~~
//    This module implements the PubControl tests.
//    :authors: Konstantin Bokarius.
//    :copyright: (c) 2015 by Fanout, Inc.
//    :license: MIT, see LICENSE for more details.

package pubcontrol

import ("testing"
        "github.com/stretchr/testify/assert")

func TestPcInitialize(t *testing.T) {
    pc := NewPubControl(nil)
    assert.Equal(t, len(pc.clients), 0)
    pc = NewPubControl([]map[string]interface{} {
            map[string]interface{} {
            "uri": "uri",
            "iss": "hello", 
            "key": "key"}})
    assert.Equal(t, len(pc.clients), 1)
}

func TestPcAddAndRemoveAllClients(t *testing.T) {
    pc := NewPubControl(nil)
    pcc := NewPubControlClient("uri")
    pc.AddClient(pcc)
    pc.AddClient(pcc)
    assert.Equal(t, len(pc.clients), 2)
    assert.Equal(t, pc.clients[0], pcc)
    assert.Equal(t, pc.clients[1], pcc)
    pc.RemoveAllClients()
    assert.Equal(t, len(pc.clients), 0)
}

func TestApplyConfig(t *testing.T) {
    pc := NewPubControl(nil)
    pc.ApplyConfig([]map[string]interface{} {
            map[string]interface{} {
            "uri": "uri",
            "iss": "hello", 
            "key": "key"},
            map[string]interface{} {
            "uri": "uri2",
            "iss": "hello2", 
            "key": "key2"}})
    claim := make(map[string]interface{})
    claim["iss"] = "hello"
    claim2 := make(map[string]interface{})
    claim2["iss"] = "hello2"
    assert.Equal(t, pc.clients[0].uri, "uri")
    assert.Equal(t, pc.clients[0].authJwtClaim, claim)
    assert.Equal(t, pc.clients[0].authJwtKey, []byte("key"))
    assert.Equal(t, pc.clients[1].uri, "uri2")
    assert.Equal(t, pc.clients[1].authJwtClaim, claim2)
    assert.Equal(t, pc.clients[1].authJwtKey, []byte("key2"))
    pc = NewPubControl(nil)
    pc.ApplyConfig([]map[string]interface{} {
            map[string]interface{} {
            "uri": "uri"}})
    assert.Equal(t, pc.clients[0].uri, "uri")
    assert.Equal(t, pc.clients[0].authJwtClaim, map[string]interface{}(nil))
    assert.Equal(t, pc.clients[0].authJwtKey, []byte(nil))
}

var publishResults1 []interface{} = nil
func publish1(pcc *PubControlClient, channel string, item *Item) error {
    publishResults1 = append(publishResults1, channel, item)
    return nil
}
var publishResults2 []interface{} = nil
func publish2(pcc *PubControlClient, channel string, item *Item) error {
    publishResults2 = append(publishResults2, channel, item)
    return nil
}

func TestPcPublish(t *testing.T) {
    publishResults1 = nil
    publishResults2 = nil
    formats := make([]Formatter, 0)
    formats = append(formats, fmt1a)
    item := NewItem(formats, "id", "prev-id")
    pc := NewPubControl(nil)
    pcc := NewPubControlClient("uri")
    pcc.publish = publish1
    pc.AddClient(pcc)
    pcc = NewPubControlClient("uri")
    pcc.publish = publish2
    pc.AddClient(pcc)
    pc.Publish("chan", item)
    assert.Equal(t, publishResults1[0], "chan")
    assert.Equal(t, publishResults1[1], item)
    assert.Equal(t, publishResults2[0], "chan")
    assert.Equal(t, publishResults2[1], item)
}
