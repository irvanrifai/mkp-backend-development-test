# Dokumentasi Rancangan Database - Platform Tiket Bioskop Nasional

Dokumen ini menjelaskan struktur data, relasi antar-tabel, serta strategi optimalisasi yang diimplementasikan pada skema database PostgreSQL untuk platform pembelian tiket bioskop nasional.

---

## 1. Ringkasan Arsitektur & Keunggulan Desain

Rancangan database ini dibuat dengan mempertimbangkan beban kerja tinggi (*high-concurrency*) serta fleksibilitas bisnis jangka panjang melalui beberapa strategi utama:

* **Polymorphic Transactions:** Struktur tabel transaksi tidak mengunci secara kaku ke satu entitas (seperti jadwal bioskop), melainkan menggunakan relasi *polymorphic* agar di masa depan platform dapat menjual produk lain (seperti *Food & Beverages*, *Merchandise*, dll) tanpa mengubah skema tabel utama.
* **Enterprise Soft Delete:** Penggunaan kolom `deleted_at` pada tabel master untuk mengamankan data dari penghapusan permanen akibat ketidaksengajaan (*accidental data loss*).
* **Pencatatan Audit Finansial Komprehensif:** Pemisahan komponen biaya (`subtotal`, `admin_fee`, `discount`, `total_price`) untuk mempermudah proses rekonsiliasi keuangan oleh tim *finance*.

---

## 2. Kamus Data & Penjelasan Tabel

### A. Tabel `users`

Menyimpan informasi data akun pengguna/pelanggan.

* **`username` & `email`:** Diatur dengan *constraint* `UNIQUE` untuk memastikan otentikasi bersifat tunggal dan mencegah duplikasi akun. Kolom `username` mempermudah variasi login pada aplikasi *mobile*.
* **`deleted_at`:** Mengimplementasikan *soft delete* agar histori penonaktifan akun pelanggan tetap tercatat.

### B. Tabel `movies`

Menyimpan data master film yang tersedia di jaringan bioskop.

* **`poster_url`:** Menyimpan referensi tautan gambar poster film untuk kebutuhan visualisasi di frontend aplikasi.
* **`duration_minutes`:** Menggunakan tipe data `INT` untuk menyimpan durasi dalam hitungan menit guna memudahkan kalkulasi waktu jeda antar-jadwal tayang di backend.

### C. Tabel `studios`

Menyimpan data ruangan studio di setiap cabang bioskop.

* **`branch_name`:** Menyimpan informasi lokasi cabang bioskop (Contoh: "MKP Cinema Grand Indonesia").
* **`studio_type`:** Membedakan tipe studio (Contoh: 'REGULAR', 'IMAX', 'PREMIERE') karena penentuan tarif dasar bioskop di Indonesia umumnya melekat pada jenis fasilitas studio.

### D. Tabel `seats`

Menyimpan denah atau daftar master kursi yang tersedia di setiap studio.

* **`studio_id (ON DELETE CASCADE)`:** Jika sebuah studio dibongkar/dihapus secara permanen dari sistem, maka seluruh master kursi di dalam studio tersebut akan otomatis terhapus secara berantai.
* **`status`:** Mengakomodasi skenario *Broken Seat Management*. Jika kursi mengalami kerusakan fisik, status diubah menjadi `BROKEN/INACTIVE` sehingga backend akan mengecualikan kursi ini dari daftar opsi yang bisa dipilih oleh pengguna.

### E. Tabel `schedules`

Tabel perantara (*Many-to-Many*) yang mempertemukan film, studio, jam tayang, dan harga tiket.

* **`movie_id` & `studio_id` (ON DELETE RESTRICT):** Menggunakan aturan `RESTRICT` karena master data film dan studio menerapkan *soft delete*. Aturan ini mencegah penghapusan data master secara permanen selama data tersebut masih aktif terikat pada jadwal tayang tertentu.
* **`price_per_ticket`:** Menggunakan tipe data `NUMERIC(10, 2)` untuk menjamin akurasi nilai desimal mata uang dan menghindari *floating-point error* saat kalkulasi keuangan.

### F. Tabel `transactions`

Inti dari sistem pencatatan transaksi pembayaran.

