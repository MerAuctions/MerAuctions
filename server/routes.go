package server

import (
	"encoding/json"
	"html/template"
	"io/ioutil"
	"log"
	"os"
	"time"

	"github.com/MerAuctions/MerAuctions/data"
	"github.com/MerAuctions/MerAuctions/models"
	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

var (
	authMiddleware *jwt.GinJWTMiddleware
	jwtIdentityKey = "id"

	cred  models.Credentials
	conf  *oauth2.Config
	state string
	store = sessions.NewCookieStore([]byte("secret"))
)

func setUpOauth2() {
	file, err := ioutil.ReadFile("/Users/vamshi.teja/go/src/github.com/MerAuctions/MerAuctions/server/creds.json")
	if err != nil {
		log.Printf("File error: %v\n", err)
		os.Exit(1)
	}
	json.Unmarshal(file, &cred)

	conf = &oauth2.Config{
		ClientID:     cred.Cid,
		ClientSecret: cred.Csecret,
		RedirectURL:  "http://me.mydomain.com:8081/auth",
		Scopes: []string{
			"https://www.googleapis.com/auth/userinfo.email", // You have to select your own scope from here -> https://developers.google.com/identity/protocols/googlescopes#google_sign-in
		},
		Endpoint: google.Endpoint,
	}
}

func setUpJWT() {
	// the jwt middleware
	var err error
	authMiddleware, err = jwt.New(&jwt.GinJWTMiddleware{
		Realm:          "MerAuction API",
		Key:            []byte("secret key"),
		IdentityKey:    jwtIdentityKey,
		Timeout:        time.Hour,
		MaxRefresh:     time.Hour,
		SendCookie:     true,
		SecureCookie:   false,
		CookieHTTPOnly: true,
		CookieName:     "token",
		TokenLookup:    "cookie:token",
		PayloadFunc: func(data interface{}) jwt.MapClaims {
			if v, ok := data.(*models.User); ok {
				return jwt.MapClaims{
					jwtIdentityKey: v.UserID,
				}
			}
			return jwt.MapClaims{}
		},
		IdentityHandler: func(c *gin.Context) interface{} {
			claims := jwt.ExtractClaims(c)
			return &models.User{
				UserID: claims[jwtIdentityKey].(string),
			}
		},
		Authenticator: func(c *gin.Context) (interface{}, error) {
			log.Println(c.Request.URL)
			loginVals := new(struct {
				Username string `form:"username" json:"username" binding:"required"`
				Password string `form:"password" json:"password" binding:"required"`
			})
			if err := c.ShouldBind(loginVals); err != nil {
				return "", jwt.ErrMissingLoginValues
			}
			userID := loginVals.Username
			password := loginVals.Password

			if user, err := data.GetUserById(userID); err != nil {
				log.Printf("Error in getting user details for '%s' from DB: %s", userID, err.Error())
				return "", jwt.ErrFailedAuthentication
			} else {
				if user == nil {
					log.Printf("User object is nil for user '%s'", userID)
					return "", jwt.ErrFailedAuthentication
				}
				if user.Password != password {
					log.Printf("Incorrect password for user '%s'", userID)
					return "", jwt.ErrFailedAuthentication
				}
				log.Printf("Successfully logged in user '%s'", userID)
				return user, nil
			}
		},
		Unauthorized: func(c *gin.Context, code int, message string) {
			c.JSON(code, gin.H{
				"code":    code,
				"message": message,
			})
		},
	})
	if err != nil {
		panic("failed to create jwt middleware")
	}
	log.Println("Successfully set up JWT")
}

func formatAuctionIDAsHexString(auctionID primitive.ObjectID) string {
	return auctionID.Hex()
}

func setupRoutes(router *gin.Engine) {
	setUpOauth2()
	setUpJWT()
	router.Use(sessions.Sessions("goquestsession", store))

	router.GET("/glogin", googleloginHandler)
	router.GET("/auth", googleauthHandler)

	router.POST("/login", authMiddleware.LoginHandler)
	router.POST("/users", addNewUser) //handle signing up

	router.SetFuncMap(template.FuncMap{
		"formatAuctionIDAsHexString": formatAuctionIDAsHexString,
	})

	if mode := gin.Mode(); mode == gin.TestMode {
		router.LoadHTMLGlob("./../templates/**/*")
	} else {
		router.LoadHTMLGlob("templates/**/*")
	}
	router.Static("/js", "./static/js")
	router.Static("/css", "./static/css")
	router.Static("/fonts", "./static/fonts")

	router.GET("/hello", hello)
	router.GET("/", getAllAuctions)
	router.GET("/auctions/:auction_id", getAuctionsByID)
	router.GET("/auctions/:auction_id/bids", getBidsAuctionsById)
	router.GET("/auctions/:auction_id/result", getResultByAuctionId)

	auth := router.Group("")
	auth.Use(authMiddleware.MiddlewareFunc())
	auth.GET("/refresh_token", authMiddleware.RefreshHandler)

	auth.POST("/auctions/:auction_id/bids", addBidAuctionIdByUserId)
	router.Run("me.mydomain.com:8081")
}
