package lectures

import (
	"net/http"
	cf "orientation-training-api/configs"
	cm "orientation-training-api/internal/common"
	rp "orientation-training-api/internal/interfaces/repository"
	param "orientation-training-api/internal/interfaces/requestparams"
	response "orientation-training-api/internal/interfaces/response"
	m "orientation-training-api/internal/models"
	cld "orientation-training-api/internal/platform/cloud"
	"orientation-training-api/internal/platform/youtube"
	"os"

	valid "github.com/asaskevich/govalidator"
	"github.com/labstack/echo/v4"
)

type LectureController struct {
	cm.BaseController

	ModuleRepo       rp.ModuleRepository
	ModuleItemRepo   rp.ModuleItemRepository
	CourseRepo       rp.CourseRepository
	UserProgressRepo rp.UserProgressRepository
	QuizRepo         rp.QuizRepository
	Cloud            cld.StorageUtility
}

func NewLectureController(
	logger echo.Logger,
	moduleRepo rp.ModuleRepository,
	moduleItemRepo rp.ModuleItemRepository,
	courseRepo rp.CourseRepository,
	upRepo rp.UserProgressRepository,
	quizRepo rp.QuizRepository,
	cloud cld.StorageUtility) (ctr *LectureController) {

	ctr = &LectureController{
		cm.BaseController{},
		moduleRepo,
		moduleItemRepo,
		courseRepo,
		upRepo,
		quizRepo,
		cloud,
	}
	ctr.Init(logger)
	return
}

