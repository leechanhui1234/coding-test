package main

import (
	"context"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Person struct {
	ID       string `form:"id" json:"id"`
	Contents string `form:"contents" json:"contents"`
}

var person = []Person{}

func remove(s int) []Person {
	return append(person[:s], person[s+1:]...)
}

func insertArray(collection *mongo.Collection, ctx context.Context) {
	filter := bson.D{{}} //모든 데이터 조회

	cur, _ := collection.Find(ctx, filter)

	for cur.Next(ctx) {
		t := Person{}
		err := cur.Decode(&t)

		if err != nil {
			log.Fatal(err)
			return
		}

		person = append(person, t)
	}

}

func selectDataAll(collection *mongo.Collection, r *gin.Engine, ctx context.Context) {
	r.GET("/person", func(c *gin.Context) {
		c.IndentedJSON(http.StatusOK, person)
	})
}

func selectData(collection *mongo.Collection, r *gin.Engine, ctx context.Context) {
	r.GET("/person/:id", func(c *gin.Context) {
		id, isExist := c.Params.Get("id")

		if !isExist {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "ID value is not exist",
			})

			return
		}

		filter := bson.D{{Key: "ID", Value: id}}
		cur, _ := collection.Find(ctx, filter)

		filterperson := []Person{}

		for cur.Next(ctx) {
			t := Person{}
			err := cur.Decode(&t)
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
				return
			}
			filterperson = append(filterperson, t)
		}

		if len(filterperson) == 0 {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "there is an empty data",
			})
		} else {
			c.IndentedJSON(http.StatusOK, filterperson)
		}
	})
}

func insertData(collection *mongo.Collection, r *gin.Engine, ctx context.Context) {
	r.POST("/person", func(c *gin.Context) {
		var newPerson Person

		if err := c.BindJSON(&newPerson); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})

			return
		}

		person = append(person, newPerson)

		data := bson.D{
			{Key: "ID", Value: newPerson.ID},
			{Key: "Contents", Value: newPerson.Contents},
		}

		_, e := collection.InsertOne(ctx, data)

		if e != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": e.Error()})
			return
		}

		c.IndentedJSON(http.StatusOK, person)
	})
}

func deleteData(collection *mongo.Collection, r *gin.Engine, ctx context.Context) {
	r.DELETE("/person/:id", func(c *gin.Context) {
		id, isExist := c.Params.Get("id")

		if !isExist {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "ID value is not exist",
			})

			return
		}

		for i, v := range person {
			if id == v.ID {
				if len(person) > 1 {
					person = remove(i)
				} else {
					person = []Person{}
				}

				filter := bson.D{
					{Key: "ID", Value: v.ID},
					{Key: "Contents", Value: v.Contents},
				}

				_, e := collection.DeleteOne(ctx, filter)

				if e != nil {
					c.JSON(http.StatusBadRequest, gin.H{"error": e.Error()})
					return
				}

				c.IndentedJSON(http.StatusOK, person)
				return
			}
		}

		c.JSON(http.StatusNotFound, gin.H{
			"error": "Not Found",
		})
	})
}

func updateData(collection *mongo.Collection, r *gin.Engine, ctx context.Context) {
	r.PATCH("/person/:id", func(c *gin.Context) {
		id, isExist := c.Params.Get("id")

		if !isExist {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "ID value is not exist",
			})

			return
		}

		var newPerson Person

		if err := c.Bind(&newPerson); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})

			return
		}

		for i, v := range person {
			if id == v.ID {
				person[i].Contents = newPerson.Contents

				filter := bson.D{{Key: "ID", Value: id}}

				update := bson.D{
					{"$set", bson.D{
						{Key: "Contents", Value: newPerson.Contents},
					}},
				}

				_, e := collection.UpdateOne(ctx, filter, update)

				if e != nil {
					c.JSON(http.StatusBadRequest, gin.H{"error": e.Error()})
					return
				}

				c.IndentedJSON(http.StatusOK, person)
				return
			}
		}

		c.JSON(http.StatusNotFound, gin.H{
			"error": "Not Found",
		})
	})
}

func main() {
	r := gin.Default()

	ctx, _ := context.WithCancel(context.Background())

	// Set client options
	clientOptions := options.Client().ApplyURI(`mongodb://leechanhui:qwer1234@localhost:20000/?connect=direct`)

	clientOptions.SetAuth(options.Credential{
		Username: "leechanhui",
		Password: "qwer1234",
	})

	// Connect to MongoDB
	client, err := mongo.Connect(ctx, clientOptions)

	if err != nil {
		log.Fatal(err)
		return
	}

	// Check the connection
	err = client.Ping(ctx, nil)

	if err != nil {
		log.Fatal(err)
		return
	}

	collection := client.Database("board").Collection("people")

	insertArray(collection, ctx)
	selectDataAll(collection, r, ctx)
	selectData(collection, r, ctx)
	insertData(collection, r, ctx)
	deleteData(collection, r, ctx)
	updateData(collection, r, ctx)

	r.Run("localhost:3000")
}
