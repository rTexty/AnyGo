package main

import "fmt"

// 1. Целевой интерфейс (Target)
type SmartDevice interface {
	TurnOn()
	TurnOff()
}

// 2. Адаптируемая структура (Adaptee) - несовместима с SmartDevice
type LegacyTV struct {
	brand string
}

func (tv *LegacyTV) ApplyVoltage() {
	fmt.Printf("Подано напряжение на %s. Лампы нагреваются, телевизор работает.\n", tv.brand)
}

func (tv *LegacyTV) CutVoltage() {
	fmt.Printf("Напряжение отключено. Экран %s с треском погас.\n", tv.brand)
}

type LegacyTVAdatper struct{
	legacyTV* LegacyTV
}

func (wrap *LegacyTVAdatper) TurnOn() {
	wrap.legacyTV.ApplyVoltage()
	fmt.Println("Обертка включен", wrap.legacyTV)
}
func (wrap *LegacyTVAdatper) TurnOff() {
	
	wrap.legacyTV.CutVoltage()
	fmt.Println("Телевизор выключен", wrap.legacyTV)
}
