package impl

import (
	"github.com/baothaihcmut/BiBox/libs/pkg/logger"
	"github.com/baothaihcmut/Bibox/storage-app/internal/common/mongo"
	"github.com/baothaihcmut/Bibox/storage-app/internal/common/storage"
	"github.com/baothaihcmut/Bibox/storage-app/internal/modules/files/interactors"
	"github.com/baothaihcmut/Bibox/storage-app/internal/modules/files/services"
	structureService "github.com/baothaihcmut/Bibox/storage-app/internal/modules/files/services"
	notificationService "github.com/baothaihcmut/Bibox/storage-app/internal/modules/notification/services"

	filePermissionRepo "github.com/baothaihcmut/Bibox/storage-app/internal/modules/file_permission/repositories"
	permissionService "github.com/baothaihcmut/Bibox/storage-app/internal/modules/file_permission/services"

	fileRepo "github.com/baothaihcmut/Bibox/storage-app/internal/modules/files/repositories"

	tagRepo "github.com/baothaihcmut/Bibox/storage-app/internal/modules/tags/repositories"
	userRepo "github.com/baothaihcmut/Bibox/storage-app/internal/modules/users/repositories"
)

type FileInteractorImpl struct {
	userRepo                  userRepo.UserRepository
	fileRepo                  fileRepo.FileRepository
	tagRepo                   tagRepo.TagRepository
	logger                    logger.Logger
	storageService            storage.StorageService
	mongoService              mongo.MongoService
	filePermission            permissionService.PermissionService
	filePermissionRepo        filePermissionRepo.FilePermissionRepository
	fileStructureService      structureService.FileStructureService
	notificationService       notificationService.NotificationService
	fileUploadProgressService services.FileUploadProgressService
}

func NewFileInteractor(
	userRepo userRepo.UserRepository,
	tagRepo tagRepo.TagRepository,
	fileRepo fileRepo.FileRepository,
	filePermission permissionService.PermissionService,
	filePermissionRepo filePermissionRepo.FilePermissionRepository,
	fileStructureService structureService.FileStructureService,
	notificationService notificationService.NotificationService,
	fileUploadProgressService services.FileUploadProgressService,
	logger logger.Logger,
	storageService storage.StorageService,
	mongoService mongo.MongoService,

) interactors.FileInteractor {
	return &FileInteractorImpl{
		userRepo:                  userRepo,
		fileRepo:                  fileRepo,
		tagRepo:                   tagRepo,
		logger:                    logger,
		storageService:            storageService,
		mongoService:              mongoService,
		filePermission:            filePermission,
		filePermissionRepo:        filePermissionRepo,
		fileStructureService:      fileStructureService,
		notificationService:       notificationService,
		fileUploadProgressService: fileUploadProgressService,
	}
}
