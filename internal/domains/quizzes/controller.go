package quizzes

import (
	"net/http"
	cf "orientation-training-api/configs"
	cm "orientation-training-api/internal/common"
	rp "orientation-training-api/internal/interfaces/repository"
	param "orientation-training-api/internal/interfaces/requestparams"
	m "orientation-training-api/internal/models"

	valid "github.com/asaskevich/govalidator"
	"github.com/labstack/echo/v4"
)

type QuizController struct {
	cm.BaseController
	QuizRepo rp.QuizRepository
}

func NewQuizController(logger echo.Logger, quizRepo rp.QuizRepository) (ctr *QuizController) {
	ctr = &QuizController{cm.BaseController{}, quizRepo}
	ctr.Init(logger)
	return
}

// CreateQuiz creates a new quiz
func (ctr *QuizController) CreateQuiz(c echo.Context) error {
	// Only managers can create quizzes
	userProfile := c.Get("user_profile").(m.User)
	if userProfile.RoleID != cf.ManagerRoleID {
		return c.JSON(http.StatusForbidden, cf.JsonResponse{
			Status:  cf.FailResponseCode,
			Message: "Only managers can create quizzes",
		})
	}

	createQuizParams := new(param.CreateQuizParams)
	if err := c.Bind(createQuizParams); err != nil {
		ctr.Logger.Errorf("Failed to bind params: %v", err)
		return c.JSON(http.StatusOK, cf.JsonResponse{
			Status:  cf.FailResponseCode,
			Message: "Invalid params",
			Data:    err,
		})
	}

	if _, err := valid.ValidateStruct(createQuizParams); err != nil {
		ctr.Logger.Errorf("Validation failed: %v", err)
		return c.JSON(http.StatusOK, cf.JsonResponse{
			Status:  cf.FailResponseCode,
			Message: err.Error(),
		})
	}

	quiz := &m.Quiz{
		Title:      createQuizParams.Title,
		Difficulty: createQuizParams.Difficulty,
		TotalScore: createQuizParams.TotalScore,
		TimeLimit:  createQuizParams.TimeLimit,
	}

	err := ctr.QuizRepo.SaveQuiz(quiz)
	if err != nil {
		ctr.Logger.Errorf("Failed to create quiz: %v", err)
		return c.JSON(http.StatusInternalServerError, cf.JsonResponse{
			Status:  cf.FailResponseCode,
			Message: "Failed to create quiz",
		})
	}

	return c.JSON(http.StatusOK, cf.JsonResponse{
		Status:  cf.SuccessResponseCode,
		Message: "Quiz created successfully",
		Data:    quiz,
	})
}

// UpdateQuiz updates an existing quiz
func (ctr *QuizController) UpdateQuiz(c echo.Context) error {
	// Only managers can update quizzes
	userProfile := c.Get("user_profile").(m.User)
	if userProfile.RoleID != cf.ManagerRoleID {
		return c.JSON(http.StatusForbidden, cf.JsonResponse{
			Status:  cf.FailResponseCode,
			Message: "Only managers can update quizzes",
		})
	}

	updateQuizParams := new(param.UpdateQuizParams)
	if err := c.Bind(updateQuizParams); err != nil {
		ctr.Logger.Errorf("Failed to bind params: %v", err)
		return c.JSON(http.StatusOK, cf.JsonResponse{
			Status:  cf.FailResponseCode,
			Message: "Invalid params",
			Data:    err,
		})
	}

	if _, err := valid.ValidateStruct(updateQuizParams); err != nil {
		ctr.Logger.Errorf("Validation failed: %v", err)
		return c.JSON(http.StatusOK, cf.JsonResponse{
			Status:  cf.FailResponseCode,
			Message: err.Error(),
		})
	}

	// Check if quiz exists
	existingQuiz, err := ctr.QuizRepo.GetQuizByID(updateQuizParams.ID)
	if err != nil {
		ctr.Logger.Errorf("Quiz not found: %v", err)
		return c.JSON(http.StatusNotFound, cf.JsonResponse{
			Status:  cf.FailResponseCode,
			Message: "Quiz not found",
		})
	}

	// Update quiz properties
	existingQuiz.Title = updateQuizParams.Title
	existingQuiz.Difficulty = updateQuizParams.Difficulty
	existingQuiz.TotalScore = updateQuizParams.TotalScore
	existingQuiz.TimeLimit = updateQuizParams.TimeLimit

	err = ctr.QuizRepo.SaveQuiz(&existingQuiz)
	if err != nil {
		ctr.Logger.Errorf("Failed to update quiz: %v", err)
		return c.JSON(http.StatusInternalServerError, cf.JsonResponse{
			Status:  cf.FailResponseCode,
			Message: "Failed to update quiz",
		})
	}

	return c.JSON(http.StatusOK, cf.JsonResponse{
		Status:  cf.SuccessResponseCode,
		Message: "Quiz updated successfully",
		Data:    existingQuiz,
	})
}

