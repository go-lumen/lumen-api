package controllers

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/go-lumen/lumen-api/helpers"
	"github.com/go-lumen/lumen-api/models"
	"github.com/go-lumen/lumen-api/store"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"net/http"
	"reflect"
)

// BaseController hold common controller actions
type BaseController struct{}

// AbortWithError abort current request with a standardized output
func (BaseController) AbortWithError(c *gin.Context, apiError helpers.Error) {
	_ = c.AbortWithError(apiError.HTTPCode, apiError)
	logrus.WithError(apiError.Trace).WithField("path", c.FullPath()).Debugln("request aborted")
}

// ErrorInternal is a shortcut to test and abort if there if an internal error.
// Returns true if the caller should return.
func (bc BaseController) ErrorInternal(c *gin.Context, err error) bool {
	if err != nil {
		bc.AbortWithError(c, helpers.ErrorInternal(err))
		return true
	}
	return false
}

// Error is a shortcut to test and abort if there if an error.
// Param message is optional and override error message.
// Returns true if the caller should return.
func (bc BaseController) Error(c *gin.Context, err error, apiError helpers.NoCtxError, message ...string) bool {
	if err != nil {
		renderedError := apiError(err)
		if len(message) == 1 {
			renderedError.Message = message[0]
		}

		bc.AbortWithError(c, renderedError)
		return true
	}
	return false
}

// BindJSONError is a shortcut for BindJSON and AbortWithError in case of error.
// Returns true if the caller should return.
func (bc BaseController) BindJSONError(c *gin.Context, obj interface{}) bool {
	if err := c.BindJSON(obj); err != nil {
		bc.AbortWithError(c, helpers.ErrorInvalidInput(err))
		return true
	}
	return false
}

// LoggedUser returns logged user, group and a bool indicating is someone is logged in.
// If no used authenticated, current request is aborted with ErrorUserUnauthorized.
func (bc BaseController) LoggedUser(c *gin.Context) (store.User, store.Group, bool) {
	storedUser, userExists := c.Get(store.CurrentUserKey)
	if !userExists || storedUser == nil {
		bc.AbortWithError(c, helpers.ErrorUserUnauthorized)
		return nil, nil, false
	}
	user := storedUser.(*models.User)

	storedUserGroup, userGroupExists := c.Get(store.CurrentUserGroupKey)
	if !userGroupExists || storedUserGroup == nil {
		bc.AbortWithError(c, helpers.ErrorUserUnauthorized)
		return user, nil, false
	}
	return user, storedUserGroup.(*models.Group), true
}

// ParamID is a shortcut for `bson.M{"_id": c.Param("id")}`
func (bc BaseController) ParamID(c *gin.Context) bson.M {
	return bson.M{"_id": c.Param("id")}
}

// ShouldBeLogged is an helper method to check if the IsLogged flag and send an http error code accordingly
func (bc BaseController) ShouldBeLogged(ctx *store.Context) bool {
	if !ctx.IsLogged {
		bc.AbortWithError(ctx.C, helpers.ErrorUserUnauthorized)
		return false
	}
	return true
}

func sliceToBaseModelSlice(slice interface{}) []store.Model {
	s := reflect.ValueOf(slice)
	ret := make([]store.Model, s.Len())
	for i := 0; i < s.Len(); i++ {
		ret[i] = s.Index(i).Interface().(store.Model)
	}
	return ret
}

// CRUDController implements generic CRUD actions
type CRUDController struct {
	BaseController
	modelType reflect.Type
}

// GetModel implements generic find one by ID
func (cc CRUDController) GetModel(c *gin.Context) {
	ctx := store.AuthContext(c)
	if !cc.ShouldBeLogged(ctx) {
		return
	}

	// create an empty model
	model := reflect.New(cc.modelType).Interface().(store.Model)
	if cc.Error(c, ctx.Store.Find(ctx, cc.ParamID(c), model), helpers.ErrorResourceNotFound) {
		return
	}

	/*if user, group, ok := cc.LoggedUser(c); ok {
		if model.CanBeRead(user, group) {
			c.JSON(http.StatusOK, model)
			return
		}

		cc.AbortWithError(c, helpers.ErrorUserUnauthorized)
	}*/
}

