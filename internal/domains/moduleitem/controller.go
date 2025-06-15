package moduleitem

import (
	"net/http"
	"net/url"
	cf "orientation-training-api/configs"
	cm "orientation-training-api/internal/common"
	rp "orientation-training-api/internal/interfaces/repository"
	param "orientation-training-api/internal/interfaces/requestparams"
	cld "orientation-training-api/internal/platform/cloud"
	"orientation-training-api/internal/platform/utils"
	"orientation-training-api/internal/platform/youtube"
	"strconv"
	"strings"
	"time"

	valid "github.com/asaskevich/govalidator"

	"github.com/go-pg/pg/v9"
	"github.com/labstack/echo/v4"
)

type ModuleItemController struct {
	cm.BaseController

	ModuleItemRepo rp.ModuleItemRepository
	QuizRepo       rp.QuizRepository
	cloud          cld.StorageUtility
}

func NewModuleItemController(logger echo.Logger, moduleItemRepo rp.ModuleItemRepository, quizRepo rp.QuizRepository, cloud cld.StorageUtility) (ctr *ModuleItemController) {
	ctr = &ModuleItemController{cm.BaseController{}, moduleItemRepo, quizRepo, cloud}
	ctr.Init(logger)
	return
}

// GetModuleItemList : get list of moduleItems(by moduleName keyword)
// Params : echo.Context
// Returns : return error
func (ctr *ModuleItemController) GetModuleItemList(c echo.Context) error {
	// userProfile := c.Get("user_profile").(m.User)
	moduleItemListParams := new(param.ModuleItemListParams)

	if err := c.Bind(moduleItemListParams); err != nil {
		return c.JSON(http.StatusOK, cf.JsonResponse{
			Status:  cf.FailResponseCode,
			Message: "Invalid Params",
			Data:    err,
		})
	}

	if _, err := valid.ValidateStruct(moduleItemListParams); err != nil {
		return c.JSON(http.StatusOK, cf.JsonResponse{
			Status:  cf.FailResponseCode,
			Message: err.Error(),
		})
	}

	moduleItems, _, err := ctr.ModuleItemRepo.GetModuleItems(moduleItemListParams)

	if err != nil {
		if err.Error() == pg.ErrNoRows.Error() {
			return c.JSON(http.StatusOK, cf.JsonResponse{
				Status:  cf.FailResponseCode,
				Message: "Get module item list failed",
			})
		}

		return c.JSON(http.StatusInternalServerError, cf.JsonResponse{
			Status:  cf.FailResponseCode,
			Message: "System Error",
		})
	}

	return c.JSON(http.StatusOK, cf.JsonResponse{
		Status:  cf.SuccessResponseCode,
		Message: "Success",
		Data: map[string]interface{}{
			"moduleItems": moduleItems,
		},
	})
}

