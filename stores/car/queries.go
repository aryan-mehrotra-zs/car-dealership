package car

const (
	insertCar        = "INSERT INTO cars (id,model,year_of_manufacture,brand,fuel_type,engine_id) VALUES (?,?,?,?,?,?)"
	getCars          = "SELECT * FROM cars;"
	getCarsWithBrand = "SELECT * FROM cars WHERE brand=?;"
	getCar           = "SELECT * FROM cars WHERE id = ?;"
	updateCar        = "UPDATE cars SET `model=?,year_of_manufacture=?,brand=?,fuel_type=?,engine_id=?` WHERE id=?"
	deleteCar        = "DELETE FROM cars WHERE id = ?;"
)
