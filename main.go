package main

import (
	"fmt"
	_ "image/png"
	"reflect"
	"strings"
)

type Vehicle struct {
	Manufacturer  Manufacturer
	VehicleNumber string
}

type Manufacturer struct {
	ManufacturerName string
}

func exportValueFromField(data interface{}, index string) string {
	indexArray := strings.Split(index, ".")
	r := reflect.ValueOf(data)
	for _, i := range indexArray {
		if r.FieldByName(i).Kind() == reflect.Struct {
			r = reflect.ValueOf(r.FieldByName(i).Interface())
		} else {
			r = r.FieldByName(i)
		}
	}
	return fmt.Sprintf("%v", r)
}

func main() {

	var vehicle = Vehicle{
		Manufacturer: Manufacturer{
			ManufacturerName: "hello",
		},
	}

	value := exportValueFromField(vehicle, "Manufacturer.ManufacturerName")
	fmt.Println(value)
}

// 70203421
// noK.uqo.3.ixo

// 14016
// specificky simbol 5757