// AddModuleItem : add new ModuleItem to database
// Params : echo.Context
// Returns : return error
func (ctr *ModuleItemController) AddModuleItem(c echo.Context) error {
	createModuleItemParams := new(param.CreateModuleItemParams)

	if err := c.Bind(createModuleItemParams); err != nil {
		return c.JSON(http.StatusOK, cf.JsonResponse{
			Status:  cf.FailResponseCode,
			Message: "Invalid Params",
			Data:    err,
		})
	}

	if _, err := valid.ValidateStruct(createModuleItemParams); err != nil {
		return c.JSON(http.StatusOK, cf.JsonResponse{
			Status:  cf.FailResponseCode,
			Message: err.Error(),
		})
	}

	if createModuleItemParams.ItemType == "" ||
		(createModuleItemParams.ItemType != "video" &&
			createModuleItemParams.ItemType != "file" &&
			createModuleItemParams.ItemType != "quiz" &&
			createModuleItemParams.ItemType != "slide") { // Thêm slide
		return c.JSON(http.StatusBadRequest, cf.JsonResponse{
			Status:  cf.FailResponseCode,
			Message: "Invalid item_type. Allowed values: video, file, quiz, slide",
		})
	}

	moduleItemListParams := &param.ModuleItemListParams{
		ModuleID: createModuleItemParams.ModuleID,
	}
	_, totalItems, err := ctr.ModuleItemRepo.GetModuleItems(moduleItemListParams)
	if err != nil {
		ctr.Logger.Warnf("Error fetching existing module items: %v", err)
		return c.JSON(http.StatusBadRequest, cf.JsonResponse{
			Status:  cf.FailResponseCode,
			Message: "Error fetching existing module items",
		})
	}

	if totalItems == 0 {
		createModuleItemParams.Position = 1
	} else {
		createModuleItemParams.Position = totalItems + 1
	}

	var quizID int = 0

	if createModuleItemParams.ItemType == "video" {
		if !valid.IsURL(createModuleItemParams.Resource) {
			return c.JSON(http.StatusBadRequest, cf.JsonResponse{
				Status:  cf.FailResponseCode,
				Message: "Invalid URL format for video",
			})
		}

		parsedURL, err := url.Parse(createModuleItemParams.Resource)
		if err != nil {
			return c.JSON(http.StatusBadRequest, cf.JsonResponse{
				Status:  cf.FailResponseCode,
				Message: "Failed to parse video URL",
			})
		}

		queryParams := parsedURL.Query()
		videoId := queryParams.Get("v")
		if videoId == "" {
			return c.JSON(http.StatusBadRequest, cf.JsonResponse{
				Status:  cf.FailResponseCode,
				Message: "Invalid YouTube video URL: missing 'v' parameter",
			})
		}

		ytService := youtube.NewYouTubeService()
		videoInfo, err := ytService.GetVideoDetails(videoId)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, cf.JsonResponse{
				Status:  cf.FailResponseCode,
				Message: "Failed to fetch video details",
			})
		}

		requiredTimeInSeconds := utils.CalculateRequiredTime(videoInfo.Duration)

		createModuleItemParams.RequiredTime = requiredTimeInSeconds
		createModuleItemParams.Resource = videoId
	} else if createModuleItemParams.ItemType == "file" || createModuleItemParams.ItemType == "slide" {
		parts := strings.SplitN(createModuleItemParams.Resource, ",", 2)
		if len(parts) != 2 {
			return c.JSON(http.StatusOK, cf.JsonResponse{
				Status:  cf.FailResponseCode,
				Message: "Invalid File Format",
			})
		}

		mimeType := parts[0]
		base64Data := parts[1]

		formatFile := ""
		if strings.HasPrefix(mimeType, "data:application/") {
			formatFile = strings.TrimPrefix(mimeType, "data:application/")
			formatFile = strings.Split(formatFile, ";")[0]
		}

		if formatFile == "" {
			return c.JSON(http.StatusOK, cf.JsonResponse{
				Status:  cf.FailResponseCode,
				Message: "Invalid File Format",
			})
		}
		// Nếu là slide, có thể kiểm tra định dạng file slide (ví dụ: pdf, ppt, pptx)
		if createModuleItemParams.ItemType == "slide" {
			allowedSlides := []string{
				"pdf",               // PDF
				"ppt",               // PowerPoint 97-2003
				"pptx",              // PowerPoint hiện đại
				"vnd.ms-powerpoint", // MIME cho ppt
				"vnd.openxmlformats-officedocument.presentationml.presentation", // MIME cho pptx
				"vnd.oasis.opendocument.presentation",                           // ODP (OpenDocument Presentation)
				"x-pdf",                                                         // Một số trình duyệt gửi pdf là x-pdf
				"vnd.ms-powerpoint.presentation.macroenabled.12",                // pptm (PowerPoint Macro-Enabled Presentation)
				"vnd.ms-powerpoint.slideshow.macroenabled.12",                   // ppsm
				"vnd.ms-powerpoint.slideshow.macroEnabled.12",                   // ppsm (viết hoa khác)
				"vnd.ms-powerpoint.addin.macroenabled.12",                       // ppam
				"vnd.ms-powerpoint.template.macroenabled.12",                    // potm
				"vnd.openxmlformats-officedocument.presentationml.template",     // potx
				"vnd.openxmlformats-officedocument.presentationml.slideshow",    // ppsx
			}
			if _, check := utils.FindStringInArray(allowedSlides, formatFile); !check {
				return c.JSON(http.StatusOK, cf.JsonResponse{
					Status:  cf.FailResponseCode,
					Message: "Slide file not allowed",
				})
			}
		} else {
			if _, check := utils.FindStringInArray(cf.AllowFormatFileList, formatFile); !check {
				return c.JSON(http.StatusOK, cf.JsonResponse{
					Status:  cf.FailResponseCode,
					Message: "File not allowed",
				})
			}
		} // Tạo tên file bao gồm phần mở rộng để đảm bảo định dạng gốc
		millisecondTimeNow := int(time.Now().UnixNano() / int64(time.Millisecond))
		fileExt := ""

		// Xác định phần mở rộng phù hợp dựa trên MIME type
		// Xử lý cả file thường và slide để đảm bảo định dạng đúng
		if createModuleItemParams.ItemType == "file" || createModuleItemParams.ItemType == "slide" {
			// Kiểm tra định dạng file từ chuỗi mime type đầy đủ
			fullMimeType := mimeType
			// Xác định phần mở rộng dựa trên MIME type
			switch formatFile {
			// Định dạng file thông thường
			case "pdf", "x-pdf":
				fileExt = ".pdf"
			case "vnd.openxmlformats-officedocument.spreadsheetml.sheet":
				fileExt = ".xlsx"
			case "vnd.openxmlformats-officedocument.wordprocessingml.document":
				fileExt = ".docx"

			// Định dạng slide
			case "ppt", "vnd.ms-powerpoint":
				fileExt = ".ppt"
			case "pptx", "vnd.openxmlformats-officedocument.presentationml.presentation":
				fileExt = ".pptx"
			case "vnd.oasis.opendocument.presentation":
				fileExt = ".odp"
			case "vnd.ms-powerpoint.presentation.macroenabled.12":
				fileExt = ".pptm"
			case "vnd.ms-powerpoint.slideshow.macroenabled.12", "vnd.ms-powerpoint.slideshow.macroEnabled.12":
				fileExt = ".ppsm"
			case "vnd.openxmlformats-officedocument.presentationml.template":
				fileExt = ".potx"
			case "vnd.openxmlformats-officedocument.presentationml.slideshow":
				fileExt = ".ppsx"
			case "vnd.ms-powerpoint.slideshow":
				fileExt = ".pps"
			default:
				// Nếu không nhận dạng được format từ MIME type, thử đoán từ tên file gốc
				if strings.Contains(fullMimeType, "pdf") {
					fileExt = ".pdf"
				} else if strings.Contains(fullMimeType, "powerpoint") ||
					strings.Contains(fullMimeType, "presentation") {
					fileExt = ".pptx" // Mặc định cho PowerPoint hiện đại
				} else {
					fileExt = ".pptx" // Mặc định
					ctr.Logger.Warnf("Không nhận dạng được định dạng slide chính xác, dùng .pptx mặc định. MIME: %s", fullMimeType)
				}
			}

			// Log để theo dõi
			ctr.Logger.Infof("Slide file extension determined: %s from format: %s", fileExt, formatFile)
		}

		fileName := strconv.Itoa(createModuleItemParams.ModuleID) + "_" + strconv.Itoa(millisecondTimeNow) + fileExt

		// Đảm bảo header Content-Type được thiết lập chính xác
		err := ctr.cloud.UploadFileToCloud(
			base64Data,
			fileName,
			cf.FileFolderGCS,
		)
		if err != nil {
			ctr.Logger.Error(err)
			return c.JSON(http.StatusInternalServerError, cf.JsonResponse{
				Status:  cf.FailResponseCode,
				Message: "Failed to upload file to cloud",
			})
		}

		createModuleItemParams.Resource = fileName
	} else if createModuleItemParams.ItemType == "quiz" {
		if createModuleItemParams.QuizData == nil {
			return c.JSON(http.StatusBadRequest, cf.JsonResponse{
				Status:  cf.FailResponseCode,
				Message: "Quiz data is required for quiz item type",
			})
		}

		if createModuleItemParams.QuizData.QuestionType != cf.QuesMultipleChoice && createModuleItemParams.QuizData.QuestionType != cf.QuesEssay {
			return c.JSON(http.StatusBadRequest, cf.JsonResponse{
				Status:  cf.FailResponseCode,
				Message: "Invalid question type. Allowed values: 1 (multiple choice), 2 (essay)",
			})
		}

		if len(createModuleItemParams.QuizData.Questions) == 0 {
			return c.JSON(http.StatusBadRequest, cf.JsonResponse{
				Status:  cf.FailResponseCode,
				Message: "Quiz must have at least one question",
			})
		}

		totalWeight := 0.0
		for _, question := range createModuleItemParams.QuizData.Questions {
			totalWeight += question.Weight

			if createModuleItemParams.QuizData.QuestionType == cf.QuesMultipleChoice {
				if len(question.Options) < 2 {
					return c.JSON(http.StatusBadRequest, cf.JsonResponse{
						Status:  cf.FailResponseCode,
						Message: "Multiple choice questions must have at least two options",
					})
				}

				hasCorrectAnswer := false
				for _, option := range question.Options {
					if option.IsCorrect {
						hasCorrectAnswer = true
						break
					}
				}

				if !hasCorrectAnswer {
					return c.JSON(http.StatusBadRequest, cf.JsonResponse{
						Status:  cf.FailResponseCode,
						Message: "Multiple choice questions must have at least one correct answer",
					})
				}
			}
		}

		if totalWeight < 0.99 || totalWeight > 1.01 {
			return c.JSON(http.StatusBadRequest, cf.JsonResponse{
				Status:  cf.FailResponseCode,
				Message: "Question weights must sum to 1.0",
			})
		}

		var err error
		quizID, err = ctr.QuizRepo.CreateQuizWithQuestionsAndAnswers(
			createModuleItemParams.QuizData,
			createModuleItemParams.Title,
		)

		ctr.Logger.Infof("Quiz ID: %d", quizID)

		if err != nil {
			return c.JSON(http.StatusInternalServerError, cf.JsonResponse{
				Status:  cf.FailResponseCode,
				Message: "Failed to create quiz: " + err.Error(),
			})
		}

		createModuleItemParams.QuizID = quizID
		createModuleItemParams.Resource = ""
		createModuleItemParams.RequiredTime = 0
	}

	savedItem, err := ctr.ModuleItemRepo.SaveModuleItem(createModuleItemParams)
	if err != nil {
		ctr.Logger.Error(err)

		if quizID > 0 {
			if deleteErr := ctr.QuizRepo.DeleteQuiz(quizID); deleteErr != nil {
				ctr.Logger.Errorf("Failed to clean up quiz after module item creation failed: %v", deleteErr)
			}
		}

		return c.JSON(http.StatusInternalServerError, cf.JsonResponse{
			Status:  cf.FailResponseCode,
			Message: "Failed to save Module Item to database",
		})
	}

	moduleItemResponse := map[string]interface{}{
		"id":       savedItem.ID,
		"type":     savedItem.ItemType,
		"title":    savedItem.Title,
		"position": savedItem.Position,
	}

	if savedItem.ItemType == "quiz" {
		moduleItemResponse["quiz_id"] = quizID
	} else {
		moduleItemResponse["resource"] = savedItem.Resource
		moduleItemResponse["required_time"] = savedItem.RequiredTime
	}

	return c.JSON(http.StatusOK, cf.JsonResponse{
		Status:  cf.SuccessResponseCode,
		Message: "Module Item Created Successfully",
		Data:    moduleItemResponse,
	})
}

