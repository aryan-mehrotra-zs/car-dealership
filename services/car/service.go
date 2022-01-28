package car

import (
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

func checkCar(car models.Car) error {
	brands := map[string]bool{"tesla": true, "porsche": true, "bmw": true, "mercedes": true, "ferrari": true}

	if car.Model == "" {
		return errors.InvalidParam{}
	}

	if car.YearOfManufacture < 1866 || car.YearOfManufacture > 2022 {
		return errors.InvalidParam{}
	}

	if _, ok := brands[strings.ToLower(car.Brand)]; !ok {
		return errors.InvalidParam{}
	}

	if car.FuelType < 0 || car.FuelType > 3 {
		return errors.InvalidParam{}
	}

	if car.Engine.Displacement > 0 && car.Engine.NCylinder > 0 && car.Engine.Range > 0 {
		return errors.InvalidParam{}
	}

	if car.Engine.Displacement < 0 && car.Engine.NCylinder < 0 && car.Engine.Range < 0 {
		return errors.InvalidParam{}
	}

	if car.Engine.Displacement == 0 && car.Engine.NCylinder == 0 && car.Engine.Range == 0 {
		return errors.InvalidParam{}
	}

	return nil

}

func (s service) Create(car models.Car) (models.Car, error) {
	err := checkCar(car)
	if err != nil {
		return models.Car{}, err
	}

	id := uuid.New()
	car.ID = id
	car.Engine.ID = id

	id, err = s.engine.Create(car.Engine)
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

func (s service) GetAll(filter filters.Car) ([]models.Car, error) {
	cars, err := s.car.GetAll(filter)
	if err != nil {
		return []models.Car{}, errors.DB{}
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
