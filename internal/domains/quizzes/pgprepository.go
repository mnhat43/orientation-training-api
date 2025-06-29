package quizzes

import (
	cf "orientation-training-api/configs"
	cm "orientation-training-api/internal/common"
	param "orientation-training-api/internal/interfaces/requestparams"
	m "orientation-training-api/internal/models"

	"github.com/go-pg/pg/v9"
	"github.com/labstack/echo/v4"
)

type PgQuizRepository struct {
	cm.AppRepository
}

func NewPgQuizRepository(logger echo.Logger) (repo *PgQuizRepository) {
	repo = &PgQuizRepository{}
	repo.Init(logger)
	return
}

// GetQuizByID fetches quiz data by ID
func (repo *PgQuizRepository) GetQuizByID(quizID int) (m.Quiz, error) {
	quiz := m.Quiz{}

	err := repo.DB.Model(&quiz).
		Where("id = ?", quizID).
		Where("deleted_at IS NULL").
		Select()

	if err != nil {
		repo.Logger.Errorf("Error fetching quiz with ID %d: %v", quizID, err)
		return quiz, err
	}

	return quiz, nil
}

// GetQuizList retrieves quizzes with pagination
func (repo *PgQuizRepository) GetQuizList(params *param.QuizListParams) ([]m.Quiz, int, error) {
	var quizzes []m.Quiz
	var count int

	query := repo.DB.Model(&quizzes).
		Where("deleted_at IS NULL")

	if params.Title != "" {
		query = query.Where("title ILIKE ?", "%"+params.Title+"%")
	}

	// Get total count
	count, err := query.Count()
	if err != nil {
		repo.Logger.Errorf("Error counting quizzes: %v", err)
		return nil, 0, err
	}

	// Apply pagination
	if params.RowPerPage > 0 {
		offset := (params.CurrentPage - 1) * params.RowPerPage
		query = query.Limit(params.RowPerPage).Offset(offset)
	}

	// Apply ordering
	query = query.Order("created_at DESC")

	err = query.Select()
	if err != nil {
		repo.Logger.Errorf("Error fetching quiz list: %v", err)
		return nil, 0, err
	}

	return quizzes, count, nil
}

// SaveQuiz creates or updates a quiz
func (repo *PgQuizRepository) SaveQuiz(quiz *m.Quiz) error {
	if quiz.ID == 0 {
		// Create new quiz
		_, err := repo.DB.Model(quiz).Insert()
		if err != nil {
			repo.Logger.Errorf("Error creating quiz: %v", err)
			return err
		}
	} else {
		// Update existing quiz
		_, err := repo.DB.Model(quiz).
			WherePK().
			Where("deleted_at IS NULL").
			Update()
		if err != nil {
			repo.Logger.Errorf("Error updating quiz with ID %d: %v", quiz.ID, err)
			return err
		}
	}
	return nil
}

// DeleteQuiz soft deletes a quiz
func (repo *PgQuizRepository) DeleteQuiz(quizID int) error {
	_, err := repo.DB.Model((*m.Quiz)(nil)).
		Set("deleted_at = NOW()").
		Where("id = ?", quizID).
		Where("deleted_at IS NULL").
		Update()

	if err != nil {
		repo.Logger.Errorf("Error deleting quiz with ID %d: %v", quizID, err)
		return err
	}
	return nil
}

// GetQuizQuestionsWithAnswers fetches all questions and answers for a quiz
func (repo *PgQuizRepository) GetQuizQuestionsWithAnswers(quizID int) ([]m.QuizQuestion, error) {
	// First, fetch all questions for this quiz
	var questions []m.QuizQuestion

	err := repo.DB.Model(&questions).
		Where("quiz_id = ?", quizID).
		Where("deleted_at IS NULL").
		Order("id ASC").
		Select()

	if err != nil {
		repo.Logger.Errorf("Error fetching questions for quiz ID %d: %v", quizID, err)
		return nil, err
	}

	// Then, for each question, fetch its answers
	for i := range questions {
		var answers []m.QuizAnswer
		err := repo.DB.Model(&answers).
			Where("quiz_question_id = ?", questions[i].ID).
			Where("deleted_at IS NULL").
			Order("id ASC").
			Select()

		if err != nil {
			repo.Logger.Errorf("Error fetching answers for question ID %d: %v", questions[i].ID, err)
			continue
		}

		questions[i].Answers = answers
	}

	return questions, nil
}

