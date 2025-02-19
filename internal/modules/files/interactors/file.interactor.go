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
	"github.com/baothaihcmut/Storage-app/internal/modules/files/services"
	tagRepo "github.com/baothaihcmut/Storage-app/internal/modules/tags/repositories"
	userModel "github.com/baothaihcmut/Storage-app/internal/modules/users/models"
	userRepo "github.com/baothaihcmut/Storage-app/internal/modules/users/repositories"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var ALLOW_FILE_SORT_FIELD = []string{"created_at", "updated_at", "opened_at"}

type FileInteractor interface {
	CreatFile(context.Context, *presenters.CreateFileInput) (*presenters.CreateFileOutput, error)
	UploadedFile(context.Context, *presenters.UploadedFileInput) (*presenters.UploadedFileOutput, error)
	FindAllFileOfUser(ctx context.Context, input *presenters.FindFileOfUserInput) (*presenters.FindFileOfUserOuput, error)
	GetFirstPageOfFiles(ctx context.Context, input *presenters.GetFirstPageInput) (*presenters.GetFirstPageOutput, error)
}

type FileInteractorImpl struct {
	userRepo         userRepo.UserRepository
	fileRepo         fileRepo.FileRepository
	tagRepo          tagRepo.TagRepository
	logger           logger.Logger
	storageService   storage.StorageService
	mongoService     mongo.MongoService
	firstPageService services.FirstPageService
}

func (f *FileInteractorImpl) CreatFile(ctx context.Context, input *presenters.CreateFileInput) (*presenters.CreateFileOutput, error) {
	//get user context
	userContext := ctx.Value(constant.UserContext).(*commonModel.UserContext)
	//tranform id
	userId, err := primitive.ObjectIDFromHex(userContext.Id)
	if err != nil {
		return nil, exception.ErrInvalidObjectId
	}
	var parentFileId *primitive.ObjectID
	if input.ParentFolderID != nil {
		id, err := primitive.ObjectIDFromHex(*input.ParentFolderID)
		if err != nil {
			return nil, exception.ErrInvalidObjectId
		}
		parentFileId = &id
	}
	tagIds := make([]primitive.ObjectID, len(input.TagIDs))
	for idx, tagId := range input.TagIDs {
		id, err := primitive.ObjectIDFromHex(tagId)
		if err != nil {
			return nil, exception.ErrInvalidObjectId
		}
		tagIds[idx] = id
	}
	//check phase
	checkWg := sync.WaitGroup{}
	checkErr := make(chan error, 1)
	userCh := make(chan *userModel.User, 1)
	ctx, cancelCheck := context.WithCancel(ctx)
	defer cancelCheck()
	//check user permission

	//check size of user
	if !input.IsFolder {
		checkWg.Add(1)
		go func() {
			defer checkWg.Done()
			user, err := f.userRepo.FindUserById(ctx, userId)
			if err != nil {
				select {
				case <-ctx.Done():
					return
				default:
				}
				cancelCheck()
				checkErr <- err
				return
			}
			if user == nil {
				cancelCheck()
				checkErr <- exception.ErrUserNotFound
				return
			}
			err = user.IncreStorageSize(input.StorageDetail.Size)
			if err != nil {
				cancelCheck()
				checkErr <- err
				return
			}
			//push user to save phase
			userCh <- user
		}()
	}

	//check if tags exists
	for _, tagId := range tagIds {
		checkWg.Add(1)
		go func() {
			defer checkWg.Done()
			tagExist, err := f.tagRepo.FindTagById(ctx, tagId)
			if err != nil {
				select {
				case <-ctx.Done():
					return
				default:
				}
				cancelCheck()
				checkErr <- err
				return
			}
			if tagExist == nil {
				cancelCheck()
				checkErr <- exception.ErrTagNotExist
				return
			}
		}()
	}

	if parentFileId != nil {
		checkWg.Add(1)
		go func() {
			defer checkWg.Done()
			parentFile, err := f.fileRepo.FindFileById(ctx, *parentFileId, false)
			if err != nil {
				if err == context.Canceled {
					return
				}
				cancelCheck()
				checkErr <- err
				return
			}
			if parentFile == nil {
				cancelCheck()
				checkErr <- exception.ErrParenFileNotExist
				return
			}
		}()
	}
	//wait for all check routine
	checkWg.Wait()
	close(userCh)
	close(checkErr)
	select {
	case err = <-checkErr:
		return nil, err
	default:
	}
	//get user result
	user := <-userCh

	//init file object
	//if file is not folder init storage detail
	var storageArg *models.FileStorageDetailArg
	if !input.IsFolder {
		storageArg = &models.FileStorageDetailArg{
			Size:            input.StorageDetail.Size,
			MimeType:        input.StorageDetail.MimeType,
			StorageProvider: f.storageService.GetStorageProviderName(),
			StorageBucket:   f.storageService.GetStorageBucket(),
		}
	}

	file := models.NewFile(
		user.ID,
		input.Name,
		parentFileId,
		input.Description,
		input.Password,
		input.IsFolder,
		input.HasPassword,
		input.IsSecure,
		tagIds,
		storageArg,
	)

	//save phase
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
	if !file.IsFolder {
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
			f.logger.Info(ctx, map[string]interface{}{
				"user_new_size": user.CurrentStorageSize,
			}, "User storage size updated")
		}()
	}
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
		f.logger.Info(ctx, map[string]interface{}{
			"file_id": file.ID.Hex(),
		}, "File created")
	}()
	//update parent file routine
	wgSave.Wait()
	close(errSave)
	select {
	case err = <-errSave:
		return nil, err
	default:
	}

	output := &presenters.CreateFileOutput{
		FileOutput: presenters.MapFileToFileOutput(file),
	}
	//get presign url for put object
	if !file.IsFolder {
		url, err := f.storageService.GetPresignUrl(ctx, storage.GetPresignUrlArg{
			Method: storage.PresignUrlPutMethod,
			Key:    file.StorageDetail.StorageKey,
		})
		if err != nil {
			return nil, err
		}
		f.logger.Info(ctx, map[string]interface{}{
			"url":     url,
			"key":     file.StorageDetail.StorageKey,
			"bucket":  file.StorageDetail.StorageBucket,
			"file_id": file.ID,
		}, "Presign url for put generated")
		output.PutObjectUrl = url
		output.UrlExpiry = 3
	}

	return output, nil

}

