const express = require('express');
const mongoose = require('mongoose');
const cors = require('cors');

// Pastikan file model ini ada
const Lokasi = require('./models/lokasi');

const app = express();
const PORT = process.env.PORT || 3000;

app.use(cors());
app.use(express.json());

// --- BAGIAN KONEKSI DATABASE (SUDAH DIPERBAIKI) ---
// Kita ganti localhost dengan link MongoDB Atlas Anda.
// Saya tambahkan nama database '/ujianSIG' di akhir link agar rapi.
const connectionString = 'mongodb+srv://sriwahyuni_db_user:EgZ2GXRliZQ1TYA7@cluster23.gnmjc2n.mongodb.net/ujianSIG';

mongoose.connect(connectionString)
  .then(() => console.log('✅ BERHASIL Konek ke MongoDB Atlas (Cloud)'))
  .catch(err => console.log('❌ Gagal Terkoneksi ke MongoDB:', err));

// --- ROUTES ---

// 1. GET (Ambil Data)
app.get('/api/lokasi', async (req, res) => {
    try {
        if (req.query.id) {
            const data = await Lokasi.findById(req.query.id);
            res.json(data);
        } else {
            const data = await Lokasi.find();
            res.json(data);
        }
    } catch (err) {
        res.status(500).json({ message: err.message });
    }
});

// 2. POST (Tambah Data)
app.post('/api/lokasi', async (req, res) => {
    try {
        const lokasiBaru = new Lokasi(req.body);
        const saved = await lokasiBaru.save();
        res.status(201).json(saved);
    } catch (err) {
        res.status(400).json({ message: err.message });
    }
});

// 3. PUT (Update Data)
app.put('/api/lokasi', async (req, res) => {
    try {
        const updated = await Lokasi.findByIdAndUpdate(req.query.id, req.body, { new: true });
        res.json(updated);
    } catch (err) {
        res.status(400).json({ message: err.message });
    }
});

// 4. DELETE (Hapus Data)
app.delete('/api/lokasi', async (req, res) => {
    try {
        await Lokasi.findByIdAndDelete(req.query.id);
        res.json({ message: "Berhasil dihapus" });
    } catch (err) {
        res.status(500).json({ message: err.message });
    }
});

app.listen(PORT, () => console.log(`🚀 Server jalan di http://localhost:${PORT}`));