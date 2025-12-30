package templatestore

import (
	"context"

	"github.com/aritradevelops/billbharat/backend/notification/internal/persistence/dao"
	"github.com/aritradevelops/billbharat/backend/notification/internal/persistence/repository"
)

type DbTemplateStore struct {
	repo repository.Repository
}

func NewDbTemplateStore(repo repository.Repository) TemplateStorage {
	return &DbTemplateStore{
		repo: repo,
	}
}

func (d *DbTemplateStore) CreateTemplate(ctx context.Context, template dao.Template) (dao.Template, error) {
	return d.repo.CreateTemplate(ctx, template)
}

func (d *DbTemplateStore) FindTemplate(ctx context.Context, params FindTemplateParams) (dao.Template, error) {
	return d.repo.FindTemplate(ctx, repository.FindTemplateParams(params))
}
