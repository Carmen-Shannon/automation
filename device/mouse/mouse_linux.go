//go:build linux
// +build linux

package mouse

import (
	"fmt"

	"github.com/BurntSushi/xgb"
	"github.com/BurntSushi/xgb/xproto"
	linux "github.com/Carmen-Shannon/automation/tools/_linux"
)

var xConn *xgb.Conn

func initXGB() error {
	var err error
	xConn, err = xgb.NewConn()
	return err
}

func (m *mouse) doMouseMove(x, y int32) error {
	if xConn == nil {
		if err := initXGB(); err != nil {
			return err
		}
	}
	root := xproto.Setup(xConn).DefaultScreen(xConn).Root
	xproto.WarpPointer(xConn, 0, root, 0, 0, 0, 0, int16(x), int16(y))
	return nil
}

func doGetMousePosition() (int32, int32, error) {
	x, y, err := linux.ExecuteXdotoolGetMousePosition()
	if err != nil {
		return 0, 0, fmt.Errorf("failed to get mouse position: %w", err)
	}
	return x, y, nil
}

func (m *mouse) doMouseClick(btn int, duration int) error {
	err := linux.ExecuteXdotoolClick(btn, duration)
	if err != nil {
		return err
	}
	return nil
}
