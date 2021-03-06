package useCase

import "hello/server/domain"

type UserRepository interface {
	CreateUser(domain.User) error
	CreateGood(domain.Good) error
	SetIcon(string, string) error
	UpdateUser(email string, u domain.User) error
	CreatePost(string, domain.Post) error
	DeletePost(string, domain.Post) error
	UpdatePost(string, *domain.Post) error
	DeleteAllByUserEMail(string) error
	DeleteGoodByPostID(postID uint, userID uint) (err error)
	ReturnUserBYEMail(string) (domain.User, error)
	ReturnUserPostByName(string) (domain.User, error)
	ReturnUserAndPostByPostID(uint) (domain.User, error)
	ReturnAllUserPost() ([]domain.User, []domain.Post, []domain.Good, error)
	ReturnGoodedPostByWord(word string) (returnUsers []domain.User, returnPosts []domain.Post, returnGoods []domain.Good, err error)
	ReturnGoodedPost(userID uint) ([]domain.User, []domain.Post, []domain.Good, error)
	GuestLogin() (domain.User, error)
}
