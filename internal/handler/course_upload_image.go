package handler

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

func UploadCourseImageHandler(c *fiber.Ctx) error {
	//* 1. Generate UUID / cek course_id
	var courseID string
	isUpdate := false

	//? 1. Logika Identifikasi: Create atau Update?
	inputID := c.FormValue("course_id")

	if inputID != "" {
		//? Jika ada inputID, berarti UPDATE
		isUpdate = true
		courseID = inputID
	} else {
		//? Jika TIDAK ADA inputID, berarti CREATE
		courseID = uuid.NewString()
		isUpdate = false
	}

	if courseID == "" {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "course_id is required",
		})
	}

	//* 2. Ambil File
	file, err := c.FormFile("image")
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "image data not found",
		})
	}

	//* 3. Susun Path Folder
	//? ./storage/1234444/course/
	folderPath := fmt.Sprintf("./storage/%s/course", courseID)

	//* 4. CEK & BUAT FOLDER (os.MkdirAll)
	if !isUpdate { //? buat folder utk CREATE saja
		// 0755 adalah permission standar (read/write untuk owner, read untuk lainnya)
		err = os.MkdirAll(folderPath, 0755)
		if err != nil {
			return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
				"success": false,
				"message": "failed to create directory",
			})
		}
	}

	//? validasi gambar
	//? memeriksa ekstensi file (validasi extensi file)
	ext := strings.ToLower(filepath.Ext(file.Filename))

	allowedExts := map[string]bool{
		".jpg":  true,
		".jpeg": true,
		".png":  true,
		".webp": true,
	}

	if !allowedExts[ext] {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "image extention not allowed (jpg, jpeg, png, webp)",
		})
	}

	//? validasi content type
	contentType := file.Header.Get("Content-Type")

	allowedContentType := map[string]bool{
		"image/jpg":  true,
		"image/jpeg": true,
		"image/png":  true,
		"image/webp": true,
	}

	if !allowedContentType[contentType] {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "Content Type is not allowed",
		})
	}

	//* course_1623232.png (membuat format imagge name)
	timestamp := time.Now().UnixNano()
	fileName := fmt.Sprintf("course_%d%s", timestamp, filepath.Ext(file.Filename))

	uploadPath := folderPath + "/" + fileName
	// c.SaveFile(file, "./storage/course/course.jpeg")
	err = c.SaveFile(file, uploadPath)

	//return error jika ada
	if err != nil {
		fmt.Println(err)

		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "internal server error",
		})
	}

	return c.JSON(fiber.Map{
		"success":   true,
		"message":   "Upload success",
		"course_id": courseID,
		"file_name": fileName,
	})
}
