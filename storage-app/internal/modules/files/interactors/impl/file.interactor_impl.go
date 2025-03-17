package impl

import (
	"github.com/baothaihcmut/Bibox/storage-app/internal/common/logger"
	"github.com/baothaihcmut/Bibox/storage-app/internal/common/mongo"
	"github.com/baothaihcmut/Bibox/storage-app/internal/common/storage"
	"github.com/baothaihcmut/Bibox/storage-app/internal/modules/files/interactors"
	"github.com/baothaihcmut/Bibox/storage-app/internal/modules/files/services"

	filePermissionRepo "github.com/baothaihcmut/Bibox/storage-app/internal/modules/file_permission/repositories"
	permissionService "github.com/baothaihcmut/Bibox/storage-app/internal/modules/file_permission/services"

	fileRepo "github.com/baothaihcmut/Bibox/storage-app/internal/modules/files/repositories"

	tagRepo "github.com/baothaihcmut/Bibox/storage-app/internal/modules/tags/repositories"
	userRepo "github.com/baothaihcmut/Bibox/storage-app/internal/modules/users/repositories"
)

type FileInteractorImpl struct {
	userRepo             userRepo.UserRepository
	fileRepo             fileRepo.FileRepository
	tagRepo              tagRepo.TagRepository
	logger               logger.Logger
	storageService       storage.StorageService
	mongoService         mongo.MongoService
	filePermission       permissionService.PermissionService
	filePermissionRepo   filePermissionRepo.FilePermissionRepository
	fileStructureService services.FileStructureService
}

func NewFileInteractor(
	userRepo userRepo.UserRepository,
	tagRepo tagRepo.TagRepository,
	fileRepo fileRepo.FileRepository,
	filePermission permissionService.PermissionService,
	filePermissionRepo filePermissionRepo.FilePermissionRepository,
	fileStructureService services.FileStructureService,
	logger logger.Logger,
	storageService storage.StorageService,
	mongoService mongo.MongoService,

) interactors.FileInteractor {
	return &FileInteractorImpl{
		userRepo:             userRepo,
		fileRepo:             fileRepo,
		tagRepo:              tagRepo,
		logger:               logger,
		storageService:       storageService,
		mongoService:         mongoService,
		filePermission:       filePermission,
		filePermissionRepo:   filePermissionRepo,
		fileStructureService: fileStructureService,
	}
}
