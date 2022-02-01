package car

import (
	"strconv"
	"strings"

	"github.com/google/uuid"

	"github.com/amehrotra/car-dealership/errors"
	"github.com/amehrotra/car-dealership/filters"
	"github.com/amehrotra/car-dealership/models"
	"github.com/amehrotra/car-dealership/stores"
)

type service struct {
	engine stores.Engine
	car    stores.Car
}

func New(engine stores.Engine, car stores.Car) service {
	return service{engine: engine, car: car}
}

// Create validates car information and sends data to store
func (s service) Create(car models.Car) (models.Car, error) {
	if err := checkCar(car); err != nil {
		return models.Car{}, err
	}

	if err := checkEngine(car.Engine); err != nil {
		return models.Car{}, err
	}

	id := uuid.New()
	car.ID = id
	car.Engine.ID = id

	id, err := s.engine.Create(car.Engine)
	if err != nil {
		return models.Car{}, err
	}

	engine, err := s.engine.GetByID(id)
	if err != nil {
		return models.Car{}, err
	}

	car.Engine = engine

	id, err = s.car.Create(car)
	if err != nil {
		return models.Car{}, err
	}

	car, err = s.GetByID(id)
	if err != nil {
		return models.Car{}, err
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
func (s service) GetByID(id uuid.UUID) (models.Car, error) {
	car, err := s.car.GetByID(id)
	if err != nil {
		return models.Car{}, err
	}

	engine, err := s.engine.GetByID(id)
	if err != nil {
		return models.Car{}, err
	}

	car.Engine = engine

	return car, nil
}

// Update updates the engine followed by car
func (s service) Update(car models.Car) (models.Car, error) {
	err := s.engine.Update(car.Engine)
	if err != nil {
		return models.Car{}, err
	}

	err = checkCar(car)
	if err != nil {
		return models.Car{}, err
	}

	err = s.car.Update(car)
	if err != nil {
		return models.Car{}, err
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
func checkCar(car models.Car) error {
	switch {
	case car.Model == "":
		return errors.InvalidParam{Param: car.Model}
	case car.ManufactureYear < 1866 || car.ManufactureYear > 2022:
		return errors.InvalidParam{Param: strconv.Itoa(car.ManufactureYear)}
	case checkBrand(car.Brand) != nil:
		return errors.InvalidParam{Param: car.Brand}
	case car.FuelType < 0 || car.FuelType > 3:
		return errors.InvalidParam{}
	default:
		return nil
	}
}

// checkEngine validates engine configuration
func checkEngine(engine models.Engine) error {
	switch {
	case engine.Displacement > 0 && engine.NCylinder > 0 && engine.Range > 0:
		return errors.InvalidParam{Param: strconv.Itoa(engine.Displacement)}
	case engine.Displacement < 0 && engine.NCylinder < 0 && engine.Range < 0:
		return errors.InvalidParam{Param: strconv.Itoa(engine.Range)}
	case engine.Displacement == 0 && engine.NCylinder == 0 && engine.Range == 0:
		return errors.InvalidParam{Param: strconv.Itoa(engine.NCylinder)}
	default:
		return nil
	}
}

// checkBrand validates the brand name
func checkBrand(brand string) error {
	brands := map[string]bool{"tesla": true, "porsche": true, "bmw": true, "mercedes": true, "ferrari": true}
	if _, ok := brands[strings.ToLower(brand)]; !ok {
		return errors.InvalidParam{}
	}

	return nil
}
