package helpers

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/gin-gonic/gin"
)

type ParamsToCheck struct {
	Key   string
	Value any
}

func ItemExists(collectionName *mongo.Collection, paramToCheck ParamsToCheck, ginCtx *gin.Context) bson.M {
	var data bson.M
	var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
	err := collectionName.FindOne(ctx, bson.M{paramToCheck.Key: paramToCheck.Value}).Decode(&data)
	if err != nil {
		defer cancel()
		return data
	}
	defer cancel()
	return data
}
