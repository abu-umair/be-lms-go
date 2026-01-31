package main

import (
	"log"
	"mime"
	"net/http"
	"os"
	"path"

	"github.com/abu-umair/be-lms-go/internal/handler"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

func handleGetFileName(c *fiber.Ctx) error {
	courseID := c.Params("course_id")
	fileNameParam := c.Params("filename")
	filePath := path.Join("storage", courseID, "course", fileNameParam)
	if _, err := os.Stat(filePath); err != nil {
		if os.IsNotExist(err) {
			return c.Status(http.StatusNotFound).SendString("Not Found")
		}
		log.Println(err)
		return c.Status(http.StatusInternalServerError).SendString("Internal Server Error")
	}

	//? membuka file
	file, err := os.Open(filePath)
	if err != nil {
		log.Println(err)
		return c.Status(http.StatusInternalServerError).SendString("Internal Server Error")
	}

	//? mengirim file sbg response
	ext := path.Ext(filePath)
	mimeType := mime.TypeByExtension(ext)

	c.Set("Content-Type", mimeType) //?konversi agar tampilan gambar sesuai (dinamis)
	// c.Set("Content-Type", "image/png") //?konversi agar tampilan gambar sesuai (blm dinamis)
	return c.SendStream(file)
}

func main() {
	app := fiber.New()
	app.Use(cors.New())

	app.Get("/storage/:course_id/course/:filename", handleGetFileName)

	app.Post("/course/upload", handler.UploadCourseImageHandler)

	app.Listen(":3000")

}
