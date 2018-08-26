package controllers

// 注册用户
const UserRegisterNicknameLengthError = 4402
const UserRegisterNicknameUniqueError = 4403
const UserRegisterPasswordLengthError = 4404
const UserRegisterEmailFormatError = 4405
const UserRegisterEmailUniqueError = 4406
const PasswordNotEqual = 4408

// 用户登录
const UserLoginError = 4407

// 用户未登录
const Unauthenticated = 4444

// 上传文件超过限制
const UploadFileExceedLimit = 4448

// 服务器内部错误
const ServerInternalError = 5000
