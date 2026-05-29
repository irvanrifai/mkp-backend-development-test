# Analisis dan Penjelasan System Design - Platform Tiket Bioskop Nasional

Dokumen ini menjelaskan pendekatan arsitektur sistem untuk menangani tantangan *high-concurrency* (lonjakan trafik massal), manajemen ketersediaan kursi, serta skenario penanganan kegagalan transaksi dan pembatalan sepihak oleh bioskop.

---

## 1. Solusi Sistem Pemilihan Tempat Duduk (*High Concurrency Seat Locking*)

### Tantangan Utama

Pada platform bioskop skala nasional, momen perilisan film populer (*box office premier*) sering kali memicu lonjakan trafik ekstrem. Ribuan pengguna dapat menekan tombol untuk memilih nomor kursi yang sama (misalnya kursi `E11`) di milidetik yang sama. Jika backend langsung melakukan operasi baca-tulis ke database relational (PostgreSQL), akan terjadi *database locking/bottleneck* yang memperlambat sistem, atau bahkan menyebabkan *double booking* (satu kursi terbeli oleh dua orang).

### Solusi Arsitektur

Sistem ini memisahkan lapisan penanganan kuncian kursi sementara (*temporary lock*) dari database utama dengan memanfaatkan **Redis Distributed Caching**.

1. **In-Memory Seat Validation:** Saat pengguna memilih kursi dan menekan tombol "Lanjutkan ke Pembayaran", backend tidak langsung memeriksa PostgreSQL. Backend akan mengecek status kuncian kursi tersebut di dalam Redis cluster menggunakan operasi atomik `SETNX` (Set if Not Exists) atau algoritma **Redlock**.
2. **TTL-Based Temporary Locking:** Jika kursi masih kosong, Redis akan membuat kuncian sementara berupa key-value pair dengan mekanisme **TTL (Time-to-Live) selama 10 menit**.
* *Key Format:* `lock:schedule:{schedule_id}:seat:{seat_id}`
* Proses ini berjalan di memori (*in-memory*), sehingga latensi pengecekan sangat rendah meskipun diakses puluhan ribu orang secara bersamaan.


3. **Isolation Principle:** Pengguna lain yang mencoba memilih kursi yang sama dalam rentang waktu tersebut akan langsung ditolak oleh backend berdasarkan kuncian di Redis, tanpa sempat membebani PostgreSQL dengan query *read/write*.

---

## 2. Mekanisme Pencatatan Transaksi dan Restok Tiket Otomatis

Pencatatan data dilakukan secara terstruktur demi menjaga performa penulisan database utama tetap stabil:

### A. Alur Pencatatan (Checkout Sukses)

Begitu kursi berhasil dikunci di Redis, backend segera melakukan operasi penulisan ke database PostgreSQL dalam satu kesatuan blok transaksi database (*Database Transaction / ACID Validation*):

1. Membuat baris baru pada tabel `transactions` dengan status `payment_status = 'PENDING'`.
2. Mencatat daftar nomor kursi yang dipesan ke tabel `ticket_seats`.
3. Menghasilkan nomor invoice unik dan meminta *Payment Gateway* untuk menerbitkan instruksi pembayaran (seperti Virtual Account atau QRIS).

### B. Alur Restok Tiket Otomatis

Sistem tiket ini dirancang agar **tidak perlu menghapus data** pada tabel `ticket_seats` ketika transaksi batal atau kedaluwarsa. Proses *restock* (pengembalian tiket ke pasar) berjalan secara pasif dan aktif melalui dua mekanisme:

1. **Mekanisme Pasif (Redis TTL Timeout):** Jika pengguna tidak menyelesaikan pembayaran dalam waktu 10 menit, kunci sementara di Redis akan otomatis terhapus (*expired*).
2. **Mekanisme Aktif (Cron Job/Worker & PG Callback):** Sebuah *background worker* di backend akan memeriksa tabel `transactions` secara berkala atau menerima *callback* kegagalan langsung dari *Payment Gateway*. Status transaksi yang melewati batas waktu akan diperbarui dari `PENDING` menjadi `EXPIRED` atau `FAILED`.

