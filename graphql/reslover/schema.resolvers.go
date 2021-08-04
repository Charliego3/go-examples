package reslover

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"

	"github.com/whimthen/temp/graphql/generated"
	"github.com/whimthen/temp/graphql/model"
)

func (r *queryResolver) MeetUps(ctx context.Context) ([]*model.Meetup, error) {
	panic(fmt.Errorf("not implemented"))
}

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

type queryResolver struct{ *Resolver }