func (ctr *LectureController) GetLectureList(c echo.Context) error {
	userProfile := c.Get("user_profile").(m.User)
	lectureListParams := new(param.LectureListParams)

	if err := c.Bind(lectureListParams); err != nil {
		ctr.Logger.Errorf("Failed to bind params: %v", err)
		return c.JSON(http.StatusOK, cf.JsonResponse{
			Status:  cf.FailResponseCode,
			Message: "Invalid Params",
			Data:    err,
		})
	}

	if _, err := valid.ValidateStruct(lectureListParams); err != nil {
		ctr.Logger.Errorf("Validation failed: %v", err)
		return c.JSON(http.StatusOK, cf.JsonResponse{
			Status:  cf.FailResponseCode,
			Message: err.Error(),
		})
	}

	userProgress, err := ctr.UserProgressRepo.GetSingleUserProgress(userProfile.ID, lectureListParams.CourseID)
	if err != nil {
		ctr.Logger.Errorf("Failed to fetch user progress: %v", err)
		return c.JSON(http.StatusInternalServerError, cf.JsonResponse{
			Status:  cf.FailResponseCode,
			Message: "Failed to fetch user progress",
		})
	}

	currentModulePosition, currentModuleItemPosition := userProgress.ModulePosition, userProgress.ModuleItemPosition

	moduleListParams := &param.ModuleListParams{
		CourseID: lectureListParams.CourseID,
	}
	modules, _, err := ctr.ModuleRepo.GetModules(moduleListParams)
	if err != nil {
		ctr.Logger.Errorf("Failed to fetch modules: %v", err)
		return c.JSON(http.StatusInternalServerError, cf.JsonResponse{
			Status:  cf.FailResponseCode,
			Message: "Failed to fetch modules for the course",
		})
	}

	var allModuleItems []m.ModuleItem
	for _, module := range modules {
		moduleItems, err := ctr.ModuleItemRepo.GetModuleItemsByModuleID(module.ID)
		if err != nil {
			ctr.Logger.Errorf("Failed to fetch module items for module ID %d: %v", module.ID, err)
			return c.JSON(http.StatusInternalServerError, cf.JsonResponse{
				Status:  cf.FailResponseCode,
				Message: "Failed to fetch module items",
			})
		}
		allModuleItems = append(allModuleItems, moduleItems...)
	}

	moduleResponses := []response.LectureModuleResponse{}

	for _, module := range modules {
		moduleResponse := response.LectureModuleResponse{
			ModuleID:       module.ID,
			ModuleTitle:    module.Title,
			ModulePosition: module.Position,
			Duration:       module.Duration,
			Lectures:       []response.LectureItemResponse{},
		}

		for _, item := range allModuleItems {
			if item.ModuleID != module.ID {
				continue
			}

			isUnlocked := false
			if module.Position < currentModulePosition ||
				(module.Position == currentModulePosition && item.Position <= currentModuleItemPosition) {
				isUnlocked = true
			}

			lectureItem := response.LectureItemResponse{
				ModuleItemID:       item.ID,
				ModuleItemTitle:    item.Title,
				ModuleItemPosition: item.Position,
				ItemType:           item.ItemType,
				Unlocked:           isUnlocked,
			}

			if item.ItemType == "video" {
				videoID := item.Resource
				ytService := youtube.NewYouTubeService()

				videoInfo, err := ytService.GetVideoDetails(videoID)
				if err != nil {
					ctr.Logger.Errorf("Failed to fetch video details for video ID %s: %v. Module ID: %d, ModuleItem ID: %d, Error details: %+v",
						videoID, err, module.ID, item.ID, err)
					continue
				}

				videoContent := response.VideoContentResponse{
					VideoID:      videoID,
					Duration:     videoInfo.Duration,
					RequiredTime: item.RequiredTime,
					Thumbnail:    videoInfo.ThumbnailURL,
					PublishedAt:  videoInfo.PublishedAt,
				}
				lectureItem.Content = videoContent
			} else if item.ItemType == "file" || item.ItemType == "slide" {
				var filePath string
				if item.Resource != "" {
					filePath = "https://storage.googleapis.com/" + os.Getenv("GOOGLE_STORAGE_BUCKET") + "/" + cf.FileFolderGCS + item.Resource
				}

				fileContent := response.FileContentResponse{
					FilePath:     filePath,
					Duration:     item.RequiredTime,
					RequiredTime: item.RequiredTime,
				}
				lectureItem.Content = fileContent
			} else if item.ItemType == "quiz" && item.QuizID > 0 {
				quiz, err := ctr.QuizRepo.GetQuizByID(item.QuizID)
				if err != nil {
					ctr.Logger.Errorf("Failed to fetch quiz details for quiz ID %d: %v", item.QuizID, err)
					continue
				}

				questions, err := ctr.QuizRepo.GetQuizQuestionsWithAnswers(item.QuizID)
				if err != nil {
					ctr.Logger.Errorf("Failed to fetch quiz questions for quiz ID %d: %v", item.QuizID, err)
					continue
				}

				quizContent := response.QuizContentResponse{
					QuizID:     quiz.ID,
					QuizTitle:  item.Title,
					Difficulty: cf.DifficultyLabels[quiz.Difficulty],
					TotalScore: quiz.TotalScore,
					TimeLimit:  quiz.TimeLimit,
				}

				if len(questions) > 0 && questions[0].QuestionType == cf.QuestionTypeEssay {
					quizContent.QuizType = "essay"
					quizContent.Questions = []response.QuizQuestionResponse{}

					for _, q := range questions {
						questionResponse := response.QuizQuestionResponse{
							QuestionID:   q.ID,
							QuestionText: q.QuestionText,
							Points:       q.Weight * quiz.TotalScore,
							Options:      []response.QuizOptionResponse{},
						}

						quizContent.Questions = append(quizContent.Questions, questionResponse)
					}
				} else {
					quizContent.QuizType = "multiple_choice"
					quizContent.Questions = []response.QuizQuestionResponse{}

					for _, q := range questions {
						questionResponse := response.QuizQuestionResponse{
							QuestionID:    q.ID,
							QuestionText:  q.QuestionText,
							AllowMultiple: q.IsMultipleCorrect,
							Points:        q.Weight * quiz.TotalScore,
							Options:       []response.QuizOptionResponse{},
						}

						for _, a := range q.Answers {
							option := response.QuizOptionResponse{
								ID:   a.ID,
								Text: a.AnswerText,
							}
							questionResponse.Options = append(questionResponse.Options, option)
						}

						quizContent.Questions = append(quizContent.Questions, questionResponse)
					}
				}

				lectureItem.Content = quizContent
			}

			moduleResponse.Lectures = append(moduleResponse.Lectures, lectureItem)
		}

		moduleResponses = append(moduleResponses, moduleResponse)
	}

	return c.JSON(http.StatusOK, cf.JsonResponse{
		Status:  cf.SuccessResponseCode,
		Message: "Success",
		Data:    moduleResponses,
	})
}
