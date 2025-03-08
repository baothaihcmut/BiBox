package interactors

import (
	"context"
	"sync"
	"time"

	"slices"

	"github.com/baothaihcmut/Bibox/storage-app/internal/common/constant"
	"github.com/baothaihcmut/Bibox/storage-app/internal/common/enums"
	"github.com/baothaihcmut/Bibox/storage-app/internal/common/exception"
	"github.com/baothaihcmut/Bibox/storage-app/internal/common/logger"
	commonModel "github.com/baothaihcmut/Bibox/storage-app/internal/common/models"
	"github.com/baothaihcmut/Bibox/storage-app/internal/common/mongo"
	"github.com/baothaihcmut/Bibox/storage-app/internal/common/response"
	"github.com/baothaihcmut/Bibox/storage-app/internal/common/storage"
	permissionModel "github.com/baothaihcmut/Bibox/storage-app/internal/modules/file_permission/models"
	filePermissionRepo "github.com/baothaihcmut/Bibox/storage-app/internal/modules/file_permission/repositories"
	"github.com/baothaihcmut/Bibox/storage-app/internal/modules/file_permission/services"
	"github.com/baothaihcmut/Bibox/storage-app/internal/modules/files/models"

	permissionPresenter "github.com/baothaihcmut/Bibox/storage-app/internal/modules/file_permission/presenters"
	"github.com/baothaihcmut/Bibox/storage-app/internal/modules/files/presenters"
	fileRepo "github.com/baothaihcmut/Bibox/storage-app/internal/modules/files/repositories"
	tagModel "github.com/baothaihcmut/Bibox/storage-app/internal/modules/tags/models"
	tagPresenter "github.com/baothaihcmut/Bibox/storage-app/internal/modules/tags/presenters"

	tagRepo "github.com/baothaihcmut/Bibox/storage-app/internal/modules/tags/repositories"
	userModel "github.com/baothaihcmut/Bibox/storage-app/internal/modules/users/models"
	userRepo "github.com/baothaihcmut/Bibox/storage-app/internal/modules/users/repositories"
	"github.com/samber/lo"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var ALLOW_FILE_SORT_FIELD = []string{"created_at", "updated_at", "opened_at"}

type FileInteractor interface {
	CreatFile(context.Context, *presenters.CreateFileInput) (*presenters.CreateFileOutput, error)
	UploadedFile(context.Context, *presenters.UploadedFileInput) (*presenters.UploadedFileOutput, error)
	FindAllFileOfUser(context.Context, *presenters.FindFileOfUserInput) (*presenters.FindFileOfUserOuput, error)
	GetFileMetaData(context.Context, *presenters.GetFileMetaDataInput) (*presenters.GetFileMetaDataOuput, error)
	GetFileTags(context.Context, *presenters.GetFileTagsInput) (*presenters.GetFileTagsOutput, error)
	GetFilePermissions(context.Context, *presenters.GetFilePermissionInput) (*presenters.GetFilePermissionOuput, error)
	GetFileDownloadUrl(context.Context, *presenters.GetFileDownloadUrlInput) (*presenters.GetFileDownloadUrlOutput, error)
	GetFileStructure(context.Context, *presenters.GetFileStructureInput) (*presenters.GetFileStructrueOuput, error)
	GetAllSubFileOfFolder(context.Context, *presenters.GetSubFileOfFolderInput) (*presenters.GetSubFileOfFolderOutput, error)
}

type FileInteractorImpl struct {
	userRepo           userRepo.UserRepository
	fileRepo           fileRepo.FileRepository
	tagRepo            tagRepo.TagRepository
	logger             logger.Logger
	storageService     storage.StorageService
	mongoService       mongo.MongoService
	filePermission     services.PermissionService
	filePermissionRepo filePermissionRepo.FilePermissionRepository
}

// GetAllSubFileOfFolder implements FileInteractor.
func (f *FileInteractorImpl) GetAllSubFileOfFolder(ctx context.Context, input *presenters.GetSubFileOfFolderInput) (*presenters.GetSubFileOfFolderOutput, error) {
	userCtx := ctx.Value(constant.UserContext).(*commonModel.UserContext)
	userId, _ := primitive.ObjectIDFromHex(userCtx.Id)
	fileId, _ := primitive.ObjectIDFromHex(input.Id)
	_, err := f.checkFilePermission(ctx, fileId, userId)
	if err != nil {
		return nil, err
	}
	//check if sort field is allowed
	args := fileRepo.FindFileOfUserArg{
		ParentFolderId: &fileId,
		IsFolder:       input.IsFolder,
		Offset:         input.Offset,
		Limit:          input.Limit,
		IsAsc:          input.IsAsc,
		PermssionLimit: 4,
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

// GetFileStructure implements FileInteractor.
func (f *FileInteractorImpl) GetFileStructure(ctx context.Context, input *presenters.GetFileStructureInput) (*presenters.GetFileStructrueOuput, error) {
	fileId, _ := primitive.ObjectIDFromHex(input.Id)
	subFiles, err := f.fileRepo.GetSubFileRecursive(ctx, fileId, []primitive.ObjectID{})
	if err != nil {
		return nil, err
	}
	return &presenters.GetFileStructrueOuput{
		SubFiles: lo.Map(subFiles, func(item *models.File, _ int) *presenters.FileOutput {
			return presenters.MapFileToFileOutput(item)
		}),
	}, nil

}

func (f *FileInteractorImpl) checkFilePermission(ctx context.Context, fileId primitive.ObjectID, userId primitive.ObjectID) (*models.File, error) {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()
	fileCh := make(chan *models.File, 1)
	errCh := make(chan error, 1)
	wg := sync.WaitGroup{}
	wg.Add(1)
	//check file exist
	go func() {
		defer wg.Done()
		file, err := f.fileRepo.FindFileById(ctx, fileId, false)
		if err != nil {
			select {
			case <-ctx.Done():
				return
			default:
			}
			cancel()
			errCh <- err
			return
		}
		if file == nil {
			cancel()
			errCh <- exception.ErrFileNotFound
			return
		}
		fileCh <- file
	}()
	wg.Add(1)
	go func() {
		defer wg.Done()
		hasPermission, err := f.filePermission.CheckPermission(ctx, fileId, userId, enums.ViewPermission)
		if err != nil {
			select {
			case <-ctx.Done():
				return
			default:
			}
			cancel()
			errCh <- err
			return
		}
		if !hasPermission {
			errCh <- exception.ErrPermissionDenied
			return
		}
	}()
	wg.Wait()
	select {
	case err := <-errCh:
		return nil, err
	default:
		return <-fileCh, nil
	}
}

// GetFileDownloadUrl implements FileInteractor.
func (f *FileInteractorImpl) GetFileDownloadUrl(ctx context.Context, input *presenters.GetFileDownloadUrlInput) (*presenters.GetFileDownloadUrlOutput, error) {
	//check permission and file find
	userCtx := ctx.Value(constant.UserContext).(*commonModel.UserContext)
	userId, _ := primitive.ObjectIDFromHex(userCtx.Id)
	fileId, _ := primitive.ObjectIDFromHex(input.Id)
	file, err := f.checkFilePermission(ctx, fileId, userId)
	if err != nil {
		return nil, err
	}
	//check if file is folder
	if file.IsFolder {
		return nil, exception.ErrFileIsFolder
	}
	//get url
	url, err := f.storageService.GetPresignUrl(ctx, storage.GetPresignUrlArg{
		Method:      storage.PresignUrlGetMethod,
		Key:         file.StorageDetail.StorageKey,
		ContentType: file.StorageDetail.MimeType,
		Expiry:      3 * time.Hour,
		Preview:     input.Preview,
	})
	if err != nil {
		return nil, err
	}
	return &presenters.GetFileDownloadUrlOutput{
		FileName:    file.Name,
		Url:         url,
		Expiry:      1,
		Method:      "GET",
		ContentType: file.StorageDetail.MimeType,
	}, nil
}

// GetFilePermissions implements FileInteractor.
func (f *FileInteractorImpl) GetFilePermissions(ctx context.Context, input *presenters.GetFilePermissionInput) (*presenters.GetFilePermissionOuput, error) {
	userCtx := ctx.Value(constant.UserContext).(*commonModel.UserContext)
	userId, _ := primitive.ObjectIDFromHex(userCtx.Id)
	fileId, _ := primitive.ObjectIDFromHex(input.Id)
	file, err := f.checkFilePermission(ctx, fileId, userId)
	if err != nil {
		return nil, err
	}
	permission, err := f.filePermissionRepo.GetPermissionOfFileWithUserInfo(ctx, file.ID)
	if err != nil {
		return nil, err
	}
	return &presenters.GetFilePermissionOuput{
		Permissions: lo.Map(permission, func(item *permissionModel.FilePermissionWithUser, _ int) *presenters.FilePermssionWithUserOutput {
			return &presenters.FilePermssionWithUserOutput{
				FilePermissionOuput: permissionPresenter.MapToOuput(item.FilePermission),
				User: &presenters.FilePermissionUserInfo{
					Email:     item.User.Email,
					FirstName: item.User.FirstName,
					LastName:  item.User.LastName,
					Image:     item.User.Image,
				},
			}
		}),
	}, nil
}

// GetFileTags implements FileInteractor.
func (f *FileInteractorImpl) GetFileTags(ctx context.Context, input *presenters.GetFileTagsInput) (*presenters.GetFileTagsOutput, error) {
	//check tag exist and check user have permission with file
	userCtx := ctx.Value(constant.UserContext).(*commonModel.UserContext)
	userId, _ := primitive.ObjectIDFromHex(userCtx.Id)
	fileId, _ := primitive.ObjectIDFromHex(input.Id)

	//check permission
	file, err := f.checkFilePermission(ctx, fileId, userId)
	if err != nil {
		return nil, err
	}
	tags, err := f.tagRepo.FindAllTagInList(ctx, file.TagIDs)
	if err != nil {
		return nil, err
	}
	//map to ouput

	return &presenters.GetFileTagsOutput{
		Tags: lo.Map(tags, func(item *tagModel.Tag, _ int) *tagPresenter.TagOutput {
			return tagPresenter.MaptoOuput(item)
		}),
	}, nil
}

// GetFileMetaData implements FileInteractor.
func (f *FileInteractorImpl) GetFileMetaData(ctx context.Context, input *presenters.GetFileMetaDataInput) (*presenters.GetFileMetaDataOuput, error) {
	userCtx := ctx.Value(constant.UserContext).(*commonModel.UserContext)
	userId, _ := primitive.ObjectIDFromHex(userCtx.Id)
	fileId, _ := primitive.ObjectIDFromHex(input.Id)

	//check permission
	file, err := f.checkFilePermission(ctx, fileId, userId)
	if err != nil {
		return nil, err
	}
	return &presenters.GetFileMetaDataOuput{
		FileOutput: presenters.MapFileToFileOutput(file),
	}, nil
}

func (f *FileInteractorImpl) CreatFile(ctx context.Context, input *presenters.CreateFileInput) (*presenters.CreateFileOutput, error) {
	//get user context
	userContext := ctx.Value(constant.UserContext).(*commonModel.UserContext)
	//tranform id
	userId, err := primitive.ObjectIDFromHex(userContext.Id)
	if err != nil {
		return nil, exception.ErrInvalidObjectId
	}
	//check phase
	checkWg := sync.WaitGroup{}
	checkErr := make(chan error, 1)
	userCh := make(chan *userModel.User, 1)
	ctx, cancelCheck := context.WithCancel(ctx)
	defer cancelCheck()
	//check user permission

	//check size of user
	if !input.IsFolder && input.StorageDetail != nil {
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
	for _, tagId := range input.TagIDs {
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

	if input.ParentFolderID != nil {
		checkWg.Add(1)
		go func() {
			defer checkWg.Done()
			parentFile, err := f.fileRepo.FindFileById(ctx, *input.ParentFolderID, false)
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
			if parentFile == nil {
				cancelCheck()
				checkErr <- exception.ErrParenFileNotExist
				return
			}

			//check user permssion in this folder
			permission, err := f.filePermission.CheckPermission(ctx, *input.ParentFolderID, userId, enums.EditPermission)
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
			if !permission {
				cancelCheck()
				checkErr <- exception.ErrPermissionDenied
				return
			}
		}()
	}
	//wait for all check routine
	checkWg.Wait()
	// close(userCh)
	// close(checkErr)
	select {
	case err = <-checkErr:
		return nil, err
	default:
	}

	//init file object
	//if file is not folder init storage detail
	var storageArg *models.FileStorageDetailArg
	if !input.IsFolder {
		storageArg = &models.FileStorageDetailArg{
			Size:            input.StorageDetail.Size,
			MimeType:        enums.MapToMimeType(input.Name, input.StorageDetail.MimeType),
			StorageProvider: f.storageService.GetStorageProviderName(),
			StorageBucket:   f.storageService.GetStorageBucket(),
		}
	}

	file := models.NewFile(
		userId,
		input.Name,
		input.ParentFolderID,
		input.Description,
		input.Password,
		input.IsFolder,
		input.HasPassword,
		input.IsSecure,
		input.TagIDs,
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
		f.mongoService.EndTransansaction(ctx, session)
	}()
	wgSave := sync.WaitGroup{}
	errSave := make(chan error, 1)
	doneCreateFileCh := make(chan struct{}, 1)
	//cancel context when have err
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()
	if !file.IsFolder {
		wgSave.Add(1)
		go func() {
			defer wgSave.Done()
			user := <-userCh
			err = f.userRepo.UpdateUserStorageSize(ctx, user)
			if err != nil {
				select {
				case <-ctx.Done():
					return
				default:
				}
				cancel()
				errSave <- err
			}
			f.logger.Info(ctx, map[string]any{
				"user_new_size": user.CurrentStorageSize,
			}, "User storage size updated")
		}()
	}
	//create file
	wgSave.Add(1)
	go func() {
		defer wgSave.Done()
		err = f.fileRepo.CreateFile(ctx, file)
		if err != nil {
			select {
			case <-ctx.Done():
				return
			default:
			}
			cancel()
			errSave <- err
		}
		doneCreateFileCh <- struct{}{}
		f.logger.Info(ctx, map[string]any{
			"file_id": file.ID.Hex(),
		}, "File created")
	}()
	//create  permission
	wgSave.Add(1)
	go func() {
		defer wgSave.Done()
		//get all permission of parent folder
		permssions := make([]*permissionModel.FilePermission, 0)
		if input.ParentFolderID != nil {
			//if file is in folder extend all permision
			parentFilePermission, err := f.filePermissionRepo.GetPermissionOfFile(ctx, *input.ParentFolderID)
			if err != nil {
				select {
				case <-ctx.Done():
					return
				default:
				}
				cancel()
				errSave <- err
			}
			for _, permission := range parentFilePermission {
				//with owner change it to edit permission

				permssions = append(permssions, permissionModel.NewFilePermission(
					file.ID,
					permission.UserID,
					permission.PermissionType,
					permission.CanShare,
					permission.ExpireAt,
				))
			}
		} else {
			//else append edit permission for user
			permssions = append(permssions, permissionModel.NewFilePermission(
				file.ID,
				userId,
				enums.EditPermission,
				true,
				nil,
			))
		}
		<-doneCreateFileCh
		err = f.filePermissionRepo.BulkCreatePermission(ctx, permssions)
		if err != nil {
			select {
			case <-ctx.Done():
				return
			default:
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
	err = f.mongoService.CommitTransaction(ctx, session)
	if err != nil {
		f.logger.Errorf(ctx, nil, "Error commit transaction mongo:", err)
		return nil, err
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
		f.logger.Info(ctx, map[string]any{
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
	err = f.fileRepo.UpdateFile(ctx, file)
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

	return &presenters.FindFileOfUserOuput{
		Data:       fileOutputs,
		Pagination: pagination,
	}, nil
}

func NewFileInteractor(
	userRepo userRepo.UserRepository,
	tagRepo tagRepo.TagRepository,
	fileRepo fileRepo.FileRepository,
	filePermission services.PermissionService,
	filePermissionRepo filePermissionRepo.FilePermissionRepository,
	logger logger.Logger,
	storageService storage.StorageService,
	mongoService mongo.MongoService,

) FileInteractor {
	return &FileInteractorImpl{
		userRepo:           userRepo,
		fileRepo:           fileRepo,
		tagRepo:            tagRepo,
		logger:             logger,
		storageService:     storageService,
		mongoService:       mongoService,
		filePermission:     filePermission,
		filePermissionRepo: filePermissionRepo,
	}
}
