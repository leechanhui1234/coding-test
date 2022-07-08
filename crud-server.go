package main
 
import (
    "context"
    "fmt"
	"log"
	"net/http"
    "go.mongodb.org/mongo-driver/bson"
    "go.mongodb.org/mongo-driver/mongo"
    "go.mongodb.org/mongo-driver/mongo/options"
	"github.com/gin-gonic/gin"
)

type Person struct {
    ID string			`form:"id"`
	Contents string		`form:"contents"`
}

var person = [] Person{}

func remove(s int) []Person {
    return append(person[:s], person[s+1:]...)
}

func selectDataAll(collection *mongo.Collection, r *gin.Engine){
	filter := bson.D{{}} //모든 데이터 조회

	cur, _ := collection.Find(context.TODO(), filter)

	for cur.Next(context.TODO()){
		t := Person{}
		err := cur.Decode(&t)
		if err != nil {
			fmt.Println(err)
		}
		person = append(person, t)
	}

	r.LoadHTMLGlob("template/*")

	r.GET("/", func(c *gin.Context){
		c.HTML(http.StatusOK, "main_page.html", gin.H{
			"persons": person,
		})
	})
}

func insertData(collection *mongo.Collection, r *gin.Engine) {
	r.POST("/insertdata", func(c *gin.Context){
		id := c.PostForm("id")
		contents := c.PostForm("contents")

		d := Person{ID: id, Contents: contents,}

		data := bson.D{
			{Key:"ID",Value:id},
			{Key:"Contents",Value:contents},
		}
	 
		_, e := collection.InsertOne(context.TODO(), data)
		person = append(person, d)
		
		if e != nil{
			fmt.Println(e)
		}

		c.HTML(http.StatusOK, "about.html", gin.H{})
	})
}

func deleteData(collection *mongo.Collection, r *gin.Engine) {
	r.POST("/deletedata", func(c *gin.Context){
		id := c.PostForm("id")
		contents := c.PostForm("contents")
		
		filter := bson.D{
			{Key: "ID", Value: id},
			{Key: "Contents", Value: contents},
		}

    	_, e := collection.DeleteOne(context.TODO(), filter)
		
		for i, v := range person {
			if v.ID == id && v.Contents == contents {
				person = remove(i)
				break
			}
		}
		
		if e != nil{
			fmt.Println(e)
		}

		c.HTML(http.StatusOK, "about.html", gin.H{})
	})
}

func updateData(collection *mongo.Collection, r *gin.Engine) {
	r.POST("/updatedata", func(c *gin.Context){
		id := c.PostForm("id")
		contents := c.PostForm("contents")
		
		filter := bson.D{{Key: "ID", Value: id}}
 
		update := bson.D{
			{"$set", bson.D{
				{Key: "Contents", Value: contents},
			}},
		}
	
		_, e := collection.UpdateOne(context.TODO(), filter, update)
		
		if e != nil{
			fmt.Println(e)
			return
		}

		for i, v := range person {
			if v.ID == id {
				person[i].Contents = contents
				break
			}
		}

		c.HTML(http.StatusOK, "about.html", gin.H{})
	})
}

func selectData(collection *mongo.Collection, r *gin.Engine){
	r.POST("/selectdata", func(c *gin.Context){
		id := c.PostForm("id")
		filter := bson.D{{Key: "ID", Value: id}}
		cur, _ := collection.Find(context.TODO(), filter)
		
		filterperson := []Person{}

		for cur.Next(context.TODO()){
			t := Person{}
			err := cur.Decode(&t)
			if err != nil {
				fmt.Println(err)
			}
			filterperson = append(filterperson, t)
		}

		c.HTML(http.StatusOK, "query.html", gin.H{
			"persons": filterperson,
		})
	})
}
 
func CheckError(e error) {
    if e != nil {
        fmt.Println(e)
    }
}

func main() {
	r := gin.Default()

    // Set client options
    clientOptions := options.Client().ApplyURI("mongodb://localhost:20000")
	
	clientOptions.SetAuth(options.Credential{
		Username: "leechanhui",
		Password: "qwer1234",
	})
 
    // Connect to MongoDB
    client, err := mongo.Connect(context.TODO(), clientOptions)

	if err != nil {
		log.Fatal(err)
	}
 
    // Check the connection
    err = client.Ping(context.TODO(), nil)

	if err != nil {
		log.Fatal(err)
	}

    collection := client.Database("board").Collection("people")
	
	selectDataAll(collection, r)
	insertData(collection, r)
	deleteData(collection, r)
	updateData(collection, r)
	selectData(collection, r)

	r.Run("localhost:3000")
}

