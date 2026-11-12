package versioning

import (
	"fmt"

	"github.com/POSIdev-community/aiproj/versions/v1_10"
	"github.com/POSIdev-community/aiproj/versions/v1_11"
	"github.com/POSIdev-community/aiproj/versions/v1_12"
	"github.com/POSIdev-community/aiproj/versions/v1_8"
	"github.com/POSIdev-community/aiproj/versions/v1_9"
)

func DefaultEngine() (*Engine, error) {
	engine, err := NewEngine(
		[]VersionSpec{
			{Version: "1.12", Validate: v1_12.Validate},
			{Version: "1.11", Validate: v1_11.Validate},
			{Version: "1.10", Validate: v1_10.Validate},
			{Version: "1.9", Validate: v1_9.Validate},
			{Version: "1.8", Validate: v1_8.Validate},
		},
		[]Migrator{
			NewFunctionMigrator("1.8", "1.9", v1_8.MigrateToV19),
			NewFunctionMigrator("1.9", "1.10", v1_9.MigrateToV110),
			NewFunctionMigrator("1.10", "1.11", v1_10.MigrateToV111),
			NewFunctionMigrator("1.11", "1.12", v1_11.MigrateToV112),
		},
	)
	if err != nil {
		return nil, fmt.Errorf("build aiproj version engine: %v", err)
	}

	return engine, nil
}
