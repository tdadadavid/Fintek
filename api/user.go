package api

import (
	"context"
	"database/sql"
	db "github/tdadadavid/fingreat/db/sqlc"
	"github/tdadadavid/fingreat/utils"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type User struct {
	server *Server
}

type UserParams struct {
	Email string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type UpdateUserNameRequest struct {
	Username string `json:"username" binding:"required"`
}

type UserResponse struct {
	Email string `json:"email"`
  Id int64 `json:"id"`
	Username string `json:"username"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (u UserResponse) toUserResponse(user *db.User) *UserResponse {
	return &UserResponse {
		Email: user.Email,
		Id: user.ID,
		Username: user.Username.String,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}
}



/////////////////////////////////
// Iniitalizes the router 
// for users.
////////////////////////////////
func (u User) router(server * Server) {
	u.server = server;

	serverGroup := server.router.Group("/users", AuthenticatedMiddleware());

	serverGroup.GET("", u.listUsers)
	serverGroup.GET("me", u.getProfile)
	serverGroup.PATCH("username", u.updateProfile)
}



/////////////////////////////////
// GET /users
// Fethces all users
////////////////////////////////
func (u *User) listUsers(ctx *gin.Context) {
  args := db.ListUsersParams {
		Offset: 0,
		Limit: 10,
	}

	users, err := u.server.queries.ListUsers(context.Background(), args)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"errror": err.Error() })
	}

	newUser := []UserResponse{}

	for _, value := range users {
		user := UserResponse{}.toUserResponse(&value)
		newUser = append(newUser, *user);
	}

	ctx.JSON(http.StatusOK, gin.H{"users": newUser})
}


/////////////////////////////////
// POST /me
// Fethces user profile.
////////////////////////////////
func (u *User) getProfile(ctx *gin.Context) {
	userID, err := utils.GetActiveUser(ctx);
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return;
	}

	foundUser, err := u.server.queries.GetUserByID(context.Background(), userID);
	if err == sql.ErrNoRows {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Incorrect user specifications" })
		return
	} else if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"user": UserResponse{}.toUserResponse(&foundUser),
		"message": "User profile.",
	})
}


/////////////////////////////////
// POST /username
// Updates user profile name.
////////////////////////////////
func (u *User) updateProfile(ctx *gin.Context) {
	userID, err := utils.GetActiveUser(ctx);
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return;
	}

	updateUsernameReq := new(UpdateUserNameRequest)

	if err = ctx.ShouldBindJSON(&updateUsernameReq); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return;
	}

	updateUserNameArgs := db.UpdateUserNameParams{
		ID: userID,
		UpdatedAt: time.Now(),
		Username: sql.NullString{
			String: updateUsernameReq.Username,
			Valid: true,
		},
	}

	updatedUser, err := u.server.queries.UpdateUserName(context.Background(), updateUserNameArgs);
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error() })
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"user": UserResponse{}.toUserResponse(&updatedUser), 
		"message": "Username successfully updated.",
	})
}