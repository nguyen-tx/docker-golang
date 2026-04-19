package types

// ── Common ────────────────────────────────────────────────────────────────────

type MessageResponse struct {
	Message string `json:"message"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}

// ── Auth ──────────────────────────────────────────────────────────────────────

type LoginResponse struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
}
