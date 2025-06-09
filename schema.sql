CREATE DATABASE IF NOT EXISTS toko_kue;
USE toko_kue;

CREATE TABLE IF NOT EXISTS kue (
    id INT AUTO_INCREMENT PRIMARY KEY,
    nama VARCHAR(100) NOT NULL,
    produksi_harian INT NOT NULL,
    keuntungan_diinginkan DECIMAL(10, 2) NOT NULL,
    harga_terakhir DECIMAL(10, 2) NOT NULL
);

CREATE TABLE IF NOT EXISTS bahan (
    id INT AUTO_INCREMENT PRIMARY KEY,
    nama VARCHAR(100) NOT NULL,
    harga DECIMAL(10, 2) NOT NULL
);

CREATE TABLE IF NOT EXISTS kue_bahan (
    kue_id INT,
    bahan_id INT,
    jumlah INT NOT NULL,
    PRIMARY KEY (kue_id, bahan_id),
    FOREIGN KEY (kue_id) REFERENCES kue(id) ON DELETE CASCADE,
    FOREIGN KEY (bahan_id) REFERENCES bahan(id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS tenaga_kerja (
    id INT AUTO_INCREMENT PRIMARY KEY,
    nama VARCHAR(100) NOT NULL,
    biaya_harian DECIMAL(10, 2) NOT NULL
)

CREATE TABLE IF NOT EXISTS riwayat_hpp (
    id INT AUTO_INCREMENT PRIMARY KEY,
    kue_id INT,
    harga DECIMAL(10,2),
    created_at DATETIME
)

-- Insert data into kue table
INSERT INTO kue (nama, produksi_harian, keuntungan_diinginkan, harga_terakhir) VALUES
('Nastar', 100, 30000.00, 150000.00),
('Black Forest', 25, 50000.00, 250000.00),
('Lapis Legit', 15, 75000.00, 300000.00),
('Donat', 200, 2000.00, 5000.00),
('Brownies', 50, 20000.00, 80000.00);

-- Insert data into bahan table
INSERT INTO bahan (nama, harga) VALUES
('Tepung Terigu', 12000.00),
('Gula Pasir', 15000.00),
('Telur', 25000.00),
('Mentega', 18000.00),
('Coklat', 35000.00),
('Keju', 30000.00),
('Nanas', 10000.00),
('Susu', 20000.00);

-- Insert data into kue_bahan table
INSERT INTO kue_bahan (kue_id, bahan_id, jumlah) VALUES
(1, 1, 500), -- Nastar: Tepung
(1, 2, 200), -- Nastar: Gula
(1, 3, 8),   -- Nastar: Telur
(1, 4, 250), -- Nastar: Mentega
(1, 7, 3),   -- Nastar: Nanas
(2, 1, 400), -- Black Forest: Tepung
(2, 2, 300), -- Black Forest: Gula
(2, 3, 10),  -- Black Forest: Telur
(2, 5, 300), -- Black Forest: Coklat
(2, 8, 500); -- Black Forest: Susu

-- Insert data into tenaga_kerja table
INSERT INTO tenaga_kerja (nama, biaya_harian) VALUES
('Ahmad', 120000.00),
('Budi', 150000.00),
('Cindy', 135000.00),
('Diana', 140000.00),
('Eko', 125000.00);

-- Insert data into riwayat_hpp table
INSERT INTO riwayat_hpp (kue_id, harga, created_at) VALUES
(1, 125000.00, '2023-01-15 10:00:00'),
(1, 130000.00, '2023-02-20 11:30:00'),
(1, 150000.00, '2023-03-25 09:15:00'),
(2, 230000.00, '2023-01-10 13:45:00'),
(2, 250000.00, '2023-03-05 14:20:00'),
(3, 280000.00, '2023-02-01 08:30:00'),
(3, 300000.00, '2023-03-15 09:45:00');