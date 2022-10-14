package nodes

type GenericResponse struct {
	Data    interface{} `json:"data,omitempty"`
	Error   interface{} `json:"error,omitempty"`
	Success bool        `json:"success"`
}

type ErrorResponse struct {
	Message string `json:"message"`
	Code    int    `json:"code"`
	Status  int    `json:"-"`
}

var (
	ErrorInvalidAuthorizationToken = ErrorResponse{Code: 100, Message: "invalid authorization token", Status: 403}
	ErrorUserNotAuthorized         = ErrorResponse{Code: 101, Message: "you are not authorized to perform this action", Status: 403}
	ErrorInvalidOAuthCode          = ErrorResponse{Code: 102, Message: "invalid oauth code", Status: 400}
	ErrorUserNotFound              = ErrorResponse{Code: 103, Message: "user not found", Status: 404}
	ErrorUserBanned                = ErrorResponse{Code: 104, Message: "user banned", Status: 403}

	ErrorModNotFound     = ErrorResponse{Code: 200, Message: "mod not found", Status: 404}
	ErrorFailedModUpload = ErrorResponse{Code: 201, Message: "failed to upload mod", Status: 500}

	ErrorVersionNotFound = ErrorResponse{Code: 300, Message: "version not found", Status: 404}
)

func GenericUserError(err error) *ErrorResponse {
	if err == nil {
		return &ErrorResponse{
			Code:    1,
			Status:  400,
			Message: "unknown error",
		}
	}

	return &ErrorResponse{
		Code:    1,
		Status:  400,
		Message: err.Error(),
	}
}
