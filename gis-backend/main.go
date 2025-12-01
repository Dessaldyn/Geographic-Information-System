package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// --- 1. STRUKTUR DATA (MODEL) ---
// Kita harus mendefinisikan bentuk JSON-nya secara ketat di Go

type GeoJSON struct {
	Type        string    `json:"type" bson:"type"`
	Coordinates []float64 `json:"coordinates" bson:"coordinates"`
}

type Lokasi struct {
	ID        primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Nama      string             `json:"nama" bson:"nama"`
	Kategori  string             `json:"kategori" bson:"kategori"`
	Deskripsi string             `json:"deskripsi" bson:"deskripsi"`
	Koordinat GeoJSON            `json:"koordinat" bson:"koordinat"`
}

var collection *mongo.Collection

// --- 2. KONEKSI DATABASE ---
func connectDB() {
	// Ganti dengan Connection String MongoDB Atlas kamu
	// Tips: Sebaiknya pakai Environment Variable di production
	mongoURI := os.Getenv("MONGO_URI")
	if mongoURI == "" {
		// Fallback untuk lokal jika tidak ada env (Ganti dengan link Atlasmu jika mau run lokal)
		mongoURI = "mongodb+srv://sriwahyuni_db_user:EgZ2GXRliZQ1TYA7@cluster23.gnmjc2n.mongodb.net/ujianSIG"
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(mongoURI))
	if err != nil {
		log.Fatal(err)
	}

	// Cek koneksi
	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatal("Gagal konek MongoDB:", err)
	}

	log.Println("✅ BERHASIL Konek ke MongoDB (via Golang)")
	collection = client.Database("ujianSIG").Collection("lokasis")
}

func main() {
	// Koneksi DB
	connectDB()

	// Setup Router (Gin)
	r := gin.Default()

	// Setup CORS (Agar frontend bisa akses)
	r.Use(cors.Default())

	// --- 3. ROUTES API ---

	// GET: Ambil Data (Semua atau Satu by ID)
	r.GET("/api/lokasi", func(c *gin.Context) {
		idParam := c.Query("id")

		if idParam != "" {
			// Jika ada ?id=xxx, ambil satu
			objID, _ := primitive.ObjectIDFromHex(idParam)
			var lokasi Lokasi
			err := collection.FindOne(context.Background(), bson.M{"_id": objID}).Decode(&lokasi)
			if err != nil {
				c.JSON(http.StatusNotFound, gin.H{"error": "Data tidak ditemukan"})
				return
			}
			c.JSON(http.StatusOK, lokasi)
		} else {
			// Jika tidak ada ID, ambil semua
			cursor, err := collection.Find(context.Background(), bson.M{})
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
			var lokasis []Lokasi
			if err = cursor.All(context.Background(), &lokasis); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
			// Agar return array kosong [] bukan null jika data kosong
			if lokasis == nil {
				lokasis = []Lokasi{}
			}
			c.JSON(http.StatusOK, lokasis)
		}
	})

	// POST: Tambah Data
	r.POST("/api/lokasi", func(c *gin.Context) {
		var lokasi Lokasi
		// Validasi JSON yang masuk
		if err := c.ShouldBindJSON(&lokasi); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Pastikan GeoJSON valid
		lokasi.Koordinat.Type = "Point"
		lokasi.ID = primitive.NewObjectID() // Generate ID baru

		_, err := collection.InsertOne(context.Background(), lokasi)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusCreated, lokasi)
	})

	// PUT: Update Data
	r.PUT("/api/lokasi", func(c *gin.Context) {
		idParam := c.Query("id")
		objID, err := primitive.ObjectIDFromHex(idParam)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "ID tidak valid"})
			return
		}

		var updateData Lokasi
		if err := c.ShouldBindJSON(&updateData); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		update := bson.M{
			"$set": bson.M{
				"nama":      updateData.Nama,
				"kategori":  updateData.Kategori,
				"deskripsi": updateData.Deskripsi,
				"koordinat": updateData.Koordinat,
			},
		}

		_, err = collection.UpdateOne(context.Background(), bson.M{"_id": objID}, update)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		// Kembalikan data yang sudah diupdate (termasuk ID lama)
		updateData.ID = objID
		c.JSON(http.StatusOK, updateData)
	})

	// DELETE: Hapus Data
	r.DELETE("/api/lokasi", func(c *gin.Context) {
		idParam := c.Query("id")
		objID, err := primitive.ObjectIDFromHex(idParam)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "ID tidak valid"})
			return
		}

		_, err = collection.DeleteOne(context.Background(), bson.M{"_id": objID})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Berhasil dihapus"})
	})

	// Jalankan Server (Support Port ENV untuk Vercel/Render)
	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}
	r.Run(":" + port)
}