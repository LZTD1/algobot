package backoffice

import "time"

type Student struct {
	ID              int         `json:"id"`
	FirstName       string      `json:"firstName"`
	LastName        string      `json:"lastName"`
	FullName        string      `json:"fullName"`
	ParentName      string      `json:"parentName"`
	Email           string      `json:"email"`
	HasLaptop       int         `json:"hasLaptop"`
	Phone           string      `json:"phone"`
	Age             int         `json:"age"`
	BirthDate       time.Time   `json:"birthDate"`
	CreatedAt       time.Time   `json:"createdAt"`
	UpdatedAt       time.Time   `json:"updatedAt"`
	DeletedAt       interface{} `json:"deletedAt"`
	HasBranchAccess bool        `json:"hasBranchAccess"`
	Username        string      `json:"username"`
	Password        string      `json:"password"`
	LastGroup       Group       `json:"lastGroup"`
	Groups          []Group     `json:"groups"`
	Links           Links       `json:"_links"`
}
