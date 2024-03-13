package main

import (
	"fmt"
	"time"
)

type RawData struct {
	Type            string
	UserId          int
	Side            string
	RefDate         string
	ProductTypeName string
	SecurityName    string
	Quantity        int
	Price           int
}

type Broker struct {
	Name     string
	BrokerId int
}

type CutOffDate struct {
	Trades    time.Time
	Treasury  time.Time
	Statement time.Time
	Transfer  time.Time
}

type BrokerInterface interface {
	GetTrades(cutOffDate CutOffDate) []RawData
	GetTransfer(cutOffDate CutOffDate) []RawData
}

type BrokerA struct {
	Broker
}

type BrokerB struct {
	Broker
}

func (broker BrokerA) GetTrades(cutOffDate CutOffDate) []RawData {
	return []RawData{
		{Type: "Trade", UserId: 1, Side: "Buy", RefDate: "2024-03-12", ProductTypeName: "Stock", SecurityName: "PETR4", Quantity: 10, Price: 100},
		{Type: "Trade", UserId: 2, Side: "Buy", RefDate: "2024-03-13", ProductTypeName: "Stock", SecurityName: "PETR4", Quantity: 20, Price: 200},
	}
}

func (broker BrokerA) GetTransfer(cutOffDate CutOffDate) []RawData {
	return []RawData{
		{Type: "Transfer", UserId: 1, Side: "Buy", RefDate: "2024-03-12", ProductTypeName: "Stock", SecurityName: "PETR4", Quantity: 30, Price: 300},
		{Type: "Transfer", UserId: 2, Side: "Sell", RefDate: "2024-03-13", ProductTypeName: "Stock", SecurityName: "PETR4", Quantity: 40, Price: 400},
	}
}

func (broker BrokerB) GetTrades(cutOffDate CutOffDate) []RawData {
	return []RawData{
		{Type: "Trade", UserId: 1, Side: "Buy", RefDate: "2024-03-12", ProductTypeName: "Stock", SecurityName: "ITSA4", Quantity: 10, Price: 100},
		{Type: "Trade", UserId: 2, Side: "Sell", RefDate: "2024-03-13", ProductTypeName: "Stock", SecurityName: "ITSA4", Quantity: 20, Price: 200},
	}
}

func (broker BrokerB) GetTransfer(cutOffDate CutOffDate) []RawData {
	return []RawData{
		{Type: "Transfer", UserId: 1, Side: "Buy", RefDate: "2024-03-12", ProductTypeName: "Stock", SecurityName: "ITSA4", Quantity: 30, Price: 300},
		{Type: "Transfer", UserId: 2, Side: "Sell", RefDate: "2024-03-13", ProductTypeName: "Stock", SecurityName: "ITSA4", Quantity: 40, Price: 400},
	}
}

func main() {
	brokerMap := make(map[int]func() BrokerInterface)
	brokerMap[123] = func() BrokerInterface { return BrokerA{Broker{BrokerId: 628, Name: "BrokerA"}} }
	brokerMap[321] = func() BrokerInterface { return BrokerB{Broker{BrokerId: 666, Name: "BrokerB"}} }
	brokerId := 123

	if brokerFunc, ok := brokerMap[brokerId]; ok {
		broker := brokerFunc()
		cutOffDate := CutOffDate{
			Trades:    time.Date(2024, 3, 12, 0, 0, 0, 0, time.UTC),
			Treasury:  time.Date(2024, 3, 11, 0, 0, 0, 0, time.UTC),
			Statement: time.Date(2024, 3, 10, 0, 0, 0, 0, time.UTC),
			Transfer:  time.Date(2024, 3, 9, 0, 0, 0, 0, time.UTC),
		}
		fmt.Println("Trades:", broker.GetTrades(cutOffDate))
		fmt.Println("Transfers:", broker.GetTransfer(cutOffDate))
	}
}
