-- Modulo de prestamos de laptops

SET NAMES utf8mb4;

-- Catalogo de laptops

CREATE TABLE IF NOT EXISTS laptops (
    id INT AUTO_INCREMENT PRIMARY KEY,
    brand VARCHAR(100) NOT NULL,
    model VARCHAR(100) NOT NULL,
    serial_number VARCHAR(100) NOT NULL UNIQUE,
    total_units INT NOT NULL DEFAULT 1,
    available_units INT NOT NULL DEFAULT 1,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Prestamos de laptops

CREATE TABLE IF NOT EXISTS laptop_loans (
    id INT AUTO_INCREMENT PRIMARY KEY,
    user_id INT NOT NULL,
    laptop_id INT NOT NULL,
    loaned_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    due_at TIMESTAMP NOT NULL,
    returned_at TIMESTAMP NULL,
    status ENUM('active', 'returned') NOT NULL DEFAULT 'active',
    FOREIGN KEY (user_id) REFERENCES users(id),
    FOREIGN KEY (laptop_id) REFERENCES laptops(id)
);

-- Indices para consultas frecuentes

CREATE INDEX idx_laptop_loans_user
ON laptop_loans(user_id);

CREATE INDEX idx_laptop_loans_laptop
ON laptop_loans(laptop_id);

-- Datos de ejemplo

INSERT INTO laptops (
    brand,
    model,
    serial_number,
    total_units,
    available_units
) VALUES
('Dell', 'Latitude 5420', 'DELL-001', 5, 5),
('HP', 'ProBook 440', 'HP-001', 3, 3),
('Lenovo', 'ThinkPad E14', 'LEN-001', 4, 4);