// SaveQuizQuestion saves a quiz question and its answers
func (repo *PgQuizRepository) SaveQuizQuestion(question *m.QuizQuestion, answers []m.QuizAnswer) error {
	// Begin transaction
	tx, err := repo.DB.Begin()
	if err != nil {
		repo.Logger.Errorf("Error beginning transaction: %v", err)
		return err
	}
	defer tx.Rollback()

	// Save question
	if question.ID == 0 {
		// Create new question
		_, err = tx.Model(question).Insert()
	} else {
		// Update existing question
		_, err = tx.Model(question).
			WherePK().
			Where("deleted_at IS NULL").
			Update()
	}

	if err != nil {
		repo.Logger.Errorf("Error saving question: %v", err)
		return err
	}

	// Delete existing answers if updating
	if question.ID > 0 {
		_, err = tx.Model((*m.QuizAnswer)(nil)).
			Where("quiz_question_id = ?", question.ID).
			Delete()

		if err != nil {
			repo.Logger.Errorf("Error deleting existing answers: %v", err)
			return err
		}
	}

	// Insert new answers
	for i := range answers {
		answers[i].QuizQuestionID = question.ID
		_, err = tx.Model(&answers[i]).Insert()
		if err != nil {
			repo.Logger.Errorf("Error inserting answer: %v", err)
			return err
		}
	}

	// Commit transaction
	err = tx.Commit()
	if err != nil {
		repo.Logger.Errorf("Error committing transaction: %v", err)
		return err
	}

	return nil
}

// GetMaxQuizAttempt gets the maximum attempt number for a specific user and quiz
func (repo *PgQuizRepository) GetMaxQuizAttempt(userID, quizID int) (int, error) {
	var maxAttempt int

	_, err := repo.DB.Query(pg.Scan(&maxAttempt),
		"SELECT COALESCE(MAX(attempt), 0) FROM quiz_submissions WHERE user_id = ? AND quiz_id = ? AND deleted_at IS NULL",
		userID, quizID)

	if err != nil {
		repo.Logger.Errorf("Error getting max attempt for user %d and quiz %d: %v", userID, quizID, err)
		return 0, err
	}

	return maxAttempt, nil
}

// SaveQuizSubmission records a user's submission for a quiz question
func (repo *PgQuizRepository) SaveQuizSubmission(submission *m.QuizSubmission) error {
	temp := m.QuizSubmission{
		UserID:            submission.UserID,
		QuizID:            submission.QuizID,
		QuizQuestionID:    submission.QuizQuestionID,
		AnswerText:        submission.AnswerText,
		SelectedAnswerIds: submission.SelectedAnswerIds,
		Score:             submission.Score,
		Attempt:           submission.Attempt,
		Reviewed:          submission.Reviewed,
		SubmittedAt:       submission.SubmittedAt,
	}

	_, err := repo.DB.Model(&temp).Insert()
	if err != nil {
		repo.Logger.Errorf("Error saving quiz submission: %v", err)
		return err
	}

	submission.ID = temp.ID
	return nil
}

// GetQuizSubmissionsByUser gets all submissions for a quiz by a specific user
func (repo *PgQuizRepository) GetQuizSubmissionsByUser(userID, quizID int) ([]m.QuizSubmission, error) {
	var submissions []m.QuizSubmission

	err := repo.DB.Model(&submissions).
		Where("user_id = ?", userID).
		Where("quiz_id = ?", quizID).
		Where("deleted_at IS NULL").
		Order("created_at ASC").
		Select()

	if err != nil {
		repo.Logger.Errorf("Error fetching quiz submissions for user %d and quiz %d: %v", userID, quizID, err)
		return nil, err
	}

	return submissions, nil
}