func (cc CRUDController) makeSlice() reflect.Value {
	slice := reflect.MakeSlice(reflect.SliceOf(reflect.PtrTo(cc.modelType)), 0, 0)
	results := reflect.New(slice.Type())
	results.Elem().Set(slice)
	return results
}

// fetchModels returns a list of models with query filters
// Automatically set errors in gin context.
func (cc CRUDController) fetchModels(c *gin.Context) ([]store.Model, error) {
	ctx := store.AuthContext(c)
	if !cc.ShouldBeLogged(ctx) {
		return nil, errors.New("should be logged")
	}

	results := cc.makeSlice().Elem()
	filters := bson.M{} // TODO: generic way to get filter from query args ?
	if cc.Error(c, ctx.Store.FindAll(ctx, filters, results.Addr().Interface()), helpers.ErrorResourceNotFound) {
		return nil, errors.New("not found")
	}
	//resTyped := sliceToBaseModelSlice(results.Interface())
	readableModels := make([]store.Model, 0)
	/*if user, group, ok := cc.LoggedUser(c); ok {
		for _, m := range resTyped {
			if m.CanBeRead(user, group) {
				readableModels = append(readableModels, m)
			}
		}
	}*/
	return readableModels, nil
}

// GetModels implements generic find several with filters
func (cc CRUDController) GetModels(c *gin.Context) {
	results, err := cc.fetchModels(c)
	if err != nil {
		return
	}
	c.JSON(http.StatusOK, results)
}

// CreateModel implements generic model creation
func (cc CRUDController) CreateModel(c *gin.Context) {
	ctx := store.AuthContext(c)
	if !cc.ShouldBeLogged(ctx) {
		return
	}

	model := reflect.New(cc.modelType).Interface().(store.Model)
	if cc.BindJSONError(c, model) {
		return
	}

	/*if user, group, ok := cc.LoggedUser(c); ok && model.CanBeCreated(user, group) {
		if cc.ErrorInternal(c, ctx.Store.Create(ctx, model)) {
			return
		}
		c.JSON(http.StatusCreated, model)
	} else {
		cc.AbortWithError(c, helpers.ErrorUserUnauthorized)
	}*/
}

// DeleteModel implements generic delete one by ID
func (cc CRUDController) DeleteModel(c *gin.Context) {
	ctx := store.AuthContext(c)
	if !cc.ShouldBeLogged(ctx) {
		return
	}

	model := reflect.New(cc.modelType).Interface().(store.Model)
	if cc.Error(c, ctx.Store.Find(ctx, cc.ParamID(c), model), helpers.ErrorResourceNotFound) {
		return
	}

	/*if user, group, ok := cc.LoggedUser(c); ok {
		if model.CanBeDeleted(user, group) {
			if cc.ErrorInternal(c, ctx.Store.Delete(ctx, c.Param("id"), model)) {
				return
			}
			c.JSON(http.StatusOK, model)
			return
		}

		cc.AbortWithError(c, helpers.ErrorUserUnauthorized)
	}*/
}

// UpdateModel implements generic update one
func (cc CRUDController) UpdateModel(c *gin.Context) {
	ctx := store.AuthContext(c)
	if !cc.ShouldBeLogged(ctx) {
		return
	}

	model := reflect.New(cc.modelType).Interface().(store.Model)
	if cc.Error(c, ctx.Store.Find(ctx, cc.ParamID(c), model), helpers.ErrorResourceNotFound) {
		return
	}

	newModel := reflect.New(cc.modelType).Interface().(store.Model)
	if cc.BindJSONError(c, newModel) {
		return
	}

	/*if user, group, ok := cc.LoggedUser(c); ok {
		if model.CanBeUpdated(user, group) {
			if cc.ErrorInternal(c, ctx.Store.Update(ctx, store.ID(c.Param("id")), newModel)) {
				return
			}
			c.JSON(http.StatusOK, newModel)
			return
		}

		cc.AbortWithError(c, helpers.ErrorUserUnauthorized)
	}*/
}
