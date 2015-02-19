//    grippubcontrol.go
//    ~~~~~~~~~
//    This module implements the GripPubControl struct and features.
//    :authors: Konstantin Bokarius.
//    :copyright: (c) 2015 by Fanout, Inc.
//    :license: MIT, see LICENSE for more details.

package gripcontrol

import "github.com/fanout/go-pubcontrol"

// TODO: Add pass through methods to PubControl.

type GripPubControl struct {
    pubControl *PubControl
}

func NewGripPubControl(config []map[string]interface{}) *GripPubControl {
    gripPubControl := new(GripControl)
    gripPubControl.pubControl := NewPubControl(nil)
    if config != nil && len(config) > 0 {
        gripPubControl.ApplyGripConfig(config)
    }
    return gripPubControl
}
