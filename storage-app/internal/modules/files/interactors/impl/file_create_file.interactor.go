package impl

import (
	"context"
	"sync"
	"time"

	"github.com/baothaihcmut/Bibox/storage-app/internal/common/constant"
	"github.com/baothaihcmut/Bibox/storage-app/internal/common/enums"
	"github.com/baothaihcmut/Bibox/storage-app/internal/common/exception"
	commonModel "github.com/baothaihcmut/Bibox/storage-app/internal/common/models"
	"github.com/baothaihcmut/Bibox/storage-app/internal/common/response"
	"github.com/google/uuid"

	"github.com/baothaihcmut/Bibox/storage-app/internal/common/storage"
	permissionModel "github.com/baothaihcmut/Bibox/storage-app/internal/modules/file_permission/models"
	"github.com/baothaihcmut/Bibox/storage-app/internal/modules/notification/services"

	"github.com/baothaihcmut/Bibox/storage-app/internal/modules/files/models"
	userModel "github.com/baothaihcmut/Bibox/storage-app/internal/modules/users/models"

	"github.com/baothaihcmut/Bibox/storage-app/internal/modules/files/presenters"
	"github.com/samber/lo"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

const lockKey = "file:upload:lock"

func (f *FileInteractorImpl) CreatFile(ctx context.Context, input *presenters.CreateFileInput) (_ *presenters.CreateFileOutput, err error) {
	//get user context
	userContext := ctx.Value(constant.UserContext).(*commonModel.UserContext)
	//tranform id
	userId, err := primitive.ObjectIDFromHex(userContext.Id)
	if err != nil {
		return nil, exception.ErrInvalidObjectId
	}
	if !input.IsFolder && input.StorageDetail == nil {
		return nil, exception.ErrMissStorageDetail
	}
	//check phase
	checkWg := sync.WaitGroup{}
	checkErr := make(chan error, 1)
	userCh := make(chan *userModel.User, 1)
	parentFolderCh := make(chan *models.File, 1)
	allParentFoldersCh := make(chan []*models.File, 1)
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
		parentExistCh := make(chan *models.File, 1)
		checkWg.Add(1)
		go func() {
			defer checkWg.Done()
			parentFile, err := f.fileRepo.FindFileById(ctx, *input.ParentFolderID)
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
				parentExistCh <- nil
				checkErr <- exception.ErrParenFileNotExist
				return
			}
			parentExistCh <- parentFile
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
			parentFolderCh <- parentFile
		}()
		//get all parent
		if !input.IsFolder {
			go func() {
				parentExist := <-parentExistCh
				if parentExist != nil {
					parentFolders, err := f.fileRepo.FindAllParentFolder(ctx, *input.ParentFolderID)
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
					parentFolders = append(parentFolders, parentExist)
					allParentFoldersCh <- parentFolders
				}
			}()
		}
	}
	//wait for all check routine
	checkWg.Wait()
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
		input.IsFolder,
		input.TagIDs,
		storageArg,
	)

	//acquire lock
	var fileLockValue string
	if !file.IsFolder {
		fileLockKey := lockKey + file.ID.Hex()
		fileLockValue = uuid.New().String()
		ok, err := f.distrutedLockService.AcquireLock(ctx, fileLockKey, fileLockValue, 3, time.Second*3, time.Minute*30)
		if err != nil {
			return nil, err
		}
		if !ok {
			return nil, exception.ErrFileIsUploading
		}
		defer func(err error) {
			if err != nil {
				if err := f.distrutedLockService.ReleaseLock(ctx, fileLockKey, fileLockValue); err != nil {
					f.logger.Errorf(ctx, nil, "Error release lock: ", err)
				}
			}
		}(err)
	}

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
	//update all parent size
	if input.ParentFolderID != nil && input.StorageDetail != nil && !input.IsFolder {
		wgSave.Add(1)
		go func() {
			defer wgSave.Done()
			parentFolders := <-allParentFoldersCh
			for _, parentFolder := range parentFolders {
				parentFolder.IncrementTotalSize(input.StorageDetail.Size)
			}
			if err := f.fileRepo.BulkUpdateFile(ctx, parentFolders); err != nil {
				select {
				case <-ctx.Done():
					return
				default:
				}
				cancel()
				errSave <- err
			}
			//send nofication
			err := f.notificationService.SendNotificationFileUploaded(
				ctx,
				lo.Map(parentFolders, func(item *models.File, _ int) services.SendNotificationFileUploadedArg {
					return services.SendNotificationFileUploadedArg{
						FileOwnerId:  file.OwnerID,
						UserUploadId: userId,
						FileId:       file.ID,
					}
				}))
			if err != nil {
				f.logger.Errorf(ctx, nil, "Error send notification:", err)
			}
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
			parentFolder := <-parentFolderCh
			//if file is in folder extend all permision
			parentFilePermission, err := f.filePermissionRepo.FindPermssionByFileId(ctx, parentFolder.ID)
			if err != nil {
				select {
				case <-ctx.Done():
					return
				default:
				}
				cancel()
				errSave <- err
			}
			targetPermission := lo.Map(
				lo.Filter(parentFilePermission, func(item *permissionModel.FilePermission, _ int) bool {
					return item.UserID != userId
				}),
				func(item *permissionModel.FilePermission, _ int) *permissionModel.FilePermission {
					return permissionModel.NewFilePermission(
						file.ID,
						item.UserID,
						item.FilePermissionType,
						item.CanShare,
						item.ExpireAt,
					)
				})
			permssions = append(permssions, targetPermission...)
			//append owner edit of parent
			if userId != parentFolder.OwnerID {
				permssions = append(permssions, permissionModel.NewFilePermission(
					file.ID,
					parentFolder.OwnerID,
					enums.EditPermission,
					true,
					nil,
				))
			}
		}
		//else append edit permission for user
		permssions = append(permssions, permissionModel.NewFilePermission(
			file.ID,
			userId,
			enums.EditPermission,
			true,
			nil,
		))
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
		FileOutput: response.MapFileToFileOutput(file),
	}
	//get presign url for put object and store to cache

	if !file.IsFolder {
		if err := f.fileUploadProgressService.StartUpload(ctx, file.ID.Hex(), file.TotalSize); err != nil {
			return nil, err
		}
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
		output.UploadLockValue = fileLockValue
	}

	return output, nil

}
