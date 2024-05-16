-- Insert sample categories
INSERT INTO category (name) VALUES
                                ('SUV'),
                                ('Sedan'),
                                ('Truck'),
                                ('Coupe');

-- Insert sample cars
INSERT INTO cars (model, brand, year, price, color, isUsed, userId, categoryName) VALUES
                                                                                      ('Range Rover', 'Land Rover', 2020, 80000, 'Black', FALSE, 46, 'SUV')