// GetQuizList retrieves a list of quizzes with pagination
func (ctr *QuizController) GetQuizList(c echo.Context) error {
	quizListParams := new(param.QuizListParams)
	if err := c.Bind(quizListParams); err != nil {
		ctr.Logger.Errorf("Failed to bind params: %v", err)
		return c.JSON(http.StatusOK, cf.JsonResponse{
			Status:  cf.FailResponseCode,
			Message: "Invalid params",
			Data:    err,
		})
	}

	// Set default current page if not provided
	if quizListParams.CurrentPage <= 0 {
		quizListParams.CurrentPage = 1
	}

	// Set default row per page if not provided
	if quizListParams.RowPerPage <= 0 {
		quizListParams.RowPerPage = 10
	}

	quizzes, total, err := ctr.QuizRepo.GetQuizList(quizListParams)
	if err != nil {
		ctr.Logger.Errorf("Failed to fetch quiz list: %v", err)
		return c.JSON(http.StatusInternalServerError, cf.JsonResponse{
			Status:  cf.FailResponseCode,
			Message: "Failed to fetch quiz list",
		})
	}

	return c.JSON(http.StatusOK, cf.JsonResponse{
		Status:  cf.SuccessResponseCode,
		Message: "Quiz list retrieved successfully",
		Data: map[string]interface{}{
			"quizzes":      quizzes,
			"total":        total,
			"current_page": quizListParams.CurrentPage,
			"row_per_page": quizListParams.RowPerPage,
		},
	})
}

// GetQuizDetail retrieves a quiz with its questions and answers
func (ctr *QuizController) GetQuizDetail(c echo.Context) error {
	getQuizParams := new(param.GetQuizParams)
	if err := c.Bind(getQuizParams); err != nil {
		ctr.Logger.Errorf("Failed to bind params: %v", err)
		return c.JSON(http.StatusOK, cf.JsonResponse{
			Status:  cf.FailResponseCode,
			Message: "Invalid params",
			Data:    err,
		})
	}

	if _, err := valid.ValidateStruct(getQuizParams); err != nil {
		ctr.Logger.Errorf("Validation failed: %v", err)
		return c.JSON(http.StatusOK, cf.JsonResponse{
			Status:  cf.FailResponseCode,
			Message: err.Error(),
		})
	}

	// Get quiz details
	quiz, err := ctr.QuizRepo.GetQuizByID(getQuizParams.QuizID)
	if err != nil {
		ctr.Logger.Errorf("Failed to fetch quiz: %v", err)
		return c.JSON(http.StatusNotFound, cf.JsonResponse{
			Status:  cf.FailResponseCode,
			Message: "Quiz not found",
		})
	}

	// Get questions with answers
	questions, err := ctr.QuizRepo.GetQuizQuestionsWithAnswers(getQuizParams.QuizID)
	if err != nil {
		ctr.Logger.Errorf("Failed to fetch quiz questions: %v", err)
		return c.JSON(http.StatusInternalServerError, cf.JsonResponse{
			Status:  cf.FailResponseCode,
			Message: "Failed to fetch quiz questions",
		})
	}

	// For non-admin users, remove the "is_correct" field from answers
	userProfile := c.Get("user_profile").(m.User)
	if userProfile.RoleID != cf.ManagerRoleID {
		for i := range questions {
			for j := range questions[i].Answers {
				questions[i].Answers[j].IsCorrect = false
			}
		}
	}

	return c.JSON(http.StatusOK, cf.JsonResponse{
		Status:  cf.SuccessResponseCode,
		Message: "Quiz details retrieved successfully",
		Data: map[string]interface{}{
			"quiz":      quiz,
			"questions": questions,
		},
	})
}

