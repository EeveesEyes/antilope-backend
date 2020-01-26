package models

type Secret struct {
	Id             int
	AuthorizedUser int
	Information    string
}
