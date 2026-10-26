package godhan

import (
	"github.com/DragoHex/godhan/models"
	"context"
	"encoding/csv"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"

	"github.com/shopspring/decimal"
)

// DefaultScripMasterURL is the public URL for Dhan's daily security master CSV dump
const DefaultScripMasterURL = "https://images.dhan.co/api-data/api-scrip-master.csv"

// models.Instrument represents a tradable security master record from DhanHQ

// ParseInstruments parses a CSV stream of Dhan security master records row-by-row
func ParseInstruments(r io.Reader) ([]models.Instrument, error) {
	reader := csv.NewReader(r)
	reader.FieldsPerRecord = -1 // Allow variable fields if trailing commas exist
	reader.ReuseRecord = true   // Minimize GC allocations

	header, err := reader.Read()
	if err != nil {
		return nil, fmt.Errorf("failed to read instrument CSV header: %w", err)
	}

	colMap := make(map[string]int)
	for i, col := range header {
		colMap[strings.TrimSpace(col)] = i
	}

	exchIdx, hasExch := colMap["SEM_EXM_EXCH_ID"]
	segIdx, hasSeg := colMap["SEM_SEGMENT"]
	secIdx, hasSec := colMap["SEM_SMST_SECURITY_ID"]
	symIdx, hasSym := colMap["SEM_TRADING_SYMBOL"]
	lotIdx, hasLot := colMap["SEM_LOT_UNITS"]
	instIdx := colMap["SEM_INSTRUMENT_NAME"]
	expIdx := colMap["SEM_EXPIRY_DATE"]
	custIdx := colMap["SEM_CUSTOM_SYMBOL"]
	strikeIdx, hasStrike := colMap["SEM_STRIKE_PRICE"]

	if !hasExch || !hasSeg || !hasSec || !hasSym || !hasLot {
		return nil, fmt.Errorf("%w: missing required headers in instrument CSV", ErrInvalidInput)
	}

	var instruments []models.Instrument
	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("error reading CSV line: %w", err)
		}

		lotSize, _ := strconv.ParseInt(strings.TrimSpace(record[lotIdx]), 10, 32)
		if lotSize <= 0 {
			lotSize = 1
		}

		var strike decimal.Decimal
		if hasStrike && strikeIdx < len(record) {
			strike, _ = decimal.NewFromString(strings.TrimSpace(record[strikeIdx]))
		}

		inst := models.Instrument{
			Exchange:      record[exchIdx],
			Segment:       record[segIdx],
			SecurityID:    record[secIdx],
			TradingSymbol: record[symIdx],
			LotSize:       int32(lotSize),
			StrikePrice:   strike,
		}

		if instIdx < len(record) {
			inst.InstrumentName = record[instIdx]
		}
		if expIdx < len(record) {
			inst.ExpiryDate = record[expIdx]
		}
		if custIdx < len(record) {
			inst.CustomSymbol = record[custIdx]
		}

		instruments = append(instruments, inst)
	}

	return instruments, nil
}

// GetInstruments streams and parses the official DhanHQ daily security master CSV
func (c *Client) GetInstruments(ctx context.Context) ([]models.Instrument, error) {
	if c.configErr != nil {
		return nil, c.configErr
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, DefaultScripMasterURL, nil)
	if err != nil {
		return nil, err
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("%w: failed to download instrument master: %v", ErrBrokerUnavailable, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, parseErrorResponse(resp.StatusCode, nil)
	}

	return ParseInstruments(resp.Body)
}
