package controllers

// 注册用户
const UserRegisterNicknameLengthError = 4402
const UserRegisterNicknameUniqueError = 4403
const UserRegisterPasswordLengthError = 4404
const UserRegisterEmailFormatError = 4405
const UserRegisterEmailUniqueError = 4406
const PasswordNotEqual = 4408
const OldPasswordInputError = 4409

// 用户登录
const UserLoginError = 4407

// 用户未登录
const Unauthenticated = 4444

// 上传文件超过限制
const UploadFileExceedLimit = 4448

// 服务器内部错误
const ServerInternalError = 5000

// 创建帖子
const ArticleTitleLengthError = 4410
const ArticleTitleContentError = 4411
const ArticleTopicNotExists = 4412

// 关注用户
const FollowUserNotExists = 4413
const FollowUserSelf = 4414
const FollowExists = 4415

// 查询帖子
const ArticleIdParamError = 4416

// 评论
const CommentIdParamError = 4417