// CreateQuizWithQuestionsAndAnswers handles the multi-step creation of a quiz
// with questions and answers in a single transaction
func (repo *PgQuizRepository) CreateQuizWithQuestionsAndAnswers(quizData *param.QuizData, title string) (int, error) {
	var quizID int = 0

	txErr := repo.DB.RunInTransaction(func(tx *pg.Tx) error {
		quiz := &m.Quiz{
			Title:      title,
			Difficulty: quizData.Difficulty,
			TotalScore: quizData.TotalScore,
			TimeLimit:  quizData.TimeLimit,
		}

		if _, err := tx.Model(quiz).Insert(); err != nil {
			repo.Logger.Errorf("Error creating quiz: %v", err)
			return err
		}

		quizID = quiz.ID

		for _, questionData := range quizData.Questions {
			question := &m.QuizQuestion{
				QuizID:            quiz.ID,
				QuestionType:      quizData.QuestionType,
				QuestionText:      questionData.QuestionText,
				Weight:            questionData.Weight,
				IsMultipleCorrect: questionData.AllowMultiple,
			}

			if _, err := tx.Model(question).Insert(); err != nil {
				repo.Logger.Errorf("Error creating question: %v", err)
				return err
			}

			if quizData.QuestionType == cf.QuesMultipleChoice {
				for _, optionData := range questionData.Options {
					answer := &m.QuizAnswer{
						QuizQuestionID: question.ID,
						AnswerText:     optionData.AnswersText,
						IsCorrect:      optionData.IsCorrect,
					}

					if _, err := tx.Model(answer).Insert(); err != nil {
						repo.Logger.Errorf("Error creating answer: %v", err)
						return err
					}
				}
			}
		}

		return nil
	})

	if txErr != nil {
		repo.Logger.Errorf("Transaction failed during quiz creation: %v", txErr)
		return 0, txErr
	}
	return quizID, nil
}

// GetEssaySubmissionsPendingReview retrieves essay submissions that need review
func (repo *PgQuizRepository) GetEssaySubmissionsPendingReview() ([]m.QuizSubmission, error) {
	var submissions []m.QuizSubmission

	err := repo.DB.Model(&submissions).
		Relation("User").
		Relation("User.UserProfile").
		Relation("QuizQuestion").
		Relation("Quiz").
		Where("\"user\".role_id = ?", cf.EmployeeRoleID).
		Where("quiz_question.question_type = ?", cf.QuesEssay).
		Where("quiz_submission.reviewed = false").
		Where("(quiz_submission.feedback IS NULL OR quiz_submission.feedback = '')").
		Where("quiz_submission.deleted_at IS NULL").
		Order("quiz_submission.user_id ASC").
		Order("quiz_submission.submitted_at DESC").
		Select()

	if err != nil {
		repo.Logger.Errorf("Error fetching essay submissions pending review: %v", err)
		return nil, err
	}

	return submissions, nil
}

// ReviewEssaySubmission updates an essay submission with review information
func (repo *PgQuizRepository) ReviewEssaySubmission(submissionID int, score float64, feedback string, reviewerID int) error {
	_, err := repo.DB.Model(&m.QuizSubmission{}).
		Set("score = ?", score).
		Set("feedback = ?", feedback).
		Set("reviewed = true").
		Set("reviewed_by = ?", reviewerID).
		Where("id = ?", submissionID).
		Where("deleted_at IS NULL").
		Update()

	if err != nil {
		repo.Logger.Errorf("Error updating essay submission review: %v", err)
		return err
	}

	return nil
}

// GetPendingEssayReviewsCountForCourse returns the count of unreviewed essay quizzes for a specific course and user
func (repo *PgQuizRepository) GetPendingEssayReviewsCountForCourse(userID int, courseID int) (int, error) {
	count, err := repo.DB.Model(&m.QuizSubmission{}).
		Join("JOIN quiz_questions qc ON qc.id = quiz_submission.quiz_question_id").
		Join("JOIN quizzes q ON q.id = quiz_submission.quiz_id").
		Join("JOIN module_items mi ON mi.quiz_id = q.id").
		Join("JOIN modules m ON m.id = mi.module_id").
		Where("quiz_submission.user_id = ?", userID).
		Where("m.course_id = ?", courseID).
		Where("qc.question_type = ?", cf.QuesEssay).
		Where("quiz_submission.reviewed = false").
		Where("quiz_submission.deleted_at IS NULL").
		Where("qc.deleted_at IS NULL").
		Where("q.deleted_at IS NULL").
		Where("mi.deleted_at IS NULL").
		Where("m.deleted_at IS NULL").
		Count()

	if err != nil {
		repo.Logger.Errorf("Error counting pending essay reviews for user %d, course %d: %v", userID, courseID, err)
		return 0, err
	}

	return count, nil
}
