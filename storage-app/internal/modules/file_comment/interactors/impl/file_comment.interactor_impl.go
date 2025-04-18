package impl

import (
	"github.com/baothaihcmut/BiBox/libs/pkg/logger"
	"github.com/baothaihcmut/Bibox/storage-app/internal/common/mongo"
	commentRepo "github.com/baothaihcmut/Bibox/storage-app/internal/modules/file_comment/repositories"
	"github.com/baothaihcmut/Bibox/storage-app/internal/modules/file_permission/services"
	fileRepo "github.com/baothaihcmut/Bibox/storage-app/internal/modules/files/repositories"
	userRepo "github.com/baothaihcmut/Bibox/storage-app/internal/modules/users/repositories"
)

type FileCommentInteractorImpl struct {
	fileRepo          fileRepo.FileRepository
	commentRepo       commentRepo.FileCommentRepository
	permissionService services.PermissionService
	userRepo          userRepo.UserRepository
	mongoService      mongo.MongoService
	logger            logger.Logger
}

func NewFileCommentInteractor(
	fileCommentRepo commentRepo.FileCommentRepository,
	fileRepo fileRepo.FileRepository,
	userRepo userRepo.UserRepository,
	filePermissionService services.PermissionService,
	mongoService mongo.MongoService,
	logger logger.Logger,
) *FileCommentInteractorImpl {
	return &FileCommentInteractorImpl{
		fileRepo:          fileRepo,
		commentRepo:       fileCommentRepo,
		permissionService: filePermissionService,
		mongoService:      mongoService,
		logger:            logger,
		userRepo:          userRepo,
	}
}
