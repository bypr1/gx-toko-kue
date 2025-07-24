Pak Abdul telah menjalankan usaha penjualan kue selama lebih dari satu tahun. Selama ini, seluruh transaksi penjualan dicatat secara manual menggunakan buku tulis. Pendekatan ini awalnya dirasa cukup, namun seiring berjalannya waktu, muncul berbagai kendala dalam pengelolaan data dan pengambilan keputusan bisnis. Buku catatan yang digunakan sering hilang atau terselip, sehingga menyulitkan Pak Abdul dalam memantau perkembangan usaha secara menyeluruh.

Selain itu, Pak Abdul belum memiliki perhitungan yang jelas dan akurat terkait harga pokok produksi (HPP) dan harga jual kuenya. Hal ini membuat proses penetapan harga sering kali dilakukan berdasarkan perkiraan, tanpa analisis yang tepat terhadap biaya bahan baku, tenaga kerja, dan keuntungan yang diinginkan.

- **Permasalahan**
    - Data penjualan dicatat secara manual dan sering kali tidak terdokumentasi dengan baik karena buku catatan mudah hilang atau terselip.
    - Sulitnya melacak riwayat transaksi penjualan dan perkembangan usaha dari waktu ke waktu.
    - Tidak adanya sistem yang dapat membantu menghitung harga pokok produksi secara otomatis.
    - Penentuan harga jual tidak berdasarkan analisis biaya dan margin keuntungan yang jelas.


### Gambaran Flow

```mermaid
flowchart TD
    A[Ingredients] -->|Dibuatkan Resep| D[Cakes]
    C[Costs] --> |Dibebankan| D
    D -->|On Create/Update| E([Update Sell Price])
    E --> D
    D -->|Dibuatkan Transaksi| T[Transactions]
```

### Fitur Sistem
- Input data bahan baku kue. (nama, deskripsi, harga per unit, unit)
- Input data kue (nama, deskripsi, margin(%)) beserta resep dan biayanya.
- Input data penjualan ID Transaksi, jumlah kue, total harga, dan tanggal penjualan.

### Rumus Menghitung Harga Jual
- **Harga Pokok Produksi (HPP)**: HPP = Total Biaya Bahan Baku + Total Biaya Tenaga Kerja + Biaya Overhead
- **Harga Jual**: Harga Jual = HPP + (HPP * Margin

### Struktur Data 
- `Cake`: Struktur data untuk menyimpan informasi kue, termasuk ID kue, nama, deskripsi, margin(%), dan harga jual.
- `CakeComponentIngredient`: Struktur data untuk menyimpan informasi bahan baku, termasuk ID bahan, nama, deskripsi, harga per unit, dan unit.
- `CakeRecipeIngredient`: Struktur data untuk menyimpan informasi resep kue, termasuk ID resep, ID kue, dan daftar bahan baku beserta takarannya.
- `CakeCost`: Struktur data untuk menyimpan informasi ongkos produksi per kue, termasuk ID kue, jenis ongkos dan ongkos.
- `Transaction`: Struktur data untuk menyimpan informasi transaksi penjualan, termasuk ID transaksi, tanggal, dan total harga.
- `TransactionCake`: Struktur data untuk menyimpan detail penjualan, termasuk ID transaksi, ID kue, jumlah yang terjual, dan harga jual per unit.
