package requestparams

// QuizOption represents an answer option for a quiz question
type QuizOption struct {
	AnswersText string `json:"answers_text"`
	IsCorrect   bool   `json:"is_correct"`
}

// QuizQuestion represents a question in a quiz
type QuizQuestion struct {
	QuestionText  string       `json:"question_text"`
	Weight        float64      `json:"weight"`
	AllowMultiple bool         `json:"allow_multiple"`
	Options       []QuizOption `json:"options"`
}

// QuizData represents the data for a quiz
type QuizData struct {
	QuestionType int            `json:"question_type"`
	Difficulty   int            `json:"difficulty"`
	TotalScore   float64        `json:"total_score"`
	TimeLimit    int            `json:"time_limit"`
	Questions    []QuizQuestion `json:"questions"`
}

type CreateModuleItemParams struct {
	Title        string    `json:"title" valid:"required"`
	ItemType     string    `json:"item_type" valid:"required"`
	Resource     string    `json:"resource"`
	Position     int       `json:"position"`
	RequiredTime int       `json:"required_time"`
	ModuleID     int       `json:"module_id" valid:"required"`
	QuizData     *QuizData `json:"quiz_data"`
	QuizID       int       `json:"-"`
}

type ModuleItemListParams struct {
	ModuleID int `json:"module_id" valid:"required"`
}

type ModuleItemIDParam struct {
	ModuleItemID int `json:"moduleItem_id" valid:"required"`
}
