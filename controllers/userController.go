package controllers

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/toshnaik/CloudBoard/initializers"
	"github.com/toshnaik/CloudBoard/models"
	"github.com/toshnaik/CloudBoard/utils"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func SignUp(c *gin.Context) {
	// get the email and password from the request body
	var body struct {
		Email string 	`json:"email" binding:"required"`
		Password string `json:"password" binding:"required"`
	}

	if c.Bind(&body) != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid body"})
		return
	}

	// check if the user already exists in the database
	var user models.User
	result := initializers.DB.Where("email = ?", body.Email).Find(&user)
	if result.RowsAffected > 0 {
		c.JSON(http.StatusConflict, gin.H{"error": "Account with this email already exists"})
		return
	}

	// hash the password using bcrypt
	hash, err := bcrypt.GenerateFromPassword([]byte(body.Password), bcrypt.DefaultCost)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error while hashing the password"})
		return
	}

	// create a new user in the database
	result = initializers.DB.Create(&models.User{Email: body.Email, Password: string(hash)})

	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error while creating the user"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User created successfully"})
}


func Login(c *gin.Context) {
	// get the email and password from the request body
	var body struct {
		Email string 	`json:"email" binding:"required"`
		Password string `json:"password" binding:"required"`
	}

	if c.Bind(&body) != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid body"})
		return
	}

	// check if the user exists in the database
	var user models.User
	result := initializers.DB.Where("email = ?", body.Email).First(&user)
	if result.RowsAffected == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Account with this email does not exist"})
		return
	}

	// compare the password with the hash
	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(body.Password))

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid password"})
		return
	}

	// create a refresh and access token
	refreshTokenDetails, err := utils.GenerateToken(fmt.Sprint(user.ID), time.Hour*24*365, os.Getenv("REFRESH_TOKEN_PRIVATE_KEY"))
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
		return
	}

	accessTokenDetails, err := utils.GenerateToken(fmt.Sprint(user.ID), time.Hour, os.Getenv("ACCESS_TOKEN_PRIVATE_KEY"))
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
		return
	}

	// save the refresh and access tokens to Redis
	err = initializers.Redis.Set(context.TODO(), refreshTokenDetails.TokenUuid, fmt.Sprint(user.ID), time.Hour*24*365).Err()
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
		return
	}
	
	err = initializers.Redis.Set(context.TODO(), accessTokenDetails.TokenUuid, fmt.Sprint(user.ID), time.Hour).Err()
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
		return
	}


	// save the access token in the cookie
	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie("access_token",
				*accessTokenDetails.Token,
				3600,
				"/",
				os.Getenv("DOMAIN"),
				false,
				true)
	// save the refresh token in the cookie
	c.SetCookie("refresh_token",
				*refreshTokenDetails.Token,
				3600*24*365,
				"/",
				os.Getenv("DOMAIN"),
				false,
				true)
	
	c.SetCookie("logged_in",
				"true",
				3600,
				"/",
				os.Getenv("DOMAIN"),
				false,
				true)

	c.JSON(http.StatusOK, gin.H{"message": "User logged in successfully",
								"access_token": accessTokenDetails.Token,})
}

func RefreshAccessToken(c *gin.Context) {
	// retrieve the refresh token from the cookie
	refresh_token, err := c.Cookie("refresh_token")
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "refresh token not found"})
		return
	}
	
	// verify the refresh token
	tokenClaims, err := utils.VerifyToken(refresh_token, os.Getenv("REFRESH_TOKEN_PUBLIC_KEY"))
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid refresh token"})
		return
	}

	// check if the refresh token is present in the Redis database
	userID, err := initializers.Redis.Get(context.TODO(), tokenClaims.TokenUuid).Result()
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "token expired"})
		return
	}

	// check if the user still exists in the database
	var user models.User
	err = initializers.DB.Where("id = ?", userID).First(&user).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "the user no longer exists"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}

	// create a new access token
	accessTokenDetails, err := utils.GenerateToken(fmt.Sprint(user.ID), time.Hour, os.Getenv("ACCESS_TOKEN_PRIVATE_KEY"))
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
		return
	}
	
	// save the access token in redis
	err = initializers.Redis.Set(context.TODO(), accessTokenDetails.TokenUuid, fmt.Sprint(user.ID), time.Hour).Err()
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
		return
	}

	// set the access token in the cookie
	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie("access_token",
				*accessTokenDetails.Token,
				3600,
				"/",
				os.Getenv("DOMAIN"),
				false,
				true)
	c.SetCookie("logged_in",
				"true",
				3600,
				"/",
				os.Getenv("DOMAIN"),
				false,
				true)

	c.JSON(http.StatusOK, gin.H{"message": "success",
								"access_token": accessTokenDetails.Token,})
}

func Logout(c *gin.Context) {
	// retrieve the refresh token from the cookie
	refresh_token, err := c.Cookie("refresh_token")
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "refresh token not found"})
		return
	}
	
	// verify the refresh token
	tokenClaims, err := utils.VerifyToken(refresh_token, os.Getenv("REFRESH_TOKEN_PUBLIC_KEY"))
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid refresh token"})
		return
	}

	accessTokenUuid, _ := c.MustGet("access_token_uuid").(string)
	
	// delete the access token and refresh token from the Redis database
	err = initializers.Redis.Del(context.TODO(), accessTokenUuid, tokenClaims.TokenUuid).Err()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie("access_token",
				"",
				-1,
				"/",
				os.Getenv("DOMAIN"),
				false,
				true)
	c.SetCookie("refresh_token",
				"",
				-1,
				"/",
				os.Getenv("DOMAIN"),
				false,
				true)
	c.SetCookie("logged_in",
				"",
				-1,
				"/",
				os.Getenv("DOMAIN"),
				false,
				true)
	c.JSON(http.StatusOK, gin.H{"message": "success"})
}