// DeleteQuiz deletes a quiz
func (ctr *QuizController) DeleteQuiz(c echo.Context) error {
	// Only managers can delete quizzes
	userProfile := c.Get("user_profile").(m.User)
	if userProfile.RoleID != cf.ManagerRoleID {
		return c.JSON(http.StatusForbidden, cf.JsonResponse{
			Status:  cf.FailResponseCode,
			Message: "Only managers can delete quizzes",
		})
	}

	deleteQuizParams := new(param.DeleteQuizParams)
	if err := c.Bind(deleteQuizParams); err != nil {
		ctr.Logger.Errorf("Failed to bind params: %v", err)
		return c.JSON(http.StatusOK, cf.JsonResponse{
			Status:  cf.FailResponseCode,
			Message: "Invalid params",
			Data:    err,
		})
	}

	if _, err := valid.ValidateStruct(deleteQuizParams); err != nil {
		ctr.Logger.Errorf("Validation failed: %v", err)
		return c.JSON(http.StatusOK, cf.JsonResponse{
			Status:  cf.FailResponseCode,
			Message: err.Error(),
		})
	}

	// Check if quiz exists
	_, err := ctr.QuizRepo.GetQuizByID(deleteQuizParams.QuizID)
	if err != nil {
		ctr.Logger.Errorf("Quiz not found: %v", err)
		return c.JSON(http.StatusNotFound, cf.JsonResponse{
			Status:  cf.FailResponseCode,
			Message: "Quiz not found",
		})
	}

	// Delete quiz
	err = ctr.QuizRepo.DeleteQuiz(deleteQuizParams.QuizID)
	if err != nil {
		ctr.Logger.Errorf("Failed to delete quiz: %v", err)
		return c.JSON(http.StatusInternalServerError, cf.JsonResponse{
			Status:  cf.FailResponseCode,
			Message: "Failed to delete quiz",
		})
	}

	return c.JSON(http.StatusOK, cf.JsonResponse{
		Status:  cf.SuccessResponseCode,
		Message: "Quiz deleted successfully",
	})
}

// CreateQuizQuestion creates or updates a quiz question with its answers
func (ctr *QuizController) CreateQuizQuestion(c echo.Context) error {
	// Only managers can create quiz questions
	userProfile := c.Get("user_profile").(m.User)
	if userProfile.RoleID != cf.ManagerRoleID {
		return c.JSON(http.StatusForbidden, cf.JsonResponse{
			Status:  cf.FailResponseCode,
			Message: "Only managers can create quiz questions",
		})
	}

	createQuestionParams := new(param.CreateQuizQuestionParams)
	if err := c.Bind(createQuestionParams); err != nil {
		ctr.Logger.Errorf("Failed to bind params: %v", err)
		return c.JSON(http.StatusOK, cf.JsonResponse{
			Status:  cf.FailResponseCode,
			Message: "Invalid params",
			Data:    err,
		})
	}

	if _, err := valid.ValidateStruct(createQuestionParams); err != nil {
		ctr.Logger.Errorf("Validation failed: %v", err)
		return c.JSON(http.StatusOK, cf.JsonResponse{
			Status:  cf.FailResponseCode,
			Message: err.Error(),
		})
	}

	// Check if quiz exists
	_, err := ctr.QuizRepo.GetQuizByID(createQuestionParams.QuizID)
	if err != nil {
		ctr.Logger.Errorf("Quiz not found: %v", err)
		return c.JSON(http.StatusNotFound, cf.JsonResponse{
			Status:  cf.FailResponseCode,
			Message: "Quiz not found",
		})
	}

	// Create question object
	question := &m.QuizQuestion{
		QuizID:            createQuestionParams.QuizID,
		QuestionType:      createQuestionParams.QuestionType,
		QuestionText:      createQuestionParams.QuestionText,
		Weight:            createQuestionParams.Weight,
		IsMultipleCorrect: createQuestionParams.IsMultipleCorrect,
	}

	// Create answer objects
	answers := make([]m.QuizAnswer, len(createQuestionParams.Answers))
	for i, answerParam := range createQuestionParams.Answers {
		answers[i] = m.QuizAnswer{
			AnswerText: answerParam.AnswerText,
			IsCorrect:  answerParam.IsCorrect,
		}
	}

	// Save question and answers
	err = ctr.QuizRepo.SaveQuizQuestion(question, answers)
	if err != nil {
		ctr.Logger.Errorf("Failed to save quiz question: %v", err)
		return c.JSON(http.StatusInternalServerError, cf.JsonResponse{
			Status:  cf.FailResponseCode,
			Message: "Failed to save quiz question",
		})
	}

	return c.JSON(http.StatusOK, cf.JsonResponse{
		Status:  cf.SuccessResponseCode,
		Message: "Quiz question created successfully",
		Data:    question,
	})
}

