package service

import (
	"context"
	"errors"
	"testing"

	"github.com/4thena-io/abyss/internal/model"
	"github.com/4thena-io/abyss/internal/service/mocks"
)

func TestDeploymentService_GetDeploymentsByApp(t *testing.T) {
	repo := mocks.NewMockDeploymentRepository(t)
	want := []model.Deployment{{ID: 1, AppID: 5}}
	repo.EXPECT().GetByApp(context.Background(), uint(5)).Return(want, nil)

	svc := NewDeploymentService(repo)

	got, err := svc.GetDeploymentsByApp(context.Background(), 5)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(got) != 1 || got[0].ID != 1 {
		t.Fatalf("unexpected deployments: %+v", got)
	}
}

func TestDeploymentService_SaveDeployment(t *testing.T) {
	t.Run("propagates repository error", func(t *testing.T) {
		repo := mocks.NewMockDeploymentRepository(t)
		d := &model.Deployment{AppID: 1}
		repo.EXPECT().Save(context.Background(), d).Return(errBoom)

		svc := NewDeploymentService(repo)

		if err := svc.SaveDeployment(context.Background(), d); !errors.Is(err, errBoom) {
			t.Fatalf("expected errBoom, got %v", err)
		}
	})

	t.Run("saves successfully", func(t *testing.T) {
		repo := mocks.NewMockDeploymentRepository(t)
		d := &model.Deployment{AppID: 1}
		repo.EXPECT().Save(context.Background(), d).Return(nil)

		svc := NewDeploymentService(repo)

		if err := svc.SaveDeployment(context.Background(), d); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	})
}
