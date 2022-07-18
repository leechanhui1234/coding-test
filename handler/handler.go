package handler

import (
	"coding-test/database"
	"coding-test/model"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
)

func SelectDataAll(c *gin.Context) { //R
	var collection, ctx = database.GetDatabase()
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
	var collection, ctx = database.GetDatabase()
	id, isExist := c.Params.Get("id")

	if !isExist {
		c.JSON(http.StatusBadRequest, gin.H{
			"status": "ID is not exist",
		})

		return
	}

	filter := bson.D{{Key: "id", Value: id}}
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

	if filterperson.ID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"status": "this is a empty data",
		})
	} else {
		c.JSON(http.StatusOK, gin.H{
			"status": "OK",
			"data":   filterperson,
		})
	}
}

func InsertData(c *gin.Context) { //C
	var collection, ctx = database.GetDatabase()
	var newPerson model.Person

	if err := c.BindJSON(&newPerson); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status": err.Error(),
		})

		return
	}

	_, e := collection.InsertOne(ctx, newPerson)

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
	var collection, ctx = database.GetDatabase()
	id, isExist := c.Params.Get("id")

	if !isExist {
		c.JSON(http.StatusBadRequest, gin.H{
			"status": "ID value is not exist",
		})

		return
	}

	filter := bson.D{
		{Key: "id", Value: id},
	}

	_, e := collection.DeleteOne(ctx, filter)

	if e != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": e.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":   "OK",
		"deleteid": id,
	})
}

func UpdateData(c *gin.Context) { //U
	var collection, ctx = database.GetDatabase()
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

	filter := bson.D{{Key: "id", Value: id}}
	update := bson.D{
		{"$set", bson.D{
			{Key: "id", Value: id},
			{Key: "title", Value: newPerson.Title},
			{Key: "contents", Value: newPerson.Contents},
		}},
	}

	_, e := collection.UpdateOne(ctx, filter, update)

	if e != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": e.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":   "OK",
		"updateid": id,
	})
}
