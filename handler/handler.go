package handler

import (
	"coding-test/database"
	"coding-test/model"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
)

func SelectDataAll(c *gin.Context) { //R
	collection, ctx := database.GetDatabase()
	filter := bson.D{{}} //모든 데이터 조회
	cur, _ := collection.Find(ctx, filter)

	var person = []model.Person{}

	for cur.Next(ctx) {
		t := model.Person{}
		err := cur.Decode(&t)

		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"status": err.Error(),
			})
			return
		}

		person = append(person, t)

	}
	c.JSON(http.StatusOK, gin.H{
		"status": "OK",
		"data":   person,
	})
}

func SelectData(c *gin.Context) { //R
	collection, ctx := database.GetDatabase()

	id, isExist := c.Params.Get("id")

	if !isExist {
		c.JSON(http.StatusBadRequest, gin.H{
			"status": "ID is not exist",
		})

		return
	}

	filter := bson.D{{Key: "ID", Value: id}}
	cur, _ := collection.Find(ctx, filter)

	filterperson := model.Person{}

	for cur.Next(ctx) {
		err := cur.Decode(&filterperson)

		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"status": err.Error(),
			})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "OK",
		"data":   filterperson,
	})
}

func InsertData(c *gin.Context) { //C
	collection, ctx := database.GetDatabase()

	var newPerson model.Person

	if err := c.BindJSON(&newPerson); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status": err.Error(),
		})

		return
	}

	data, e := collection.InsertOne(ctx, newPerson)

	fmt.Println(data)

	if e != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": e.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "OK",
		"data":   newPerson,
	})
}

func DeleteData(c *gin.Context) { //D
	collection, ctx := database.GetDatabase()
	id, isExist := c.Params.Get("id")

	if !isExist {
		c.JSON(http.StatusBadRequest, gin.H{
			"status": "ID value is not exist",
		})

		return
	}

	filter := bson.D{
		{Key: "ID", Value: id},
	}

	data, e := collection.DeleteOne(ctx, filter)

	if e != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": e.Error()})
		return
	}

	fmt.Println(data)

	c.JSON(http.StatusOK, gin.H{
		"status":   "OK",
		"deleteid": id,
	})
}

func UpdateData(c *gin.Context) { //U
	collection, ctx := database.GetDatabase()
	id, isExist := c.Params.Get("id")

	if !isExist {
		c.JSON(http.StatusBadRequest, gin.H{
			"status": "ID value is not exist",
		})

		return
	}

	var newPerson model.Person

	if err := c.Bind(&newPerson); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status": err.Error(),
		})

		return
	}

	fmt.Println(newPerson)

	filter := bson.D{{Key: "ID", Value: id}}
	update := bson.D{
		{"$set", bson.D{
			{Key: "ID", Value: id},
			{Key: "Contents", Value: newPerson.Contents},
		}},
	}

	_, e := collection.UpdateOne(ctx, filter, update)

	if e != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": e.Error()})
		return
	}

	c.IndentedJSON(http.StatusOK, gin.H{
		"status":   "OK",
		"updateid": id,
	})
}