func (f *FileInteractorImpl) UploadedFile(ctx context.Context, input *presenters.UploadedFileInput) (*presenters.UploadedFileOutput, error) {
	//check file exist
	fileId, err := primitive.ObjectIDFromHex(input.Id)
	if err != nil {
		return nil, exception.ErrInvalidObjectId
	}
	file, err := f.fileRepo.FindFileById(ctx, fileId, false)
	if err != nil {
		return nil, err
	}
	if file == nil {
		return nil, exception.ErrFileNotFound
	}
	if file.IsFolder {
		return nil, exception.ErrFileIsFolder
	}
	file.StorageDetail.IsUploaded = true
	//update db
	err = f.fileRepo.UploadedFile(ctx, file)
	if err != nil {
		return nil, err
	}
	return &presenters.UploadedFileOutput{
		FileOutput: presenters.MapFileToFileOutput(file),
	}, nil
}

func (f *FileInteractorImpl) FindAllFileOfUser(ctx context.Context, input *presenters.FindFileOfUserInput) (*presenters.FindFileOfUserOuput, error) {
	//get user context
	userContext := ctx.Value(constant.UserContext).(*commonModel.UserContext)
	//tranform id
	userId, err := primitive.ObjectIDFromHex(userContext.Id)
	if err != nil {
		return nil, exception.ErrInvalidObjectId
	}
	//check if sort field is allowed
	args := fileRepo.FindFileOfUserArg{
		IsInFolder: input.IsFolder,
		IsFolder:   input.IsFolder,
		Offset:     input.Offset,
		Limit:      input.Limit,
		IsAsc:      input.IsAsc,
	}
	//check allow sort field
	allow := false
	for _, allowField := range ALLOW_FILE_SORT_FIELD {
		if allowField == input.SortBy {
			allow = true
			break
		}
	}
	if !allow {
		return nil, exception.ErrUnAllowedSortField
	}
	args.SortBy = input.SortBy

	res, err := f.fileRepo.FindAllFileOfUser(ctx, userId, args)
	if err != nil {
		return nil, err
	}
	fileOutput := make([]*presenters.FileOutput, len(res))
	for idx, file := range res {
		fileOutput[idx] = presenters.MapFileToFileOutput(file)
	}
	return &presenters.FindFileOfUserOuput{
		Files: fileOutput,
	}, nil
}

func (f *FileInteractorImpl) GetFirstPageOfFiles(ctx context.Context, input *presenters.GetFirstPageInput) (*presenters.GetFirstPageOutput, error) {
	//check file permission
	userContext := ctx.Value(constant.UserContext).(*commonModel.UserContext)

	//get file by id
	fileId, err := primitive.ObjectIDFromHex(input.FileId)
	if err != nil {
		return nil, err
	}
	file, err := f.fileRepo.FindFileById(ctx, fileId, false)
	if file.OwnerID.Hex() != userContext.Id {
		return nil, exception.ErrUserForbiddenFile
	}
	if file.IsFolder {
		return nil, exception.ErrFileIsFolder
	}
	//get file from storage
	fileStorage, err := f.storageService.GetFile(ctx, file.StorageDetail.StorageKey)
	if err != nil {
		return nil, err
	}
	defer fileStorage.Close()
	//get first page
	img, err := f.firstPageService.GetFirstPage(ctx, fileStorage, file.StorageDetail.MimeType, input.OuputType)
	if err != nil {
		return nil, err
	}
	return &presenters.GetFirstPageOutput{Image: img}, nil

}

func NewFileInteractor(
	userRepo userRepo.UserRepository,
	tagRepo tagRepo.TagRepository,
	fileRepo fileRepo.FileRepository,
	logger logger.Logger,
	storageService storage.StorageService,
	mongoService mongo.MongoService,
	firstPageService services.FirstPageService,
) FileInteractor {
	return &FileInteractorImpl{
		userRepo:         userRepo,
		fileRepo:         fileRepo,
		tagRepo:          tagRepo,
		logger:           logger,
		storageService:   storageService,
		mongoService:     mongoService,
		firstPageService: firstPageService,
	}
}
