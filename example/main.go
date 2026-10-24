package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/DragoHex/godhan"
	"github.com/DragoHex/godhan/models"
	"github.com/DragoHex/godhan/ticker"
	"github.com/shopspring/decimal"
)

func main() {
	ctx := context.Background()

	// 1. Initialize DhanHQ Client
	client := godhan.NewClient("1000000001", "your-access-token")

	// 2. Query Account Profile & Balances
	profile, err := client.GetUserProfile(ctx)
	if err == nil {
		fmt.Printf("[+] Dhan Client Profile: %s\n", profile.DhanClientID)
	}

	funds, err := client.GetFundLimits(ctx)
	if err == nil {
		fmt.Printf("[+] Available Cash Balance: %s\n", funds.AvailableBalance)
	}

	// 3. Place a Limit Buy Order
	orderResp, err := client.PlaceOrder(ctx, models.OrderRequest{
		TransactionType: "BUY",
		ExchangeSegment: "NSE_EQ",
		ProductType:     "INTRADAY",
		OrderType:       "LIMIT",
		Validity:        "DAY",
		TradingSymbol:   "TATASTEEL",
		SecurityID:      "11536",
		Quantity:        10,
		Price:           decimal.NewFromFloat(120.50),
	})
	if err != nil {
		log.Printf("PlaceOrder (demo fallback): %v", err)
	} else {
		fmt.Printf("[+] Order Placed! ID: %s, Status: %s\n", orderResp.OrderID, orderResp.OrderStatus)
	}

	// 4. Query Open Positions & Demat Holdings
	positions, err := client.GetPositions(ctx)
	if err == nil {
		fmt.Printf("[+] Open Positions: %d\n", len(positions))
	}

	holdings, err := client.GetHoldings(ctx)
	if err == nil {
		fmt.Printf("[+] Custody Holdings: %d\n", len(holdings))
	}

	// 5. Connect to DhanHQ Single-Socket Multiplexed Feed (Quotes + Order Updates)
	tk := ticker.NewTicker("1000000001", "your-access-token")

	tk.OnConnect(func() {
		fmt.Println("[+] Connected to DhanHQ WebSocket Feed!")
		_ = tk.Subscribe([]ticker.TickerSubscription{
			{ExchangeSegment: 1, SecurityID: 11536},
		})
	})

	tk.OnTick(func(tick ticker.Tick) {
		fmt.Printf("[Tick] Security: %d, LTP: %s\n", tick.SecurityID, tick.LTP)
	})

	tk.OnOrderUpdate(func(update ticker.OrderUpdate) {
		fmt.Printf("[Order Alert] %s -> Status: %s, Filled: %d\n",
			update.OrderID, update.OrderStatus, update.FilledQty)
	})

	feedCtx, feedCancel := context.WithTimeout(ctx, 3*time.Second)
	defer feedCancel()

	go func() {
		_ = tk.ServeWithContext(feedCtx)
	}()

	time.Sleep(100 * time.Millisecond)
	_ = tk.Close()
	fmt.Println("[+] DhanHQ Example run complete.")
}
