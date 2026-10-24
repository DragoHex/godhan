package ticker

import (
	"bytes"
	"encoding/binary"
	"testing"
	"time"

	"github.com/shopspring/decimal"
)

func buildBinaryPacket(code uint8, segment uint8, securityID uint32, payload []byte) []byte {
	buf := new(bytes.Buffer)
	msgLen := uint16(8 + len(payload))
	buf.WriteByte(code)
	binary.Write(buf, binary.LittleEndian, msgLen)
	buf.WriteByte(segment)
	binary.Write(buf, binary.LittleEndian, securityID)
	buf.Write(payload)
	return buf.Bytes()
}

func TestTickerBinary_ParseLTP(t *testing.T) {
	// Build LTP payload: LTP = 1450.50, LTT = 1697356800
	payload := new(bytes.Buffer)
	binary.Write(payload, binary.LittleEndian, float32(1450.50))
	binary.Write(payload, binary.LittleEndian, uint32(1697356800))

	rawPkt := buildBinaryPacket(11, 1, 11536, payload.Bytes())

	tick, err := ParseBinaryTick(rawPkt)
	if err != nil {
		t.Fatalf("ParseBinaryTick failed: %v", err)
	}

	if tick.ResponseCode != 11 {
		t.Errorf("expected response code 11, got %d", tick.ResponseCode)
	}
	if tick.SecurityID != 11536 {
		t.Errorf("expected security ID 11536, got %d", tick.SecurityID)
	}
	expectedLTP := decimal.NewFromFloat32(1450.50)
	if !tick.LTP.Equal(expectedLTP) {
		t.Errorf("expected LTP %v, got %v", expectedLTP, tick.LTP)
	}
	expectedTime := time.Unix(1697356800, 0)
	if !tick.LastTradedTime.Equal(expectedTime) {
		t.Errorf("expected LTT %v, got %v", expectedTime, tick.LastTradedTime)
	}
}

func TestTickerBinary_ParseQuote(t *testing.T) {
	// Build Quote payload
	payload := new(bytes.Buffer)
	binary.Write(payload, binary.LittleEndian, float32(2450.25)) // LTP
	binary.Write(payload, binary.LittleEndian, uint32(50))      // LTQ
	binary.Write(payload, binary.LittleEndian, uint32(1697356800)) // LTT
	binary.Write(payload, binary.LittleEndian, float32(2445.00)) // ATP
	binary.Write(payload, binary.LittleEndian, uint32(100000))   // Volume
	binary.Write(payload, binary.LittleEndian, uint32(25000))    // Total Buy Qty
	binary.Write(payload, binary.LittleEndian, uint32(18000))    // Total Sell Qty
	binary.Write(payload, binary.LittleEndian, float32(2440.00)) // Open
	binary.Write(payload, binary.LittleEndian, float32(2465.00)) // High
	binary.Write(payload, binary.LittleEndian, float32(2435.00)) // Low
	binary.Write(payload, binary.LittleEndian, float32(2442.00)) // Close
	binary.Write(payload, binary.LittleEndian, uint32(5000))     // OI

	rawPkt := buildBinaryPacket(12, 1, 2885, payload.Bytes())

	tick, err := ParseBinaryTick(rawPkt)
	if err != nil {
		t.Fatalf("ParseBinaryTick quote failed: %v", err)
	}

	if tick.ResponseCode != 12 || tick.SecurityID != 2885 {
		t.Errorf("unexpected header: code=%d, secID=%d", tick.ResponseCode, tick.SecurityID)
	}
	if tick.Volume != 100000 || tick.OpenInterest != 5000 {
		t.Errorf("unexpected volume or OI: vol=%d, oi=%d", tick.Volume, tick.OpenInterest)
	}
}

func TestTickerBinary_ParseDepth(t *testing.T) {
	// Build 5-Depth payload (5 bids + 5 asks, each: price float32, qty uint32, orders uint16 = 10 bytes)
	payload := new(bytes.Buffer)
	for i := 1; i <= 5; i++ {
		binary.Write(payload, binary.LittleEndian, float32(100.0-float32(i)*0.05))
		binary.Write(payload, binary.LittleEndian, uint32(i*100))
		binary.Write(payload, binary.LittleEndian, uint16(i*2))
	}
	for i := 1; i <= 5; i++ {
		binary.Write(payload, binary.LittleEndian, float32(100.0+float32(i)*0.05))
		binary.Write(payload, binary.LittleEndian, uint32(i*150))
		binary.Write(payload, binary.LittleEndian, uint16(i*3))
	}

	rawPkt := buildBinaryPacket(13, 1, 500, payload.Bytes())

	tick, err := ParseBinaryTick(rawPkt)
	if err != nil {
		t.Fatalf("ParseBinaryTick depth failed: %v", err)
	}

	if tick.ResponseCode != 13 || len(tick.Depth.Bids) != 5 || len(tick.Depth.Asks) != 5 {
		t.Fatalf("unexpected depth tick: %+v", tick)
	}
	if tick.Depth.Bids[0].Quantity != 100 || tick.Depth.Asks[0].Quantity != 150 {
		t.Errorf("unexpected top of book depth: bids[0]=%+v, asks[0]=%+v", tick.Depth.Bids[0], tick.Depth.Asks[0])
	}
}

func TestTickerBinary_TruncatedPacket(t *testing.T) {
	// Less than 8 bytes header
	shortPkt := []byte{11, 16, 0, 1}
	_, err := ParseBinaryTick(shortPkt)
	if err == nil {
		t.Fatal("expected error for truncated packet, got nil")
	}

	// Truncated payload
	truncated := buildBinaryPacket(11, 1, 100, []byte{1, 2, 3})
	_, err = ParseBinaryTick(truncated)
	if err == nil {
		t.Fatal("expected error for truncated payload, got nil")
	}
}
