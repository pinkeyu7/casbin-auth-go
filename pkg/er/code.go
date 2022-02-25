package er

const (
	ErrorParamInvalid          = 400400
	UnauthorizedError          = 400401
	ForbiddenError             = 400403
	ResourceNotFoundError      = 400404
	TokenExpiredError          = 401001
	AWSInitError               = 400405
	FirebaseIdTokenError       = 400409
	DataDuplicateError         = 400410
	LimitExceededError         = 400001
	DecryptError               = 400002
	UploadFileErrUnknown       = 400900
	UploadFileErrNotExist      = 400901
	UploadFileErrSizeOverLimit = 400902
	UploadFileErrFileName      = 400903
	UploadFileErrEmpty         = 400904
	UploadFileErrTypeNotMatch  = 400905
	UploadFileErrRowOverLimit  = 400906
	UnknownError               = 500000
	DBInsertError              = 500001
	DBUpdateError              = 500002
	DBDeleteError              = 500003
	DBDuplicateKeyError        = 500004
	RedisSerError              = 500005
)
