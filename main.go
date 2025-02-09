package main

import (
	"log"
	"math/rand"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

type Config struct {
	DB_URL string
	PORT   string
}

type Ad struct {
	ID    int    `json:"id"`
	Image string `json:"image"`
	URL   string `json:"url"`
}

type ClickData struct {
	AdID      int     `json:"adId"`
	Timestamp string  `json:"timestamp"`
	VideoTime float64 `json:"videoTime"`
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal(err)
	}
	var config Config
	config.DB_URL = os.Getenv("DB_URL")
	config.PORT = os.Getenv("PORT")
	app := fiber.New()
	app.Use(cors.New())

	app.Get("/ads", getAds)

	log.Fatal(app.Listen(config.PORT))
}

func getAds(c *fiber.Ctx) error {
	ads := []Ad{
		{1, "https://scontent.fknu2-1.fna.fbcdn.net/v/t39.30808-6/476005818_922314306781136_6981070312990821632_n.jpg?_nc_cat=101&ccb=1-7&_nc_sid=cc71e4&_nc_ohc=EaUSjEMenw0Q7kNvgEfp6ww&_nc_oc=AdjhAOiUvKdnLs3ziCEvH4NBkogUV0RsPN734aFCnpNzzCbGaaamNqnKHdPqJVS6Qx0&_nc_zt=23&_nc_ht=scontent.fknu2-1.fna&_nc_gid=AagEkjoz0MVbn1yVO4UV3ms&oh=00_AYCK-qC5I9IAEn72-HI2imUUhPmMJTfFjcmtbkdiiHwpvg&oe=67AC48A1",
			"https://www.zeptonow.com/"},
		{2, "https://www.analyticssteps.com/backend/media/thumbnail/1890055/4828382_1669140152_BlinkitArtboard%201.jpg",
			"https://www.blinkit.com/"},
	}
	return c.JSON(ads[rand.Intn(len(ads)):])
}
