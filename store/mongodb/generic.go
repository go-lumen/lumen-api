package mongodb

import (
	"context"
	"encoding/json"
	"github.com/chidiwilliams/flatbson"
	"github.com/go-lumen/lumen-api/store"
	"github.com/go-lumen/lumen-api/utils"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"reflect"
)

const idName = "ID"

func setID(model store.Model, id string) {
	v := reflect.ValueOf(model).Elem().FieldByName(idName)
	if v.IsValid() && v.CanSet() {
		v.SetString(id)
	}
}

func orderToMongoOrder(order store.SortOrder) int {
	if order == store.SortAscending {
		return 1
	}
	return -1
}

// getCacheKey returns an unique cache key
func getCacheKey(filter bson.M, input interface{}) (string, error) {
	data, err := json.Marshal(filter)
	if err != nil {
		return "", errors.Wrap(err, "cannot marshall filters")
	}
	cacheKey := reflect.TypeOf(input).String() + ":" + string(data)
	return cacheKey, nil
}

// GetCollection returns mongo collection
func (db *Mngo) GetCollection(c *store.Context, model store.Model) *mongo.Collection {
	utils.EnsurePointer(model)
	mongoModel := store.EnsureGenericModel(model)
	return db.database.Collection(mongoModel.GetCollection())
}

// Create a generic model
func (db *Mngo) Create(c *store.Context, collectionName string, model store.Model) error {
	utils.EnsurePointer(model)
	mongoModel := store.EnsureGenericModel(model)
	collection := db.database.Collection(mongoModel.GetCollection())

	if creator, ok := model.(store.BeforeCreator); ok {
		if err := creator.BeforeCreate(); err != nil {
			return errors.Wrap(err, "error in BeforeCreate")
		}
	}
	if creator, ok := model.(store.BeforeCreatorWithContext); ok {
		if err := creator.BeforeCreate(c); err != nil {
			return errors.Wrap(err, "error in BeforeCreatorWithContext")
		}
	}

	res, err := collection.InsertOne(db.context, model)
	if err != nil {
		logrus.WithError(err).Errorln("cannot insert model")
		return errors.Wrap(err, "cannot insert model")
	}

	// update with inserted id
	if id, ok := res.InsertedID.(primitive.ObjectID); ok {
		setID(model, id.Hex())
	}
	if id, ok := res.InsertedID.(string); ok {
		setID(model, id)
	}

	return nil
}

// Find return a generic model
func (db *Mngo) Find(c *store.Context, filter bson.M, model store.Model, opts ...store.FindOption) error {
	optValues := store.GetFindOptions(opts...)
	utils.EnsurePointer(model)
	mongoModel := store.EnsureGenericModel(model)
	collection := db.database.Collection(mongoModel.GetCollection())

	cacheKey, err := getCacheKey(filter, model)
	if err != nil {
		return errors.Wrap(err, "cannot compute cache key")
	}

	if cached, found := c.GetCache(cacheKey); found && optValues.Cache {
		model = cached.(store.Model)
		return nil
	}

	var findOptions options.FindOneOptions
	// apply sort
	if len(optValues.SortedFields) > 0 {
		sortBson := bson.D{}
		for _, sortedField := range optValues.SortedFields {
			sortBson = append(sortBson, bson.E{Key: sortedField.Field, Value: orderToMongoOrder(sortedField.Order)})
		}
		findOptions.SetSort(sortBson)
	}

	err = collection.FindOne(db.context, filter, &findOptions).Decode(model)
	if err != nil {
		return errors.Wrap(err, "cannot find model")
	}

	c.SetCache(cacheKey, model)
	return err
}

