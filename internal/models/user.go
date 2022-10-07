package models

import (
	"time"

	"github.com/golang-jwt/jwt/v4"
)

type UserCreateParams struct {
	Username string  `validate:"required"`
	Email    string  `validate:"required,email"`
	Password string  `validate:"required"`
	Bio      *string `validate:"omitempty,max=500"`
	Avatar   *string `validate:"omitempty"`
}

type UserAuthParams struct {
	Username string `validate:"omitempty"`
	Email    string `validate:"required,email"`
	Password string `validate:"required"`
}

type UserClaims struct {
	jwt.RegisteredClaims
	Email string
	ID    int
}

type UserFriend struct {
	ID           int
	Username     string
	Avatar       string
	FriendsSince time.Time
}

// ActiveFriendRequest represents friend requests that has been sent to authenticated user
type ActiveFriendRequest struct {
	ID                int               // identity of user that sends friend request
	Username          string            // username of user
	Avatar            string            // avatar of user
	SentAt            time.Time         // datetime when request were send
	TotalPosts        int               // count of all user's posts
	AvgScore          float32           // average beers score of this user
	LatestBeers       BeersMinimized    // latest beers reviews
	FavoriteBeers     BeersMinimized    // beers with highest score by this user
	FavoriteBreweries FavoriteBreweries // breweries that have highest beer score by this user
}

type BreweryMinimized struct {
	ID   int    // brewery's identity
	Name string // brewery's name
}

type BeersMinimized struct {
	ID           int              // beers's identity
	Name         string           // beer name
	AvgUserScore float32          // average beer score of concrete user
	Brewery      BreweryMinimized // brewery that brews this beer
}

type BestBeer struct {
	ID       int     // beers's identity
	Name     string  // beer name
	AvgScore float32 // average users score of beer
}

type FavoriteBreweries struct {
	ID       int      // identity of brewery
	Name     string   // brewery name
	BestBeer BestBeer // brewery's best beer
}
