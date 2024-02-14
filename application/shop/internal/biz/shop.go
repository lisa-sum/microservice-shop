package biz

import (
	"context"
	"errors"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware/auth/jwt"
	jwt2 "github.com/golang-jwt/jwt/v4"
	"gorm.io/gorm"
	v1 "shop/api/shop/v1"
	"shop/internal/conf"
	"shop/internal/pkg/captcha"
	"shop/internal/pkg/middleware/auth"
	"time"
)

// 定义错误信息
var (
	ErrPasswordInvalid     = errors.New("password invalid")
	ErrUsernameInvalid     = errors.New("username invalid")
	ErrCaptchaInvalid      = errors.New("verification code error")
	ErrMobileInvalid       = errors.New("mobile invalid")
	ErrUserNotFound        = errors.New("user not found")
	ErrLoginFailed         = errors.New("login failed")
	ErrGenerateTokenFailed = errors.New("generate token failed")
	ErrAuthFailed          = errors.New("authentication failed")
)

// User 定义返回的数据的结构体
type User struct {
	ID        int64
	Mobile    string
	Nickname  string
	Birthday  int64
	Gender    string
	Role      int
	CreatedAt uint64 `gorm:"column:created_at"`
	UpdatedAt uint64 `gorm:"column:updated_at"`
	DeletedAt gorm.DeletedAt
	Password  string
}

//go:generate mockery --name UserRepo
type UserRepo interface {
	CreateUser(c context.Context, u *User) (*User, error)
	UserByMobile(ctx context.Context, mobile string) (*User, error)
	UserById(ctx context.Context, Id int64) (*User, error)
	CheckPassword(ctx context.Context, password, encryptedPassword string) (bool, error)
}

type UserUsecase struct {
	uRepo      UserRepo
	log        *log.Helper
	signingKey string // 这里是为了生存 token 的时候可以直接取配置文件里面的配置
}

func NewUserUsecase(repo UserRepo, logger log.Logger, conf *conf.Auth) *UserUsecase {
	helper := log.NewHelper(log.With(logger, "module", "usecase/shop"))
	return &UserUsecase{uRepo: repo, log: helper, signingKey: conf.JwtKey}
}

// GetCaptcha 验证码
func (uc *UserUsecase) GetCaptcha(ctx context.Context) (*v1.CaptchaReply, error) {
	captchaInfo, err := captcha.GetCaptcha(ctx)
	if err != nil {
		return nil, err
	}

	return &v1.CaptchaReply{
		CaptchaId: captchaInfo.CaptchaId,
		PicPath:   captchaInfo.PicPath,
	}, nil
}

func (uc *UserUsecase) UserDetailByID(ctx context.Context) (*v1.UserDetailResponse, error) {
	// 在上下文 context 中取出 claims 对象
	var uId int64
	if claims, ok := jwt.FromContext(ctx); ok {
		c := claims.(jwt2.MapClaims)
		if c["ID"] == nil {
			return nil, ErrAuthFailed
		}
		uId = int64(c["ID"].(float64))
	}

	user, err := uc.uRepo.UserById(ctx, uId)
	if err != nil {
		return nil, err
	}
	return &v1.UserDetailResponse{
		Id:       user.ID,
		Nickname: user.Nickname,
		Mobile:   user.Mobile,
	}, nil
}

func (uc *UserUsecase) PassWordLogin(ctx context.Context, req *v1.LoginReq) (*v1.RegisterReply, error) {
	// 表单验证
	if len(req.Mobile) <= 0 {
		return nil, ErrMobileInvalid
	}
	if len(req.Password) <= 0 {
		return nil, ErrUsernameInvalid
	}
	// 验证验证码是否正确
	if !captcha.Store.Verify(req.CaptchaId, req.Captcha, true) {
		return nil, ErrCaptchaInvalid
	}

	if user, err := uc.uRepo.UserByMobile(ctx, req.Mobile); err != nil {
		return nil, ErrUserNotFound
	} else {
		// 用户存在检查密码
		if passRsp, pasErr := uc.uRepo.CheckPassword(ctx, req.Password, user.Password); pasErr != nil {
			return nil, ErrPasswordInvalid
		} else {
			if passRsp {
				claims := auth.CustomClaims{
					ID:          user.ID,
					Nickname:    user.Nickname,
					AuthorityId: user.Role,
					RegisteredClaims: jwt2.RegisteredClaims{
						ExpiresAt: jwt2.NewNumericDate(time.Now().Add(24 * time.Hour)),
						IssuedAt:  jwt2.NewNumericDate(time.Now()),
						NotBefore: jwt2.NewNumericDate(time.Now()),
						// Subject:   "somebody",
						// ID:        "1",
						// Audience:  []string{"somebody_else"},
						Issuer: "Gyl",
					},
				}

				token, err := auth.CreateToken(claims, uc.signingKey)
				if err != nil {
					return nil, ErrGenerateTokenFailed
				}
				return &v1.RegisterReply{
					Id:        user.ID,
					Mobile:    user.Mobile,
					Username:  user.Nickname,
					Token:     token,
					ExpiredAt: time.Now().Unix() + 60*60*24*30,
				}, nil
			} else {
				return nil, ErrLoginFailed
			}
		}
	}
}

func (uc *UserUsecase) CreateUser(ctx context.Context, req *v1.RegisterReq) (*v1.RegisterReply, error) {
	newUser, err := NewUser(req.Mobile, req.Username, req.Password)
	if err != nil {
		return nil, err
	}
	createUser, err := uc.uRepo.CreateUser(ctx, &newUser)
	if err != nil {
		return nil, err
	}
	claims := auth.CustomClaims{
		ID:          createUser.ID,
		Nickname:    createUser.Nickname,
		AuthorityId: createUser.Role,
		RegisteredClaims: jwt2.RegisteredClaims{
			ExpiresAt: jwt2.NewNumericDate(time.Now().Add(24 * time.Hour)),
			IssuedAt:  jwt2.NewNumericDate(time.Now()),
			NotBefore: jwt2.NewNumericDate(time.Now()),
			// Subject:   "somebody",
			// ID:        "1",
			// Audience:  []string{"somebody_else"},
			Issuer: "Gyl",
		},
	}
	token, err := auth.CreateToken(claims, uc.signingKey)
	if err != nil {
		return nil, err
	}

	return &v1.RegisterReply{
		Id:        createUser.ID,
		Mobile:    createUser.Mobile,
		Username:  createUser.Nickname,
		Token:     token,
		ExpiredAt: time.Now().Unix() + 60*60*24*30,
	}, nil
}

func NewUser(mobile, username, password string) (User, error) {
	// check mobile
	if len(mobile) <= 10 {
		return User{}, ErrMobileInvalid
	}
	// check username
	if len(username) <= 3 {
		return User{}, ErrUsernameInvalid
	}
	// check password
	if len(password) <= 8 {
		return User{}, ErrPasswordInvalid
	}
	return User{
		Mobile:    mobile,
		Nickname:  username,
		Password:  password,
		CreatedAt: uint64(time.Now().Unix()),
	}, nil
}
