package server

import (
	"log"
	"time"

	"github.com/MerAuctions/MerAuctions/data"
	"github.com/MerAuctions/MerAuctions/models"
	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
)

var (
	authMiddleware *jwt.GinJWTMiddleware
	identityKey    = "id"
)

func setUpJWT() {
	// the jwt middleware
	var err error
	authMiddleware, err = jwt.New(&jwt.GinJWTMiddleware{
		Realm:          "MerAuction API",
		Key:            []byte("secret key"),
		IdentityKey:    identityKey,
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
					identityKey: string(v.UserID),
				}
			}
			return jwt.MapClaims{}
		},
		IdentityHandler: func(c *gin.Context) interface{} {
			claims := jwt.ExtractClaims(c)
			if val, ok := claims[identityKey].(string); ok == false {
				return nil
			} else {
				return &models.User{
					UserID: models.ID(val),
				}
			}
		},
		Authenticator: func(c *gin.Context) (interface{}, error) {
			log.Println(c.Request.URL)
			loginVals := new(struct {
				Username string `form:"username" json:"username" binding:"required"`
				Password string `form:"password" json:"password" binding:"required"`
			})
			if err := c.ShouldBind(loginVals); err != nil {
				log.Printf("Missing Login Values in JWT")
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
				return &user, nil
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

func setupRoutes(router *gin.Engine) {

	setUpJWT()
	router.POST("/login", authMiddleware.LoginHandler)
	router.POST("/users", addNewUser) //handle signing up

	if mode := gin.Mode(); mode == gin.TestMode {
		router.LoadHTMLGlob("./../templates/**/*")
	} else {
		router.LoadHTMLGlob("templates/**/*")
	}
	router.Static("/js", "./static/js")
	router.Static("/css", "./static/css")

	router.GET("/hello", hello)
	router.GET("/", getAllAuctions)
	router.GET("/auctions/:auction_id", getAuctionsById)
	router.GET("/auctions/:auction_id/bids", getBidsAuctionsById)
	router.GET("/auctions/:auction_id/result", getResultByAuctionId)

	auth := router.Group("")
	auth.Use(authMiddleware.MiddlewareFunc())
	auth.GET("/refresh_token", authMiddleware.RefreshHandler)

	auth.POST("/auctions/:auction_id/bids", addBidAuctionIdByUserId)
}
