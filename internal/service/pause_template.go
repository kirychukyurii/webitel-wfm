package service

import (
	"context"

	"github.com/webitel/webitel-wfm/internal/model"
	"github.com/webitel/webitel-wfm/internal/model/options"
	"github.com/webitel/webitel-wfm/internal/storage"
)

// TODO: add validation for cause.id
// for i, v := range in {
//		if err := p.store.ReadPauseCause(ctx, v.DomainId, v.Cause.Token); err != nil {
//			return err.SetDetailedError(fmt.Sprintf("items[%d].cause.id: not found: %s", i, err.appError()))
//		}
// }

type PauseTemplateManager interface {
	CreatePauseTemplate(ctx context.Context, read *options.Read, in *model.PauseTemplate) (*model.PauseTemplate, error)
	ReadPauseTemplate(ctx context.Context, read *options.Read) (*model.PauseTemplate, error)
	SearchPauseTemplate(ctx context.Context, search *options.Search) ([]*model.PauseTemplate, bool, error)
	UpdatePauseTemplate(ctx context.Context, read *options.Read, in *model.PauseTemplate) (*model.PauseTemplate, error)
	DeletePauseTemplate(ctx context.Context, read *options.Read) (int64, error)
}
type PauseTemplate struct {
	storage storage.PauseTemplateManager
}

func NewPauseTemplate(storage storage.PauseTemplateManager) *PauseTemplate {
	return &PauseTemplate{
		storage: storage,
	}
}

func (p *PauseTemplate) CreatePauseTemplate(ctx context.Context, read *options.Read, in *model.PauseTemplate) (*model.PauseTemplate, error) {
	id, err := p.storage.CreatePauseTemplate(ctx, read.User(), in)
	if err != nil {
		return nil, err
	}

	return p.ReadPauseTemplate(ctx, read.WithID(id))
}

func (p *PauseTemplate) ReadPauseTemplate(ctx context.Context, read *options.Read) (*model.PauseTemplate, error) {
	out, err := p.storage.ReadPauseTemplate(ctx, read)
	if err != nil {
		return nil, err
	}

	return out, nil
}

func (p *PauseTemplate) SearchPauseTemplate(ctx context.Context, search *options.Search) ([]*model.PauseTemplate, bool, error) {
	out, err := p.storage.SearchPauseTemplate(ctx, search)
	if err != nil {
		return nil, false, err
	}

	next, out := model.ListResult(int32(search.Size()), out)

	return out, next, nil
}

func (p *PauseTemplate) UpdatePauseTemplate(ctx context.Context, read *options.Read, in *model.PauseTemplate) (*model.PauseTemplate, error) {
	if err := p.storage.UpdatePauseTemplate(ctx, read.User(), in); err != nil {
		return nil, err
	}

	return p.ReadPauseTemplate(ctx, read)
}

func (p *PauseTemplate) DeletePauseTemplate(ctx context.Context, read *options.Read) (int64, error) {
	out, err := p.storage.DeletePauseTemplate(ctx, read)
	if err != nil {
		return 0, err
	}

	return out, nil
}
