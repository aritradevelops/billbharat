package templatestore

import (
	"context"

	"github.com/aritradevelops/billbharat/backend/notification/internal/persistence/dao"
	"github.com/aritradevelops/billbharat/backend/notification/internal/persistence/repository"
)

type TemplateStorage interface {
	CreateTemplate(ctx context.Context, template dao.Template) (dao.Template, error)
	FindTemplate(ctx context.Context, params FindTemplateParams) (dao.Template, error)
}

func NewTemplateStore(env string, repo repository.Repository) TemplateStorage {
	if env == "production" {
		return NewDbTemplateStore(repo)
	}
	return NewFSTemplateStore()
}
