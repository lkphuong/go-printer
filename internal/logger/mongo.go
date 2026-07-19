package logger

import (
	"context"
	"encoding/json"
	"go-printer/internal/constants"
	"log"
	"os"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	databaseName         = "printer_logs"
	collectionPrefix     = "logs_"
	collectionDateLayout = "20060102"
	retentionDays        = 14
	deviceConfigPath     = "./config/device.json"
	defaultLocation      = "unknown"
)

// client is held at package level so every log call reuses one connection pool.
// A nil client means Mongo is unavailable; callers must degrade to file logging.
var client *mongo.Client

type printLog struct {
	StatusCode string    `bson:"status_code"`
	HTTPStatus int       `bson:"http_status"`
	Device     string    `bson:"device"`
	Message    string    `bson:"message"`
	Timestamp  time.Time `bson:"timestamp"`
}

// Init connects to MongoDB and pings it. Any failure is logged and returned, never
// panicked, so the service keeps running with file-only logging when Mongo is down.
func Init() error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	c, err := mongo.Connect(ctx, options.Client().ApplyURI(constants.MONGODB_URI))
	if err != nil {
		log.Println("mongo connect error:", err)
		return err
	}

	if err := c.Ping(ctx, nil); err != nil {
		log.Println("mongo ping error:", err)
		return err
	}

	client = c
	log.Println("MongoDB logger initialized.")
	return nil
}

// LogPrint writes one structured entry into today's collection. It always mirrors to the
// file log first so no event is lost when Mongo is unavailable.
func LogPrint(statusCode string, httpStatus int, message string) {
	log.Printf("print log: status=%s http=%d message=%s", statusCode, httpStatus, message)

	if client == nil {
		return
	}

	entry := printLog{
		StatusCode: statusCode,
		HTTPStatus: httpStatus,
		Device:     getDeviceLocation(),
		Message:    message,
		Timestamp:  time.Now(),
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	collection := client.Database(databaseName).Collection(collectionName(time.Now()))
	if _, err := collection.InsertOne(ctx, entry); err != nil {
		log.Println("mongo insert error:", err)
	}
}

// CleanupOldCollections drops day-collections older than the retention window. It reuses
// the existing daily scheduler in app.go and contains all errors so it can never crash.
func CleanupOldCollections() {
	if client == nil {
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	db := client.Database(databaseName)
	names, err := db.ListCollectionNames(ctx, bson.M{})
	if err != nil {
		log.Println("mongo list collections error:", err)
		return
	}

	cutoff := time.Now().AddDate(0, 0, -retentionDays)
	for _, name := range names {
		if !strings.HasPrefix(name, collectionPrefix) {
			continue
		}

		datePart := strings.TrimPrefix(name, collectionPrefix)
		collDate, err := time.Parse(collectionDateLayout, datePart)
		if err != nil {
			continue
		}

		if collDate.Before(cutoff) {
			if err := db.Collection(name).Drop(ctx); err != nil {
				log.Println("mongo drop collection error:", err)
				continue
			}
			log.Printf("dropped old log collection: %s", name)
		}
	}
}

// getDeviceLocation resolves the machine location from device.json, defaulting to
// "unknown" so a missing or malformed file never blocks a log write.
func getDeviceLocation() string {
	data, err := os.ReadFile(deviceConfigPath)
	if err != nil {
		return defaultLocation
	}

	var config map[string]string
	if err := json.Unmarshal(data, &config); err != nil {
		return defaultLocation
	}

	location := config["location"]
	if location == "" {
		return defaultLocation
	}
	return location
}

func collectionName(t time.Time) string {
	return collectionPrefix + t.Format(collectionDateLayout)
}
