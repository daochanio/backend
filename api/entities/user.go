package entities

import "time"

type User struct {
	Address   string
	EnsName   string
	CreatedAt time.Time
	UpdatedAt time.Time
}

func NewUser() User {
	return User{}
}

func (u User) SetAddress(address string) User {
	u.Address = address
	return u
}

func (u User) SetEnsName(ensName string) User {
	u.EnsName = ensName
	return u
}

func (u User) SetCreatedAt(createdAt time.Time) User {
	u.CreatedAt = createdAt
	return u
}

func (u User) SetUpdatedAt(updatedAt time.Time) User {
	u.UpdatedAt = updatedAt
	return u
}

func (u User) GetAddress() string {
	return u.Address
}

func (u User) GetEnsName() string {
	return u.EnsName
}

func (u User) GetCreatedAt() time.Time {
	return u.CreatedAt
}

func (u User) GetUpdatedAt() time.Time {
	return u.UpdatedAt
}
