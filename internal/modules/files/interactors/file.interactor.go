package interactors

import (
	"context"
	"sync"

	"github.com/baothaihcmut/Storage-app/internal/common/constant"
	"github.com/baothaihcmut/Storage-app/internal/common/exception"
	"github.com/baothaihcmut/Storage-app/internal/common/logger"
	commonModel "github.com/baothaihcmut/Storage-app/internal/common/models"
	"github.com/baothaihcmut/Storage-app/internal/common/mongo"
	"github.com/baothaihcmut/Storage-app/internal/common/storage"
	"github.com/baothaihcmut/Storage-app/internal/modules/files/models"
	"github.com/baothaihcmut/Storage-app/internal/modules/files/presenters"
	fileRepo "github.com/baothaihcmut/Storage-app/internal/modules/files/repositories"
	userRepo "github.com/baothaihcmut/Storage-app/internal/modules/users/repositories"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type FileInteractor interface {
	CreatFile(context.Context, *presenters.CreateFileInput) (*presenters.CreateFileOutput, error)
}

type FileInteractorImpl struct {
	userRepo       userRepo.UserRepository
	fileRepo       fileRepo.FileRepository
	logger         logger.Logger
	storageService storage.StorageService
	mongoService   mongo.MongoService
}

func (f *FileInteractorImpl) CreatFile(ctx context.Context, input *presenters.CreateFileInput) (*presenters.CreateFileOutput, error) {
	//get user context
	userContext := ctx.Value(string(constant.UserContext)).(*commonModel.UserContext)

	//check user permission

	//check size of user
	userId, err := primitive.ObjectIDFromHex(userContext.Id)
	if err != nil {
		return nil, exception.ErrInvalidObjectId
	}
	user, err := f.userRepo.FindUserById(ctx, userId)
	if err != nil {
		return nil, err
	}
	err = user.IncreStorageSize(input.StorageDetail.Size)
	if err != nil {
		return nil, err
	}

	//check if tag exists
	var parentFileId primitive.ObjectID
	if input.ParentFolderID != nil {
		parentFileId, err = primitive.ObjectIDFromHex(*input.ParentFolderID)
		if err != nil {
			return nil, exception.ErrInvalidObjectId
		}
		parentFile, err := f.fileRepo.FindFileById(ctx, parentFileId, false)
		if err != nil {
			return nil, err
		}
		if parentFile == nil {
			return nil, exception.ErrParenFileNotExist
		}
	}
	//init file object
	file := models.NewFile(
		user.ID,
		&parentFileId,
		input.Description,
		input.Password,
		input.IsFolder,
		input.HasPassword,
		input.IsSecure,
		[]primitive.ObjectID{},
		&struct {
			Size            int
			FileType        string
			StorageProvider string
			StorageBucket   string
		}{
			Size:            input.StorageDetail.Size,
			FileType:        input.StorageDetail.Type,
			StorageProvider: f.storageService.GetStorageBucket(),
			StorageBucket:   f.storageService.GetStorageProviderName(),
		})

	//save to db
	session, err := f.mongoService.BeginTransaction(ctx)
	if err != nil {
		f.logger.Errorf(ctx, nil, "Error init transaction mongo:", err)
		return nil, err
	}
	defer func() {
		if err != nil {
			f.mongoService.RollbackTransaction(ctx, session)
		}
		f.mongoService.RollbackTransaction(ctx, session)
	}()
	wgSave := sync.WaitGroup{}
	errSave := make(chan error, 1)
	//cancel context when have err
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()
	wgSave.Add(1)
	go func() {
		defer wgSave.Done()
		err = f.userRepo.UpdateUserStorageSize(ctx, user)
		if err != nil {
			if err == context.Canceled {
				return
			}
			cancel()
			errSave <- err
		}
	}()
	wgSave.Add(1)
	go func() {
		defer wgSave.Done()
		err = f.fileRepo.CreateFile(ctx, file)
		if err != nil {
			if err == context.Canceled {
				return
			}
			cancel()
			errSave <- err
		}
	}()
	//update parent file routine
	wgSave.Wait()
	select {
	case err = <-errSave:
		return nil, err
	default:
	}
	var storageOuputDetail *presenters.StorageDetailOuput
	//get presign url for put object
	if !file.IsFolder {
		url, err := f.storageService.GetPresignUrl(ctx, storage.GetPresignUrlArg{
			Method: storage.PresignUrlPutMethod,
			Key:    file.StorageDetail.StorageKey,
		})
		if err != nil {
			return nil, err
		}
		storageOuputDetail = &presenters.StorageDetailOuput{
			Size:         file.StorageDetail.Size,
			Type:         file.StorageDetail.FileType,
			PutObjectUrl: url,
			UrlExpiry:    3,
		}
	}
	return &presenters.CreateFileOutput{
		ID:             file.ID.Hex(),
		OwnerID:        file.OwnerID.Hex(),
		IsFolder:       file.IsFolder,
		ParentFolderID: input.ParentFolderID,
		CreatedAt:      file.CreatedAt,
		UpdatedAt:      file.UpdatedAt,
		DeletedAt:      file.DeletedAt,
		OpenedAt:       file.OpenedAt,
		HasPassword:    file.HasPassword,
		Description:    file.Description,
		IsSecure:       file.IsSecure,
		TotalSize:      file.TotalSize,
		TagIDs:         input.TagIDs,
		StorageDetails: storageOuputDetail,
	}, nil

}

func NewFileInteractor(
	userRepo userRepo.UserRepository,
	fileRepo fileRepo.FileRepository,
	logger logger.Logger,
	storageService storage.StorageService,
	mongoService mongo.MongoService,
) FileInteractor {
	return &FileInteractorImpl{
		userRepo:       userRepo,
		fileRepo:       fileRepo,
		logger:         logger,
		storageService: storageService,
		mongoService:   mongoService,
	}
}
