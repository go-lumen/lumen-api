package server

import (
	"github.com/go-lumen/lumen-api/services"
	"github.com/go-lumen/lumen-api/store"
)

// TaskHandler define the shape of a cron task handler
type TaskHandler = func(ctx *store.Context) // , apiKeys models.APIKeys

// TaskSpec represents a cron task specification
type TaskSpec struct {
	CronSpec string
	Handler  TaskHandler
	Name     string
}

// tasks holds server cron task
var tasks = []TaskSpec{
	{"0 7 * * *", services.KPITask, "compute-kpis"},
}
