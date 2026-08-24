package memory

import (
	"context"
	"testing"

	"github.com/example/asset-maintenance-service/internal/domain"
)

func templateFixture() domain.InspectionTemplate {
	return domain.InspectionTemplate{ID: "tpl-1", Name: "pump", Version: 1, Active: true, Items: []domain.InspectionItem{{ID: "pressure", Name: "Pressure", Required: true}}}
}

func TestTemplateRepositoryGetIsolatesItems(t *testing.T) {
	r := NewTemplateRepository()
	if _, err := r.Create(context.Background(), templateFixture()); err != nil {
		t.Fatal(err)
	}
	got, err := r.Get(context.Background(), "tpl-1")
	if err != nil {
		t.Fatal(err)
	}
	got.Items[0].Name = "mutated"
	again, err := r.Get(context.Background(), "tpl-1")
	if err != nil {
		t.Fatal(err)
	}
	if again.Items[0].Name != "Pressure" {
		t.Fatalf("repository item changed through returned slice: %q", again.Items[0].Name)
	}
}

func TestTemplateRepositoryListIsolatesItems(t *testing.T) {
	r := NewTemplateRepository()
	if _, err := r.Create(context.Background(), templateFixture()); err != nil {
		t.Fatal(err)
	}
	list, err := r.List(context.Background(), true)
	if err != nil {
		t.Fatal(err)
	}
	list[0].Items[0].Name = "mutated"
	again, err := r.Get(context.Background(), "tpl-1")
	if err != nil {
		t.Fatal(err)
	}
	if again.Items[0].Name != "Pressure" {
		t.Fatalf("repository item changed through list slice: %q", again.Items[0].Name)
	}
}

func TestTemplateRepositoryCreateIsolatesInput(t *testing.T) {
	r := NewTemplateRepository()
	input := templateFixture()
	if _, err := r.Create(context.Background(), input); err != nil {
		t.Fatal(err)
	}
	input.Items[0].Name = "mutated"
	got, err := r.Get(context.Background(), input.ID)
	if err != nil {
		t.Fatal(err)
	}
	if got.Items[0].Name != "Pressure" {
		t.Fatalf("create retained caller slice: %q", got.Items[0].Name)
	}
}

func TestTemplateRepositoryUpdateIsolatesInput(t *testing.T) {
	r := NewTemplateRepository()
	input := templateFixture()
	if _, err := r.Create(context.Background(), input); err != nil {
		t.Fatal(err)
	}
	input.Items[0].Name = "updated"
	if _, err := r.Update(context.Background(), input); err != nil {
		t.Fatal(err)
	}
	input.Items[0].Name = "mutated after update"
	got, err := r.Get(context.Background(), input.ID)
	if err != nil {
		t.Fatal(err)
	}
	if got.Items[0].Name != "updated" {
		t.Fatalf("update retained caller slice: %q", got.Items[0].Name)
	}
}
