package interactors

import (
	"context"

	"github.com/baothaihcmut/Bibox/storage-app/internal/modules/files/presenters"
)

type FileInteractor interface {
	CreatFile(context.Context, *presenters.CreateFileInput) (*presenters.CreateFileOutput, error)
	UploadFolder(context.Context, *presenters.UploadFolderInput) (*presenters.UploadFolderOutput, error)
	UploadedFile(context.Context, *presenters.UploadedFileInput) (*presenters.UploadedFileOutput, error)
	FindAllFileOfUser(context.Context, *presenters.FindFileOfUserInput) (*presenters.FindFileOfUserOuput, error)
	AddFilePermission(context.Context, *presenters.AddFilePermissionInput) (*presenters.AddFilePermissionOutput, error)
	GetFileMetaData(context.Context, *presenters.GetFileMetaDataInput) (*presenters.GetFileMetaDataOuput, error)
	GetFileTags(context.Context, *presenters.GetFileTagsInput) (*presenters.GetFileTagsOutput, error)
	GetFilePermissions(context.Context, *presenters.GetFilePermissionInput) (*presenters.GetFilePermissionOuput, error)
	GetFileDownloadUrl(context.Context, *presenters.GetFileDownloadUrlInput) (*presenters.GetFileDownloadUrlOutput, error)
	GetAllSubFileOfFolder(context.Context, *presenters.GetSubFileOfFolderInput) (*presenters.GetSubFileOfFolderOutput, error)
	GetSubFileMetaData(context.Context, *presenters.GetSubFileMetaDataInput) (*presenters.GetSubFileMetaDataOutput, error)
}
