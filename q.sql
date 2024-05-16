-- Insert sample categories
INSERT INTO category (name) VALUES
                                ('SUV'),
                                ('Sedan'),
                                ('Truck'),
                                ('Coupe');

-- Insert sample cars
INSERT INTO cars (model, brand, year, price, color, isUsed, userId, categoryName) VALUES
                                                                                      ('Model X', 'Tesla', 2020, 80000, 'Black', FALSE, 45, 'SUV'),
                                                                                      ('Model S', 'Tesla', 2019, 75000, 'White', TRUE, 45, 'Sedan'),
                                                                                      ('F-1', 'F1', 2018, 30000, 'Red', TRUE, 45, 'Sport'),
                                                                                      ('Mustang', 'Ford', 2021, 55000, 'Blue', FALSE, 45, 'Coupe');
