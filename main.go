package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
)

const (
	userkey = "user"
)

func main() {
	r := engine()
	r.Use(gin.Logger())
	if err := engine().Run(":8080"); err != nil {
		log.Fatal("Unable to start:", err)
	}
}

func engine() *gin.Engine {
	r := gin.New()
	r.Use(sessions.Sessions("mysession", sessions.NewCookieStore([]byte("secret"))))
	r.POST("/login", login)
	r.GET("/", loginCheck)
	r.GET("/logout", logout)
	r.GET("/test", test)

	r.GET("/home", menu_home)
	r.GET("/timeline", menu_timeline)

	r.Static("dist", "./lib/dist")
	r.Static("vendors", "./lib/vendors")
	r.Static("image", "./resource/image")
	r.Static("templates", "./templates")
	r.LoadHTMLGlob("templates/*")

	//로그인이 되었을때만 보이는 요청
	private := r.Group("/private")
	private.Use(AuthRequired)
	{
		private.GET("/me", me)
		private.GET("/status", status)
	}
	return r
}
func AuthRequired(c *gin.Context) {
	session := sessions.Default(c)
	user := session.Get(userkey)
	if user == nil {
		// Abort the request with the appropriate error code
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	// Continue down the chain to handler etc
	c.Next()
}

func menu_home(c *gin.Context) {
	session := sessions.Default(c)
	user := session.Get(userkey)
	if user == nil {
		c.HTML(http.StatusOK, "welcome.html", gin.H{
			"title":   "Posts",
			"warning": "",
		})
	} else {
		c.HTML(http.StatusOK, "main.html", gin.H{
			"title":    "Posts",
			"contents": "home.html",
		})
	}
}
func menu_timeline(c *gin.Context) {
	session := sessions.Default(c)
	user := session.Get(userkey)
	if user == nil {
		c.HTML(http.StatusOK, "welcome.html", gin.H{
			"title":   "Posts",
			"warning": "",
		})
	} else {
		c.HTML(http.StatusOK, "main.html", gin.H{
			"title":    "Posts",
			"contents": "timeline.html",
		})
	}
}

func test(c *gin.Context) {
	c.HTML(http.StatusOK, "welcome.html", gin.H{
		"title": "Posts",
	})
}

// login is a handler that parses a form and checks for specific data
func login(c *gin.Context) {
	session := sessions.Default(c)
	username := c.PostForm("username")
	password := c.PostForm("password")

	// Validate form input
	if strings.Trim(username, " ") == "" || strings.Trim(password, " ") == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Parameters can't be empty"})
		return
	}

	// Check for username and password match, usually from a database
	if username != "hello" || password != "world" {
		//c.JSON(http.StatusUnauthorized, gin.H{"error": "Authentication failed"})
		c.HTML(http.StatusOK, "welcome.html", gin.H{
			"title":   "Posts",
			"warning": "id or password is not correct.",
		})
		return
	}

	// Save the username in the session
	session.Set(userkey, username) // In real world usage you'd set this to the users ID
	if err := session.Save(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save session"})
		return
	}

	//c.JSON(http.StatusOK, gin.H{"message": "Successfully authenticated user"})
	c.HTML(http.StatusOK, "main.html", gin.H{
		"title":    "Posts",
		"contents": "home.html",
	})
}

func logout(c *gin.Context) {
	session := sessions.Default(c)
	user := session.Get(userkey)
	if user == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid session token"})
		return
	}
	session.Delete(userkey)
	if err := session.Save(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save session"})
		return
	}
	//c.JSON(http.StatusOK, gin.H{"message": "Successfully logged out"})
	c.HTML(http.StatusOK, "welcome.html", gin.H{
		"title":   "Posts",
		"warning": "",
	})
}

func me(c *gin.Context) {
	session := sessions.Default(c)
	user := session.Get(userkey)
	c.JSON(http.StatusOK, gin.H{"user": user})
}

func loginCheck(c *gin.Context) {
	session := sessions.Default(c)
	user := session.Get(userkey)
	if user == nil {
		c.HTML(http.StatusOK, "welcome.html", gin.H{
			"title":   "Posts",
			"warning": "",
		})
	} else {
		c.HTML(http.StatusOK, "main.html", gin.H{
			"title":    "Posts",
			"contents": "home.html",
		})
	}

}

func status(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "You are logged in"})
}

func Time(c *gin.Context) {
	db, err := ConnectToDB()
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()

	var time string
	err = db.QueryRow("SELECT now()").Scan(&time)
	if err != nil {
		fmt.Println(err)
	}

	c.JSON(http.StatusOK, map[string]string{
		"Time": time,
	})
}

func ConnectToDB() (*sql.DB, error) {
	db, err := sql.Open("mysql", "root:teamimp0611#@tcp(35.192.24.180:3306)/teamimp")
	if err != nil {
		panic(err.Error())
	}
	if db != nil {
		db.SetMaxOpenConns(100)
		db.SetMaxIdleConns(10)
	}
	//defer db.Close()

	/* connection test
	var version string
	db.QueryRow("SELECT VERSION()").Scan(&version)
	fmt.Println("Connected to:", version)
	*/
	return db, err
}
