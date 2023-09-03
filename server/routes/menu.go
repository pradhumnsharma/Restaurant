package routes

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"server/helpers"
	"server/models"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var validateItem = validator.New()
var menuCollection *mongo.Collection = OpenCollection(Client, "menu")

// add an item
func AddMenuItem(c *gin.Context) {
	var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
	var item models.MenuItem

	if err := c.BindJSON(&item); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	validationErr := validateItem.Struct(item)
	if validationErr != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": validationErr.Error()})
		return
	}
	item.ID = primitive.NewObjectID()

	result, insertErr := menuCollection.InsertOne(ctx, item)
	if insertErr != nil {
		msg := fmt.Sprintf("order item was not created")
		c.JSON(http.StatusInternalServerError, gin.H{"error": msg})
		return
	}

	defer cancel()
	c.JSON(http.StatusOK, result)
}

func GetItemDescription(c *gin.Context) {
	itemName := c.Params.ByName("foodName")
	result := helpers.ItemExists(menuCollection, helpers.ParamsToCheck{Key: "foodname", Value: itemName}, c)
	if len(result) > 0 {
		fmt.Println("result>? ", result["description"])
		c.JSON(http.StatusOK, gin.H{"code": http.StatusOK, "message": "Item found successfully", "description": result["description"]})
	}
}

func DeleteMenuItem(c *gin.Context) {
	itemID := c.Params.ByName("id")
	docID, _ := primitive.ObjectIDFromHex(itemID)
	var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)

	result := helpers.ItemExists(menuCollection, helpers.ParamsToCheck{Key: "_id", Value: docID}, c)

	if len(result) > 0 {
		result, err := menuCollection.DeleteOne(ctx, bson.M{"_id": docID})
		fmt.Println("result = ", result)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			defer cancel()

			return
		}
		defer cancel()
		c.JSON(http.StatusOK, gin.H{"code": http.StatusOK, "message": "Item deleted successfully", "id": itemID})
	} else {
		c.Writer.WriteHeader(http.StatusNoContent)
	}
	defer cancel()
}