// DeleteModuleItem : delete module item by id
// Params : echo.Context
// Returns : object
func (ctr *ModuleItemController) DeleteModuleItem(c echo.Context) error {
	moduleItemIDParam := new(param.ModuleItemIDParam)
	if err := c.Bind(moduleItemIDParam); err != nil {
		return c.JSON(http.StatusOK, cf.JsonResponse{
			Status:  cf.FailResponseCode,
			Message: "Invalid Params",
			Data:    err,
		})
	}

	if _, err := valid.ValidateStruct(moduleItemIDParam); err != nil {
		return c.JSON(http.StatusOK, cf.JsonResponse{
			Status:  cf.FailResponseCode,
			Message: err.Error(),
		})
	}
	moduleItem, er := ctr.ModuleItemRepo.GetModuleItemByID(moduleItemIDParam.ModuleItemID)

	if er != nil {
		return c.JSON(http.StatusOK, cf.JsonResponse{
			Status:  cf.FailResponseCode,
			Message: "Module item not found",
			Data:    er,
		})
	}

	if moduleItem.ItemType == "file" || moduleItem.ItemType == "slide" {
		// Lấy tên file từ resource
		fileName := moduleItem.Resource

		// Đối với slide, tên file có thể đã được lưu với đuôi mở rộng
		// Nếu không có đuôi mở rộng, hệ thống vẫn xử lý được dựa trên tên file đã lưu trong DB
		err := ctr.cloud.DeleteFileCloud(fileName, cf.FileFolderGCS)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, cf.JsonResponse{
				Status:  cf.FailResponseCode,
				Message: "Failed to delete file from cloud",
				Data:    err,
			})
		}
	}
	err := ctr.ModuleItemRepo.DeleteModuleItem(moduleItemIDParam.ModuleItemID)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, cf.JsonResponse{
			Status:  cf.FailResponseCode,
			Message: "System Error",
		})
	}

	return c.JSON(http.StatusOK, cf.JsonResponse{
		Status:  cf.SuccessResponseCode,
		Message: "Deleted",
	})
}
