package interactors

import (
	"context"
	"errors"
	"fmt"

	"github.com/baothaihcmut/Bibox/storage-app/internal/common/constant"
	"github.com/baothaihcmut/Bibox/storage-app/internal/common/enums"
	"github.com/baothaihcmut/Bibox/storage-app/internal/common/exception"
	"github.com/baothaihcmut/Bibox/storage-app/internal/common/models"
	"github.com/baothaihcmut/Bibox/storage-app/internal/modules/file_comment/repositories"
	permissionRepo "github.com/baothaihcmut/Bibox/storage-app/internal/modules/file_permission/repositories"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var ErrPermissionDenied = errors.New("permission denied")

type CommentInteractor struct {
	Repo           repositories.CommentRepository          // changed to interface
	PermissionRepo permissionRepo.FilePermissionRepository // changed to interface
}

func NewCommentInteractor(commentRepo repositories.CommentRepository, permissionRepo permissionRepo.FilePermissionRepository) *CommentInteractor {
	return &CommentInteractor{
		Repo:           commentRepo,
		PermissionRepo: permissionRepo,
	}
}

func (ci *CommentInteractor) GetAllComments(ctx context.Context) ([]map[string]interface{}, error) {
	comments, err := ci.Repo.FetchCommentsWithUsersAndAnswers(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch comments: %w", err)
	}
	return comments, nil
}

func (ci *CommentInteractor) AddComment(ctx context.Context, fileID, content string) error {
	// get user context from token
	userContext, ok := ctx.Value(constant.UserContext).(*models.UserContext)
	if !ok {
		return exception.ErrUnauthorized
	}

	// convert fileID to ObjectID
	fileObjectID, err := primitive.ObjectIDFromHex(fileID)
	if err != nil {
		return exception.ErrInvalidObjectId
	}

	userID, err := primitive.ObjectIDFromHex(userContext.Id)
	if err != nil {
		return exception.ErrInvalidObjectId
	}

	// check user has permission to comment
	option := permissionRepo.FilterPermssionType{
		Option: permissionRepo.FilterPermssionOption(0),
		Value:  []enums.FilePermissionType{enums.CommentPermission, enums.EditPermission}, // use enums for permission types
	}
	filePermission, err := ci.PermissionRepo.GetFilePermission(ctx, fileObjectID, userID, option)
	if err != nil {
		return fmt.Errorf("failed to get file permission: %w", err)
	}
	if filePermission == nil {
		return ErrPermissionDenied
	}

	// user has permission
	err = ci.Repo.CreateComment(ctx, fileObjectID, userID, content)
	if err != nil {
		return fmt.Errorf("failed to create comment: %w", err)
	}
	return nil
}
