package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.
// Code generated by github.com/99designs/gqlgen version v0.17.36

import (
	"context"
	"errors"
	"fmt"
	"slices"
	"strings"

	"github.com/99designs/gqlgen/graphql"
	"github.com/drew-harris/lapis/graph/model"
	"github.com/drew-harris/lapis/maps"
	"github.com/google/uuid"
	"gorm.io/datatypes"
)

// Player is the resolver for the player field.
func (r *logResolver) Player(ctx context.Context, obj *model.Log) (*model.Player, error) {
	if obj.Player != nil {
		return obj.Player, nil
	}
	fmt.Println("Using long query")
	player := model.Player{}
	r.db.Where("id = ?", obj.PlayerID).First(&player)
	if r.db.Error != nil {
		return nil, r.db.Error
	}
	return &player, nil
}

// Attributes is the resolver for the attributes field.
func (r *logResolver) Attributes(ctx context.Context, obj *model.Log) (map[string]interface{}, error) {
	attributes, err := maps.ToMap(obj.Attributes)
	if err != nil {
		return nil, err
	}
	return attributes, nil
}

// Log is the resolver for the log field.
func (r *mutationResolver) Log(ctx context.Context, input model.LogInput) (*model.Log, error) {
	// Check if playerid is valid
	player := model.Player{}
	r.db.Where("name = ?", input.PlayerName).First(&player)
	if r.db.Error != nil {
		fmt.Println(r.db.Error)
		return nil, r.db.Error
	}
	if player.ID == "" {
		return nil, errors.New("Player id is not valid")
	}
	attributes, err := maps.FromMap(input.Attributes)
	if err != nil {
		return nil, err
	}
	log := model.Log{
		ID:         uuid.New().String(),
		Message:    input.Message,
		PlayerID:   player.ID,
		Attributes: attributes,
		Type:       input.Type,
	}

	r.db.Create(&log)
	if r.db.Error != nil {
		return nil, r.db.Error
	}
	return &log, nil
}

// Logs is the resolver for the logs field.
func (r *queryResolver) Logs(ctx context.Context, filter *model.LogQueryFilter, limit *model.LimitFilter) ([]model.Log, error) {
	logs := []model.Log{}
	db := r.db
	if filter != nil {
		if filter.PlayerID != nil {
			db = db.Where("player_id = ?", filter.PlayerID)
		}
		if filter.Type != nil {
			db = db.Where("type in ?", filter.Type)
		}
		if filter.HasAttribute != nil {
			db = db.Where(datatypes.JSONQuery("attributes").HasKey(*filter.HasAttribute))
		}
	}

	if limit != nil {
		if limit.Limit != nil {
			db = db.Limit(*limit.Limit)
		}
		if limit.Page != nil {
			db = db.Offset((*limit.Page - 1) * *limit.Limit)
		}
	}

	fmt.Println(strings.Join(graphql.CollectAllFields(ctx), " "))
	if slices.Contains(graphql.CollectAllFields(ctx), "player") {
		fmt.Println("Preloading player")
		db = db.Joins("Player")
	}

	result := db.Order("created_at desc").Find(&logs)
	if result.Error != nil {
		return nil, result.Error
	}
	return logs, nil
}

// Log returns LogResolver implementation.
func (r *Resolver) Log() LogResolver { return &logResolver{r} }

// Mutation returns MutationResolver implementation.
func (r *Resolver) Mutation() MutationResolver { return &mutationResolver{r} }

// Query returns QueryResolver implementation.
func (r *Resolver) Query() QueryResolver { return &queryResolver{r} }

type logResolver struct{ *Resolver }
type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
