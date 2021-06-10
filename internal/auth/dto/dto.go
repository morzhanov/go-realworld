package dto

type LoginInput struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type SignupInput struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type LoginDto struct {
	AccessToken string `json:"accessToken"`
}

type ValidateRestRequestInput struct {
	Path        string
	AccessToken string
}

type ValidateGrpcRequestInput struct {
	Method      string
	AccessToken string
}

type ValidateEventsRequestInput struct {
	Event       string
	AccessToken string
}
