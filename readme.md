# car-dealership
It allows storing of data of different cars with their engine specifications. 

# coverage
```
handler : 100%
service : 100%
stores
     --car     : 94.9%
     --engine  : 100%
```

#### linter-check `no error`


### DATABASE SETUP

SQL COMMANDS


```
CREATE DATABASE car_dealership

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