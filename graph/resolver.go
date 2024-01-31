package graph

//go:generate go run github.com/99designs/gqlgen generate

import (
	"github.com/drew-harris/lapis/standards"
	"github.com/posthog/posthog-go"
	"gorm.io/gorm"
)

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

type Resolver struct {
	db        *gorm.DB
	posthog   posthog.Client
	standards *standards.StandardsController
}

func NewResolver(db *gorm.DB, posthog posthog.Client,
	standards *standards.StandardsController,
) Resolver {
	return Resolver{
		db:        db,
		posthog:   posthog,
		standards: standards,
	}
}
