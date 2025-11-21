package main

import (
	"fmt"
	"log"
)

func main() {
	mfPath := "samples/millennium-falcon.json"
	empirePath := "samples/empire.json"

	c, err := NewC3PO(mfPath)
	if err != nil {
		log.Fatalf("NewC3PO error: %v", err)
	}

	odds, err := c.GiveMeTheOdds(empirePath)
	if err != nil {
		log.Fatalf("GiveMeTheOdds error: %v", err)
	}

	fmt.Printf("%.3f\n", odds)
}
