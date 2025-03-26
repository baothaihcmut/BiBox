package impl

import (
	"github.com/baothaihcmut/Bibox/storage-app/internal/common/logger"
	"github.com/baothaihcmut/Bibox/storage-app/internal/common/mongo"
	"github.com/baothaihcmut/Bibox/storage-app/internal/modules/files/repositories"
	tagRepo "github.com/baothaihcmut/Bibox/storage-app/internal/modules/tags/repositories"
)

type TagInteractorImpl struct {
	repo         tagRepo.TagRepository
	fileRepo     repositories.FileRepository
	logger       logger.Logger
	mongoService mongo.MongoService
}

func NewTagInteractor(repo tagRepo.TagRepository, fileRepo repositories.FileRepository, logger logger.Logger, mongoService mongo.MongoService) *TagInteractorImpl {
	return &TagInteractorImpl{
		repo:         repo,
		logger:       logger,
		mongoService: mongoService,
		fileRepo:     fileRepo,
	}
}
