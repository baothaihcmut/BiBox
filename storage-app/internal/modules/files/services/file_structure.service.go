package services

import (
	"container/list"
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
	TraverseUploadFolder(_ context.Context, folder *presenters.UploadFolderInput, ownerID primitive.ObjectID, storageProvider string, storageBucket string) ([]*FileWithPath, int, error)
}

type FileStructureServiceImpl struct {
}

func (f *FileStructureServiceImpl) BuildFileStructureTree(_ context.Context, files []*models.File) (*models.FileStructure, error) {
	targetStructure := &models.FileStructure{
		Root: []*models.FileNode{},
	}
	mapNode := make(map[primitive.ObjectID]*models.FileNode)
	for _, file := range files {
		mapNode[file.ID] = &models.FileNode{
			Value:    file,
			SubFiles: []*models.File{},
		}
	}
	for _, file := range files {
		targetNode := mapNode[file.ID]
		if file.ParentFolderID != nil {
			targetNode.SubFiles = append(mapNode[*file.ParentFolderID].SubFiles, file)
		} else {
			targetStructure.Root = append(targetStructure.Root, targetNode)
		}
	}
	return targetStructure, nil
}

func (f *FileStructureServiceImpl) ExcludeFile(_ context.Context, fileStructure *models.FileStructure, excludeFileIds []primitive.ObjectID) ([]*models.File, error) {
	//using BFS
	q := list.New()
	for _, node := range fileStructure.Root {
		q.PushBack(node)
	}
	excludeFileIdsSet := make(map[primitive.ObjectID]struct{})
	for _, fileID := range excludeFileIds {
		excludeFileIdsSet[fileID] = struct{}{}
	}
	res := make([]*models.File, 0)
	for q.Len() > 0 {
		targetNode := q.Remove(q.Front()).(*models.FileNode)
		if _, exist := excludeFileIdsSet[targetNode.Value.ID]; !exist {
			for _, subFile := range targetNode.SubFiles {
				q.PushBack(subFile)
			}
		}
		res = append(res, targetNode.Value)
	}
	return res, nil
}

func (f *FileStructureServiceImpl) TraverseUploadFolder(_ context.Context, folder *presenters.UploadFolderInput, ownerID primitive.ObjectID, storageProvider string, storageBucket string) ([]*FileWithPath, int, error) {
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
			folder.Data.Password,
			folder.Data.IsFolder,
			folder.Data.HasPassword,
			folder.Data.IsSecure,
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
	return files, totalSize, nil
}
func NewFileStructureService() FileStructureService {
	return &FileStructureServiceImpl{}
}
