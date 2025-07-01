package controllers

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	//"os"

	"auth2_google/internal/config"
	"auth2_google/internal/models"
	"auth2_google/internal/utils"
	"auth2_google/pkg/database"
	"net/http"
	//"net/url"

	"github.com/gin-gonic/gin"
)

func GoogleLogin(c *gin.Context) {
     state:= generateRandomState() 

	 //c.SetCookie("oauth_state", state, 600, "/", "", false, true)
     
	 url:=config.GoogleOAuthConfig.AuthCodeURL(state)

	 c.JSON(http.StatusOK, gin.H {
		 "auth_url" : url, 
		 "message":  "Redirect to this URL to login with Google",
		 "state":    state,
	 })

}

func GoogleCallback(c *gin.Context) {
	//  state:=c.Query("state")
	//  cookie, err := c.Cookie("oauth_state")
    //  if err != nil || state != cookie {
    //     c.JSON(http.StatusBadRequest, gin.H{
    //         "error": "Invalid state parameter - possible CSRF attack",
    //     })
    //     return
    //  }
	 code := c.Query("code")
     if code == "" {
        c.JSON(http.StatusBadRequest, gin.H{
            "error": "Authorization code not provided by Google",
        })
        return
     }
	 token, err := config.GoogleOAuthConfig.Exchange(c, code)
     if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{
            "error": "Failed to exchange code for token: " + err.Error(),
        })
        return
     }
	 userInfo, err := getUserInfoFromGoogle(token.AccessToken)
     if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{
            "error": "Failed to get user info from Google: " + err.Error(),
        })
        return
    }

    // Step 5: Save or update user in our database
     user, err := saveOrUpdateUser(userInfo)
     if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{
            "error": "Failed to save user to database: " + err.Error(),
        })
        return
    }
    jwtToken, err := utils.GenerateJWT(user.ID, user.Email, user.Name)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{
            "error": "Failed to generate authentication token: " + err.Error(),
        })
        return
    }
    c.JSON(http.StatusOK, gin.H{
        "message": "Successfully authenticated with Google!",
        "user": gin.H{
            "id":      user.ID,
            "name":    user.Name,
            "email":   user.Email,
            "picture": user.Picture,
        },
        "token":   jwtToken,
        "expires": "7 days",
    })

    // // 🆕 NEW: Store JWT in secure cookie
    // c.SetCookie("auth_token", jwtToken, 3600*24*7, "/", "", false, true)

    // // 🆕 UPDATED: Return YOUR JWT token instead of Google's
    // c.JSON(http.StatusOK, gin.H{
    //     "message": "Successfully authenticated with Google!",
    //     "user":    user,
    //     "token":   jwtToken, // ✅ Your JWT token (lasts 7 days)
    //     "expires": "7 days",
    // })
}

func generateRandomState() string {
    bytes := make([]byte, 32)
    rand.Read(bytes)
    return base64.URLEncoding.EncodeToString(bytes)
}

func getUserInfoFromGoogle(accessToken string) (*models.GoogleUserInfo, error) {
    // Call Google's userinfo API
    url := fmt.Sprintf("https://www.googleapis.com/oauth2/v2/userinfo?access_token=%s", accessToken)
    
    resp, err := http.Get(url)
    if err != nil {
        return nil, fmt.Errorf("failed to call Google API: %v", err)
    }
    defer resp.Body.Close()

    // Read response body
    body, err := io.ReadAll(resp.Body)
    if err != nil {
        return nil, fmt.Errorf("failed to read Google API response: %v", err)
    }

    // Parse JSON response
    var userInfo models.GoogleUserInfo
    if err := json.Unmarshal(body, &userInfo); err != nil {
        return nil, fmt.Errorf("failed to parse Google API response: %v", err)
    }

    return &userInfo, nil
}
func saveOrUpdateUser(googleUser *models.GoogleUserInfo) (*models.User, error) {
    var user models.User
    
    // Try to find existing user by Google ID
    result := database.DB.Where("google_id = ?", googleUser.ID).First(&user)
    
    if result.Error != nil {
        // User doesn't exist, create new one
        user = models.User{
            GoogleID: googleUser.ID,
            Email:    googleUser.Email,
            Name:     googleUser.Name,
            Picture:  googleUser.Picture,
        }
        
        if err := database.DB.Create(&user).Error; err != nil {
            return nil, fmt.Errorf("failed to create user: %v", err)
        }
        
        fmt.Printf("✅ New user created: %s (%s)\n", user.Name, user.Email)
    } else {
    // User exists, check if any data has changed before updating
    needsUpdate := user.Email != googleUser.Email || 
                   user.Name != googleUser.Name || 
                   user.Picture != googleUser.Picture
    
    if needsUpdate {
        user.Email = googleUser.Email
        user.Name = googleUser.Name
        user.Picture = googleUser.Picture
        
        if err := database.DB.Save(&user).Error; err != nil {
            return nil, fmt.Errorf("failed to update user: %v", err)
        }
        
        fmt.Printf("✅ User updated: %s (%s)\n", user.Name, user.Email)
    } else {
        fmt.Printf("✅ User login: %s (%s) - no updates needed\n", user.Name, user.Email)
    }
    }
    
    return &user, nil
}