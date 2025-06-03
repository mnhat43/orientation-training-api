package repository

import (
	param "orientation-training-api/internal/interfaces/requestparams"
	m "orientation-training-api/internal/models"
)

type QuizRepository interface {
	GetQuizByID(quizID int) (m.Quiz, error)
	GetQuizList(params *param.QuizListParams) ([]m.Quiz, int, error)
	SaveQuiz(quiz *m.Quiz) error
	DeleteQuiz(quizID int) error
	GetQuizQuestionsWithAnswers(quizID int) ([]m.QuizQuestion, error)
	SaveQuizQuestion(question *m.QuizQuestion, answers []m.QuizAnswer) error
	SaveQuizSubmission(submission *m.QuizSubmission) error
	GetQuizSubmissionsByUser(userID int, quizID int) ([]m.QuizSubmission, error)
	CreateQuizWithQuestionsAndAnswers(quizData *param.QuizData, title string) (int, error)
	GetMaxQuizAttempt(userID int, quizID int) (int, error)
	GetEssaySubmissionsPendingReview() ([]m.QuizSubmission, error)
	ReviewEssaySubmission(submissionID int, score float64, feedback string, reviewerID int) error
	GetPendingEssayReviewsCountForCourse(userID int, courseID int) (int, error)
}
