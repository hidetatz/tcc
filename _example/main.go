package main

import (
	"fmt"
	"log"

	"github.com/yagi5/tcc"
)

var (
	db = &FakeDB{
		flight: flight{StockSeatCount: uint64(3)},
		hotel:  hotel{StockRoomCount: uint64(1)},
	}

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

type flight struct {
	StockSeatCount    uint64
	ReservedSeatCount uint64
}

type hotel struct {
	StockRoomCount    uint64
	ReservedRoomCount uint64
}

// FakeDB represents a database for example
type FakeDB struct {
	flight flight
	hotel  hotel
}

func (f *FakeDB) tryReserveFlightSeat() error {
	if f.flight.StockSeatCount == 0 {
		return fmt.Errorf("no seat")
	}
	f.flight.StockSeatCount--
	return nil
}

func (f *FakeDB) confirmFlightSeatReservation() error {
	f.flight.ReservedSeatCount++
	return nil
}

func (f *FakeDB) cancelFlightSeat() error {
	f.flight.StockSeatCount++
	return nil
}

func (f *FakeDB) tryReserveHotelRoom() error {
	if f.hotel.StockRoomCount == 0 {
		return fmt.Errorf("no room")
	}
	f.hotel.StockRoomCount--
	return nil
}

func (f *FakeDB) confirmHotelRoomReservation() error {
	f.hotel.ReservedRoomCount++
	return nil
}

func (f *FakeDB) cancelHotelRoom() error {
	f.hotel.StockRoomCount++
	return nil
}

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
	orchestrator := tcc.NewOrchestrator([]*tcc.Service{flightService, hotelService}, tcc.WithMaxRetries(1))
	err := orchestrator.Orchestrate()
	if err != nil {
		log.Printf("error happened in 2nd reservation: %s", err)
	}
	tccErr := err.(*tcc.Error)
	log.Printf("tccErr.Error: %v", tccErr.Error())
	log.Printf("tccErr.FailedPhase == ErrTryFailed: %v", tccErr.FailedPhase() == tcc.ErrTryFailed)
	log.Printf("tccErr.ServiceName: %v", tccErr.ServiceName())
}
