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

	"github.com/baothaihcmut/Bibox/storage-app/internal/common/storage"
	permissionModel "github.com/baothaihcmut/Bibox/storage-app/internal/modules/file_permission/models"
	"github.com/baothaihcmut/Bibox/storage-app/internal/modules/file_permission/repositories"
	"github.com/baothaihcmut/Bibox/storage-app/internal/modules/files/models"
	userModel "github.com/baothaihcmut/Bibox/storage-app/internal/modules/users/models"

	"github.com/baothaihcmut/Bibox/storage-app/internal/modules/files/presenters"
	"github.com/baothaihcmut/Bibox/storage-app/internal/modules/files/services"
	"github.com/samber/lo"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// UploadFolder implements FileInteractor.
func (f *FileInteractorImpl) UploadFolder(ctx context.Context, folder *presenters.UploadFolderInput) (*presenters.UploadFolderOutput, error) {
	//get file list from foler tree
	//get user context
	userContext := ctx.Value(constant.UserContext).(*commonModel.UserContext)
	//tranform id
	userId, err := primitive.ObjectIDFromHex(userContext.Id)
	if err != nil {
		return nil, exception.ErrInvalidObjectId
	}
	files, totalSize := f.fileStructureService.TraverseUploadFolder(ctx, folder, userId, f.storageService.GetStorageProviderName(), f.storageService.GetStorageBucket())

	//check permission and size
	wgCheck := sync.WaitGroup{}
	ctx, cancelCheck := context.WithCancel(ctx)
	defer cancelCheck()
	errCheck := make(chan error, 1)
	userCh := make(chan *userModel.User, 1)
	parentFolderCh := make(chan *models.File, 1)
	allParentCh := make(chan []*models.File, 1)
	//check size
	if totalSize > 0 {
		wgCheck.Add(1)
		go func() {
			defer wgCheck.Done()
			user, err := f.userRepo.FindUserById(ctx, userId)
			if err != nil {
				select {
				case <-ctx.Done():
					return
				default:
				}
				cancelCheck()
				errCheck <- err
				return
			}
			err = user.IncreStorageSize(totalSize)
			if err != nil {
				cancelCheck()
				errCheck <- err
				return
			}
			userCh <- user
		}()
	}
	if folder.Data.ParentFolderID != nil {
		//check if parent folder exist
		parentFolderExistCh := make(chan *models.File, 1)
		wgCheck.Add(1)
		go func() {
			parentFolder, err := f.fileRepo.FindFileById(ctx, *folder.Data.ParentFolderID, false)
			if err != nil {
				select {
				case <-ctx.Done():
					return
				default:
				}
				cancelCheck()
				errCheck <- err
				return
			}
			if parentFolder == nil {
				cancelCheck()
				parentFolderExistCh <- nil
				errCheck <- exception.ErrParenFileNotExist
				return
			}
			parentFolderExistCh <- parentFolder
			parentFolderCh <- parentFolder
		}()
		//get parent of folder
		go func() {
			parentFolder := <-parentFolderExistCh
			if parentFolder != nil {
				allParentFolder, err := f.fileRepo.FindAllParentFolder(ctx, parentFolder.ID)
				if err != nil {
					select {
					case <-ctx.Done():
						return
					default:
					}
					cancelCheck()
					errCheck <- err
					return
				}
				allParentFolder = append(allParentFolder, parentFolder)
				allParentCh <- allParentFolder
			}
		}()
	}

	//store to db
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
	//save phase
	wgSave := sync.WaitGroup{}
	ctx, cancelSave := context.WithCancel(ctx)
	errSave := make(chan error, 1)
	defer cancelSave()
	wgSave.Add(1)
	go func() {
		defer wgSave.Done()
		err = f.fileRepo.BulkCreateFiles(ctx, lo.Map(files, func(item *services.FileWithPath, _ int) *models.File {
			return item.File
		}))
		if err != nil {
			select {
			case <-ctx.Done():
				return
			default:
			}
			cancelSave()
			errSave <- err
			return
		}
	}()
	//update user size
	if totalSize > 0 {
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
				cancelSave()
				errSave <- err
			}
			f.logger.Info(ctx, map[string]any{
				"user_new_size": user.CurrentStorageSize,
			}, "User storage size updated")
		}()
	}
	//update file size of all parent
	if folder.Data.ParentFolderID != nil && totalSize > 0 {
		wgSave.Add(1)
		go func() {
			defer wgSave.Done()
			allParentFolder := <-allParentCh
			for _, parentFolder := range allParentFolder {
				parentFolder.IncrementTotalSize(totalSize)
			}
			err = f.fileRepo.BulkUpdateTotalSize(ctx, allParentFolder)
			if err != nil {
				select {
				case <-ctx.Done():
					return
				default:
				}
				cancelSave()
				errSave <- err
			}
		}()
	}
	//extend permission of parent
	wgSave.Add(1)
	go func() {
		defer wgSave.Done()
		permissions := make([]*permissionModel.FilePermission, 0)

		if folder.Data.ParentFolderID != nil {
			parentFolder := <-parentFolderCh
			parentPermissions, err := f.filePermissionRepo.GetFilePermissions(ctx, repositories.GetPermissionArg{
				FileId:             &parentFolder.ID,
				FilePermissionType: enums.GetPermissionTypePointer(enums.FilePermissionType(enums.ViewPermission)),
			})
			if err != nil {
				select {
				case <-ctx.Done():
					return
				default:
				}
				cancelSave()
				errSave <- err
			}
			for _, file := range files {
				//add parent view permission
				targetPermission := lo.Map(parentPermissions, func(item *permissionModel.FilePermission, _ int) *permissionModel.FilePermission {
					return permissionModel.NewFilePermission(
						file.ID,
						item.UserID,
						item.FilePermissionType,

						item.CanShare,
						item.ExpireAt,
					)
				})
				permissions = append(permissions, targetPermission...)
				// add owner permission edit
				if parentFolder.OwnerID != userId {
					permissions = append(permissions, permissionModel.NewFilePermission(
						file.ID,
						parentFolder.OwnerID,
						enums.EditPermission,
						true,
						nil,
					))
				}
			}
		}
		userPermissions := lo.Map(files, func(item *services.FileWithPath, _ int) *permissionModel.FilePermission {
			return permissionModel.NewFilePermission(
				item.ID,
				userId,
				enums.EditPermission,
				true,
				nil,
			)
		})
		permissions = append(permissions, userPermissions...)
		err = f.filePermissionRepo.BulkCreatePermission(ctx, permissions)
		if err != nil {
			select {
			case <-ctx.Done():
				return
			default:
			}
			cancelSave()
			errSave <- err
		}
	}()
	wgSave.Wait()
	select {
	case err = <-errSave:
		return nil, err
	default:
	}
	//storage phase
	//respone
	wg := sync.WaitGroup{}
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()
	errCh := make(chan error, 1)
	resFile := make([]*presenters.FileWithPathOutput, len(files))
	for idx, file := range files {
		wg.Add(1)
		go func() {
			defer wg.Done()
			resFile[idx] = &presenters.FileWithPathOutput{
				FileOutput: response.MapFileToFileOutput(file.File),

				Path: file.Path,
			}

			if !file.IsFolder && file.StorageDetail != nil {
				//get presign url
				presignUrl, err := f.storageService.GetPresignUrl(ctx, storage.GetPresignUrlArg{
					Method:      storage.PresignUrlPutMethod,
					Key:         file.StorageDetail.StorageKey,
					ContentType: file.StorageDetail.MimeType,
					Expiry:      30 * time.Minute,
				})
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
				resFile[idx].PutObjectUrl = presignUrl
				resFile[idx].UrlExpiry = 30
			}
		}()
	}
	wg.Wait()
	select {
	case err = <-errCh:
		return nil, err
	default:
		return &presenters.UploadFolderOutput{
			Files: resFile,
		}, nil
	}

}