// SubmitFullQuiz submits all answers for a quiz at once
func (ctr *QuizController) SubmitFullQuiz(c echo.Context) error {
	userProfile := c.Get("user_profile").(m.User)
	submitParams := new(param.SubmitFullQuizParams)

	if err := c.Bind(submitParams); err != nil {
		ctr.Logger.Errorf("Failed to bind params: %v", err)
		return c.JSON(http.StatusOK, cf.JsonResponse{
			Status:  cf.FailResponseCode,
			Message: "Invalid params",
			Data:    err,
		})
	}

	if _, err := valid.ValidateStruct(submitParams); err != nil {
		ctr.Logger.Errorf("Validation failed: %v", err)
		return c.JSON(http.StatusOK, cf.JsonResponse{
			Status:  cf.FailResponseCode,
			Message: err.Error(),
		})
	}

	maxAttempt, err := ctr.QuizRepo.GetMaxQuizAttempt(userProfile.ID, submitParams.QuizID)
	if err != nil {
		ctr.Logger.Errorf("Error getting max attempt: %v", err)
		return c.JSON(http.StatusInternalServerError, cf.JsonResponse{
			Status:  cf.FailResponseCode,
			Message: "Failed to get attempt history",
		})
	}

	currentAttempt := maxAttempt + 1

	questions, err := ctr.QuizRepo.GetQuizQuestionsWithAnswers(submitParams.QuizID)
	if err != nil {
		ctr.Logger.Errorf("Failed to fetch quiz questions: %v", err)
		return c.JSON(http.StatusInternalServerError, cf.JsonResponse{
			Status:  cf.FailResponseCode,
			Message: "Failed to fetch quiz questions",
		})
	}

	quiz, err := ctr.QuizRepo.GetQuizByID(submitParams.QuizID)
	if err != nil {
		ctr.Logger.Errorf("Failed to fetch quiz details: %v", err)
		return c.JSON(http.StatusInternalServerError, cf.JsonResponse{
			Status:  cf.FailResponseCode,
			Message: "Failed to fetch quiz details",
		})
	}

	questionsMap := make(map[int]m.QuizQuestion)
	for _, q := range questions {
		questionsMap[q.ID] = q
	}

	totalScore := 0.0
	submissionDetails := []map[string]interface{}{}
	essaySubmissions := []map[string]interface{}{}
	hasEssayQuestions := false

	for _, answer := range submitParams.Answers {
		question, exists := questionsMap[answer.QuestionID]
		if !exists {
			ctr.Logger.Warnf("Question ID %d not found in quiz %d", answer.QuestionID, submitParams.QuizID)
			continue
		}

		score := 0.0
		isCorrect := false
		correctAnswersList := []int{}
		reviewed := false

		if question.QuestionType == cf.QuestionTypeMultipleChoice {
			correctAnswerIDs := make(map[int]bool)
			for _, a := range question.Answers {
				if a.IsCorrect {
					correctAnswerIDs[a.ID] = true
					correctAnswersList = append(correctAnswersList, a.ID)
				}
			}

			isCorrect = true
			if len(answer.SelectedAnswerIds) != len(correctAnswerIDs) {
				isCorrect = false
			} else {
				for _, selectedID := range answer.SelectedAnswerIds {
					if !correctAnswerIDs[selectedID] {
						isCorrect = false
						break
					}
				}
			}

			if isCorrect {
				score = question.Weight * quiz.TotalScore
			}

			reviewed = true
		} else if question.QuestionType == cf.QuestionTypeEssay {
			hasEssayQuestions = true
			score = 0 // Initially 0, to be updated by manager
			reviewed = false
		}

		submission := &m.QuizSubmission{
			UserID:            userProfile.ID,
			QuizID:            submitParams.QuizID,
			QuizQuestionID:    answer.QuestionID,
			AnswerText:        answer.AnswerText,
			SelectedAnswerIds: answer.SelectedAnswerIds,
			Score:             score,
			Attempt:           currentAttempt,
			Reviewed:          reviewed,
		}

		err = ctr.QuizRepo.SaveQuizSubmission(submission)
		if err != nil {
			ctr.Logger.Errorf("Failed to save quiz submission for question %d: %v", answer.QuestionID, err)
			continue
		}

		if question.QuestionType == cf.QuestionTypeMultipleChoice {
			submissionDetail := map[string]interface{}{
				"question_id":      answer.QuestionID,
				"score":            score,
				"is_correct":       isCorrect,
				"selected_answers": answer.SelectedAnswerIds,
				"correct_answers":  correctAnswersList,
				"reviewed":         reviewed,
				"attempt":          currentAttempt,
				"explanation":      question.Explanation,
			}
			submissionDetails = append(submissionDetails, submissionDetail)
		} else if question.QuestionType == cf.QuestionTypeEssay {
			essaySubmission := map[string]interface{}{
				"question_id": answer.QuestionID,
				"reviewed":    reviewed,
				"attempt":     currentAttempt,
			}
			essaySubmissions = append(essaySubmissions, essaySubmission)
		}

		totalScore += score
	}

	// Determine if passed (usually 70% is passing)
	passThreshold := quiz.TotalScore * 0.7
	passed := totalScore >= passThreshold

	responseData := map[string]interface{}{
		"quiz_id":        submitParams.QuizID,
		"total_score":    totalScore,
		"max_score":      quiz.TotalScore,
		"passed":         passed,
		"pass_threshold": passThreshold,
		"submissions":    submissionDetails,
		"attempt":        currentAttempt,
	}

	if hasEssayQuestions {
		responseData["passed"] = true
		responseData["essay_submissions"] = essaySubmissions
	}

	return c.JSON(http.StatusOK, cf.JsonResponse{
		Status:  cf.SuccessResponseCode,
		Message: "Quiz answers submitted successfully",
		Data:    responseData,
	})
}

