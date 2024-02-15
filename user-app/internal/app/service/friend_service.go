package service

import (
	"context"

	"github.com/esmailemami/chess/shared/database/psql"
	"github.com/esmailemami/chess/shared/errs"
	"github.com/esmailemami/chess/shared/util/dbutil"
	appmodels "github.com/esmailemami/chess/user/internal/app/models"
	"github.com/esmailemami/chess/user/internal/models"
	"github.com/esmailemami/chess/user/internal/util"
	"github.com/google/uuid"
)

type FriendService struct {
}

func NewFriendService() *FriendService {
	return &FriendService{}
}

func (f *FriendService) MakeFriend(ctx context.Context, currentUserID, friendID uuid.UUID) (*models.Friend, error) {
	db := psql.DBContext(ctx)

	if dbutil.Exists(&models.Friend{}, "user_id = ? AND friend_id = ?", currentUserID, friendID) {
		return nil, errs.BadRequestErr().Msg("already as your friend")
	}

	model := &models.Friend{
		UserID:   currentUserID,
		FriendID: friendID,
	}

	if err := db.Create(model).Error; err != nil {
		return nil, errs.InternalServerErr().WithError(err)
	}

	return model, nil
}

func (f *FriendService) RemoveFriend(ctx context.Context, currentUserID, friendID uuid.UUID) error {
	db := psql.DBContext(ctx)

	var friend models.Friend

	if err := db.First(&friend, "user_id = ? AND friend_id = ?", currentUserID, friendID); err != nil {
		return errs.BadRequestErr().Msg("user is not your friend")
	}

	if err := db.Delete(&friend).Error; err != nil {
		return errs.InternalServerErr().WithError(err)
	}

	return nil
}

func (f *FriendService) GetFriends(ctx context.Context, userID uuid.UUID, params *appmodels.FriendQueryParams) ([]appmodels.FriendOutPutModel, error) {
	db := psql.DBContext(ctx)

	var friends []appmodels.FriendOutPutModel

	qry := db.Model(&models.Friend{}).
		Joins("INNER JOIN public.user u on u.id = friend.friend_id").
		Where("friend.user_id = ?", userID)

	qry = dbutil.Filter(qry, params.SearchTerm, "u.first_name", "u.last_name", "u.username").
		Select("u.id, u.first_name,u.last_name,u.username,u.profile")

	if err := qry.Find(&friends).Error; err != nil {
		return nil, errs.InternalServerErr().WithError(err)
	}

	// full the profile path
	for i := 0; i < len(friends); i++ {
		friends[i].Profile = util.FilePathPrefix(friends[i].Profile)
	}

	return friends, nil
}
