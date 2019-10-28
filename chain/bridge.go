package chain

import (
	"errors"
	"net/url"
	"strings"

	"github.com/wzshiming/bridge"
)

type BridgeChain map[string]bridge.Bridger

func (b BridgeChain) Bridge(dialer bridge.Dialer, addr string) (bridge.Dialer, bridge.ListenConfig, error) {
	return b.BridgeChain(dialer, strings.Split(addr, ">")...)
}

func (b BridgeChain) BridgeChain(dialer bridge.Dialer, addrs ...string) (bridge.Dialer, bridge.ListenConfig, error) {
	if len(addrs) == 0 {
		return dialer, nil, nil
	}
	addr := addrs[0]
	d, l, err := b.bridge(dialer, addr)
	if err != nil {
		return nil, nil, err
	}
	addrs = addrs[1:]
	if len(addrs) == 0 {
		return d, l, nil
	}
	return b.BridgeChain(d, addrs...)
}

func (b BridgeChain) bridge(dialer bridge.Dialer, addr string) (bridge.Dialer, bridge.ListenConfig, error) {
	ur, err := url.Parse(addr)
	if err != nil {
		return nil, nil, err
	}
	dial, ok := b[ur.Scheme]
	if !ok {
		return nil, nil, errors.New("not define " + addr)
	}
	return dial.Bridge(dialer, addr)
}
