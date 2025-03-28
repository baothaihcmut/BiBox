package impl

import (
	"context"
	"slices"

	"github.com/baothaihcmut/Bibox/storage-app/internal/common/constant"
	"github.com/baothaihcmut/Bibox/storage-app/internal/common/enums"
	"github.com/baothaihcmut/Bibox/storage-app/internal/common/exception"
	commonModel "github.com/baothaihcmut/Bibox/storage-app/internal/common/models"
	"github.com/baothaihcmut/Bibox/storage-app/internal/common/response"
	"github.com/baothaihcmut/Bibox/storage-app/internal/modules/files/presenters"
	"github.com/baothaihcmut/Bibox/storage-app/internal/modules/files/repositories"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (f *FileInteractorImpl) GetAllFileOfUser(ctx context.Context, input *presenters.GetAllFileOfUserInput) (*presenters.GetAllFileOfUserOuput, error) {
	//get user context
	userContext := ctx.Value(constant.UserContext).(*commonModel.UserContext)
	//tranform id
	userId, err := primitive.ObjectIDFromHex(userContext.Id)
	if err != nil {
		return nil, exception.ErrInvalidObjectId
	}

	//check if sort field is allowed
	args := repositories.FindFileWithPermissionArg{
		IsFolder:  input.IsFolder,
		Offset:    input.Offset,
		Limit:     input.Limit,
		IsAsc:     input.IsAsc,
		OwnerId:   &userId,
		UserId:    userId,
		IsDeleted: input.IsDeleted,
	}
	//check allow sort field
	if !slices.Contains(ALLOW_FILE_SORT_FIELD, input.SortBy) {
		return nil, exception.ErrUnAllowedSortField
	}
	args.SortBy = input.SortBy

	if input.FileType != nil && input.IsFolder != nil && !*input.IsFolder {
		fileType := enums.MapToMimeType("", *input.FileType)
		args.FileType = &fileType
	}

	data, count, err := f.fileRepo.FindFileWithPermssionAndCount(ctx, args)
	if err != nil {
		return nil, err
	}
	fileOutputs := make([]*response.FileWithPermissionOutput, len(data))
	for idx, file := range data {
		permissionOfFile := make([]*response.PermissionOfFileOuput, len(file.Permissions))
		for j, permission := range file.Permissions {
			permissionOfFile[j] = &response.PermissionOfFileOuput{
				UserID:             permission.UserID,
				FilePermissionType: permission.FilePermissionType,
				UserImage:          permission.UserImage,
			}
		}
		fileOutputs[idx] = &response.FileWithPermissionOutput{
			FileOutput:         response.MapFileToFileOutput(&file.File),
			Permissions:        permissionOfFile,
			FilePermissionType: file.FilePermissionType,
		}
	}

	return &presenters.GetAllFileOfUserOuput{
		Data:       fileOutputs,
		Pagination: response.InitPaginationResponse(int(count), input.Limit, input.Offset),
	}, nil
}
