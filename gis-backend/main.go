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
	// Ambil MONGO_URI dari Environment Variable (Settingan di Vercel)
	mongoURI := os.Getenv("MONGO_URI")
	if mongoURI == "" {
		// Fallback ke link Atlas kamu jika dijalankan di lokal laptop
		mongoURI = "mongodb+srv://sriwahyuni_db_user:EgZ2GXRliZQ1TYA7@cluster23.gnmjc2n.mongodb.net/ujianSIG"
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(mongoURI))
	if err != nil {
		log.Fatal(err)
	}

	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatal("Gagal konek MongoDB:", err)
	}

	log.Println("✅ BERHASIL Konek ke MongoDB")
	collection = client.Database("ujianSIG").Collection("lokasis")
}

func main() {
	connectDB()

	r := gin.Default()

	// --- PERBAIKAN CORS DISINI ---
	// Kita izinkan semua origin (*) agar Frontend di Localhost maupun GitHub Pages bisa masuk
	config := cors.DefaultConfig()
	config.AllowAllOrigins = true
	config.AllowMethods = []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}
	config.AllowHeaders = []string{"Origin", "Content-Type", "Accept"}
	
	r.Use(cors.New(config))

	// --- 3. ROUTES API ---

	// GET: Ambil Data
	r.GET("/api/lokasi", func(c *gin.Context) {
		idParam := c.Query("id")

		if idParam != "" {
			objID, err := primitive.ObjectIDFromHex(idParam)
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": "ID tidak valid"})
				return
			}
			var lokasi Lokasi
			err = collection.FindOne(context.Background(), bson.M{"_id": objID}).Decode(&lokasi)
			if err != nil {
				c.JSON(http.StatusNotFound, gin.H{"error": "Data tidak ditemukan"})
				return
			}
			c.JSON(http.StatusOK, lokasi)
		} else {
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
			if lokasis == nil {
				lokasis = []Lokasi{}
			}
			c.JSON(http.StatusOK, lokasis)
		}
	})

	// POST: Tambah Data
	r.POST("/api/lokasi", func(c *gin.Context) {
		var lokasi Lokasi
		if err := c.ShouldBindJSON(&lokasi); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		lokasi.Koordinat.Type = "Point"
		lokasi.ID = primitive.NewObjectID()

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

	// Jalankan Server
	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}
	r.Run(":" + port)
}