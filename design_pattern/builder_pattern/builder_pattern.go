/**
 * @Author: Sean
 * @Date: 2022/5/17 15:13
 */

package builder_pattern

/**
separate the construction of a complex object from representation
设计思想:
	*Builder interface (包含1.build_XX method 返回的是builder接口, 2.get_XX 返回对象)
	*父struct
	*Director struct,属性为Builder,实现Construct()和SetBuilder()方法
	*不同的子Struct组合,实现接口builder
*/

type Wheels int

// 定义父struct
type Vehicle struct {
	Wheels    Wheels
	Seats     int
	Structure string
}

// builder interface
type Builder interface {
	SetWheels() Builder
	SetSeats() Builder
	SetStructure() Builder
	GetVehicle() Vehicle
}

// Director
type Director struct {
	builder Builder
}

func (director *Director) Construct() {
	director.builder.SetWheels().SetSeats().SetStructure() //链式调用
}

func (director *Director) SetBuilder(builder Builder) {
	director.builder = builder
}

// Car struct
type Car struct {
	vehicle Vehicle
}

// 实现集成builder
func (car *Car) SetWheels() Builder {
	car.vehicle.Wheels = 4
	return car
}
func (car *Car) SetSeats() Builder {
	car.vehicle.Seats = 4
	return car
}
func (car *Car) SetStructure() Builder {
	car.vehicle.Structure = "Car"
	return car
}
func (car *Car) GetVehicle() Vehicle {
	return car.vehicle
}