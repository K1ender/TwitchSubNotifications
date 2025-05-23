package types

type UserAccessToken string
type UserRefreshToken string

type Response[T any] struct {
	Data []T `json:"data"`
}
