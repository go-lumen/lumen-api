package models

import "go.mongodb.org/mongo-driver/bson"

type QueryParams struct {
	Limit     uint32
	Order     int8
	StartTime uint32
	EndTime   uint32
}

// All returns an empty filter (for semantic)
func All() bson.M { return bson.M{} }

// ByID returns a by ID filter
func ByID(key string) bson.M { return bson.M{"_id": key} }

// ByGroupID returns a group_id filter
func ByGroupID(id string) bson.M { return bson.M{"group_id": id} }
