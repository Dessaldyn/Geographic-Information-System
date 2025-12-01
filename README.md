🗺️ Sistem Informasi Geografis (SIG) - Web GIS

Aplikasi pemetaan interaktif berbasis web untuk mengelola data lokasi spasial (Latitude & Longitude). Dibangun dengan arsitektur Fullstack JavaScript yang memisahkan Frontend dan Backend secara modular.

🛠️ Teknologi yang Digunakan (Tech Stack)

Project ini dibangun menggunakan teknologi modern MERN Stack (Minus React, diganti Vanilla JS + Leaflet):

Frontend

Backend

Database

✨ Fitur Utama

✅ Peta Interaktif: Menggunakan OpenStreetMap & Leaflet.js.
✅ Input Lokasi Otomatis: Klik di peta, koordinat (Lat/Long) langsung terisi.
✅ CRUD Data:

Create: Tambah titik lokasi baru.

Read: Menampilkan semua marker yang tersimpan di Database.

Update: Edit data lokasi.

Delete: Hapus titik lokasi.
✅ Database Cloud: Terintegrasi dengan MongoDB Atlas.
✅ Responsive Design: Tampilan rapi di Desktop maupun Mobile.

🚀 Cara Menjalankan (Instalasi Lokal)

Ikuti langkah ini untuk menjalankan proyek di komputer kamu:

1. Clone Repository
```
git clone [https://github.com/username-kamu/nama-repo.git](https://github.com/username-kamu/nama-repo.git)
cd nama-repo
```

2. Setup Backend

Masuk ke folder backend dan install dependency:
```
cd backend
npm install
```

Jalankan server:
```
npm start
# Server akan berjalan di http://localhost:3000
```

3. Setup Frontend

Buka folder frontend.

Edit file main.js, pastikan apiUrl mengarah ke backend lokal:

const apiUrl = 'http://localhost:3000/api/lokasi';


Buka file index.html di browser (atau gunakan Live Server di VS Code).

🔌 API Endpoints (Dokumentasi Backend)

Method

Endpoint

Deskripsi
```
GET

/api/lokasi
```
Mengambil semua data lokasi (GeoJSON).
```
GET

/api/lokasi?id={id}
```
Mengambil satu data lokasi berdasarkan ID.
```
POST

/api/lokasi
```
Menambahkan lokasi baru.
```
PUT

/api/lokasi?id={id}
```
Mengupdate data lokasi.
```
DELETE

/api/lokasi?id={id}
```
Menghapus data lokasi.

📂 Struktur Folder

📦 Project-Root
 ┣ 📂 backend          # Server Side (Node.js & Express)
 ┃ ┣ 📂 models         # Schema Database Mongoose
 ┃ ┣ 📜 server.js      # Entry point server
 ┃ ┗ 📜 package.json
 ┗ 📂 frontend         # Client Side (Interface)
   ┣ 📜 index.html     # Halaman Utama
   ┣ 📜 style.css      # Styling (Modern UI)
   ┗ 📜 main.js        # Logic Leaflet & Fetch API


👤 Author

Dibuat untuk memenuhi tugas Sistem Informasi Geografis.

Nama: Hafizh Fakhri Muharram

NPM : 714230031

GitHub: @Dessaldyn

Jangan lupa kasih ⭐️ (Star) jika project ini bermanfaat!