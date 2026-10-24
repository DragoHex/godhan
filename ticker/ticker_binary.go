package ticker


import (
	"bytes"
	"encoding/binary"
	"fmt"
	"time"

	"github.com/shopspring/decimal"
)

// Ticker Modes
const (
	TickerModeLTP   = 15 // Code 15 subscribe request
	TickerModeQuote = 16 // Code 16 subscribe request
	TickerModeDepth = 17 // Code 17 subscribe request
)

// DepthLevel represents a single order book bid or ask tier
type DepthLevel struct {
	Price    decimal.Decimal `json:"price"`
	Quantity uint32          `json:"quantity"`
	Orders   uint16          `json:"orders"`
}

// MarketDepth represents 5-level bids and asks
type MarketDepth struct {
	Bids []DepthLevel `json:"bids"`
	Asks []DepthLevel `json:"asks"`
}

// Tick represents an unpacked market feed tick
type Tick struct {
	ResponseCode    uint8           `json:"responseCode"`
	MessageLength   uint16          `json:"messageLength"`
	ExchangeSegment uint8           `json:"exchangeSegment"`
	SecurityID      uint32          `json:"securityId"`
	LTP             decimal.Decimal `json:"ltp"`
	LastTradedTime  time.Time       `json:"lastTradedTime"`
	LastTradedQty   uint32          `json:"lastTradedQty"`
	AvgPrice        decimal.Decimal `json:"avgPrice"`
	Volume          uint32          `json:"volume"`
	TotalBuyQty     uint32          `json:"totalBuyQty"`
	TotalSellQty    uint32          `json:"totalSellQty"`
	Open            decimal.Decimal `json:"open"`
	High            decimal.Decimal `json:"high"`
	Low             decimal.Decimal `json:"low"`
	Close           decimal.Decimal `json:"close"`
	OpenInterest    uint32          `json:"openInterest"`
	Depth           MarketDepth     `json:"depth"`
}

// ParseBinaryTick decodes a single Dhan binary feed packet
func ParseBinaryTick(raw []byte) (*Tick, error) {
	if len(raw) < 8 {
		return nil, fmt.Errorf("%w: packet shorter than 8-byte header: got %d bytes", ErrInvalidInput, len(raw))
	}

	reader := bytes.NewReader(raw)

	var code uint8
	var msgLen uint16
	var segment uint8
	var securityID uint32

	if err := binary.Read(reader, binary.LittleEndian, &code); err != nil {
		return nil, err
	}
	if err := binary.Read(reader, binary.LittleEndian, &msgLen); err != nil {
		return nil, err
	}
	if err := binary.Read(reader, binary.LittleEndian, &segment); err != nil {
		return nil, err
	}
	if err := binary.Read(reader, binary.LittleEndian, &securityID); err != nil {
		return nil, err
	}

	if len(raw) < int(msgLen) {
		return nil, fmt.Errorf("%w: packet truncated: expected %d bytes, got %d", ErrInvalidInput, msgLen, len(raw))
	}

	tick := &Tick{
		ResponseCode:    code,
		MessageLength:   msgLen,
		ExchangeSegment: segment,
		SecurityID:      securityID,
	}

	switch code {
	case 11: // LTP Packet (8 bytes payload)
		var ltpFloat float32
		var lttUnix uint32
		if err := binary.Read(reader, binary.LittleEndian, &ltpFloat); err != nil {
			return nil, err
		}
		if err := binary.Read(reader, binary.LittleEndian, &lttUnix); err != nil {
			return nil, err
		}
		tick.LTP = decimal.NewFromFloat32(ltpFloat)
		tick.LastTradedTime = time.Unix(int64(lttUnix), 0)

	case 12: // Quote Packet (50 bytes payload)
		var ltp, atp, open, high, low, close float32
		var ltq, ltt, vol, buyQty, sellQty, oi uint32

		_ = binary.Read(reader, binary.LittleEndian, &ltp)
		_ = binary.Read(reader, binary.LittleEndian, &ltq)
		_ = binary.Read(reader, binary.LittleEndian, &ltt)
		_ = binary.Read(reader, binary.LittleEndian, &atp)
		_ = binary.Read(reader, binary.LittleEndian, &vol)
		_ = binary.Read(reader, binary.LittleEndian, &buyQty)
		_ = binary.Read(reader, binary.LittleEndian, &sellQty)
		_ = binary.Read(reader, binary.LittleEndian, &open)
		_ = binary.Read(reader, binary.LittleEndian, &high)
		_ = binary.Read(reader, binary.LittleEndian, &low)
		_ = binary.Read(reader, binary.LittleEndian, &close)
		_ = binary.Read(reader, binary.LittleEndian, &oi)

		tick.LTP = decimal.NewFromFloat32(ltp)
		tick.LastTradedQty = ltq
		tick.LastTradedTime = time.Unix(int64(ltt), 0)
		tick.AvgPrice = decimal.NewFromFloat32(atp)
		tick.Volume = vol
		tick.TotalBuyQty = buyQty
		tick.TotalSellQty = sellQty
		tick.Open = decimal.NewFromFloat32(open)
		tick.High = decimal.NewFromFloat32(high)
		tick.Low = decimal.NewFromFloat32(low)
		tick.Close = decimal.NewFromFloat32(close)
		tick.OpenInterest = oi

	case 13: // Full Depth Packet (100 bytes payload: 5 bids + 5 asks)
		bids := make([]DepthLevel, 5)
		asks := make([]DepthLevel, 5)

		for i := 0; i < 5; i++ {
			var price float32
			var qty uint32
			var orders uint16
			_ = binary.Read(reader, binary.LittleEndian, &price)
			_ = binary.Read(reader, binary.LittleEndian, &qty)
			_ = binary.Read(reader, binary.LittleEndian, &orders)
			bids[i] = DepthLevel{
				Price:    decimal.NewFromFloat32(price),
				Quantity: qty,
				Orders:   orders,
			}
		}

		for i := 0; i < 5; i++ {
			var price float32
			var qty uint32
			var orders uint16
			_ = binary.Read(reader, binary.LittleEndian, &price)
			_ = binary.Read(reader, binary.LittleEndian, &qty)
			_ = binary.Read(reader, binary.LittleEndian, &orders)
			asks[i] = DepthLevel{
				Price:    decimal.NewFromFloat32(price),
				Quantity: qty,
				Orders:   orders,
			}
		}

		tick.Depth = MarketDepth{
			Bids: bids,
			Asks: asks,
		}

	default:
		// Unsupported or custom code (e.g. heartbeat 14/15/21)
	}

	return tick, nil
}
