package mode

import "github.com/rqpt/picker"

var modes = []string{
	"New",
	"Existing",
}

func Select() (string, error) {
	return picker.Run(modes)
}
