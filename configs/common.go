package configs

// Common const
const (
	FormatDate                  = "2006-01-02 15:04:05"
	FormatDateNoSec             = "2006-01-02 15:04"
	FormatDateDisplay           = "2006/01/02"
	FormatDateDatabase          = "2006-01-02"
	PostgreCharacterDisplayDate = "YYYY/MM/DD"
	AvatarFolderGCS             = "images/avatar/"
	ThumbnailFolderGCS          = "images/thumbnail/"
	FileFolderGCS               = "files/"
	VideoFolderGCS              = "videos/"

	// Course category
	Onboarding = 1
	Company    = 2
	Technical  = 3
	Soft       = 4
	Compliance = 5

	// Quiz difficulty levels
	DifficultyEasy   = 1
	DifficultyMedium = 2
	DifficultyHard   = 3

	// Quiz question types
	QuestionTypeMultipleChoice = 1
	QuestionTypeEssay          = 2
)

var CourseCategoryList = map[int]string{
	Onboarding: "Onboarding",
	Company:    "Company",
	Technical:  "Technical",
	Soft:       "Soft",
	Compliance: "Compliance",
}

var DifficultyLabels = map[int]string{
	DifficultyEasy:   "Easy",
	DifficultyMedium: "Medium",
	DifficultyHard:   "Hard",
}

var QuestionTypeLabels = map[int]string{
	QuestionTypeMultipleChoice: "Multiple Choice",
	QuestionTypeEssay:          "Essay",
}

// AllowFormatImageList format image allow
var AllowFormatImageList = []string{
	"png",
	"jpg",
	"jpeg",
	"gif",
}

// AllowFormatFileList format file allow
var AllowFormatFileList = []string{
	"pdf",
	"vnd.openxmlformats-officedocument.spreadsheetml.sheet",
	"vnd.openxmlformats-officedocument.wordprocessingml.document",
}
