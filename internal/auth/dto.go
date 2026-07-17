package auth

type MeResponse struct {
	FirstName      string `json:"firstName,omitempty"`
	Role           string `json:"role"`
	VerifiedStatus string `json:"verifiedStatus,omitempty"`
	Balance        any    `json:"balance,omitempty"`
	LogoURL        string `json:"logoUrl,omitempty"`

	NeedAgreementRedirect bool `json:"needAgreementRedirect,omitempty"`
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginResponse struct {
	AccessToken string `json:"accessToken"`
}
