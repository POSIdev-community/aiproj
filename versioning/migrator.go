package versioning

import (
	"fmt"
)

type FunctionMigrator struct {
	fromVersion string
	toVersion   string
	migrateFn   func(data []byte) ([]byte, error)
}

func NewFunctionMigrator(
	fromVersion, toVersion string,
	migrateFn func(data []byte) ([]byte, error),
) *FunctionMigrator {
	return &FunctionMigrator{
		fromVersion: fromVersion,
		toVersion:   toVersion,
		migrateFn:   migrateFn,
	}
}

func (m *FunctionMigrator) From() string {
	return m.fromVersion
}

func (m *FunctionMigrator) To() string {
	return m.toVersion
}

func (m *FunctionMigrator) Migrate(data []byte) ([]byte, error) {
	if m.migrateFn == nil {
		return nil, fmt.Errorf("migrate function is nil")
	}

	migrated, err := m.migrateFn(data)
	if err != nil {
		return nil, fmt.Errorf("run migrate function: %w", err)
	}

	return migrated, nil
}
