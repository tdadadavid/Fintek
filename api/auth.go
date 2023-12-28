package api

import (
	"context"
	"database/sql"
	db "github/tdadadavid/fingreat/db/sqlc"
	utils "github/tdadadavid/fingreat/utils"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/lib/pq"
)


type Auth struct {
	server *Server
}

/////////////////////////////////
// Initalizes the router 
// for Auth.
////////////////////////////////
func (a Auth) router(server *Server) {
	a.server = server;

	serverGroup := server.router.Group("/auth");

	serverGroup.POST("/sign-up", a.register)
	serverGroup.POST("/login", a.login)
	serverGroup.POST("/logout", a.logout)
}

/////////////////////////////////
// POST /users
// Creates user
////////////////////////////////
func (a *Auth) register(ctx *gin.Context) {
	var user UserParams

	if err := ctx.ShouldBindJSON(&user); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	hashPassword, err := utils.GenerateHashedPassword(user.Password)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": err.Error() })
		return
	}

	args := db.CreateUserParams{
		Email: user.Email,
		HashedPassword: hashPassword,
	}

	newUser, err := a.server.queries.CreateUser(context.Background(), args)
	if err != nil {

		if pgError, ok := err.(*pq.Error); ok {
			if pgError.Code == "23505" {
				ctx.JSON(http.StatusBadRequest, gin.H{"error": "Account with specifications already exists.", "status": "Failed"});
				return
			}
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": err.Error() });
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"user": UserResponse{}.toUserResponse(&newUser) })
}



/////////////////////////////////
// POST /login 
// for Auth.
////////////////////////////////
func (a Auth) login(ctx *gin.Context)  {
	user := new(UserParams)

	if err := ctx.ShouldBindJSON(&user); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error() })
		return
	}

	foundUser, err := a.server.queries.GetUserByEmail(context.Background(), user.Email)
	if err == sql.ErrNoRows {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Incorrect user specifications" })
		return
	} else if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error() })
		return
	}

	// "Incorrect user specifications"
	if err := utils.VerifyPassword(user.Password, foundUser.HashedPassword); err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error":"Incorrect user specifications" })
		return
	}

	token, err := tokenController.CreateToken(foundUser.ID);

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error() })
		return
	}

	ctx.JSON(http.StatusOK, gin.H{ 
		"user": UserResponse{}.toUserResponse(&foundUser),
		"token": token,
	})
}

func (a Auth) logout(ctx *gin.Context) {

}