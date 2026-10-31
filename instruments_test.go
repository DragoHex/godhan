package godhan

import (
	"strings"
	"testing"
)

func TestParseInstruments_LotSizeParsing(t *testing.T) {
	csvData := `SEM_EXM_EXCH_ID,SEM_SEGMENT,SEM_SMST_SECURITY_ID,SEM_INSTRUMENT_NAME,SEM_EXPIRY_CODE,SEM_TRADING_SYMBOL,SEM_LOT_UNITS,SEM_CUSTOM_SYMBOL,SEM_EXPIRY_DATE,SEM_STRIKE_PRICE,SEM_OPTION_TYPE,SEM_TICK_SIZE,SEM_EXPIRY_FLAG,SEM_EXCH_INSTRUMENT_TYPE,SEM_SERIES,SM_SYMBOL_NAME
NSE,D,53376,OPTIDX,0,NIFTY-Nov2026-21000-PE,65.0,NIFTY 03 NOV 21000 PUT,2026-11-03 14:30:00,21000.00000,PE,5.0000,W,OP,,
NSE,D,53377,OPTIDX,0,NIFTY-Nov2026-21000-CE,50,NIFTY 03 NOV 21000 CALL,2026-11-03 14:30:00,21000.00000,CE,5.0000,W,OP,,
NSE,D,53378,OPTIDX,0,NIFTY-Nov2026-21000-PE,0.0,NIFTY 03 NOV 21000 PUT,2026-11-03 14:30:00,21000.00000,PE,5.0000,W,OP,,
NSE,E,11536,EQUITY,0,TATASTEEL,1.0,TATA STEEL,,,,5.0000,NA,ES,B,TATA STEEL
`
	insts, err := ParseInstruments(strings.NewReader(csvData))
	if err != nil {
		t.Fatalf("ParseInstruments unexpected error: %v", err)
	}

	if len(insts) != 4 {
		t.Fatalf("expected 4 instruments, got %d", len(insts))
	}

	// 1. Float lot size "65.0" should parse as 65
	if insts[0].LotSize != 65 {
		t.Errorf("insts[0].LotSize = %d, want 65", insts[0].LotSize)
	}
	// 2. Integer lot size "50" should parse as 50
	if insts[1].LotSize != 50 {
		t.Errorf("insts[1].LotSize = %d, want 50", insts[1].LotSize)
	}
	// 3. Zero lot size "0.0" should fall back to 1
	if insts[2].LotSize != 1 {
		t.Errorf("insts[2].LotSize = %d, want 1", insts[2].LotSize)
	}
	// 4. Equity lot size "1.0" should parse as 1
	if insts[3].LotSize != 1 {
		t.Errorf("insts[3].LotSize = %d, want 1", insts[3].LotSize)
	}
}
