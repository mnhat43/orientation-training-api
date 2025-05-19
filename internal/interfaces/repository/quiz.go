package repository

import (
	param "orientation-training-api/internal/interfaces/requestparams"
	m "orientation-training-api/internal/models"
)

// QuizRepository defines methods for accessing quiz data
type QuizRepository interface {
	GetQuizByID(quizID int) (m.Quiz, error)
	GetQuizList(params *param.QuizListParams) ([]m.Quiz, int, error)
	SaveQuiz(quiz *m.Quiz) error
	DeleteQuiz(quizID int) error
	GetQuizQuestionsWithAnswers(quizID int) ([]m.QuizQuestion, error)
	SaveQuizQuestion(question *m.QuizQuestion, answers []m.QuizAnswer) error
	SaveQuizSubmission(submission *m.QuizSubmission) error
	GetQuizSubmissionsByUser(userID, quizID int) ([]m.QuizSubmission, error)
}
