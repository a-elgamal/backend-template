// Package test provides mocking capability for stored package
package test

import (
	"context"

	"github.com/stretchr/testify/mock"
	"alielgamal.com/myservice/internal/stored"
)

// Store provides mocking capability for store.
type Store[T any] struct {
	mock.Mock
}

// Add a new Stored item with a specific id and content
func (m *Store[T]) Add(ctx context.Context, creator string, id string, content T) (*stored.Stored[T], error) {
	args := m.Called(ctx, creator, id, content)
	return args.Get(0).(*stored.Stored[T]), args.Error(1)
}

// Patch Updates a single attribute in the content. This method doesn't check the attribute existence but guarantees
// that content stored is still a valid. If it is not, the patch operation will fail without impacting storage.
func (m *Store[T]) Patch(ctx context.Context, updater string, id string, attributes map[string]any) (*stored.Stored[T], error) {
	args := m.Called(ctx, updater, id, attributes)
	return args.Get(0).(*stored.Stored[T]), args.Error(1)
}

// Get finds a storable by its id
func (m *Store[T]) Get(ctx context.Context, id string) (*stored.Stored[T], error) {
	args := m.Called(ctx, id)
	return args.Get(0).(*stored.Stored[T]), args.Error(1)
}

// List returns all items that fill certain all conditions (AND operator between the conditions).
// if no conditions are passed, all stored items are returned.
func (m *Store[T]) List(ctx context.Context, conds ...stored.Condition) ([]stored.Stored[T], error) {
	allArgs := []any{ctx}
	for _, c := range conds {
		allArgs = append(allArgs, c)
	}
	args := m.Called(allArgs...)
	return args.Get(0).([]stored.Stored[T]), args.Error(1)
}
