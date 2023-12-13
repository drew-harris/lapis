package graph

//go:generate go run github.com/99designs/gqlgen generate

import (
	"github.com/posthog/posthog-go"
	"gorm.io/gorm"
)

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

type Resolver struct {
	db      *gorm.DB
	posthog posthog.Client
}

func NewResolver(db *gorm.DB, posthog posthog.Client) Resolver {
	return Resolver{
		db:      db,
		posthog: posthog,
	}
}
