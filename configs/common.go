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
)

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
