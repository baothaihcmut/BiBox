package services

import (
	"context"
	"fmt"

	"github.com/baothaihcmut/Bibox/storage-app/internal/common/enums"
	permissionModel "github.com/baothaihcmut/Bibox/storage-app/internal/modules/file_permission/models"
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
	ExtractPermissionFromFileStructure(ctx context.Context, userId primitive.ObjectID, fileStructure *models.FileStructure, addtionPermssions []presenters.AdditionFilePermission) []*permissionModel.FilePermission
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

func (f *FileStructureServiceImpl) ExtractPermissionFromFileStructure(ctx context.Context, userId primitive.ObjectID, fileStructure *models.FileStructure, addtionPermssions []presenters.AdditionFilePermission) []*permissionModel.FilePermission {
	//map for addition permission
	mapAdditionPermission := make(map[primitive.ObjectID]enums.PermissionType)
	for _, additionPermission := range addtionPermssions {
		mapAdditionPermission[additionPermission.FileId] = enums.PermissionType(additionPermission.PermissionType)
	}
	var traverseFileStructure func(*models.FileNode, enums.FilePermissionType) []*permissionModel.FilePermission
	traverseFileStructure = func(root *models.FileNode, permissionType enums.FilePermissionType) []*permissionModel.FilePermission {
		if overwritePermission, exist := mapAdditionPermission[root.Value.ID]; exist {
			permissionType = enums.FilePermissionType(overwritePermission)
		}
		permissions := []*permissionModel.FilePermission{permissionModel.NewFilePermission(
			root.Value.ID,
			userId,
			permissionType,
			true,
			nil,
		),
		}
		for _, subFile := range root.SubFiles {
			subPermissions := traverseFileStructure(subFile, permissionType)
			permissions = append(permissions, subPermissions...)
		}
		return permissions
	}
	permissions := make([]*permissionModel.FilePermission, 0)
	for _, root := range fileStructure.Root {
		permissions = append(permissions, traverseFileStructure(root, enums.ViewPermission)...)
	}
	return permissions
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
