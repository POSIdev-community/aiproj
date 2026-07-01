package versioning

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
)

var (
	ErrNoSupportedVersions = errors.New("no supported versions")
	ErrNoValidVersion      = errors.New("aiproj version was not detected")
	ErrMigrationPath       = errors.New("migration path was not found")
	ErrUnsupportedVersion  = errors.New("unsupported aiproj version")
)

type Validator func(data []byte) error

type VersionSpec struct {
	Version  string
	Validate Validator
}

type Migrator interface {
	From() string
	To() string
	Migrate(data []byte) ([]byte, error)
}

type Engine struct {
	orderedVersions []VersionSpec
	latestVersion   string
	migratorsByFrom map[string]Migrator
	indexByVersion  map[string]int
}

type baseAIProj struct {
	Version string `json:"Version"`
}

func NewEngine(orderedVersions []VersionSpec, migrators []Migrator) (*Engine, error) {
	if len(orderedVersions) == 0 {
		return nil, ErrNoSupportedVersions
	}

	indexByVersion := make(map[string]int, len(orderedVersions))

	for idx, spec := range orderedVersions {
		if spec.Version == "" {
			return nil, fmt.Errorf("version spec[%d]: empty version", idx)
		}

		if spec.Validate == nil {
			return nil, fmt.Errorf("version spec[%d] %s: validator is nil", idx, spec.Version)
		}

		if _, exists := indexByVersion[spec.Version]; exists {
			return nil, fmt.Errorf("duplicate version spec: %s", spec.Version)
		}

		indexByVersion[spec.Version] = idx
	}

	migratorsByFrom := make(map[string]Migrator, len(migrators))

	for _, migrator := range migrators {
		if migrator == nil {
			return nil, errors.New("nil migrator")
		}

		if _, exists := indexByVersion[migrator.From()]; !exists {
			return nil, fmt.Errorf("migrator source version is unsupported: %s", migrator.From())
		}

		if _, exists := indexByVersion[migrator.To()]; !exists {
			return nil, fmt.Errorf("migrator target version is unsupported: %s", migrator.To())
		}

		if _, exists := migratorsByFrom[migrator.From()]; exists {
			return nil, fmt.Errorf("duplicate migrator for version: %s", migrator.From())
		}

		migratorsByFrom[migrator.From()] = migrator
	}

	return &Engine{
		orderedVersions: orderedVersions,
		latestVersion:   orderedVersions[0].Version,
		migratorsByFrom: migratorsByFrom,
		indexByVersion:  indexByVersion,
	}, nil
}

func (e *Engine) LatestVersion() string {
	return e.latestVersion
}

func (e *Engine) Detect(data []byte) (string, error) {
	baseVersion := extractVersion(data)
	if baseVersion != "" {
		if idx, ok := e.indexByVersion[baseVersion]; ok {
			spec := e.orderedVersions[idx]
			if err := spec.Validate(data); err != nil {
				return "", fmt.Errorf("validate explicit version %s: %w", baseVersion, err)
			}

			return baseVersion, nil
		}

		return "", fmt.Errorf("%w: %s", ErrUnsupportedVersion, baseVersion)
	}

	var validationErrors []string

	for _, spec := range e.orderedVersions {
		if err := spec.Validate(data); err == nil {
			return spec.Version, nil
		} else {
			validationErrors = append(validationErrors, fmt.Sprintf("%s: %v", spec.Version, err))
		}
	}

	return "", fmt.Errorf("%w: %s", ErrNoValidVersion, strings.Join(validationErrors, " | "))
}

func extractVersion(data []byte) string {
	var base baseAIProj

	decoder := json.NewDecoder(bytes.NewReader(data))
	decoder.UseNumber()
	if err := decoder.Decode(&base); err != nil {
		return ""
	}

	return base.Version
}

func (e *Engine) MigrateToLatest(data []byte, fromVersion string) ([]byte, error) {
	return e.MigrateTo(data, fromVersion, e.latestVersion)
}

func (e *Engine) MigrateTo(data []byte, fromVersion, toVersion string) ([]byte, error) {
	startIdx, ok := e.indexByVersion[fromVersion]
	if !ok {
		return nil, fmt.Errorf("%w: %s", ErrMigrationPath, fromVersion)
	}

	targetIdx, ok := e.indexByVersion[toVersion]
	if !ok {
		return nil, fmt.Errorf("%w: %s", ErrMigrationPath, toVersion)
	}

	if startIdx < targetIdx {
		return nil, fmt.Errorf("%w: cannot downgrade %s -> %s", ErrMigrationPath, fromVersion, toVersion)
	}

	if startIdx == targetIdx {
		if err := e.orderedVersions[startIdx].Validate(data); err != nil {
			return nil, fmt.Errorf("validate data as %s: %w", fromVersion, err)
		}

		return data, nil
	}

	if err := e.orderedVersions[startIdx].Validate(data); err != nil {
		return nil, fmt.Errorf("validate data as %s: %w", fromVersion, err)
	}

	migrated := data
	currentVersion := fromVersion

	for idx := startIdx; idx > targetIdx; idx-- {
		nextVersion := e.orderedVersions[idx-1].Version
		migrator, exists := e.migratorsByFrom[currentVersion]
		if !exists || migrator.To() != nextVersion {
			return nil, fmt.Errorf("%w: %s -> %s", ErrMigrationPath, currentVersion, nextVersion)
		}

		var err error
		migrated, err = migrator.Migrate(migrated)
		if err != nil {
			return nil, fmt.Errorf("migrate %s -> %s: %w", currentVersion, nextVersion, err)
		}

		if err := e.orderedVersions[idx-1].Validate(migrated); err != nil {
			return nil, fmt.Errorf("validate migrated data as %s: %w", nextVersion, err)
		}

		currentVersion = nextVersion
	}

	return migrated, nil
}
