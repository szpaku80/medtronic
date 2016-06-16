package main

import (
	"fmt"
	"log"
	"time"

	"github.com/ecc1/medtronic"
)

func main() {
	pump := medtronic.Open()
	result := pump.CarbRatios()
	if pump.Error() != nil {
		log.Fatal(pump.Error())
	}
	fmt.Printf("%+v\n", result)
	s := result.CarbRatioAt(time.Now())
	fmt.Printf("Current carb ratio: %d\n", s.Carbs)
}