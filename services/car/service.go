package car

import (
	"strings"

	"github.com/google/uuid"

	"github.com/amehrotra/car-dealership/errors"
	"github.com/amehrotra/car-dealership/filters"
	"github.com/amehrotra/car-dealership/models"
	"github.com/amehrotra/car-dealership/services"
	"github.com/amehrotra/car-dealership/stores"
)

type service struct {
	engine stores.Engine
	car    stores.Car
}

func New(engine stores.Engine, car stores.Car) services.Car {
	return service{engine: engine, car: car}
}

// Create validates car information and sends data to store
func (s service) Create(car *models.Car) (*models.Car, error) {
	if err := checkCar(car); err != nil {
		return nil, err
	}

	if err := checkEngine(car.Engine); err != nil {
		return nil, err
	}

	id := uuid.New()
	car.ID = id
	car.Engine.ID = id

	id, err := s.engine.Create(&car.Engine)
	if err != nil {
		return nil, err
	}

	engine, err := s.engine.GetByID(id)
	if err != nil {
		return nil, err
	}

	car.Engine = engine

	id, err = s.car.Create(car)
	if err != nil {
		return nil, err
	}

	car, err = s.GetByID(id)
	if err != nil {
		return nil, err
	}

	return car, nil
}

// GetAll based on filter extracts data from store about cars
func (s service) GetAll(filter filters.Car) ([]models.Car, error) {
	// validate brand from filter
	switch {
	case filter.Brand == "":
		break
	case checkBrand(filter.Brand) != nil:
		return nil, errors.InvalidParam{}
	default:
		break
	}

	cars, err := s.car.GetAll(filter)
	if err != nil {
		return nil, err
	}

	if filter.Engine {
		for i, car := range cars {
			engine, err := s.engine.GetByID(car.ID)
			if err != nil {
				return nil, err
			}

			cars[i].Engine = engine
		}
	}

	return cars, nil
}

// GetByID based on ID provided extracts data from store about car
func (s service) GetByID(id uuid.UUID) (*models.Car, error) {
	car, err := s.car.GetByID(id)
	if err != nil {
		return nil, err
	}

	engine, err := s.engine.GetByID(id)
	if err != nil {
		return nil, err
	}

	car.Engine = engine

	return &car, nil
}

// Update updates the engine followed by car
func (s service) Update(car *models.Car) (*models.Car, error) {
	err := s.engine.Update(&car.Engine)
	if err != nil {
		return nil, err
	}

	err = checkCar(car)
	if err != nil {
		return nil, err
	}

	err = s.car.Update(car)
	if err != nil {
		return nil, err
	}

	return car, nil
}

// Delete deletes the car from store
func (s service) Delete(id uuid.UUID) error {
	err := s.car.Delete(id)
	if err != nil {
		return err
	}

	err = s.engine.Delete(id)
	if err != nil {
		return err
	}

	return nil
}

// checkCar validates the all parameters of the car
func checkCar(car *models.Car) error {
	switch {
	case car.Model == "":
		return errors.InvalidParam{Param: []string{"model"}}
	case car.ManufactureYear < 1866 || car.ManufactureYear > 2022:
		return errors.InvalidParam{Param: []string{"yearOfManufacture"}}
	case checkBrand(car.Brand) != nil:
		return errors.InvalidParam{Param: []string{"brand"}}
	case car.FuelType < 0 || car.FuelType > 3:
		return errors.InvalidParam{Param: []string{"fuelType"}}
	default:
		return nil
	}
}

// checkEngine validates engine configuration
func checkEngine(engine models.Engine) error {
	equal := func(x int) bool {
		return x == 0
	}

	greater := func(x int) bool {
		return x > 0
	}

	if validate(engine, equal) || validate(engine, greater) {
		return errors.InvalidParam{Param: []string{"displacement", "noOfCylinder", "range"}}
	}

	params := make([]string, 0)

	switch {
	case engine.Displacement < 0:
		params = append(params, "displacement")
	case engine.NCylinder < 0:
		params = append(params, "noOfCylinder")
	case engine.Range < 0:
		params = append(params, "range")
	}

	if len(params) > 0 {
		return errors.InvalidParam{Param: params}
	}

	return nil
}

func validate(e models.Engine, compare func(int) bool) bool {
	return compare(e.Displacement) && compare(e.NCylinder) && compare(e.Range)
}

// checkBrand validates the brand name
func checkBrand(brand string) error {
	brands := map[string]bool{"tesla": true, "porsche": true, "bmw": true, "mercedes": true, "ferrari": true}
	if _, ok := brands[strings.ToLower(brand)]; !ok {
		return errors.InvalidParam{}
	}

	return nil
}
