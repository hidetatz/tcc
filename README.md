# tcc - Tiny and minimal TCC (Try - Confirm - Cancel) pattern implementation

[![Build Status](https://travis-ci.org/dty1er/tcc.svg?branch=master)](https://travis-ci.org/dty1er/tcc)
[![Coverage Status](https://coveralls.io/repos/github/dty1er/tcc/badge.svg?branch=master)](https://coveralls.io/github/dty1er/tcc?branch=master)
[![Go Report Card](https://goreportcard.com/badge/github.com/dty1er/tcc)](https://goreportcard.com/report/github.com/dty1er/tcc)
[![GoDoc](https://godoc.org/github.com/dty1er/tcc?status.svg)](https://godoc.org/github.com/dty1er/tcc)

### Abstract

In distributed system, transactions over network (_distributed transaction_) is a hard thing.
There are some patterns to achieve distributed transactions, the one of them is called `TCC`.
This library enables to implement TCC easily.

### Usage

Working sample is in [_example](./_example) directory, or you can try the code in [Go Playground](https://play.golang.org/p/w5_ar85dmGx).

```go
package main

import (
	"log"

	"github.com/dty1er/tcc"
)

var (
	flightService = tcc.NewService(
		"flight reservation",
		db.tryReserveFlightSeat,
		db.confirmFlightSeatReservation,
		db.cancelFlightSeat,
	)

	hotelService = tcc.NewService(
		"hotel reservation",
		db.tryReserveHotelRoom,
		db.confirmHotelRoomReservation,
		db.cancelHotelRoom,
	)
)

func main() {
	doFirstReservation(db)
	doSecondReservation(db)
}

func doFirstReservation(db *FakeDB) {
	orchestrator := tcc.NewOrchestrator([]*tcc.Service{flightService, hotelService}, tcc.WithMaxRetries(1))
	err := orchestrator.Orchestrate()
	if err != nil {
		log.Printf("error happened in 1st reservation: %s", err)
	}
}

func doSecondReservation(db *FakeDB) {
	// In second reservation, flight seat is not enough
	// Please refer to working example
	orchestrator := tcc.NewOrchestrator([]*tcc.Service{flightService, hotelService}, tcc.WithMaxRetries(1))
	err := orchestrator.Orchestrate()
	if err != nil {
		log.Printf("error happened in 2nd reservation: %s", err)
	}

	// When error is returned, it can be casted into *tcc.Error
	tccErr := err.(*tcc.Error)
	log.Printf("tccErr.Error: %v", tccErr.Error())
	log.Printf("tccErr.FailedPhase == ErrTryFailed: %v", tccErr.FailedPhase() == tcc.ErrTryFailed)
	log.Printf("tccErr.ServiceName: %v", tccErr.ServiceName())
}
```

### Documents

Described in [GoDoc](https://godoc.org/github.com/dty1er/tcc).

### Ref

References for TCC pattern.

[Eventual Data Consistency Solution in ServiceComb - part 3](https://servicecomb.apache.org/docs/distributed_saga_3/)
[Transactions for the REST of Us](https://dzone.com/articles/transactions-for-the-rest-of-us)

### License

MIT
