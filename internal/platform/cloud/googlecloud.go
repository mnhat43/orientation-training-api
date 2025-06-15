package cloud

import (
	"context"
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"net/url"
	"os"
	"strings"

	"cloud.google.com/go/storage"
	"github.com/labstack/echo/v4"
)

// GcsStorage struct
type GcsStorage struct {
	Logger echo.Logger
	Ctx    context.Context
	Bucket *storage.BucketHandle
}

// NewGcsStorage export
func NewGcsStorage(logger echo.Logger) (cloud *GcsStorage) {
	cloud = &GcsStorage{}
	cloud.Init(logger)
	return
}

// Init new connect
func (cloud *GcsStorage) Init(logger echo.Logger) {
	cloud.Logger = logger
	cloud.Ctx = context.Background()
	client, err := storage.NewClient(cloud.Ctx)
	if err != nil {
		cloud.Logger.Error(err)
	}

	cloud.Bucket = client.Bucket(os.Getenv("GOOGLE_STORAGE_BUCKET"))
}

// GetFileByFileName show image
func (cloud *GcsStorage) GetFileByFileName(fileName string, directoryCloud string) ([]byte, error) {
	linkCloud := directoryCloud + fileName
	reader, err := cloud.Bucket.Object(linkCloud).NewReader(cloud.Ctx)
	if err != nil {
		cloud.Logger.Error(err)
		return nil, err
	}
	defer reader.Close()

	data, err := ioutil.ReadAll(reader)
	if err != nil {
		cloud.Logger.Error(err)
		return nil, err
	}

	return data, nil
}

// UploadFileToCloud : upload file to cloud
// Params      : directory file name
// Returns     : error or nil
func (cloud *GcsStorage) UploadFileToCloud(file string, fileName string, directoryCloud string) error {
	linkCloud := directoryCloud + fileName

	// Thiết lập content type dựa trên phần mở rộng file
	contentType := "application/octet-stream" // Mặc định

	// Xử lý các loại file phổ biến, bao gồm file thông thường và slide
	lowerFileName := strings.ToLower(fileName)
	if strings.HasSuffix(lowerFileName, ".pdf") {
		contentType = "application/pdf"
	} else if strings.HasSuffix(lowerFileName, ".xlsx") {
		contentType = "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet"
	} else if strings.HasSuffix(lowerFileName, ".docx") {
		contentType = "application/vnd.openxmlformats-officedocument.wordprocessingml.document"
	} else if strings.HasSuffix(lowerFileName, ".ppt") {
		contentType = "application/vnd.ms-powerpoint"
	} else if strings.HasSuffix(lowerFileName, ".pptx") {
		contentType = "application/vnd.openxmlformats-officedocument.presentationml.presentation"
	} else if strings.HasSuffix(lowerFileName, ".odp") {
		contentType = "application/vnd.oasis.opendocument.presentation"
	} else if strings.HasSuffix(lowerFileName, ".pptm") {
		contentType = "application/vnd.ms-powerpoint.presentation.macroenabled.12"
	} else if strings.HasSuffix(lowerFileName, ".ppsm") {
		contentType = "application/vnd.ms-powerpoint.slideshow.macroenabled.12"
	} else if strings.HasSuffix(lowerFileName, ".potx") {
		contentType = "application/vnd.openxmlformats-officedocument.presentationml.template"
	} else if strings.HasSuffix(lowerFileName, ".ppsx") {
		contentType = "application/vnd.openxmlformats-officedocument.presentationml.slideshow"
	} else if strings.HasSuffix(lowerFileName, ".pps") {
		contentType = "application/vnd.ms-powerpoint.slideshow"
	}

	// Tạo một đối tượng writer với các thuộc tính tùy chỉnh
	wc := cloud.Bucket.Object(linkCloud).NewWriter(cloud.Ctx)
	wc.ContentType = contentType
	// Đặt Cache-Control để tối ưu hiệu suất
	wc.CacheControl = "public, max-age=86400"

	var err error
	var dataBytes []byte

	// Xử lý base64 đúng cách
	// Có thể có nhiều định dạng base64 khác nhau (standard, URL-safe)
	dataBytes, err = base64.StdEncoding.DecodeString(file)
	if err != nil {
		// Thử với URL-safe encoding nếu standard encoding không thành công
		dataBytes, err = base64.URLEncoding.DecodeString(file)
		if err != nil {
			dataBytes, err = base64.RawStdEncoding.DecodeString(file)
			if err != nil {
				dataBytes, err = base64.RawURLEncoding.DecodeString(file)
				if err != nil {
					cloud.Logger.Errorf("Failed to decode base64 data: %v", err)
					return err
				}
			}
		}
	}

	// Ghi dữ liệu đã giải mã vào cloud storage
	if _, err := wc.Write(dataBytes); err != nil {
		cloud.Logger.Error(err)
		return err
	}

	// Đảm bảo đóng writer để hoàn thành quá trình tải lên
	if err := wc.Close(); err != nil {
		cloud.Logger.Error(err)
		return err
	}

	if err := wc.Close(); err != nil {
		cloud.Logger.Error(err)
		return err
	}
	// [END upload_file]
	return nil
}

