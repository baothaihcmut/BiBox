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

var ALLOW_FILE_SORT_FIELD = []string{"created_at", "updated_at", "opened_at"}

// GetAllSubFileOfFolder implements FileInteractor.
func (f *FileInteractorImpl) GetAllSubFileOfFolder(ctx context.Context, input *presenters.GetSubFileOfFolderInput) (*presenters.GetSubFileOfFolderOutput, error) {
	userCtx := ctx.Value(constant.UserContext).(*commonModel.UserContext)
	userId, _ := primitive.ObjectIDFromHex(userCtx.Id)
	fileId, _ := primitive.ObjectIDFromHex(input.Id)
	_, err := f.checkFilePermission(ctx, fileId, userId, enums.ViewPermission)
	if err != nil {
		return nil, err
	}
	//check if sort field is allowed
	args := repositories.FindFileWithPermissionArg{
		ParentFolderId: &fileId,
		IsFolder:       input.IsFolder,
		Offset:         input.Offset,
		Limit:          input.Limit,
		IsAsc:          input.IsAsc,
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
	pagination := response.PaginationResponse{
		Offset:  input.Offset,
		Limit:   input.Limit,
		Total:   count,
		HasNext: false,
		HasPrev: false,
	}
	if input.Offset+input.Limit < int(count) {
		nextOffset := input.Offset + input.Limit
		pagination.HasNext = true
		pagination.NextOffset = &nextOffset
	}
	if input.Offset > 0 {
		prevOffset := input.Offset - input.Limit
		pagination.HasPrev = true
		pagination.PrevOffset = &prevOffset
	}
	return &presenters.GetSubFileOfFolderOutput{
		Data:       fileOutputs,
		Pagination: pagination,
	}, nil
}
