package errors

type ErrorResponse struct {
	Error              string `json:"error"`
	Description        string `json:"description,omitempty"`
	XRateLimitResetsIn int64  `json:"limitResetsIn,omitempty"`
}

var (
	DatabaseError             = &ErrorResponse{Error: "Internal server error", Description: "database"}
	ShortenedURLNotFoundError = &ErrorResponse{Error: "URL not found"}
	InvalidParametersError    = &ErrorResponse{Error: "Invalid parameters in request"}
	URLInUseError             = &ErrorResponse{Error: "URL already in use"}
	InvalidURLError           = &ErrorResponse{Error: "Invalid URL"}
	InvalidDomainError        = &ErrorResponse{Error: "Invalid domain"}
)
