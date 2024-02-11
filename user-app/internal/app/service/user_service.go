package service

import (
	"context"
	"time"

	"github.com/esmailemami/chess/shared/database/psql"
	"github.com/esmailemami/chess/shared/errs"

	sharedModels "github.com/esmailemami/chess/shared/models"
	"github.com/esmailemami/chess/shared/service"
	appModels "github.com/esmailemami/chess/user/internal/app/models"
	"github.com/esmailemami/chess/user/internal/models"
	"github.com/esmailemami/chess/user/internal/util"
	"github.com/esmailemami/chess/user/pkg/consts"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	service.BaseService[sharedModels.User]
}

func NewUserService() *UserService {
	return &UserService{}
}

func (u *UserService) GetProfile(ctx context.Context, id uuid.UUID) (*appModels.UserProfileOutPutModel, error) {
	db := psql.DBContext(ctx).Table(`"user" u`).
		Joins("INNER JOIN role r ON r.id = u.role_id").
		Select("u.id, u.first_name, u.last_name, u.mobile, u.username, u.role_id, r.name as role_name, u.profile").
		Where("u.deleted_at is null").
		Where("u.id = ?", id)

	var resp appModels.UserProfileOutPutModel

	if err := db.Find(&resp).Error; err != nil {
		return nil, errs.InternalServerErr().WithError(err)
	}

	// set prefix of files
	resp.Profile = util.FilePathPrefix(resp.Profile)

	return &resp, nil
}

func (u *UserService) ChangePassword(ctx context.Context, id uuid.UUID, req *appModels.UserChangePasswordInputModel) error {
	if err := req.Validate(); err != nil {
		return errs.ValidationErr(err)
	}

	var user models.User

	db := psql.DBContext(ctx)

	if err := db.First(&user, id).Error; err != nil {
		return errs.NotFoundErr()
	}

	if bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)) != nil {
		return errs.BadRequestErr().Msg(consts.PasswordMismatch)
	}

	newPass, err := bcrypt.GenerateFromPassword([]byte(req.NewPassword), bcrypt.DefaultCost)

	if err != nil {
		return errs.InternalServerErr().WithError(err)
	}

	user.Password = string(newPass)

	if err := db.Save(&user).Error; err != nil {
		return errs.InternalServerErr().WithError(err)
	}

	return nil
}

func (u *UserService) ChangeProfile(ctx context.Context, id uuid.UUID, req *appModels.UserChangeProfileInputModel) error {
	if err := req.Validate(); err != nil {
		return errs.ValidationErr(err)
	}

	var user models.User

	db := psql.DBContext(ctx)

	if err := db.First(&user, id).Error; err != nil {
		return errs.NotFoundErr()
	}

	req.MergeDBModel(&user)

	if err := db.Save(&user).Error; err != nil {
		return errs.InternalServerErr().WithError(err)
	}

	return nil
}

func (u *UserService) UpdateLastConnection(ctx context.Context, userID uuid.UUID, lastConnection time.Time) error {
	db := psql.DBContext(ctx)

	if err := db.Model(&models.User{}).Where("id = ?", userID).
		UpdateColumn("last_connection", lastConnection).Error; err != nil {
		return errs.InternalServerErr().WithError(err)
	}

	return nil
}

func (u *UserService) UpdateProfile(ctx context.Context, userID uuid.UUID, profile string) error {
	db := psql.DBContext(ctx)

	if err := db.Model(&models.User{}).Where("id = ?", userID).
		UpdateColumn("profile", profile).Error; err != nil {
		return errs.InternalServerErr().WithError(err)
	}

	return nil
}
