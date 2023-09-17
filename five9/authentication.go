package five9

// type authentication interface {
// 	AddFive9Authentication(r *http.Request)
// }

type AuthenticationRestAPI struct {
	Username string
	Password string
	TokenID  AuthenticationTokenID
	FarmID   FarmID
}
