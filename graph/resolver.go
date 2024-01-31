package graph

//go:generate go run github.com/99designs/gqlgen generate

import (
	"github.com/drew-harris/lapis/logging"
	"github.com/drew-harris/lapis/realtime"
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
	logging   *logging.LoggingService
	standards *standards.StandardsController
	realtime  *realtime.Service
}

func NewResolver(db *gorm.DB, posthog posthog.Client,
	standards *standards.StandardsController, realtime *realtime.Service, logging *logging.LoggingService,
) Resolver {
	return Resolver{
		db:        db,
		posthog:   posthog,
		standards: standards,
		realtime:  realtime,
		logging:   logging,
	}
}
