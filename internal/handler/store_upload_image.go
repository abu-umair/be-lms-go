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

func UploadStoreImageHandler(c *fiber.Ctx) error {
	//* 1. Generate UUID / cek store_id
	var storeID string
	isUpdate := false

	//? 1. Logika Identifikasi: Create atau Update?
	inputID := c.FormValue("store_id")

	if inputID != "" {
		//? Jika ada inputID, berarti UPDATE
		isUpdate = true
		storeID = inputID
	} else {
		//? Jika TIDAK ADA inputID, berarti CREATE
		storeID = uuid.NewString()
		isUpdate = false
	}

	if storeID == "" {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "store_id is required",
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
	//? ./storage/1234444/store/
	folderPath := fmt.Sprintf("./storage/%s/store", storeID)

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

	//* store_1623232.png (membuat format imagge name)
	timestamp := time.Now().UnixNano()
	fileName := fmt.Sprintf("store_%d%s", timestamp, filepath.Ext(file.Filename))

	uploadPath := folderPath + "/" + fileName
	// c.SaveFile(file, "./storage/store/store.jpeg")
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
		"store_id":  storeID,
		"file_name": fileName,
	})
}
