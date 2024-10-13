CREATE TABLE Dishes (
    name VARCHAR(30) UNIQUE NOT NULL ,
    calorie INT,
    protein INT,
    fat INT,
    carbohydrate INT
);

CREATE TABLE SelectedDished (
    name VARCHAR(30) REFERENCES Dishes(name),
    date DATE DEFAULT CURRENT_DATE
)