// GetQuizResults retrieves quiz results for a user
func (ctr *QuizController) GetQuizResults(c echo.Context) error {
	userProfile := c.Get("user_profile").(m.User)
	getResultsParams := new(param.GetQuizResultsParams)

	if err := c.Bind(getResultsParams); err != nil {
		ctr.Logger.Errorf("Failed to bind params: %v", err)
		return c.JSON(http.StatusOK, cf.JsonResponse{
			Status:  cf.FailResponseCode,
			Message: "Invalid params",
			Data:    err,
		})
	}

	if _, err := valid.ValidateStruct(getResultsParams); err != nil {
		ctr.Logger.Errorf("Validation failed: %v", err)
		return c.JSON(http.StatusOK, cf.JsonResponse{
			Status:  cf.FailResponseCode,
			Message: err.Error(),
		})
	}

	targetUserID := userProfile.ID
	if userProfile.RoleID == cf.ManagerRoleID && getResultsParams.UserID > 0 {
		targetUserID = getResultsParams.UserID
	}

	maxAttempt, err := ctr.QuizRepo.GetMaxQuizAttempt(targetUserID, getResultsParams.QuizID)
	if err != nil {
		ctr.Logger.Errorf("Error getting max attempt: %v", err)
		return c.JSON(http.StatusInternalServerError, cf.JsonResponse{
			Status:  cf.FailResponseCode,
			Message: "Failed to get attempt history",
		})
	}

	if maxAttempt == 0 {
		return c.JSON(http.StatusOK, cf.JsonResponse{
			Status:  cf.SuccessResponseCode,
			Message: "No quiz attempts found",
			Data: map[string]interface{}{
				"passed": false,
				"results": map[string]interface{}{
					"quiz_id": getResultsParams.QuizID,
					"answers": []map[string]interface{}{},
				},
			},
		})
	}

	quiz, err := ctr.QuizRepo.GetQuizByID(getResultsParams.QuizID)
	if err != nil {
		ctr.Logger.Errorf("Quiz not found: %v", err)
		return c.JSON(http.StatusNotFound, cf.JsonResponse{
			Status:  cf.FailResponseCode,
			Message: "Quiz not found",
		})
	}

	questions, err := ctr.QuizRepo.GetQuizQuestionsWithAnswers(getResultsParams.QuizID)
	if err != nil {
		ctr.Logger.Errorf("Failed to fetch quiz questions: %v", err)
		return c.JSON(http.StatusInternalServerError, cf.JsonResponse{
			Status:  cf.FailResponseCode,
			Message: "Failed to fetch quiz questions",
		})
	}

	questionsMap := make(map[int]m.QuizQuestion)
	for _, q := range questions {
		questionsMap[q.ID] = q
	}

	submissions, err := ctr.QuizRepo.GetQuizSubmissionsByUser(targetUserID, getResultsParams.QuizID)
	if err != nil {
		ctr.Logger.Errorf("Failed to fetch quiz submissions: %v", err)
		return c.JSON(http.StatusInternalServerError, cf.JsonResponse{
			Status:  cf.FailResponseCode,
			Message: "Failed to fetch quiz submissions",
		})
	}

	latestAttemptSubmissions := []m.QuizSubmission{}
	for _, submission := range submissions {
		if submission.Attempt == maxAttempt {
			latestAttemptSubmissions = append(latestAttemptSubmissions, submission)
		}
	}

	totalScore := 0.0
	hasEssayQuestions := false
	allEssaysReviewed := true
	answers := []map[string]interface{}{}

	for _, submission := range latestAttemptSubmissions {
		question, exists := questionsMap[submission.QuizQuestionID]
		if !exists {
			continue
		}

		if question.QuestionType == cf.QuestionTypeEssay {
			hasEssayQuestions = true
			if !submission.Reviewed {
				allEssaysReviewed = false
			}

			answerData := map[string]interface{}{
				"question_id": submission.QuizQuestionID,
				"answer_text": submission.AnswerText,
			}
			answers = append(answers, answerData)
		} else if question.QuestionType == cf.QuestionTypeMultipleChoice {
			isCorrect := false
			correctAnswerIDs := []int{}

			for _, a := range question.Answers {
				if a.IsCorrect {
					correctAnswerIDs = append(correctAnswerIDs, a.ID)
				}
			}

			if len(submission.SelectedAnswerIds) == len(correctAnswerIDs) {
				isCorrect = true
				selectedMap := make(map[int]bool)
				for _, id := range submission.SelectedAnswerIds {
					selectedMap[id] = true
				}

				for _, id := range correctAnswerIDs {
					if !selectedMap[id] {
						isCorrect = false
						break
					}
				}
			}

			answerData := map[string]interface{}{
				"question_id":         submission.QuizQuestionID,
				"selected_answer_ids": submission.SelectedAnswerIds,
				"is_correct":          isCorrect,
				"correct_answer_ids":  correctAnswerIDs,
				"explanation":         question.Explanation,
				"points":              question.Weight * quiz.TotalScore,
			}
			answers = append(answers, answerData)
		}

		totalScore += submission.Score
	}

	passThreshold := quiz.TotalScore * 0.7
	passed := totalScore >= passThreshold || hasEssayQuestions

	resultsData := map[string]interface{}{
		"quiz_id": getResultsParams.QuizID,
		"answers": answers,
		"attempt": maxAttempt,
	}

	if !hasEssayQuestions {
		// Case 1: Multiple choice quiz
		resultsData["total_score"] = quiz.TotalScore
		resultsData["user_score"] = totalScore
	} else if hasEssayQuestions && !allEssaysReviewed {
		// Case 2: Essay quiz not yet reviewed
		// No additional fields needed
	} else if hasEssayQuestions && allEssaysReviewed {
		// Case 3: Essay quiz with review and feedback
		resultsData["user_score"] = totalScore

		for _, submission := range latestAttemptSubmissions {
			if submission.Feedback != "" {
				resultsData["feedback"] = submission.Feedback
				break
			}
		}
	}

	return c.JSON(http.StatusOK, cf.JsonResponse{
		Status:  cf.SuccessResponseCode,
		Message: "Quiz results retrieved successfully",
		Data: map[string]interface{}{
			"passed":  passed,
			"results": resultsData,
		},
	})
}
