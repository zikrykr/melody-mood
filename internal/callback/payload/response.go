package payload

type SessionResponse struct {
	SessionID string `json:"session_id"`
	ExpiresIn int64  `json:"expires_in"`
}
