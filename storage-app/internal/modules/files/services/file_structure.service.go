package services

import (
	"container/list"
	"context"

	"github.com/baothaihcmut/Bibox/storage-app/internal/modules/files/models"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type FileStructureService interface{}

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
