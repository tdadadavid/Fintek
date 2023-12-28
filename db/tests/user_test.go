package test

import (
	"context"
	db "github/tdadadavid/fingreat/db/sqlc"
	"github/tdadadavid/fingreat/utils"
	"log"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)


func TestCreateUser(t *testing.T) {

	defer cleanUp()

	hashPassword  := genFakePassword();
	arg := genFakeUserData(hashPassword)

	user, err := genRandUser(arg);
	if err != nil {
		log.Fatal("Error creating random user")
	}

	assert.NoError(t, err);
	assert.NotEmpty(t, user)

	assert.Equal(t, user.Email, arg.Email)
	assert.Equal(t, user.HashedPassword, arg.HashedPassword)
	assert.WithinDuration(t, user.CreatedAt, time.Now(), 2*time.Second)
	assert.WithinDuration(t, user.UpdatedAt, time.Now(), 2*time.Second)

	// ensure uniqueness of email.
	duplicateUser, err := testQuery.CreateUser(context.Background(), arg)
	assert.Error(t, err)
	assert.Empty(t, duplicateUser)

} 


func TestUpdateUser(t *testing.T) {

	defer cleanUp()

	newHhashPassword  := genFakePassword();
	arg := genFakeUserData(newHhashPassword)

	user, err := genRandUser(arg);
	if err != nil {
		log.Fatal("Error creating random user")
	}
	
	updateArg := db.UpdateUserPasswordParams {
		HashedPassword: newHhashPassword,
		ID: user.ID,
	}

	updatedUser, err := testQuery.UpdateUserPassword(context.Background(), updateArg)


	assert.NoError(t, err);
	assert.NotEmpty(t, updatedUser)

	assert.Equal(t, updatedUser.HashedPassword, updateArg.HashedPassword)
	assert.WithinDuration(t, user.UpdatedAt, time.Now(), 10*time.Second)
}

func TestGetUserByEmail(t *testing.T) {

	defer cleanUp()

  newHhashPassword  := genFakePassword();
	arg := genFakeUserData(newHhashPassword)

	user, err := genRandUser(arg);
	if err != nil {
		log.Fatal("Error creating random user")
	}

	foundUser, err := testQuery.GetUserByEmail(context.Background(), arg.Email);

	assert.NoError(t, err);
	assert.NotEmpty(t, foundUser)

  assert.Equal(t, user.Email, foundUser.Email)
	assert.Equal(t, user.HashedPassword, foundUser.HashedPassword)
}


func TestGetUserByID(t *testing.T) {

	defer cleanUp()

	newHashPassword  := genFakePassword();
	arg := genFakeUserData(newHashPassword)

	user, err := genRandUser(arg);
	if err != nil {
		log.Fatal("Error creating random user")
	}

	foundUser, err := testQuery.GetUserByID(context.Background(), user.ID);

	assert.NoError(t, err);
	assert.NotEmpty(t, foundUser)

  assert.Equal(t, user.Email, foundUser.Email)
	assert.Equal(t, user.HashedPassword, foundUser.HashedPassword)
}

func TestGetAllUsers(t *testing.T) {
	defer cleanUp()
	
	var waitGroup sync.WaitGroup
	for i := 0; i < 100; i++ {
		waitGroup.Add(1)
		go func() {
			defer waitGroup.Done()
			createRandomUser()
		}()
	}

	waitGroup.Wait()

	args := db.ListUsersParams {
		Limit: 10,
		Offset: 0,
	}

	allUsers, err := testQuery.ListUsers(context.Background(), args)
  if err != nil {
		log.Fatal("Error fetching all users")
	}

	assert.NoError(t, err);
	// assert.NotEmpty(t, allUsers)

	assert.NotNil(t, allUsers);

	args2 := db.ListUsersParams {
		Limit: 10,
		Offset: 200,
	}

	allUsers2, err := testQuery.ListUsers(context.Background(), args2)
  if err != nil {
		log.Fatal("Error fetching all users")
	}

	assert.NoError(t, err);

  assert.Equal(t, len(allUsers2), 0);
}

func TestDeleteUserByID(t *testing.T) {
	defer cleanUp()

	newHhashPassword  := genFakePassword();
	arg := genFakeUserData(newHhashPassword)

	user, err := genRandUser(arg);
	if err != nil {
		log.Fatal("Error creating random user.")
	}


	err = testQuery.DeleteUser(context.Background(), user.ID); 

	assert.NoError(t, err);


	foundUser, err := testQuery.GetUserByID(context.Background(), user.ID);

	assert.Error(t, err);
	assert.Empty(t, foundUser);
}


///////////////////////////////
// Helper functions.         
//////////////////////////////

func cleanUp() {
	err := testQuery.DeleteAllUsers(context.Background())
	if err != nil {
		log.Fatalf("Error deleting all users {%v}", err)
	}
}

func createRandomUser() db.User {
	password := genFakePassword()
	userData := genFakeUserData(password)
	
	user, err := genRandUser(userData)
	if err != nil {
		log.Fatal("Error generating random user")
	}

	return user;
}

func genRandUser(args db.CreateUserParams) (db.User,error) {

	user, err := testQuery.CreateUser(context.Background(), args)
	if err != nil {
		log.Printf("%v", err)
		return user, err;
	}

	return user, nil;
}

func genFakePassword() (string) {
	hashPassword, err := utils.GenerateHashedPassword(utils.RandomString(8))
	if err != nil {
	  log.Fatal("Error generating hash password")
	}
	return string(hashPassword)
}

func genFakeUserData(hashedPassword string) db.CreateUserParams {
  return db.CreateUserParams {
		Email: utils.RandomEmail(),
		HashedPassword: hashedPassword,
	}
}