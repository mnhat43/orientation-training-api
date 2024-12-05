package cloud

import (
	"context"
	"fmt"
	"io"
	"os"

	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
	"github.com/labstack/echo/v4"
)

// CloudinaryStorage struct
type CloudinaryStorage struct {
	Logger echo.Logger
	Ctx    context.Context
	Cloud  *cloudinary.Cloudinary
}

// NewCloudinaryStorage export
func NewCloudinaryStorage(logger echo.Logger) (cloud *CloudinaryStorage) {
	cloud = &CloudinaryStorage{}
	cloud.Init(logger)
	return
}

// Init new connect
func (cloud *CloudinaryStorage) Init(logger echo.Logger) {
	cloud.Logger = logger
	cloud.Ctx = context.Background()

	// Tạo một client Cloudinary
	cloudinaryURL := os.Getenv("CLOUDINARY_URL")
	cld, err := cloudinary.NewFromURL(cloudinaryURL)
	if err != nil {
		cloud.Logger.Error(err)
	}
	cloud.Cloud = cld
}

func (cloud *CloudinaryStorage) GetFileByFileName(fileName string, directoryCloud string) string {
	publicID := directoryCloud + fileName
	cloudinaryName := os.Getenv("CLOUDINARY_NAME")
	secureURL := fmt.Sprintf("https://res.cloudinary.com/%s/image/upload/%s", cloudinaryName, publicID)

	return secureURL
}

func (cloud *CloudinaryStorage) UploadFileToCloud(file io.Reader, fileName string, directoryCloud string) error {
	_, err := cloud.Cloud.Upload.Upload(cloud.Ctx, file, uploader.UploadParams{
		PublicID: fileName,
		Folder:   directoryCloud,
	})
	if err != nil {
		cloud.Logger.Error(err)
		return err
	}

	return nil
}

func (cloud *CloudinaryStorage) DeleteFileCloud(fileName string, directoryCloud string) error {
	publicID := directoryCloud + fileName

	_, err := cloud.Cloud.Upload.Destroy(cloud.Ctx, uploader.DestroyParams{
		PublicID: publicID,
	})

	if err != nil {
		cloud.Logger.Error(err)
		return err
	}

	return nil
}
