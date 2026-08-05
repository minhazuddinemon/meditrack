-- Seed Doctor
INSERT INTO Doctor (doc_id, name, specialization, contact)
VALUES (1, 'Dr. Sarah Smith', 'General Physician', '+8801700000000')
ON DUPLICATE KEY UPDATE name = name;

-- Seed Master Medicine Catalog
INSERT INTO Medicine (medicine_name, medicine_type, manufacturer) VALUES
('Napa 500mg', 'Tablet', 'Beximco'),
('Seclo 20mg', 'Capsule', 'Square'),
('Histacin', 'Syrup', 'Incepta')
ON DUPLICATE KEY UPDATE manufacturer = manufacturer;

-- Seed Diagnostic Tests
INSERT INTO Test (test_name, test_fee) VALUES
('CBC (Complete Blood Count)', 400.00),
('Blood Grouping', 150.00),
('Dengue NS1 Antigen', 800.00)
ON DUPLICATE KEY UPDATE test_fee = test_fee;
