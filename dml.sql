-- DML (Data Manipulation Language) - Seed Data
-- Run this after DDL to populate initial data

-- Insert initial car data
INSERT INTO cars (name, availability, stock_availability, rental_costs, category, description, image_url)
VALUES
    -- SEDAN (Economic & Comfortable)
    ('Toyota Camry', true, 5, 50.00, 'Sedan', 'Comfortable mid-size sedan perfect for city driving', 'https://images.unsplash.com/photo-1621007947382-bb3c3994e3fb?q=80&w=800'),
    ('Honda Civic', true, 8, 45.00, 'Sedan', 'Reliable and fuel-efficient compact sedan', 'https://images.unsplash.com/photo-1596733430284-f7437764b1a9?q=80&w=800'),
    ('Mazda 6', true, 3, 55.00, 'Sedan', 'Premium feel with exceptional handling', 'https://images.unsplash.com/photo-1552519507-da3b142c6e3d?q=80&w=800'),
    ('Audi A4', true, 2, 85.00, 'Sedan', 'German engineering at its finest', 'https://images.unsplash.com/photo-1606148644561-9c98471c998c?q=80&w=800'),
    
    -- SUV (Versatile & Strong)
    ('Honda CR-V', true, 3, 75.00, 'SUV', 'Spacious SUV with excellent fuel economy', 'https://images.unsplash.com/photo-1568844293986-8d0400bd4745?q=80&w=800'),
    ('Toyota Fortuner', true, 4, 95.00, 'SUV', 'Rugged 7-seater SUV for all terrains', 'https://images.unsplash.com/photo-1533473359331-0135ef1b58bf?q=80&w=800'),
    ('Range Rover Evoque', true, 2, 180.00, 'SUV', 'Compact luxury SUV with iconic design', 'https://images.unsplash.com/photo-1506521781263-d8422e82f27a?q=80&w=800'),
    ('Mazda CX-5', true, 4, 70.00, 'SUV', 'Stylish crossover with premium interior', 'https://images.unsplash.com/photo-1621255554313-05f3e9618037?q=80&w=800'),
    ('Jeep Wrangler', true, 3, 110.00, 'SUV', 'The ultimate off-road adventure vehicle', 'https://images.unsplash.com/photo-1533473359331-0135ef1b58bf?q=80&w=800'),
    ('Toyota Land Cruiser', true, 1, 250.00, 'SUV', 'Indestructible luxury off-roader', 'https://images.unsplash.com/photo-1605515298946-d062f2e9da53?q=80&w=800'),
    
    -- LUXURY (Prestige & Status)
    ('BMW 3 Series', true, 2, 120.00, 'Luxury', 'Premium luxury sedan with sporty performance', 'https://images.unsplash.com/photo-1555215695-3004980ad54e?q=80&w=800'),
    ('Mercedes-Benz C-Class', true, 2, 130.00, 'Luxury', 'The pinnacle of luxury and comfort in its class', 'https://images.unsplash.com/photo-1618843479313-40f8afb4b4d8?q=80&w=800'),
    ('Rolls Royce Ghost', true, 1, 1200.00, 'Luxury', 'The ultimate expression of automotive luxury', 'https://images.unsplash.com/photo-1631214503951-375126d4704b?q=80&w=800'),
    ('Bentley Continental GT', true, 1, 950.00, 'Luxury', 'Hand-crafted luxury grand tourer', 'https://images.unsplash.com/photo-1621135802920-133df287f89c?q=80&w=800'),
    
    -- SPORTS (Speed & Thrills)
    ('Ford Mustang', true, 2, 150.00, 'Sports', 'Iconic American muscle car with a V8 engine', 'https://images.unsplash.com/photo-1584345604482-8135a2153242?q=80&w=800'),
    ('Porsche 911', true, 1, 350.00, 'Sports', 'Legendary sports car with precision handling', 'https://images.unsplash.com/photo-1503376780353-7e6692767b70?q=80&w=800'),
    ('Nissan GT-R', true, 1, 300.00, 'Sports', 'The Godzilla of the road, unmatched performance', 'https://images.unsplash.com/photo-1614162692292-7ac56d7f7f1e?q=80&w=800'),
    ('Lamborghini Huracan', true, 1, 1500.00, 'Sports', 'V10 power with stunning Italian design', 'https://images.unsplash.com/photo-1511919884226-fd3cad34687c?q=80&w=800'),
    
    -- ELECTRIC (Sustainable Future)
    ('Tesla Model 3', true, 4, 100.00, 'Electric', 'All-electric sedan with cutting-edge tech', 'https://images.unsplash.com/photo-1560958089-b8a1929cea89?q=80&w=800'),
    ('Hyundai IONIQ 5', true, 3, 110.00, 'Electric', 'Futuristic electric crossover with fast charging', 'https://images.unsplash.com/photo-1669023030485-573b6a75aa64?q=80&w=800'),
    ('Porsche Taycan', true, 1, 280.00, 'Electric', 'Electric performance that only Porsche can deliver', 'https://images.unsplash.com/photo-1614200179396-2bdb77ebf81b?q=80&w=800'),
    ('Lucid Air', true, 2, 220.00, 'Electric', 'Unmatched range and sophisticated luxury', 'https://images.unsplash.com/photo-1633513364214-4a6547169001?q=80&w=800'),
    
    -- MPV & MINIVAN (Family First)
    ('Toyota Innova', true, 10, 65.00, 'MPV', 'The ultimate family car for long road trips', 'https://images.unsplash.com/photo-1605810230434-7631ac76ec81?q=80&w=800'),
    ('Mitsubishi Xpander', true, 6, 55.00, 'MPV', 'Stylish and practical 7-seater for the family', 'https://images.unsplash.com/photo-1533473359331-0135ef1b58bf?q=80&w=800'),
    ('Kia Carnival', true, 3, 90.00, 'MPV', 'Spacious luxury people mover', 'https://images.unsplash.com/photo-1549317661-bd32c8ce0db2?q=80&w=800'),
    
    -- HATCHBACK (Compact & Fun)
    ('Toyota Yaris', true, 12, 35.00, 'Hatchback', 'Perfect compact car for tight city streets', 'https://images.unsplash.com/photo-1629897048514-3dd7414fe72a?q=80&w=800'),
    ('Volkswagen Golf', true, 5, 45.00, 'Hatchback', 'The classic versatile hatchback', 'https://images.unsplash.com/photo-1541899481282-d53bffe3c35d?q=80&w=800'),
    ('Mini Cooper', true, 3, 65.00, 'Hatchback', 'Fun, stylish, and iconic city car', 'https://images.unsplash.com/photo-1526726538690-5cbf95642cb4?q=80&w=800'),
    
    -- PICKUP TRUCK (Working Hard)
    ('Ford F-150', true, 2, 120.00, 'Pickup', 'The best-selling truck, capable of anything', 'https://images.unsplash.com/photo-1583121274602-3e2820c69888?q=80&w=800'),
    ('Toyota Hilux', true, 5, 85.00, 'Pickup', 'Famous for its legendary durability', 'https://images.unsplash.com/photo-1594270141959-1588122046bc?q=80&w=800')
