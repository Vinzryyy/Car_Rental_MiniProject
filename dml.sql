-- DML (Data Manipulation Language) - Seed Data
-- Run this after DDL to populate initial data

-- Insert initial car data
INSERT INTO cars (name, availability, stock_availability, rental_costs, category, description, image_url)
VALUES
    ('Toyota Camry', true, 5, 50.00, 'Sedan', 'Comfortable mid-size sedan perfect for city driving', 'https://example.com/camry.jpg'),
    ('Honda CR-V', true, 3, 75.00, 'SUV', 'Spacious SUV with excellent fuel economy', 'https://example.com/crv.jpg'),
    ('BMW 3 Series', true, 2, 120.00, 'Luxury', 'Premium luxury sedan with sporty performance', 'https://example.com/bmw3.jpg'),
    ('Ford Mustang', true, 2, 150.00, 'Sports', 'Iconic American muscle car', 'https://example.com/mustang.jpg'),
    ('Tesla Model 3', true, 4, 100.00, 'Electric', 'All-electric sedan with autopilot features', 'https://example.com/model3.jpg')
ON CONFLICT DO NOTHING;
