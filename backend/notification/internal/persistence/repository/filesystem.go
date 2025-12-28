package repository

import (
	"context"
	"embed"
	"encoding/json"
	"fmt"
	"mime"
	"path"

	"github.com/aritradeveops/billbharat/backend/notification/internal/persistence/dao"
	"github.com/google/uuid"
)

//go:embed templates/**
var templateFs embed.FS
var preferred = map[string]string{
	"text/html":        ".html",
	"application/json": ".json",
	"text/plain":       ".txt",
}

func PreferredExt(mimeType string) string {
	if ext, ok := preferred[mimeType]; ok {
		return ext
	}
	exts, _ := mime.ExtensionsByType(mimeType)
	if len(exts) > 0 {
		return exts[0]
	}
	return ""
}

// this is only for running in local
// and storing the templates in filesystem
type FilesystemRepository struct {
	fs embed.FS
}

func NewFilesystemRepository() Repository {
	return &FilesystemRepository{
		fs: templateFs,
	}
}

// NOTE: scope is not considered in the path, for local filesystem
// all will be using default scope
func (t FindTemplateParams) directory() string {
	return path.Join("templates", string(t.Event), string(t.Channel))
}

func (t FindTemplateParams) contentPath() (string, error) {
	ext, ok := preferred[t.Mimetype]
	if !ok {
		return "", fmt.Errorf("unsupported mimetype: %s", t.Mimetype)
	}
	return path.Join(t.directory(), fmt.Sprintf("%s%s", t.Locale, ext)), nil
}

func (t FindTemplateParams) metadataPath() string {
	return path.Join(t.directory(), fmt.Sprintf("%s%s", t.Locale, ".json"))
}

// CreateTemplate implements Repository.
func (f *FilesystemRepository) CreateTemplate(ctx context.Context, template *dao.Template) (*dao.Template, error) {
	return nil, fmt.Errorf("not supported")
}

func (f *FilesystemRepository) FindTemplate(ctx context.Context, params FindTemplateParams) (*dao.Template, error) {
	contentPath, err := params.contentPath()
	if err != nil {
		return nil, err
	}
	content, err := f.fs.ReadFile(contentPath)
	if err != nil {
		return nil, err
	}

	metadataPath := params.metadataPath()
	metadataBytes, err := f.fs.ReadFile(metadataPath)
	if err != nil {
		return nil, err
	}

	metadata := struct {
		Subject string `json:"subject"`
	}{}

	err = json.Unmarshal(metadataBytes, &metadata)
	if err != nil {
		return nil, err
	}

	return &dao.Template{
		ID:       params.directory(),
		UID:      uuid.New(),
		Event:    params.Event,
		Channel:  params.Channel,
		Locale:   params.Locale,
		Subject:  metadata.Subject,
		Body:     string(content),
		Mimetype: params.Mimetype,
		Scope:    params.Scope,
	}, nil
}
