//go:generate go run github.com/99designs/gqlgen

package graph

import (
	"context"
	"fmt"
)

var (
	TodoA = &Todo{
		ID:   "0be25fcf-20e6-4a6d-b0f9-7804224ef20e",
		Text: "Todo A",
		Done: false,
		User: &User{
			ID:   "123",
			Name: "Name A",
		},
	}

	TodoB = &Todo{
		ID:   "53314dfd-35a8-4665-b528-20d06ba549be",
		Text: "Todo B",
		Done: false,
		User: &User{
			ID:   "456",
			Name: "Name B",
		},
	}

	TodoC = &Todo{
		ID:   "bb5893c7-eb8d-496b-ae36-6002bbe50b7e",
		Text: "Todo C",
		Done: false,
		User: &User{
			ID:   "789",
			Name: "Name C",
		},
	}
)

type Resolver struct{}

func (r *Resolver) Mutation() MutationResolver {
	return &mutationResolver{r}
}
func (r *Resolver) Query() QueryResolver {
	return &queryResolver{r}
}

type mutationResolver struct{ *Resolver }

func (r *mutationResolver) CreateTodo(ctx context.Context, input NewTodo) (*Todo, error) {
	return &Todo{
		ID:   "bf7f59b7-a147-43a5-9f80-00a2d4ed6880",
		Text: input.Text,
		Done: false,
		User: &User{
			ID:   input.UserID,
			Name: fmt.Sprintf("Test%v", input.UserID),
		},
	}, nil
}

type queryResolver struct{ *Resolver }

func (r *queryResolver) Todos(ctx context.Context) ([]*Todo, error) {
	return []*Todo{TodoA, TodoB, TodoC}, nil
}
