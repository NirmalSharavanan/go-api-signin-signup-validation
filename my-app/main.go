package main
 
import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gin-gonic/gin"
	"github.com/dgrijalva/jwt-go"
	"net/http"
	"os"
	"time"
	"log"
	"fmt"
)

var db *sql.DB
var err error

var (
  router = gin.Default()
)
 
func main() {
  db, err = sql.Open("mysql", "root:Little1515.@tcp(127.0.0.1:3306)/hrdb")
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		panic(err.Error())
  }
  
  router.POST("/login", Login)
  router.POST("/signup", Signup)
  log.Fatal(router.Run(":8080"))
}
type User struct {
  id int `json:"id"`
  UserName string `json:"username"`
  PassW string `json:"password"`
  FullName string `json:"fullname"`
}

func Signup(c *gin.Context) {
  var creds User
  if err := c.ShouldBindJSON(&creds); err != nil {
     c.JSON(http.StatusUnprocessableEntity, "Invalid json provided")
     return
  }

  var databaseUsername string

  err := db.QueryRow("SELECT UserName FROM Employee WHERE UserName=?", creds.UserName).Scan(&databaseUsername)

	switch {
  case err == sql.ErrNoRows:
    

    _, err = db.Exec("INSERT INTO Employee(id, UserName, PassW, FullName) VALUES(?, ?, ?, ?)", creds.id, creds.UserName, creds.PassW, creds.FullName)
    fmt.Println(err)
		if err != nil {
      c.JSON(http.StatusUnprocessableEntity, err.Error())
			return
		}
    c.JSON(http.StatusOK, "Insert Successful")
  }
}

func Login(c *gin.Context) {
  var creds User
  if err := c.ShouldBindJSON(&creds); err != nil {
     c.JSON(http.StatusUnprocessableEntity, "Invalid json provided")
     return
  }
  fmt.Println(c)
  fmt.Println(creds)

  var databaseId int
	var databaseUsername string
	var databasePassword string

	err := db.QueryRow("SELECT id, UserName, PassW FROM Employee WHERE UserName=?", creds.UserName).Scan(&databaseId, &databaseUsername, &databasePassword)
	
	// if there is an error opening the connection, handle it
	if err != nil {
		panic(err.Error())
  }
  fmt.Println(databaseUsername)
  fmt.Println(databasePassword)

  //compare the user from the request, with the one we defined:
  if databaseUsername != creds.UserName || databasePassword != creds.PassW {
     c.JSON(http.StatusUnauthorized, "Please provide valid login details")
     return
  }
  token, err := CreateToken(databaseId)
  if err != nil {
     c.JSON(http.StatusUnprocessableEntity, err.Error())
     return
  }
  c.JSON(http.StatusOK, token)
}

func CreateToken(userId int) (string, error) {
  var err error
  //Creating Access Token
  os.Setenv("ACCESS_SECRET", "jdnfksdmfksd") //this should be in an env file
  atClaims := jwt.MapClaims{}
  atClaims["authorized"] = true
  atClaims["user_id"] = userId
  atClaims["exp"] = time.Now().Add(time.Minute * 15).Unix()
  at := jwt.NewWithClaims(jwt.SigningMethodHS256, atClaims)
  token, err := at.SignedString([]byte(os.Getenv("ACCESS_SECRET")))
  if err != nil {
     return "", err
  }
  return token, nil
}