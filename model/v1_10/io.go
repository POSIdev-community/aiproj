package v1_10

import (
	"github.com/POSIdev-community/aiproj/model"
)

func ReadFromBytes(data []byte) (AIProj, error) {
	return model.ReadFromBytes[AIProj](data)
}