// DeleteFileCloud : delete image in cloud
// Params      : file name
// Returns     : error or nil
func (cloud *GcsStorage) DeleteFileCloud(fileName string, directoryCloud string) error {
	linkCloud := directoryCloud + fileName
	o := cloud.Bucket.Object(linkCloud)

	if err := o.Delete(cloud.Ctx); err != nil {
		cloud.Logger.Error(err)
		return err
	}

	return nil
}

// GetURL generates a public URL for a in the cloud storage
func (cloud *GcsStorage) GetURL(fileName string, directoryCloud string) string {
	// Tạo URL chuẩn cho file trên GCS với các tham số bổ sung để đảm bảo xử lý đúng
	baseURL := "https://storage.googleapis.com/" + os.Getenv("GOOGLE_STORAGE_BUCKET") + "/" + directoryCloud + fileName

	// Sử dụng chuỗi hoa thường để kiểm tra an toàn hơn
	lowerFileName := strings.ToLower(fileName)

	// Thêm tham số cho các file để đảm bảo hiển thị đúng
	if strings.HasSuffix(lowerFileName, ".pdf") ||
		strings.HasSuffix(lowerFileName, ".docx") || strings.HasSuffix(lowerFileName, ".xlsx") ||
		strings.HasSuffix(lowerFileName, ".ppt") || strings.HasSuffix(lowerFileName, ".pptx") ||
		strings.HasSuffix(lowerFileName, ".pptm") || strings.HasSuffix(lowerFileName, ".ppsx") ||
		strings.HasSuffix(lowerFileName, ".pps") || strings.HasSuffix(lowerFileName, ".potx") ||
		strings.HasSuffix(lowerFileName, ".odp") || strings.HasSuffix(lowerFileName, ".ppsm") {
		// Thêm các header cần thiết để đảm bảo trình duyệt xử lý file PowerPoint đúng cách
		return baseURL + "?response-content-disposition=attachment%3B%20filename%3D" + url.QueryEscape(fileName) +
			"&response-content-type=" + url.QueryEscape(getContentTypeFromFileName(fileName))
	}

	return baseURL
}

// getContentTypeFromFileName trả về MIME type dựa trên tên file
func getContentTypeFromFileName(fileName string) string {
	fmt.Printf("Getting content type for file: %s\n", fileName)
	lowerFileName := strings.ToLower(fileName)
	fmt.Printf("Lowercase file name: %s\n", lowerFileName)
	// File thông thường
	if strings.HasSuffix(lowerFileName, ".pdf") {
		return "application/pdf"
	} else if strings.HasSuffix(lowerFileName, ".xlsx") {
		return "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet"
	} else if strings.HasSuffix(lowerFileName, ".docx") {
		return "application/vnd.openxmlformats-officedocument.wordprocessingml.document"
	} else if strings.HasSuffix(lowerFileName, ".ppt") { // File slide
		return "application/vnd.ms-powerpoint"
	} else if strings.HasSuffix(lowerFileName, ".pptx") {
		return "application/vnd.openxmlformats-officedocument.presentationml.presentation"
	} else if strings.HasSuffix(lowerFileName, ".odp") {
		return "application/vnd.oasis.opendocument.presentation"
	} else if strings.HasSuffix(lowerFileName, ".pptm") {
		return "application/vnd.ms-powerpoint.presentation.macroenabled.12"
	} else if strings.HasSuffix(lowerFileName, ".ppsm") {
		return "application/vnd.ms-powerpoint.slideshow.macroenabled.12"
	} else if strings.HasSuffix(lowerFileName, ".potx") {
		return "application/vnd.openxmlformats-officedocument.presentationml.template"
	} else if strings.HasSuffix(lowerFileName, ".ppsx") {
		return "application/vnd.openxmlformats-officedocument.presentationml.slideshow"
	} else if strings.HasSuffix(lowerFileName, ".pps") {
		return "application/vnd.ms-powerpoint.slideshow"
	}
	return "application/octet-stream"
}
