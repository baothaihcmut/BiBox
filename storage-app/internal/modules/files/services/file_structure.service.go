package services

import (
	"context"
	"fmt"

	"github.com/baothaihcmut/Bibox/storage-app/internal/common/enums"
	"github.com/baothaihcmut/Bibox/storage-app/internal/modules/files/models"

	"github.com/baothaihcmut/Bibox/storage-app/internal/modules/files/presenters"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type FileWithPath struct {
	*models.File
	Path string
}

type FileStructureService interface {
	TraverseUploadFolder(_ context.Context, folder *presenters.UploadFolderInput, ownerID primitive.ObjectID, storageProvider string, storageBucket string) ([]*FileWithPath, int)
	BuildFileStructureTree(_ context.Context, files []*models.File) *models.FileStructure
}

type FileStructureServiceImpl struct {
}

func (f *FileStructureServiceImpl) BuildFileStructureTree(_ context.Context, files []*models.File) *models.FileStructure {
	targetStructure := &models.FileStructure{
		Root: []*models.FileNode{},
	}
	mapNode := make(map[primitive.ObjectID]*models.FileNode)
	for _, file := range files {
		mapNode[file.ID] = &models.FileNode{
			Value:    file,
			SubFiles: []*models.FileNode{},
		}
	}
	for _, file := range files {
		targetNode := mapNode[file.ID]
		if file.ParentFolderID != nil {
			if parentNode, exist := mapNode[*file.ParentFolderID]; exist {
				parentNode.SubFiles = append(parentNode.SubFiles, targetNode)
				continue
			}
		}
		targetStructure.Root = append(targetStructure.Root, targetNode)
	}
	return targetStructure
}

func (f *FileStructureServiceImpl) TraverseUploadFolder(_ context.Context, folder *presenters.UploadFolderInput, ownerID primitive.ObjectID, storageProvider string, storageBucket string) ([]*FileWithPath, int) {
	var traverseFunc func(*presenters.UploadFolderInput, *primitive.ObjectID, string) ([]*FileWithPath, int)
	traverseFunc = func(folder *presenters.UploadFolderInput, parentFolderID *primitive.ObjectID, path string) ([]*FileWithPath, int) {
		var storageDetail *models.FileStorageDetailArg
		if !folder.Data.IsFolder && folder.Data.StorageDetail != nil {
			storageDetail = &models.FileStorageDetailArg{
				Size:            folder.Data.StorageDetail.Size,
				MimeType:        enums.MapToMimeType(folder.Data.Name, folder.Data.StorageDetail.MimeType),
				StorageProvider: storageProvider,
				StorageBucket:   storageBucket,
			}
		}
		rootFile := models.NewFile(
			ownerID,
			folder.Data.Name,
			parentFolderID,
			folder.Data.Description,
			folder.Data.IsFolder,
			folder.Data.TagIDs,
			storageDetail,
		)
		rootNode := &FileWithPath{
			File: rootFile,
			Path: fmt.Sprintf("%s/%s", path, folder.Data.Name),
		}
		res := make([]*FileWithPath, 0)
		res = append(res, rootNode)
		if !folder.Data.IsFolder && len(folder.SubFiles) == 0 {
			return res, folder.Data.StorageDetail.Size
		}
		totalSize := 0
		for _, subFile := range folder.SubFiles {
			subNodes, subSize := traverseFunc(subFile, &rootFile.ID, rootNode.Path)
			res = append(res, subNodes...)
			totalSize += subSize
		}
		return res, totalSize
	}
	files, totalSize := traverseFunc(folder, folder.Data.ParentFolderID, "")
	return files, totalSize
}
func NewFileStructureService() FileStructureService {
	return &FileStructureServiceImpl{}
}
