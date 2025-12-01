const mongoose = require('mongoose');

const LokasiSchema = new mongoose.Schema({
    nama: { 
        type: String, 
        required: true 
    },
    kategori: { 
        type: String, 
        required: true 
    },
    deskripsi: { 
        type: String 
    },
    koordinat: {
        type: { 
            type: String, 
            enum: ['Point'], 
            required: true,
            default: 'Point'
        },
        coordinates: { 
            type: [Number], 
            required: true 
        }
    }
});

module.exports = mongoose.model('Lokasi', LokasiSchema);