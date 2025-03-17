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

func (f *FileInteractorImpl) FindAllFileOfUser(ctx context.Context, input *presenters.FindFileOfUserInput) (*presenters.FindFileOfUserOuput, error) {
	//get user context
	userContext := ctx.Value(constant.UserContext).(*commonModel.UserContext)
	//tranform id
	userId, err := primitive.ObjectIDFromHex(userContext.Id)
	if err != nil {
		return nil, exception.ErrInvalidObjectId
	}

	//check if sort field is allowed
	args := repositories.FindFileOfUserArg{
		IsFolder:       input.IsFolder,
		Offset:         input.Offset,
		Limit:          input.Limit,
		IsAsc:          input.IsAsc,
		PermssionLimit: 4,
		OwnerId:        &userId,
		UserId:         userId,
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
	fileOutputs := make([]*presenters.FileWithPermissionOutput, len(data))
	for idx, file := range data {
		permissionOfFile := make([]*presenters.PermissionOfFileOuput, len(file.Permissions))
		for j, permission := range file.Permissions {
			permissionOfFile[j] = &presenters.PermissionOfFileOuput{
				UserID:         permission.UserID,
				PermissionType: permission.PermissionType,
				UserImage:      permission.UserImage,
			}
		}
		fileOutputs[idx] = &presenters.FileWithPermissionOutput{
			FileOutput:     presenters.MapFileToFileOutput(&file.File),
			Permissions:    permissionOfFile,
			PermissionType: file.PermissionType,
		}
	}

	return &presenters.FindFileOfUserOuput{
		Data:       fileOutputs,
		Pagination: response.InitPaginationResponse(int(count), input.Limit, input.Offset),
	}, nil
}
