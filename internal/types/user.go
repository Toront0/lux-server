package types

import "time"

type AuthUser struct {
	ID int `json:"id"`
	FirstName string `json:"firstName"`
	LastName string `json:"lastName"`
	ProfileImg string `json:"profileImg"`
}

type LoginUser struct {
	ID int `json:"id"`
	FirstName string `json:"firstName"`
	LastName string `json:"lastName"`
	Password string `json:"-"`
	ProfileImg string `json:"profileImg"`
}

type User struct {
	ID int `json:"id"`
	CreatedAt time.Time `json:"createdAt"`
	FirstName string `json:"firstName"`
	LastName string `json:"lastName"`
	ProfileImg string `json:"profileImg"`
	BannerImg string `json:"bannerImg"`
	Status string `json:"status"`
	FollowersAmount int `json:"followersAmount"`
	FolloweesAmount int `json:"followeesAmount"`
	UsersRelations string `json:"usersRelations"`
}

type UserDialog struct {
	SenderID int `json:"senderId"`
	ContactFName string `json:"contactFName"`
	ContactLName string `json:"contactLName"`
	ContactPImg string `json:"contactPImg"`
	Message string `json:"message"`
	ID int `json:"id"`
	CreatedAt time.Time `json:"createdAt"`
	ReceiverID int `json:"receiverId"`
}

type UserMessage struct {
	ID int `json:"id"`
	CreatedAt time.Time `json:"createdAt"`
	ReceiverID int `json:"receiverId"`
	SenderID int `json:"senderId"`
	Message string `json:"message"`
	ProfileImg string `json:"profileImg"`
}

type UserFriend struct {
	ID int `json:"id"`
	FriendID int `json:"friendId"`
	FirstName string `json:"firstName"`
	LastName string `json:"lastName"`
	ProfileImg string `json:"profileImg"`
}

type UserPreview struct {
	ID int `json:"id"`
	UserID int `json:"userId"`
	FirstName string `json:"firstName"`
	LastName string `json:"lastName"`
	ProfileImg string `json:"profileImg"`
}

type UserSettingsData struct {
	ID int `json:"id"`
	FirstName string `json:"firstName"`
	LastName string `json:"lastName"`
	Status string `json:"status"`
	ProfileImg string `json:"profileImg"`
	BannerImg string `json:"bannerImg"`
	Email string `json:"email`
}