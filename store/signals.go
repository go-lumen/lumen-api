package store

// BeforeCreator add a BeforeCreate callback
type BeforeCreator interface {
	BeforeCreate() error
}

// BeforeCreatorWithContext add a BeforeCreate callback with context
type BeforeCreatorWithContext interface {
	BeforeCreate(ctx *Context) error
}
