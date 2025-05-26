package respond

// LectureModuleResponse represents a module with its lectures
type LectureModuleResponse struct {
	ModuleID       int                   `json:"module_id"`
	ModuleTitle    string                `json:"module_title"`
	ModulePosition int                   `json:"module_position"`
	Duration       int                   `json:"duration"`
	Lectures       []LectureItemResponse `json:"lecture"`
}

// LectureItemResponse represents a lecture item
type LectureItemResponse struct {
	ModuleItemID       int         `json:"module_item_id"`
	ModuleItemTitle    string      `json:"module_item_title"`
	ModuleItemPosition int         `json:"module_item_position"`
	ItemType           string      `json:"item_type"`
	Unlocked           bool        `json:"unlocked"`
	Content            interface{} `json:"content"`
}

// VideoContentResponse represents video content
type VideoContentResponse struct {
	VideoID      string `json:"videoId"`
	Duration     string `json:"duration"`
	RequiredTime int    `json:"required_time"`
	Thumbnail    string `json:"thumbnail"`
	PublishedAt  string `json:"publishedAt,omitempty"`
}

// FileContentResponse represents file content
type FileContentResponse struct {
	FilePath     string `json:"file_path"`
	Duration     int    `json:"duration"`
	RequiredTime int    `json:"required_time"`
}

// QuizContentResponse represents quiz content
type QuizContentResponse struct {
	QuizID     int                    `json:"quiz_id"`
	QuizType   string                 `json:"quiz_type"`
	QuizTitle  string                 `json:"quiz_title"`
	Difficulty string                 `json:"difficulty"`
	TotalScore float64                `json:"total_score"`
	TimeLimit  int                    `json:"time_limit"`
	Questions  []QuizQuestionResponse `json:"questions"`
}

// QuizQuestionResponse represents a quiz question
type QuizQuestionResponse struct {
	QuestionID    int                  `json:"question_id"`
	QuestionText  string               `json:"question_text"`
	AllowMultiple bool                 `json:"allow_multiple"`
	Points        float64              `json:"points,omitempty"`
	Options       []QuizOptionResponse `json:"options"`
}

// QuizOptionResponse represents a quiz option
type QuizOptionResponse struct {
	ID   int    `json:"id"`
	Text string `json:"text"`
}
