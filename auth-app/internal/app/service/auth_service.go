package service

import (
	"context"
	"time"

	appmodels "github.com/esmailemami/chess/auth/internal/app/models"
	"github.com/esmailemami/chess/auth/internal/consts"
	"github.com/esmailemami/chess/auth/internal/keys"
	"github.com/esmailemami/chess/auth/internal/models"

	"github.com/esmailemami/chess/shared/database/psql"
	"github.com/esmailemami/chess/shared/errs"
	"github.com/esmailemami/chess/shared/logging"
	sharedmodels "github.com/esmailemami/chess/shared/models"
	"github.com/google/uuid"
	"github.com/lestrrat-go/jwx/v2/jwa"
	"github.com/lestrrat-go/jwx/v2/jwt"
	"github.com/spf13/viper"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	userService *UserService
}

func NewAuthService(userService *UserService) *AuthService {
	return &AuthService{
		userService: userService,
	}
}

func (*AuthService) NewToken(claims map[string]any) (jwt.Token, error) {
	expiresIn := viper.GetInt("app.token_expires_in")
	if expiresIn == 0 {
		expiresIn = 1
	}

	builder := jwt.NewBuilder().
		Claim(jwt.IssuedAtKey, time.Now().Add(time.Second*-1).UTC()).
		Claim(jwt.NotBeforeKey, time.Now().Add(time.Second*-1).UTC()).
		Claim(jwt.ExpirationKey, time.Now().UTC().Add(time.Duration(expiresIn)*time.Hour)).
		Claim(jwt.JwtIDKey, uuid.New().String())

	for k, v := range claims {
		builder.Claim(k, v)
	}

	token, err := builder.Build()

	if err != nil {
		logging.ErrorE("failed to build token", err)
		return nil, err
	}

	return token, nil
}

func (*AuthService) SignToken(token jwt.Token) (string, error) {
	pk, err := keys.GetPrivateKey()
	if err != nil {
		return "", err
	}
	bts, err := jwt.Sign(token, jwt.WithKey(jwa.RS256, pk))

	if err != nil {
		logging.ErrorE("failed to generate signed payload", err)
		return "", err
	}
	return string(bts), nil
}

func (a *AuthService) GetToken(claims map[string]any) (string, error) {
	token, err := a.NewToken(claims)
	if err != nil {
		return "", err
	}
	return a.SignToken(token)
}

func (a *AuthService) ParseTokenString(tokenString string, validateToken bool) (jwt.Token, error) {
	pk, err := keys.GetPrivateKey()
	if err != nil {
		return nil, err
	}

	token, err := jwt.ParseString(
		tokenString,
		jwt.WithValidate(validateToken),
		jwt.WithKey(jwa.RS256, pk.PublicKey),
	)

	if err != nil {
		logging.ErrorE("failed to parse JWT token", err)
		return nil, err
	}

	return token, nil
}

func (a *AuthService) GetUser(ctx context.Context, tokenID uuid.UUID) (*sharedmodels.User, error) {
	db := psql.DBContext(ctx)

	// load token from DB
	var authToken models.AuthToken

	if err := db.Where(`"id"=?`, tokenID).Where(`"revoked"=?`, false).
		Where(`"expires_at">?`, time.Now()).First(&authToken).Error; err != nil {
		return nil, errs.UnAuthorizedErr().WithError(err)
	}

	// load token's related user
	var user sharedmodels.User

	if err := db.Where(`"id"=?`, authToken.UserID).Preload("Role").First(&user).Error; err != nil {
		return nil, errs.UnAuthorizedErr().WithError(err)
	}

	return &user, nil
}

func (a *AuthService) Register(ctx context.Context, req *appmodels.RegisterInputModel) (*sharedmodels.User, error) {
	dbData, err := req.ToDBModel()

	if err != nil {
		return nil, errs.InternalServerErr().WithError(err)
	}

	db := psql.DBContext(ctx)

	if err := db.Model(&sharedmodels.User{}).Create(dbData).Error; err != nil {
		return nil, errs.InternalServerErr().WithError(err)
	}

	return dbData, nil
}

func (a *AuthService) Login(ctx context.Context, req *appmodels.LoginInputModel) (*appmodels.LoginOutputModel, error) {
	user, err := a.userService.GetByUsername(ctx, req.Username)
	if err != nil {
		return nil, errs.BadRequestErr().WithError(err).Msg(consts.LoginFailed)
	}

	if bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)) != nil {
		return nil, errs.BadRequestErr().Msg(consts.LoginFailed)
	}

	if !user.Enabled {
		return nil, errs.BadRequestErr().Msg(consts.UserIsDisabled)
	}

	token, err := a.NewToken(map[string]interface{}{
		"userID":    user.ID,
		"username":  user.Username,
		"firstName": user.FirstName,
		"lastName":  user.LastName,
	})

	if err != nil {
		return nil, errs.InternalServerErr().WithError(err)
	}

	tokenStr, err := a.SignToken(token)

	if err != nil {
		return nil, errs.InternalServerErr().WithError(err)
	}

	output := &appmodels.LoginOutputModel{
		Token:     tokenStr,
		ExpiresAt: token.Expiration(),
		ExpiresIn: token.Expiration().Unix() - time.Now().Unix(),
		User: appmodels.LoginOutputUserModel{
			Username:  user.Username,
			FirstName: user.FirstName,
			LastName:  user.LastName,
		},
	}

	authToken := models.AuthToken{
		UserID:    user.ID,
		ExpiresAt: output.ExpiresAt,
		Revoked:   false,
	}

	authToken.ID = uuid.MustParse(token.JwtID())

	db := psql.DBContext(ctx)
	tx := db.Begin()

	err = tx.Create(&authToken).Error
	if err != nil {
		tx.Rollback()
		return nil, errs.InternalServerErr().WithError(err)
	}
	output.TokenID = authToken.ID

	history := models.LoginHistory{
		UserID:    user.ID,
		TokenID:   &output.TokenID,
		UserAgent: &req.UserAgent,
		IP:        &req.IP,
	}

	err = tx.Create(&history).Error
	if err != nil {
		tx.Rollback()
		return nil, errs.InternalServerErr().WithError(err)
	}

	tx.Commit()

	return output, nil
}

func (a *AuthService) RevokeToken(ctx context.Context, jwtID uuid.UUID) error {
	db := psql.DBContext(ctx)

	if err := db.Model(&models.AuthToken{}).
		Where("id = ?", jwtID).
		UpdateColumn("revoked", true).Error; err != nil {
		return errs.InternalServerErr().WithError(err)
	}

	return nil
}