ON CONFLICT DO NOTHING;

-- Insert initial users (Passwords are 'password123' hashed)
INSERT INTO users (email, password, deposit_amount, role)
VALUES
    ('admin@rentalcar.com', '$2a$10$8K9O./8yT4F.mHshI6kXWOfS4Z.X9z0Zp3FjM3gK7xL5W5w5w5w5w', 1000.00, 'admin'),
    ('user@test.com', '$2a$10$8K9O./8yT4F.mHshI6kXWOfS4Z.X9z0Zp3FjM3gK7xL5W5w5w5w5w', 500.00, 'user')
ON CONFLICT DO NOTHING;

-- Insert rental histories
INSERT INTO rental_histories (user_id, car_id, rental_date, return_date, total_cost, status, payment_status)
VALUES
    ((SELECT id FROM users WHERE email = 'user@test.com'), (SELECT id FROM cars WHERE name = 'Toyota Camry'), NOW() - INTERVAL '5 days', NOW() - INTERVAL '2 days', 150.00, 'completed', 'paid'),
    ((SELECT id FROM users WHERE email = 'user@test.com'), (SELECT id FROM cars WHERE name = 'Honda CR-V'), NOW() - INTERVAL '10 days', NOW() - INTERVAL '7 days', 225.00, 'completed', 'paid'),
    ((SELECT id FROM users WHERE email = 'user@test.com'), (SELECT id FROM cars WHERE name = 'Tesla Model 3'), NOW() - INTERVAL '3 days', NOW() - INTERVAL '1 day', 200.00, 'completed', 'paid'),
    ((SELECT id FROM users WHERE email = 'user@test.com'), (SELECT id FROM cars WHERE name = 'Ford Mustang'), NOW() - INTERVAL '1 day', NULL, 450.00, 'ongoing', 'paid'),
    ((SELECT id FROM users WHERE email = 'admin@rentalcar.com'), (SELECT id FROM cars WHERE name = 'BMW 3 Series'), NOW() - INTERVAL '7 days', NOW() - INTERVAL '5 days', 240.00, 'completed', 'paid'),
    ((SELECT id FROM users WHERE email = 'admin@rentalcar.com'), (SELECT id FROM cars WHERE name = 'Porsche 911'), NOW() - INTERVAL '2 days', NULL, 700.00, 'ongoing', 'paid')
ON CONFLICT DO NOTHING;

-- Insert top-up transactions
INSERT INTO top_up_transactions (user_id, amount, status, payment_method)
VALUES
    ((SELECT id FROM users WHERE email = 'user@test.com'), 500.00, 'completed', 'credit_card'),
    ((SELECT id FROM users WHERE email = 'user@test.com'), 300.00, 'completed', 'bank_transfer'),
    ((SELECT id FROM users WHERE email = 'user@test.com'), 1000.00, 'completed', 'e_wallet'),
    ((SELECT id FROM users WHERE email = 'admin@rentalcar.com'), 2000.00, 'completed', 'bank_transfer'),
    ((SELECT id FROM users WHERE email = 'admin@rentalcar.com'), 500.00, 'pending', 'credit_card')
ON CONFLICT DO NOTHING;

-- Insert user sessions
INSERT INTO user_sessions (user_id, token, expires_at)
VALUES
    ((SELECT id FROM users WHERE email = 'user@test.com'), 'eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoidXNlciIsImV4cCI6OTk5OTk5OTk5OX0.dummy_token_for_user', NOW() + INTERVAL '24 hours'),
    ((SELECT id FROM users WHERE email = 'admin@rentalcar.com'), 'eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoiYWRtaW4iLCJleHAiOjk5OTk5OTk5OTl9.dummy_token_for_admin', NOW() + INTERVAL '24 hours')
ON CONFLICT DO NOTHING;
