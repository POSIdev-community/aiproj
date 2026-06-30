package v1_9

import (
	"github.com/POSIdev-community/aiproj/model"
)

func ReadFromBytes(data []byte) (AIProj, error) {
	return model.ReadFromBytes[AIProj](data)
}