* **`user_id (ON DELETE SET NULL)`:** Jika pengguna menghapus akun mereka secara permanen, data transaksi finansial **tidak boleh ikut terhapus** demi kebutuhan audit pembukuan keuangan. Kolom `user_id` otomatis berubah menjadi `NULL` (anonimisasi data), namun angka perputaran uang tetap tercatat utuh.
* **`transactionable_type` & `transactionable_id`:** Struktur *Polymorphic Association*. Kolom ini merekam tipe entitas yang dibeli (misal: 'schedules') beserta ID-nya secara dinamis.
* **`invoice_number`:** Diatur sebagai `UNIQUE`. PostgreSQL otomatis membuatkan indeks B-Tree pada kolom ini untuk mempercepat query pencarian saat menerima data *callback* dari *Payment Gateway*.

### G. Tabel `ticket_seats`

Menyimpan detail kursi bioskop yang dipesan untuk setiap transaksi.

* **`transaction_id (ON DELETE CASCADE)`:** Menjamin jika data transaksi dihapus/dibersihkan, kuncian kursi yang terkait otomatis terlepas.
* **`seat_id (ON DELETE RESTRICT)`:** Mencegah penghapusan master data kursi jika kursi tersebut sedang atau pernah digunakan dalam sebuah transaksi.
* **`UNIQUE(transaction_id, seat_id)`:** *Composite Unique Key* untuk memastikan sebuah kursi tidak dapat diduplikasi dalam satu invoice transaksi yang sama.

---

## 3. Logika Bisnis & Strategi Alur Data

### A. Mekanisme Cek Ketersediaan Kursi & Pool Pending

Sistem menganggap sebuah kursi **"tidak tersedia"** apabila kursi tersebut telah terikat pada transaksi yang berstatus `SUCCESS` (lunas) ATAU berstatus `PENDING` (dalam batas waktu tunggu bayar).

Proses filter ketersediaan kursi dilakukan dengan melakukan `JOIN` antara tabel `ticket_seats` dan `transactions` berdasarkan `schedule_id` terkait:

```sql
SELECT ts.seat_id
FROM ticket_seats ts
JOIN transactions t ON ts.transaction_id = t.id
WHERE t.transactionable_type = 'schedules'
  AND t.transactionable_id = :schedule_id
  AND t.payment_status IN ('PENDING', 'SUCCESS')
  AND t.deleted_at IS NULL;

```

*Hasil dari query di atas kemudian dikurangi dari total master kursi aktif di studio tersebut untuk menampilkan sisa kursi yang tersedia bagi pengguna baru.*

### B. Otomatisasi Restok Tiket (Skenario Timeout & Gagal)

Ketika pengguna melakukan *checkout*, data langsung ditulis ke tabel `transactions` (`PENDING`) dan `ticket_seats`. Sistem **tidak perlu menghapus data** di tabel `ticket_seats` apabila pengguna batal membayar atau waktu pembayaran habis (*timeout*).

Backend/worker cukup memperbarui status pada tabel utama:

```sql
UPDATE transactions SET payment_status = 'EXPIRED' WHERE id = :transaction_id;

```

Secara otomatis, pada query pengecekan ketersediaan kursi berikutnya, transaksi berstatus `EXPIRED` tidak akan tersaring. Hal ini membuat kursi tersebut **otomatis ter-restock (tersedia kembali)** bagi publik, sementara histori kegagalan transaksi tetap tersimpan dengan aman untuk analitik data internal IT.

---

## 4. Strategi Indeksasi & Optimalisasi Performa Query

Untuk menghindari *Full Table Scan* pada operasi yang sering diakses di lingkungan produksi ber-trafik tinggi, dipasang 2 indeks taktis berikut:

1. **`CREATE INDEX idx_schedules_show_time ON schedules(show_time);`**
* *Alasan:* Halaman beranda dan pencarian aplikasi akan terus-menerus melakukan query penyaringan jadwal film berdasarkan parameter waktu hari ini atau akhir pekan (`WHERE show_time BETWEEN ...`). Indeks ini memastikan pencarian jadwal berjalan dalam skala waktu konstan ($O(\log N)$).


2. **`CREATE INDEX idx_transactions_user_status ON transactions(user_id, payment_status);`**
* *Alasan:* Indeks komposit ini mengoptimalkan proses validasi alur *checkout* di backend. Sebelum mengizinkan pengguna memilih kursi baru, sistem dapat mengecek secara instan apakah pengguna tersebut masih memiliki transaksi berjalan yang menggantung (`PENDING`) tanpa membebani performa database.
