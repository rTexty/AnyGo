package main

import (
	"fmt"
	"structs"

	"go.mongodb.org/mongo-driver/mongo/address"
)

type ITransport interface {
	deliver()
}

type Truck struct{
	address string  
}

type Ship struct{
	country string 
}

type Plane struct{
	country string
}

func (p *Truck) deliver() {
	fmt.Println("Delivered by Track")
} 

func (p *Truck) deliver() {
	fmt.Println("Delivered by Track")
} 

func (p *Truck) deliver() {
	fmt.Println("Delivered by Track")
} 
