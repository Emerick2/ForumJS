package forumjs 

import (
	"time"
)

func Date(dateBrute time.Time) string {
	return dateBrute.Format("02 01 2006")
}