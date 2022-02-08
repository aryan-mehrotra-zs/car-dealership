# car-dealership
It allows storing of data of different cars with their engine specifications. 

# test-coverage
```
handler : 100%
service : 100%
stores
     --car     : 100%
     --engine  : 100%
```

#### linter-check 
```
no error
```

### Server Setup

Start Database
```
docker exec -it customer-api mysql -u root -ppassword organisation
```
Begin Server 
```
go run main.go
```


### DATABASE SETUP

Create Docker Image 
```
docker run --name car_dealership -e MYSQL_ROOT_PASSWORD=password -e MYSQL_DATABASE=car_dealership -p 3306:3306 -d mysql:latest
```

SQL Commands To Create Table
```
CREATE TABLE engines(
id varchar(36) NOT NULL,
displacement INT,
no_of_cylinder INT,
`range` INT,
PRIMARY KEY (id)
);

CREATE TABLE cars(
id varchar(36) NOT NULL,
model varchar(50) NOT NULL,
year_of_manufacture year NOT NULL,
brand varchar(50) NOT NULL,
fuel_type ENUM('petrol','diesel','electric') NOT NULL,
engine_id varchar(36) NOT NULL,
PRIMARY KEY (ID),
FOREIGN KEY (engine_id) REFERENCES engines(id)
);

```