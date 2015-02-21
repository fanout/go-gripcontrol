//    pubcontrol.go
//    ~~~~~~~~~
//    This module implements the PubControl struct and features.
//    :authors: Konstantin Bokarius.
//    :copyright: (c) 2015 by Fanout, Inc.
//    :license: MIT, see LICENSE for more details.

package pubcontrol

type PubControl struct {
    clients []*PubControlClient
}

func NewPubControl(config []map[string]interface{}) *PubControl {
    pc := new(PubControl)
    pc.clients = make([]*PubControlClient, 0)
    if config != nil && len(config) > 0 {
        pc.ApplyConfig(config)
    }
    return pc
}

func (pc *PubControl) RemoveAllClients() {
    pc.clients = make([]*PubControlClient, 0)
}

func (pc *PubControl) AddClient(pcc *PubControlClient) {
    pc.clients = append(pc.clients, pcc)
}

func (pc *PubControl) ApplyConfig(config []map[string]interface{}) {
    for _, entry := range config {
        if _, ok := entry["uri"]; !ok {
            continue
        }
        pcc := NewPubControlClient(entry["uri"].(string))
        if _, ok := entry["iss"]; ok {
            claim := make(map[string]interface{})
            claim["iss"] = entry["iss"]
            pcc.SetAuthJwt(claim, entry["key"].([]byte))
        }
        pc.clients = append(pc.clients, pcc)
    }
}

func (pc *PubControl) Publish(channel string, item *Item) error {
    for _, pcc := range pc.clients {
        err := pcc.Publish(channel, item)
        if err != nil {
            return err
        }
    }
    return nil
}