**Mengapa Kursi Otomatis Ter-restock?**
Karena query pengecekan ketersediaan kursi oleh pengguna baru di aplikasi hanya menyaring transaksi yang memiliki status `PENDING` atau `SUCCESS`. Ketika transaksi lama berubah menjadi `EXPIRED`, kursi tersebut secara otomatis **tidak akan tersaring** lagi oleh query pencarian, yang berarti kursi tersebut langsung berstatus **tersedia kembali untuk dipilih publik detik itu juga**. Histori transaksi gagal tetap aman tersimpan di database untuk kebutuhan audit internal IT tanpa mengorbankan fungsi bisnis aplikasi.

### C. Penanganan State Transaksi dan Retensi Data Third-Party (Fintech Best Practice)

Untuk menjaga efisiensi kinerja jaringan dan menghindari batasan limitasi panggilan API (Rate Limiting) dari server Payment Gateway, platform ini menerapkan prinsip *State Persistence*:

1. **Local Data Persistence:** Saat pengguna sukses melakukan checkout, backend langsung meminta instansiasi kode bayar/Virtual Account ke Payment Gateway, lalu menyimpan metadata tersebut (`payment_method`, `payment_reference`, dan `payment_url`) secara permanen ke database lokal.
2. **Anti-Polling Principle:** Ketika pengguna sempat meninggalkan aplikasi lalu kembali lagi ke halaman peninjauan pembayaran, aplikasi frontend cukup melakukan query baca lokal (`SELECT`) ke database internal untuk menampilkan nomor VA atau QR Code yang sudah ada. Backend tidak akan melakukan pemanggilan sinkronus berulang (re-hit) ke API Payment Gateway demi menghemat bandwith jaringan server.
3. **Event-Driven Reconciliation:** Pembaruan status pembayaran sepenuhnya dikendalikan secara asinkronus berbasis kejadian (Event-Driven) melalui Webhook/Callback yang dikirim oleh Payment Gateway ke server kita ketika pengguna menyelesaikan pembayaran di jaringan bank.

---

## 3. Alur Refund dan Pembatalan dari Pihak Bioskop (*Force Majeure*)

Jika terjadi kendala teknis di lapangan (misalnya proyektor studio rusak, pemadaman listrik, atau bencana alam), pihak bioskop harus membatalkan penayangan secara sepihak dan mengembalikan dana konsumen secara otomatis demi menjaga *Customer Experience*.

Sistem menangani skenario ini melalui arsitektur berbasis kejadian asinkronus (*Asynchronous Event-Driven Flow*):

```
[Admin Bioskop] ──> Ubah Status Schedule ke 'CANCELLED'
                                │
                                ▼
                     [Trigger Async Worker]
                                │
            ┌───────────────────┴─────────────────────────┐
            ▼                                             ▼
 [Scan Tabel Transactions]                      [Kirim Notifikasi Massal]
(Cari Schedules terkait & Status SUCCESS)       (Via Email / WhatsApp API)
            │
            ▼
[Bulk Request Refund ke PG API]
            │
            ▼
[Ubah Status Transaksi ke 'REFUNDED']

```

### Tahapan Eksekusi Sistem:

1. **State Mutation:** Admin bioskop mengubah status jadwal film terkait di dasbor manajemen menjadi `CANCELLED`. Perubahan ini otomatis mengaktifkan penonaktifan jadwal agar tidak ada pengguna baru yang bisa membeli tiket pada sesi tersebut.
2. **Asynchronous Background Processing:** Backend memicu sebuah *worker* (background job) secara asinkronus agar proses pembatalan massal tidak membuat aplikasi admin mengalami *freeze/timeout*.
3. **Data Scanning:** *Worker* akan menyisir tabel `transactions` untuk mencari semua transaksi yang memiliki `transactionable_type = 'schedules'`, `transactionable_id = [ID jadwal terkait]`, dan `payment_status = 'SUCCESS'` (transaksi yang sudah lunas).
4. **Automated PG Refund Interconnect:** Untuk setiap transaksi sukses yang ditemukan, backend akan menembak API *Payment Gateway* (menggunakan *Idempotency-Key* berbasis nomor invoice asli untuk mencegah *double refund*) dengan instruksi pengembalian dana penuh (*Full Refund*) ke metode pembayaran awal pengguna (E-Wallet/Virtual Account).
5. **Audit Trail and Notification:** Setelah *Payment Gateway* sukses memproses refund, backend mengubah status transaksi menjadi `REFUNDED`. Di saat yang bersamaan, sistem secara otomatis mengirimkan notifikasi resmi (melalui WhatsApp API atau Email) kepada pengguna bahwa jadwal tayang dibatalkan dan dana mereka sedang dikembalikan oleh bank pengirim.
