CREATE DATABASE IF NOT EXISTS toko_kue;
USE toko_kue;

CREATE TABLE IF NOT EXISTS kue (
    id INT AUTO_INCREMENT PRIMARY KEY,
    nama VARCHAR(100) NOT NULL,
    produksi_harian INT NOT NULL,
    harga_terakhir DECIMAL(10, 2) NOT NULL,
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