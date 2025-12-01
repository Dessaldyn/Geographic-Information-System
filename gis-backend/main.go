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

// Global Variables
var (
	app        *gin.Engine
	collection *mongo.Collection
)

// --- 1. SETUP AWAL (INIT) ---
// Fungsi init() jalan otomatis saat Vercel memulai serverless function
func init() {
	connectDB()
	setupRouter()
}

// --- 2. ENTRYPOINT VERCEL (PENTING!) ---
// Vercel mencari fungsi bernama "Handler" ini. Jangan dihapus/ubah namanya.
func Handler(w http.ResponseWriter, r *http.Request) {
	app.ServeHTTP(w, r)
}

// --- 3. ENTRYPOINT LOKAL (LAPTOP) ---
// Ini dipakai kalau kamu jalankan "go run main.go" di laptop
func main() {
	// Pastikan setup jalan jika dijalankan manual
	if app == nil {
		connectDB()
		setupRouter()
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}
	log.Println("🚀 Server Golang jalan di http://localhost:" + port)
	app.Run(":" + port)
}

// --- 4. SETUP ROUTER & CORS ---
func setupRouter() {
	// Gunakan gin.New() agar lebih ringan di serverless
	app = gin.New()
	app.Use(gin.Recovery())

	// Konfigurasi CORS (Agar Frontend di GitHub Pages bisa masuk)
	config := cors.DefaultConfig()
	config.AllowAllOrigins = true
	config.AllowMethods = []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}
	config.AllowHeaders = []string{"Origin", "Content-Type", "Accept"}
	app.Use(cors.New(config))

	// Routes
	app.GET("/api/lokasi", getLokasi)
	app.POST("/api/lokasi", createLokasi)
	app.PUT("/api/lokasi", updateLokasi)
	app.DELETE("/api/lokasi", deleteLokasi)
}

// --- 5. KONEKSI DATABASE ---
func connectDB() {
	// Cek apakah koneksi sudah ada? Kalau ada, pakai yang lama (Caching)
	if collection != nil {
		return
	}

	mongoURI := os.Getenv("MONGO_URI")
	if mongoURI == "" {
		// Fallback hardcode (Hanya untuk darurat/testing)
		mongoURI = "mongodb+srv://sriwahyuni_db_user:EgZ2GXRliZQ1TYA7@cluster23.gnmjc2n.mongodb.net/ujianSIG"
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(mongoURI))
	if err != nil {
		log.Println("❌ Gagal buat client Mongo:", err)
		return
	}

	err = client.Ping(ctx, nil)
	if err != nil {
		log.Println("❌ Gagal ping Mongo:", err)
		return
	}

	log.Println("✅ Terhubung ke MongoDB Atlas")
	collection = client.Database("ujianSIG").Collection("lokasis")
}

// --- 6. MODELS ---
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

// --- 7. CONTROLLERS (LOGIC) ---

func getLokasi(c *gin.Context) {
	if collection == nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database belum terkoneksi"})
		return
	}

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
}

func createLokasi(c *gin.Context) {
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
}

func updateLokasi(c *gin.Context) {
	idParam := c.Query("id")
	objID, err := primitive.ObjectIDFromHex(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID invalid"})
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
}

func deleteLokasi(c *gin.Context) {
	idParam := c.Query("id")
	objID, err := primitive.ObjectIDFromHex(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID invalid"})
		return
	}
	_, err = collection.DeleteOne(context.Background(), bson.M{"_id": objID})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Berhasil dihapus"})
}