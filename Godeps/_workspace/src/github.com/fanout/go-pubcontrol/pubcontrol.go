//    pubcontrol.go
//    ~~~~~~~~~
//    This module implements the PubControl struct and features.
//    :authors: Konstantin Bokarius.
//    :copyright: (c) 2015 by Fanout, Inc.
//    :license: MIT, see LICENSE for more details.

package pubcontrol

// The PubControl struct allows a consumer to manage a set of publishing
// endpoints and to publish to all of those endpoints via a single publish
// method call. A PubControl instance can be configured either using a
// hash or array of hashes containing configuration information or by
// manually adding PubControlClient instances.
type PubControl struct {
    clients []*PubControlClient
}

// Initialize with or without a configuration. A configuration can be applied
// after initialization via the apply_config method.
func NewPubControl(config []map[string]interface{}) *PubControl {
    pc := new(PubControl)
    pc.clients = make([]*PubControlClient, 0)
    if config != nil && len(config) > 0 {
        pc.ApplyConfig(config)
    }
    return pc
}

// Remove all of the configured PubControlClient instances.
func (pc *PubControl) RemoveAllClients() {
    pc.clients = make([]*PubControlClient, 0)
}

// Add the specified PubControlClient instance.
func (pc *PubControl) AddClient(pcc *PubControlClient) {
    pc.clients = append(pc.clients, pcc)
}

// Apply the specified configuration to this PubControl instance. The
// configuration object can either be a hash or an array of hashes where
// each hash corresponds to a single PubControlClient instance. Each hash
// will be parsed and a PubControlClient will be created either using just
// a URI or a URI and JWT authentication information.
func (pc *PubControl) ApplyConfig(config []map[string]interface{}) {
    for _, entry := range config {
        if _, ok := entry["uri"]; !ok {
            continue
        }
        pcc := NewPubControlClient(entry["uri"].(string))
        if _, ok := entry["iss"]; ok {
            claim := make(map[string]interface{})
            claim["iss"] = entry["iss"]
            switch entry["key"].(type) {
                case string:
                    pcc.SetAuthJwt(claim, []byte(entry["key"].(string)))
                case []byte:
                    pcc.SetAuthJwt(claim, entry["key"].([]byte))
            }
        }
        pc.clients = append(pc.clients, pcc)
    }
}

// The publish method for publishing the specified item to the specified
// channel on the configured endpoints.
func (pc *PubControl) Publish(channel string, item *Item) error {
    for _, pcc := range pc.clients {
        err := pcc.Publish(channel, item)
        if err != nil {
            return err
        }
    }
    return nil
}
