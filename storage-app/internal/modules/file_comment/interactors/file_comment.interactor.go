package interactors

import (
	"context"
	"errors"

	"github.com/baothaihcmut/Bibox/storage-app/internal/common/constant"
	"github.com/baothaihcmut/Bibox/storage-app/internal/common/exception"
	"github.com/baothaihcmut/Bibox/storage-app/internal/common/models"
	"github.com/baothaihcmut/Bibox/storage-app/internal/modules/file_comment/repositories"
	permissionRepo "github.com/baothaihcmut/Bibox/storage-app/internal/modules/file_permission/repositories"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var ErrPermissionDenied = errors.New("permission denied")

type CommentInteractor struct {
	Repo           *repositories.CommentRepository
	PermissionRepo *permissionRepo.FilePermissionRepository
}

func NewCommentInteractor(commentRepo *repositories.CommentRepository, permissionRepo *permissionRepo.FilePermissionRepository) *CommentInteractor {
	return &CommentInteractor{
		Repo:           commentRepo,
		PermissionRepo: permissionRepo,
	}
}
func (ci *CommentInteractor) GetAllComments() ([]map[string]any, error) {
	return ci.Repo.FetchComments()
}

func (ci *CommentInteractor) AddComment(ctx context.Context, fileID, content string) error {
	// Get user context from token
	userContext, ok := ctx.Value(constant.UserContext).(*models.UserContext)
	if !ok {
		return exception.ErrUnauthorized
	}

	// Convert fileID to ObjectID
	fileObjectID, err := primitive.ObjectIDFromHex(fileID)
	if err != nil {
		return exception.ErrInvalidObjectId
	}

	userID, err := primitive.ObjectIDFromHex(userContext.Id)
	if err != nil {
		return exception.ErrInvalidObjectId
	}

	// Check if user has permission to comment
	// hasPermission, err := ci.PermissionRepo.CheckUserPermission(ctx, fileObjectID, userID, []int{2, 3}) // 2 = Comment, 3 = Edit
	// if err != nil {
	// 	return err
	// }
	// if !hasPermission {
	// 	return ErrPermissionDenied // User does not have permission
	// }

	// User has permission, add comment
	return ci.Repo.CreateComment(ctx, fileObjectID, userID, content)
}
