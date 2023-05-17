package tzinit

import (
	"os"
)

func init() {
	os.Setenv("TZ", "Europe/Copenhagen") // this should make gin's timestamps in Copenhagen time
}
