package entities

import (
	"math/big"
	"time"
)

type User struct {
	address    string
	ensName    *string
	ensAvatar  *Image
	reputation *big.Int
	createdAt  time.Time
	updatedAt  *time.Time
}

type UserParams struct {
	Address    string
	EnsName    *string
	EnsAvatar  *Image
	Reputation *big.Int
	CreatedAt  time.Time
	UpdatedAt  *time.Time
}

func NewUser(params UserParams) User {
	return User{
		address:    params.Address,
		ensName:    params.EnsName,
		ensAvatar:  params.EnsAvatar,
		reputation: params.Reputation,
		createdAt:  params.CreatedAt,
		updatedAt:  params.UpdatedAt,
	}
}

func (u *User) Address() string {
	return u.address
}

func (u *User) EnsName() *string {
	return u.ensName
}

func (u *User) EnsAvatar() *Image {
	return u.ensAvatar
}

func (u *User) Reputation() *big.Int {
	return u.reputation
}

func (u *User) CreatedAt() time.Time {
	return u.createdAt
}

func (u *User) UpdatedAt() *time.Time {
	return u.updatedAt
}