// FindAll return several generic models
func (db *Mngo) FindAll(c *store.Context, filter bson.M, results interface{}, opts ...store.FindOption) error {
	utils.EnsurePointer(results)
	optValues := store.GetFindOptions(opts...)
	slice := reflect.ValueOf(results).Elem()
	modelType := reflect.TypeOf(results).Elem().Elem().Elem()
	newModel := reflect.New(modelType)
	sliceItem := newModel.Interface().(store.Model)
	mongoModel := store.EnsureGenericModel(sliceItem)
	collection := db.database.Collection(mongoModel.GetCollection())

	cacheKey, err := getCacheKey(filter, results)
	if err != nil {
		return errors.Wrap(err, "cannot compute cache key")
	}

	if cached, found := c.GetCache(cacheKey); found && optValues.Cache {
		results = cached
		return nil
	}

	// apply model default options
	var findOptions options.FindOptions
	if optioner, ok := mongoModel.(store.MongoFindAllOptioner); ok {
		optioner.ApplyOptions(&findOptions)
	}
	// apply limit
	if optValues.HasLimit {
		findOptions.SetLimit(optValues.Limit)
	}
	// apply sort
	if len(optValues.SortedFields) > 0 {
		sortBson := bson.D{}
		for _, sortedField := range optValues.SortedFields {
			sortBson = append(sortBson, bson.E{Key: sortedField.Field, Value: orderToMongoOrder(sortedField.Order)})
		}
		findOptions.SetSort(sortBson)
	}

	cur, err := collection.Find(context.TODO(), filter, &findOptions)
	if err != nil {
		logrus.WithError(err).Errorln("cannot find models")
		return errors.Wrap(err, "cannot find models")
	}

	for cur.Next(context.TODO()) {
		record := reflect.New(modelType)
		err = cur.Decode(record.Interface())
		if err != nil {
			logrus.Warnln("ErrorInternal on Decoding the document", err)
		}
		reflect.ValueOf(results).Elem().Set(reflect.Append(slice, record))
	}

	c.SetCache(cacheKey, results)
	return err
}

// Update a generic model
func (db *Mngo) Update(c *store.Context, filter bson.M, model store.Model, opts ...store.UpdateOption) error {
	optValues := store.GetUpdateOptions(opts...)
	utils.EnsurePointer(model)
	mongoModel := store.EnsureGenericModel(model)
	collection := db.database.Collection(mongoModel.GetCollection())

	// flatten field to update for mongodb driver
	fields, err := flatbson.Flatten(model)
	// filter out fields if OnlyFields is set
	if len(optValues.OnlyFields) > 0 {
		for key := range fields {
			if !utils.FindStringInSlice(key, optValues.OnlyFields) {
				delete(fields, key)
			}
		}
	}

	if err != nil {
		logrus.WithError(err).Errorln("cannot flatten model")
		return errors.Wrap(err, "cannot flatten model")
	}

	result, err := collection.UpdateOne(context.TODO(), filter,
		bson.M{"$set": fields}, options.Update().SetUpsert(optValues.CreateIfNotExists))

	if err != nil {
		logrus.WithError(err).Errorln("cannot update model")
		return errors.Wrap(err, "cannot update model")
	}

	if result.MatchedCount != 0 {
		logrus.Debugln("matched and replaced an existing document")
		return nil
	}
	if result.UpsertedCount != 0 {
		logrus.WithField("id", result.UpsertedID).Debugln("inserted a new document")
		return nil
	}

	return nil
}

// Delete a generic model
func (db *Mngo) Delete(c *store.Context, id string, model store.Model) error {
	utils.EnsurePointer(model)
	mongoModel := store.EnsureGenericModel(model)
	collection := db.database.Collection(mongoModel.GetCollection())

	_, err := collection.DeleteOne(db.context, bson.M{"_id": id})
	if err != nil {
		logrus.WithError(err).Errorln("cannot delete model")
		return errors.Wrap(err, "cannot delete model")
	}

	return nil
}

// DeleteAll a generic model
func (db *Mngo) DeleteAll(c *store.Context, filter bson.M, model store.Model) (int64, error) {
	utils.EnsurePointer(model)
	mongoModel := store.EnsureGenericModel(model)
	collection := db.database.Collection(mongoModel.GetCollection())

	res, err := collection.DeleteOne(db.context, filter)
	if err != nil {
		logrus.WithError(err).Errorln("cannot delete model")
		return 0, errors.Wrap(err, "cannot delete model")
	}

	return res.DeletedCount, nil
}
