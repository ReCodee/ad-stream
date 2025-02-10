package main

import (
	"database/sql"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/websocket/v2"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

type Config struct {
	DB_PORT      string
	DB_HOST      string
	DB_USER      string
	DB_PASSWORD  string
	DB_NAME      string
	CLICKS_TABLE string
	APP_PORT     string
}

type Ad struct {
	ID    int    `json:"id"`
	Image string `json:"image"`
	URL   string `json:"url"`
}

type AdPosition struct {
	X int `json:"x"`
	Y int `json:"y"`
}

type ClickData struct {
	AdID      int        `json:"adId"`
	Timestamp string     `json:"timestamp"`
	VideoTime float64    `json:"videoTime"`
	Position  AdPosition `json:"position"`
	HoverTime float64    `json:"hoverTime"`
}

type WSClient struct {
	conn *websocket.Conn
	ads  chan Ad
}

var (
	clients    = make(map[*websocket.Conn]*WSClient)
	register   = make(chan *WSClient)
	unregister = make(chan *WSClient)
	broadcast  = make(chan Ad)
)

var db *sql.DB
var config Config

func main() {
	var err error
	err = godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}
	config = Config{
		DB_PORT:      os.Getenv("DB_PORT"),
		DB_HOST:      os.Getenv("DB_HOST"),
		APP_PORT:     os.Getenv("APP_PORT"),
		DB_USER:      os.Getenv("DB_USER"),
		DB_PASSWORD:  os.Getenv("DB_PASSWORD"),
		DB_NAME:      os.Getenv("DB_NAME"),
		CLICKS_TABLE: os.Getenv("CLICKS_TABLE"),
	}

	db, err = sql.Open("postgres", fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", config.DB_HOST, config.DB_PORT, config.DB_USER, config.DB_PASSWORD, config.DB_NAME))
	if err != nil {
		log.Fatalf("Error connecting to database: %v", err)
	}
	err = db.Ping()
	if err != nil {
		log.Fatalf("Error pinging database: %v", err)
	}
	log.Printf("Connected to PostgreSQL DB at Port: %s", config.DB_PORT)

	defer db.Close()

	_, err = db.Exec(fmt.Sprintf("CREATE TABLE IF NOT EXISTS %s (id SERIAL PRIMARY KEY, ad_id INTEGER, timestamp TEXT, video_time REAL, x INTEGER, y INTEGER, hover_time REAL)", config.CLICKS_TABLE))
	if err != nil {
		log.Fatalf("Error creating table: %v", err)
	}

	app := fiber.New()
	app.Use(cors.New())

	app.Use("/ws", func(c *fiber.Ctx) error {
		if websocket.IsWebSocketUpgrade(c) {
			return c.Next()
		}
		return fiber.ErrUpgradeRequired
	})

	app.Get("/ws/ads", websocket.New(handleWebSocket))
	app.Post("/ads/click", logAdClick)

	go handleConnections()
	go broadcastAds()

	log.Fatal(app.Listen(config.APP_PORT))
}

func handleWebSocket(c *websocket.Conn) {
	client := &WSClient{
		conn: c,
		ads:  make(chan Ad),
	}
	log.Printf("New WebSocket client connected from: %s", c.RemoteAddr().String())
	register <- client
	log.Printf("Client registered: %s", c.RemoteAddr().String())

	defer func() {
		log.Printf("Client disconnecting: %s", c.RemoteAddr().String())
		unregister <- client
		log.Printf("Client unregistered: %s", c.RemoteAddr().String())
		client.conn.Close()
		log.Printf("Client connection closed: %s", c.RemoteAddr().String())
	}()

	for {
		if err := c.WriteJSON(<-client.ads); err != nil {
			log.Printf("Error writing to client %s: %v", c.RemoteAddr().String(), err)
			break
		}
		log.Printf("Successfully sent ad to client: %s", c.RemoteAddr().String())
	}
}

func handleConnections() {
	for {
		select {
		case client := <-register:
			clients[client.conn] = client
		case client := <-unregister:
			if client, ok := clients[client.conn]; ok {
				close(client.ads)
				delete(clients, client.conn)
			}

		case ad := <-broadcast:
			for _, client := range clients {
				select {
				case client.ads <- ad:
				default:
					log.Printf("Failed to send ad to client, channel full")
				}
			}
		}
	}
}

func broadcastAds() {
	ads := []Ad{
		{1, "https://scontent.fknu2-1.fna.fbcdn.net/v/t39.30808-6/476005818_922314306781136_6981070312990821632_n.jpg?_nc_cat=101&ccb=1-7&_nc_sid=cc71e4&_nc_ohc=EaUSjEMenw0Q7kNvgEfp6ww&_nc_oc=AdjhAOiUvKdnLs3ziCEvH4NBkogUV0RsPN734aFCnpNzzCbGaaamNqnKHdPqJVS6Qx0&_nc_zt=23&_nc_ht=scontent.fknu2-1.fna&_nc_gid=AagEkjoz0MVbn1yVO4UV3ms&oh=00_AYCK-qC5I9IAEn72-HI2imUUhPmMJTfFjcmtbkdiiHwpvg&oe=67AC48A1",
			"https://www.zeptonow.com/"},
		{2, "https://www.analyticssteps.com/backend/media/thumbnail/1890055/4828382_1669140152_BlinkitArtboard%201.jpg",
			"https://www.blinkit.com/"},
	}

	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop()

	for range ticker.C {
		broadcast <- ads[rand.Intn(len(ads))]
	}
}

func getAds(c *fiber.Ctx) error {
	ads := []Ad{
		{1, "https://scontent.fknu2-1.fna.fbcdn.net/v/t39.30808-6/476005818_922314306781136_6981070312990821632_n.jpg?_nc_cat=101&ccb=1-7&_nc_sid=cc71e4&_nc_ohc=EaUSjEMenw0Q7kNvgEfp6ww&_nc_oc=AdjhAOiUvKdnLs3ziCEvH4NBkogUV0RsPN734aFCnpNzzCbGaaamNqnKHdPqJVS6Qx0&_nc_zt=23&_nc_ht=scontent.fknu2-1.fna&_nc_gid=AagEkjoz0MVbn1yVO4UV3ms&oh=00_AYCK-qC5I9IAEn72-HI2imUUhPmMJTfFjcmtbkdiiHwpvg&oe=67AC48A1",
			"https://www.zeptonow.com/"},
		{2, "https://www.analyticssteps.com/backend/media/thumbnail/1890055/4828382_1669140152_BlinkitArtboard%201.jpg",
			"https://www.blinkit.com/"},
	}
	if err := c.JSON(ads[rand.Intn(len(ads)):]); err != nil {
		log.Printf("Error sending JSON response: %v", err)
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to send response",
		})
	}
	return nil
}

func logAdClick(c *fiber.Ctx) error {
	var data ClickData
	if err := c.BodyParser(&data); err != nil {
		log.Printf("Error parsing request body: %v", err)
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request format",
		})
	}

	_, err := db.Exec(fmt.Sprintf("INSERT INTO %s (ad_id, timestamp, video_time, x, y, hover_time) VALUES ($1, $2, $3, $4, $5, $6)", config.CLICKS_TABLE), data.AdID, data.Timestamp, data.VideoTime, data.Position.X, data.Position.Y, data.HoverTime)
	if err != nil {
		log.Printf("Error saving click data: %v", err)
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to save click data",
		})
	}

	return c.Status(http.StatusOK).JSON(fiber.Map{
		"message": "Click data saved successfully",
	})
}
