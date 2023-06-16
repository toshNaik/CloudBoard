package controllers

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/toshnaik/CloudBoard/initializers"
	"github.com/toshnaik/CloudBoard/models"
	"github.com/toshnaik/CloudBoard/utils"
	"golang.org/x/crypto/bcrypt"
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
	refreshTokenDetails, err := utils.GenerateToken(fmt.Sprint(user.ID), time.Hour*24*365, "REFRESH_TOKEN_PRIVATE_KEY")
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
		return
	}

	accessTokenDetails, err := utils.GenerateToken(fmt.Sprint(user.ID), time.Hour, "ACCESS_TOKEN_PRIVATE_KEY")
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
		return
	}

	// TODO: save the refresh token to Redis


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

	c.JSON(http.StatusOK, gin.H{"message": "User logged in successfully",
								"access_token": accessTokenDetails.Token,})
}