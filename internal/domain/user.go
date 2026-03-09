package domain

import(
	"time"
)

type User struct {
	id			UserID
	name  		string
	email		string
	iconUrl		string
	githubID 	int64
	createdAt 	time.Time
}

func NewUser(
		id UserID, name string, email string,
		iconUrl string, githubID int64, createdAt time.Time,
	) User {
	return User{
		id:			id,
		name:		name,
		email:		email,
		iconUrl: 	iconUrl,
		githubID: 	githubID,
		createdAt: 	createdAt,
	}
}

//ゲッター
func (u *User) ID() UserID {
    return u.id
}

func (u *User) Name() string {
	return  u.name
}

func (u *User) Email() string {
	return u.email
}

func (u *User) IconUrl() string {
	return u.iconUrl
}

func (u *User) GithubId() int64 {
	return u.githubID
}

func (u *User) CreatedAt() time.Time {
	return u.createdAt
}