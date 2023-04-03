package entities

import "time"

type User struct {
	address   string
	ensName   string
	createdAt time.Time
	updatedAt time.Time
}

func NewUser() User {
	return User{}
}

func (u User) SetAddress(address string) User {
	u.address = address
	return u
}

func (u User) SetEnsName(ensName string) User {
	u.ensName = ensName
	return u
}

func (u User) SetCreatedAt(createdAt time.Time) User {
	u.createdAt = createdAt
	return u
}

func (u User) SetUpdatedAt(updatedAt time.Time) User {
	u.updatedAt = updatedAt
	return u
}

func (u User) GetAddress() string {
	return u.address
}

func (u User) GetEnsName() string {
	return u.ensName
}

func (u User) GetCreatedAt() time.Time {
	return u.createdAt
}

func (u User) GetUpdatedAt() time.Time {
	return u.updatedAt
}
