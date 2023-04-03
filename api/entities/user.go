package entities

import "time"

type User struct {
	address   string
	ensName   string
	createdAt time.Time
	updatedAt time.Time
}

type UserParams struct {
	Address   string
	EnsName   string
	CreatedAt time.Time
	UpdatedAt time.Time
}

func NewUser(params UserParams) User {
	return User{
		address:   params.Address,
		ensName:   params.EnsName,
		createdAt: params.CreatedAt,
		updatedAt: params.UpdatedAt,
	}
}

func (u *User) Address() string {
	return u.address
}

func (u *User) EnsName() string {
	return u.ensName
}

func (u *User) CreatedAt() time.Time {
	return u.createdAt
}

func (u *User) UpdatedAt() time.Time {
	return u.updatedAt
}
