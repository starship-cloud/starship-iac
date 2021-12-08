package models

const(
	DefaultPageLimit int64 = 2

	DB_NAME string = "starship"
	DB_COLLECTION_USERS string = "users"
)

type PaginOption struct {
	Index int64
	Limit int64
}
