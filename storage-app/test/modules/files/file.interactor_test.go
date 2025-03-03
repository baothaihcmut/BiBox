package files_test

import (
	"context"
	"testing"

	"github.com/baothaihcmut/Bibox/storage-app/internal/common/constant"
	"github.com/baothaihcmut/Bibox/storage-app/internal/common/enums"
	commonModel "github.com/baothaihcmut/Bibox/storage-app/internal/common/models"
	"github.com/baothaihcmut/Bibox/storage-app/internal/modules/files/interactors"
	"github.com/baothaihcmut/Bibox/storage-app/internal/modules/files/models"
	"github.com/baothaihcmut/Bibox/storage-app/internal/modules/files/presenters"
	userModel "github.com/baothaihcmut/Bibox/storage-app/internal/modules/users/models"
	"github.com/baothaihcmut/Bibox/storage-app/test/mocks"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.uber.org/mock/gomock"
)

func TestFileInteractor_CreatFile(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	userId := primitive.NewObjectID()
	parentFileId := primitive.NewObjectID()

	file := models.File{
		Name:        "test",
		OwnerID:     userId,
		IsFolder:    false,
		HasPassword: false,
		Password:    nil,
		Description: "",
		IsSecure:    false,
		StorageDetail: &models.FileStorageDetail{
			Size:            100,
			MimeType:        "application/pdf",
			IsUploaded:      false,
			StorageProvider: "s3",
			StorageKey:      "test",
			StorageBucket:   "test",
		},
	}

	mockUserRepo := mocks.NewMockUserRepository(ctrl)
	mockUserRepo.EXPECT().FindUserById(gomock.Any(), gomock.Any()).Return(
		&userModel.User{
			ID:                 userId,
			CurrentStorageSize: 0,
			LimitStorageSize:   100,
		},
		nil,
	)
	mockUserRepo.EXPECT().UpdateUserStorageSize(gomock.Any(), gomock.Any()).Return(nil)
	// mockTageRepo := mocks.NewMockTagRepository(ctrl)
	mockFileRepo := mocks.NewMockFileRepository(ctrl)
	mockFileRepo.EXPECT().FindFileById(gomock.Any(), gomock.Any(), false).Return(
		&models.File{
			ID:        parentFileId,
			IsDeleted: false,
			IsFolder:  true,
		}, nil)
	mockFileRepo.EXPECT().CreateFile(gomock.Any(), gomock.Any()).Return(nil)
	mockFilePermissionService := mocks.NewMockPermissionService(ctrl)
	mockFilePermissionService.EXPECT().CreatePermssion(gomock.Any(), gomock.Any()).Return(nil)
	mockFilePermissionService.EXPECT().CheckPermission(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(true, nil)
	mockLogger := mocks.NewMockLogger(ctrl)
	mockLogger.EXPECT().Info(gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes()
	mockLogger.EXPECT().Error(gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes()
	mockLogger.EXPECT().Infof(gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes()
	mockLogger.EXPECT().Errorf(gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes()
	mockStorageService := mocks.NewMockStorageService(ctrl)
	mockStorageService.EXPECT().GetStorageBucket().Return(file.StorageDetail.StorageBucket)
	mockStorageService.EXPECT().GetStorageProviderName().Return(file.StorageDetail.StorageProvider)
	mockStorageService.EXPECT().GetPresignUrl(gomock.Any(), gomock.Any()).Return("test", nil)
	mockMongoService := mocks.NewMockMongoService(ctrl)
	mockMongoService.EXPECT().BeginTransaction(gomock.Any()).Return(nil, nil)
	mockMongoService.EXPECT().CommitTransaction(gomock.Any(), gomock.Any()).Return(nil)
	mockMongoService.EXPECT().EndTransansaction(gomock.Any(), gomock.Any())
	interactor := interactors.NewFileInteractor(
		mockUserRepo,
		nil,
		nil,
		mockFilePermissionService,
		nil,
		mockLogger,
		mockStorageService,
		mockMongoService,
	)
	ouput, err := interactor.CreatFile(
		context.WithValue(context.Background(), constant.UserContext, &commonModel.UserContext{
			Id:   file.OwnerID.Hex(),
			Role: enums.UserRole,
		}),
		&presenters.CreateFileInput{
			Name:           file.Name,
			ParentFolderID: file.ParentFolderID,
			TagIDs:         file.TagIDs,
			IsFolder:       file.IsFolder,
			HasPassword:    file.HasPassword,
			Password:       file.Password,
			Description:    file.Description,
			IsSecure:       file.IsSecure,
			StorageDetail: &struct {
				Size     int    "json:\"size\" validate:\"required\""
				MimeType string "json:\"mime_type\" validate:\"required\""
			}{
				Size:     file.StorageDetail.Size,
				MimeType: string(file.StorageDetail.MimeType),
			},
		},
	)
	assert.Nil(t, err)
	assert.NotNil(t, ouput)
	assert.Equal(t, file.Name, ouput.Name)
	assert.Equal(t, file.OwnerID.Hex(), ouput.OwnerID)
	assert.Equal(t, file.IsFolder, ouput.IsFolder)
	assert.Equal(t, file.ParentFolderID, ouput.ParentFolderID)
	assert.Equal(t, file.HasPassword, ouput.HasPassword)
	assert.Equal(t, file.IsSecure, ouput.IsSecure)
	assert.Equal(t, file.StorageDetail.Size, ouput.StorageDetails.Size)
	assert.Equal(t, file.StorageDetail.MimeType, ouput.StorageDetails.MimeType)
	assert.Equal(t, "test", ouput.PutObjectUrl)
	assert.NotNil(t, ouput.ID)
	assert.NotNil(t, ouput.CreatedAt)
	assert.NotNil(t, ouput.UpdatedAt)
	assert.Nil(t, ouput.OpenedAt)
}
