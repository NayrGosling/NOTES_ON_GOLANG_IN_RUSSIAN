Больше информации - https://refactoringguru.cn/ru/design-patterns/catalog

1. Паттерн "Фабрика" (Factory)
Что это: Способ создавать объекты, не указывая их точный тип напрямую.
Пример с машинами: представь, что у тебя есть фабрика, которая производит машины. Ты говоришь «сделай машину», а она решает, красную или синюю.
package main
import "fmt"
// Машина
type Car struct {
    color string
}
// Фабрика машин
func NewCar(color string) *Car {
    return &Car{color: color}
}
func main() {
    redCar := NewCar("красная")
    blueCar := NewCar("синяя")
    fmt.Println("У меня есть", redCar.color, "машина и", blueCar.color, "машина!")
}
Вывод: «У меня есть красная машина и синяя машина!»

2. Паттерн "Одиночка" (Singleton)
Что это: Гарантирует, что будет только один экземпляр объекта.
Пример с яблоками: у тебя есть одно большое яблоко, и все дети делят его, а не берут каждый своё.
package main
import "fmt"
type BigApple struct {
    size int
}
var oneApple *BigApple
func GetBigApple() *BigApple {
    if oneApple == nil {
        oneApple = &BigApple{size: 10}
    }
    return oneApple
}
func main() {
    apple1 := GetBigApple()
    apple2 := GetBigApple()
    fmt.Println("Яблоко 1 размер:", apple1.size)
    fmt.Println("Яблоко 2 размер:", apple2.size)
    fmt.Println("Это одно и то же яблоко!")
}
Вывод: «Размер яблока 1: 10», «Размер яблока 2: 10», «Это одно и то же яблоко!»

3. Паттерн "Наблюдатель" (Observer)
Что это такое: один объект (например, машина) сообщает другим (например, детям), что что-то случилось.
Пример с машинами: Машина сигналит, и дети бегут посмотреть.
package main
import "fmt"
type Car struct {
    watchers []func()
}
func (c *Car) AddWatcher(w func()) {
    c.watchers = append(c.watchers, w)
}
func (c *Car) Honk() {
    fmt.Println("Машина сигналит: БИП-БИП!")
    for _, w := range c.watchers {
        w()
    }
}
func main() {
    car := &Car{}
    car.AddWatcher(func() { fmt.Println("Ребёнок 1: Я слышу машину!") })
    car.AddWatcher(func() { fmt.Println("Ребёнок 2: Пойду посмотрю!") })
    car.Honk()
}
Вывод:
 «Машина сигналит: БИП-БИП!»
 «Ребёнок 1: Я слышу машину!»
 «Ребёнок 2: Пойду посмотрю!»

4. Паттерн "Стратегия" (Strategy)
Что это: Способ менять поведение объекта, как выбирать, что делать.
Пример с яблоками: ты можешь съесть яблоко целиком или разрезать его на кусочки.
package main
import "fmt"
type EatingStrategy interface {
    Eat()
}
type WholeApple struct{}
func (w *WholeApple) Eat() { fmt.Println("Съел яблоко целиком!") }
type SlicedApple struct{}
func (s *SlicedApple) Eat() { fmt.Println("Съел яблоко кусочками!") }
type Kid struct {
    strategy EatingStrategy
}
func (k *Kid) SetStrategy(s EatingStrategy) {
    k.strategy = s
}
func (k *Kid) EatApple() {
    k.strategy.Eat()
}
func main() {
    kid := &Kid{}
    kid.SetStrategy(&WholeApple{})
    kid.EatApple()
    kid.SetStrategy(&SlicedApple{})
    kid.EatApple()
}
Вывод:
"Съел яблоко целиком!"
"Съел яблоко кусочками!"

5. Паттерн "Декоратор" (Decorator)
Что это такое: добавляет новые возможности объекту, не изменяя его основную структуру.
Пример с машинами: ты берёшь простую машину и добавляешь к ней крутые колёса или громкий гудок.
package main
import "fmt"
type Car interface {
    Drive() string
}
type BasicCar struct{}
func (b *BasicCar) Drive() string {
    return "Еду на простой машине"
}
type CoolWheelsCar struct {
    car Car
}
func (c *CoolWheelsCar) Drive() string {
    return c.car.Drive() + " с крутыми колёсами!"
}
func main() {
    car := &BasicCar{}
    coolCar := &CoolWheelsCar{car: car}
    fmt.Println(car.Drive())
    fmt.Println(coolCar.Drive())
}
Вывод:
«Еду на простой машине»
«Еду на простой машине с крутыми колёсами!»

6. Паттерн "Команда" (Command)
Что это такое: упаковывает действие в объект, чтобы его можно было выполнить позже или передать.
Пример с яблоками: ты говоришь «съешь яблоко», и это задание можно выполнить сейчас или потом.
package main
import "fmt"
type Command interface {
    Execute()
}
type EatAppleCommand struct {
    name string
}
func (e *EatAppleCommand) Execute() {
    fmt.Println(e.name, "съел яблоко!")
}
type Kid struct {
    command Command
}
func (k *Kid) DoCommand() {
    k.command.Execute()
}
func main() {
    kid := &Kid{}
    eatCommand := &EatAppleCommand{name: "Вася"}
    kid.command = eatCommand
    kid.DoCommand()
}
Вывод:
"Вася съел яблоко!"

7. Паттерн "Адаптер" (Adapter)
Что это делает: позволяет двум разным объектам работать вместе, даже если они несовместимы.
Пример с машинами: у тебя есть игрушечная машинка с круглыми колёсами, а ты хочешь поставить квадратные колёса от другой игрушки.
package main
import "fmt"
// Старая машина с круглыми колёсами
type OldCar struct{}
func (o *OldCar) RollRoundWheels() string {
    return "Катаю круглые колёса"
}
// Новая машина должна использовать квадратные колёса
type Car interface {
    Roll() string
}
// Адаптер для старой машины
type OldCarAdapter struct {
    oldCar *OldCar
}
func (a *OldCarAdapter) Roll() string {
    return a.oldCar.RollRoundWheels() + ", но теперь как квадратные!"
}
func main() {
    oldCar := &OldCar{}
    adapter := &OldCarAdapter{oldCar: oldCar}
    fmt.Println(adapter.Roll())
}
Вывод:
 «Катаю круглые колёса, но теперь они квадратные!»

8. Паттерн "Цепочка обязанностей" (Chain of Responsibility)
Что это такое: передаёт запрос по цепочке объектов, пока кто-нибудь его не обработает.
Пример с яблоками: ты спрашиваешь детей, кто хочет яблоко, пока кто-нибудь не скажет «я!».
package main
import "fmt"
type Kid struct {
    name     string
    nextKid  *Kid
}
func (k *Kid) AskForApple() {
    if k.name == "Маша" {
        fmt.Println(k.name, "сказала: Я хочу яблоко!")
    } else if k.nextKid != nil {
        fmt.Println(k.name, "сказал: Не хочу.")
        k.nextKid.AskForApple()
    }
}
func main() {
    kid1 := &Kid{name: "Вася"}
    kid2 := &Kid{name: "Петя"}
    kid3 := &Kid{name: "Маша"}
    kid1.nextKid = kid2
    kid2.nextKid = kid3

    kid1.AskForApple()
}
Вывод:
 «Вася сказал: «Не хочу».»
 «Петя сказал: «Не хочу».»
 «Маша сказала: «Я хочу яблоко!»»
