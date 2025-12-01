document.addEventListener('DOMContentLoaded', () => {
    // URL Backend Lokal Kita
    // Ganti ke link backend GO yang baru
    const apiUrl = 'https://gis-backend-golang.vercel.app';
    // const apiUrl = 'http://localhost:3000/api/lokasi';

    // Inisialisasi Peta Leaflet (Default Jakarta)
    const map = L.map('map').setView([-6.200000, 106.816666], 12); 
    L.tileLayer('https://{s}.tile.openstreetmap.org/{z}/{x}/{y}.png', {
        attribution: '&copy; OpenStreetMap contributors'
    }).addTo(map);

    // Ambil Elemen HTML
    const form = document.getElementById('lokasi-form');
    const lokasiIdInput = document.getElementById('lokasi-id');
    const namaInput = document.getElementById('nama');
    const kategoriInput = document.getElementById('kategori');
    const deskripsiInput = document.getElementById('deskripsi');
    const longitudeInput = document.getElementById('longitude');
    const latitudeInput = document.getElementById('latitude');
    const submitBtn = document.getElementById('submit-btn');
    const cancelBtn = document.getElementById('cancel-btn');
    const lokasiList = document.getElementById('lokasi-list');

    let markers = {}; 

    // Fungsi Ambil Data dari Backend
    const fetchLokasi = async () => {
        try {
            const response = await fetch(apiUrl);
            const lokasis = await response.json() || [];
            
            // Bersihkan data lama
            lokasiList.innerHTML = '';
            Object.values(markers).forEach(marker => map.removeLayer(marker));
            markers = {};

            // Isi data baru
            lokasis.forEach(lokasi => {
                addLokasiToList(lokasi);
                addMarkerToMap(lokasi);
            });
        } catch (error) {
            console.error('Gagal mengambil data:', error);
        }
    };

    // Tambah ke List Kanan
    const addLokasiToList = (lokasi) => {
        const li = document.createElement('li');
        li.dataset.id = lokasi._id;
        li.innerHTML = `
            <div class="list-item-info">
                <strong>${lokasi.nama}</strong>
                <span>${lokasi.kategori}</span>
            </div>
            <div class="action-buttons">
                <button class="edit-btn" title="Edit">✏️</button>
                <button class="delete-btn" title="Hapus">🗑️</button>
            </div>
        `;
        lokasiList.appendChild(li);
    };

    // Tambah Marker ke Peta
    const addMarkerToMap = (lokasi) => {
        // Mongo nyimpan [Long, Lat], Leaflet butuh [Lat, Long] -> Kita balik
        const [lng, lat] = lokasi.koordinat.coordinates;
        const marker = L.marker([lat, lng])
            .addTo(map)
            .bindPopup(`<b>${lokasi.nama}</b><br>${lokasi.deskripsi}`);
        
        markers[lokasi._id] = marker;
    };

    // Reset Form
    const resetForm = () => {
        form.reset();
        lokasiIdInput.value = '';
        submitBtn.textContent = 'Tambah Lokasi';
        cancelBtn.classList.add('hidden');
    };

    // Isi Form untuk Edit
    const fillFormForEdit = (lokasi) => {
        lokasiIdInput.value = lokasi._id;
        namaInput.value = lokasi.nama;
        kategoriInput.value = lokasi.kategori;
        deskripsiInput.value = lokasi.deskripsi;
        longitudeInput.value = lokasi.koordinat.coordinates[0];
        latitudeInput.value = lokasi.koordinat.coordinates[1];
        submitBtn.textContent = 'Update Lokasi';
        cancelBtn.classList.remove('hidden');
    };

    // Klik Peta -> Isi Koordinat
    map.on('click', (e) => {
        longitudeInput.value = e.latlng.lng.toFixed(6);
        latitudeInput.value = e.latlng.lat.toFixed(6);
    });

    // Simpan Data (Tambah/Update)
    form.addEventListener('submit', async (e) => {
        e.preventDefault();
        const id = lokasiIdInput.value;
        const data = {
            nama: namaInput.value,
            kategori: kategoriInput.value,
            deskripsi: deskripsiInput.value,
            koordinat: {
                type: 'Point',
                coordinates: [parseFloat(longitudeInput.value), parseFloat(latitudeInput.value)]
            }
        };

        const method = id ? 'PUT' : 'POST';
        const url = id ? `${apiUrl}?id=${id}` : apiUrl;

        try {
            const response = await fetch(url, {
                method: method,
                headers: { 'Content-Type': 'application/json' },
                body: JSON.stringify(data)
            });
            if (!response.ok) throw new Error('Gagal Simpan');
            resetForm();
            fetchLokasi();
        } catch (error) {
            alert('Gagal menyimpan data.');
        }
    });

    // Handle Klik List (Edit/Hapus)
    lokasiList.addEventListener('click', async (e) => {
        const li = e.target.closest('li');
        const id = li.dataset.id;

        if (e.target.classList.contains('delete-btn')) {
            if (confirm('Hapus lokasi ini?')) {
                await fetch(`${apiUrl}?id=${id}`, { method: 'DELETE' });
                fetchLokasi();
            }
        } else if (e.target.classList.contains('edit-btn')) {
            const response = await fetch(`${apiUrl}?id=${id}`);
            const lokasi = await response.json();
            fillFormForEdit(lokasi);
        } else {
            if (markers[id]) {
                map.flyTo(markers[id].getLatLng(), 15);
                markers[id].openPopup();
            }
        }
    });
    
    cancelBtn.addEventListener('click', resetForm);
    fetchLokasi(); // Load data saat pertama buka
});