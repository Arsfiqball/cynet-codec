package resource

import (
	"context"
	"errors"
	"feature/internal/value/domain"
	"feature/internal/value/user"
)

func (s *Service) Create(ctx context.Context, patch domain.Patch, user user.Identity) (domain.Entity, error) {
	ctx, span := s.tracer.Start(ctx, "feature/internal/application/resource.Service.Create")
	defer span.End()

	ent := domain.NewEntity()

	if err := ent.Patch(patch); err != nil {
		return nil, NewError(err, err.Error(), ErrCodeInvalidEntity)
	}

	if err := ent.Validate(); err != nil {
		return nil, NewError(err, err.Error(), ErrCodeInvalidEntity)
	}

	if err := s.authorizeCreate(user, ent); err != nil {
		return nil, NewError(err, err.Error(), ErrCodeUnauthorized)
	}

	ent, err := s.repo.Create(ctx, ent)
	if err != nil {
		return nil, NewError(err, err.Error(), ErrCodeUnknown)
	}

	if err := s.event.Created(ctx, ent); err != nil {
		return nil, NewError(err, err.Error(), ErrCodeUnknown)
	}

	return ent, nil
}

func (s *Service) authorizeCreate(u user.Identity, ent domain.Entity) error {
	// TODO: implement authorization logic

	if u.Provider() == user.ProviderGuest {
		return errors.New("guest cannot create resource")
	}

	return nil
}
