package ws

import (
	"encoding/json"
	"fmt"

	"github.com/amir-the-h/okex"
	"github.com/amir-the-h/okex/events"
	"github.com/amir-the-h/okex/events/business"
	requests "github.com/amir-the-h/okex/requests/ws/business"
)

// Business
//
// https://www.okex.com/docs-v5/en/#websocket-api-public-channels
type Business struct {
	*ClientWs
	cCh chan *business.Candle
}

// NewPublic returns a pointer to a fresh Business
func NewBusiness(c *ClientWs) *Business {
	return &Business{ClientWs: c}
}

// Instruments
// The full instrument list will be pushed for the first time after subscription. Subsequently, the instruments will be pushed if there's any change to the instrument’s state (such as delivery of FUTURES, exercise of OPTION, listing of new contracts / trading pairs, trading suspension, etc.).
//
// https://www.okex.com/docs-v5/en/#websocket-api-public-channels-instruments-channel
func (c *Business) Candle(req requests.Candle, ch ...chan *business.Candle) error {
	m := okex.S2M(req)
	if len(ch) > 0 {
		c.cCh = ch[0]
	}
	return c.Subscribe(false, "business", []okex.ChannelName{"candle1D"}, m)
}

// UInstruments
//
// https://www.okex.com/docs-v5/en/#websocket-api-public-channels-instruments-channel
func (c *Business) UCandle(req requests.Candle, rCh ...bool) error {
	m := okex.S2M(req)
	if len(rCh) > 0 && rCh[0] {
		c.cCh = nil
	}
	return c.Unsubscribe(false, "business", []okex.ChannelName{"candle1D"}, m)
}

func (c *Business) Process(data []byte, e *events.Basic) bool {
	if e.Event == "" && e.Arg != nil && e.Data != nil && len(e.Data) > 0 {
		ch, ok := e.Arg.Get("channel")
		if !ok {
			return false
		}
		switch ch {
		case "candle1D":
			fmt.Println("candle1D : ", string(data))
			e := business.Candle{}
			if err := json.Unmarshal(data, &e); err != nil {
				return false
			}
			if c.cCh != nil {
				c.cCh <- &e
			}
			if c.StructuredEventChan != nil {
				c.StructuredEventChan <- e
			}
			return true
		default:
		}
	}
	return false
}
