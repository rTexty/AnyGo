
    2
    3 import (
    4     "fmt"
    5 )
    6
    7 // 1. Интерфейс Посетителя (Visitor)
    8 type Visitor interface {
    9     VisitBook(b *Book) float64
   10     VisitCoffee(c *Coffee) float64
   11     VisitTea(t *Tea) float64
   12 }
   13
   14 // 2. Интерфейс Элемента (Element)
   15 type Item interface {
   16     Accept(v Visitor) float64
   17 }
   18
   19 // Конкретные элементы
   20 type Book struct {
   21     isbn  string
   22     price float64
   23     tax   float64
   24 }
   25
   26 func (b *Book) Accept(v Visitor) float64 { return v.VisitBook(b) }
   27
   28 type Coffee struct {
   29     brand    string
   30     price    float64
   31     tax      float64
   32     discount bool
   33 }
   34
   35 func (c *Coffee) Accept(v Visitor) float64 { return v.VisitCoffee(c) }
   36
   37 type Tea struct {
   38     brand    string
   39     price    float64
   40     tax      float64
58     return cost 59 } 60
   62     cost := t.price + (t.tax * t.price)
   63     if t.discount {
   64         cost -= cost * 0.1
   65     }
   66     return cost
   67 }
   68
   69 func main() {
   70     items := []Item{
   71         &Book{"1234", 20.01, 0.08},
   72         &Book{"5678", 345.0, 0.08},
   73         &Coffee{"Espresso", 300.0, 0.092, false},
   74         &Coffee{"Starbucks", 400.0, 0.099, true},
   75         &Tea{"Curtis", 50.0, 0.003, true},
   76     }
   77
   78     visitor := &StoreVisitor{}
   79     var totalCost float64
   80
   81     for _, item := range items {
   82         totalCost += item.Accept(visitor)
   83     }
   84
   85     fmt.Printf("Total cost = %.2f\n", totalCost)
   86 }
