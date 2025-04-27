package interactors

import (
	"context"

	"github.com/baothaihcmut/Bibox/storage-app/internal/modules/files/presenters"
)

type FileInteractor interface {
	CreatFile(context.Context, *presenters.CreateFileInput) (*presenters.CreateFileOutput, error)
	UploadFolder(context.Context, *presenters.UploadFolderInput) (*presenters.UploadFolderOutput, error)
	UploadedFile(context.Context, *presenters.UploadedFileInput) (*presenters.UploadedFileOutput, error)
	GetAllFileOfUser(context.Context, *presenters.GetAllFileOfUserInput) (*presenters.GetAllFileOfUserOuput, error)
	AddFilePermission(context.Context, *presenters.AddFilePermissionInput) (*presenters.AddFilePermissionOutput, error)
	GetFileMetaData(context.Context, *presenters.GetFileMetaDataInput) (*presenters.GetFileMetaDataOuput, error)
	GetFileTags(context.Context, *presenters.GetFileTagsInput) (*presenters.GetFileTagsOutput, error)
	GetFilePermissions(context.Context, *presenters.GetFilePermissionInput) (*presenters.GetFilePermissionOuput, error)
	GetFileDownloadUrl(context.Context, *presenters.GetFileDownloadUrlInput) (*presenters.GetFileDownloadUrlOutput, error)
	GetAllSubFileOfFolder(context.Context, *presenters.GetSubFileOfFolderInput) (*presenters.GetSubFileOfFolderOutput, error)
	GetSubFileMetaData(context.Context, *presenters.GetSubFileMetaDataInput) (*presenters.GetSubFileMetaDataOutput, error)
	GetFilePermissionOfUser(context.Context, *presenters.GetFilePermissionOfUserInput) (*presenters.GetFilePermissionOfUserOutput, error)
	UpdateFilePermission(context.Context, *presenters.UpdateFilePermissionInput) (*presenters.UpdateFilePermissionOuput, error)
	DeleteFilePermission(context.Context, *presenters.DeleteFilePermissionOfUserInput) (*presenters.DeleteFilePermissionOfUserOutput, error)
	SoftDeleteFile(context.Context, *presenters.SoftDeleteFileInput) (*presenters.SoftDeleteFileOuput, error)
	RecoverFile(context.Context, *presenters.RecoverFileInput) (*presenters.RecoverFileOutput, error)
	HardDeleteFile(context.Context, *presenters.HardDeleteFileInput) (*presenters.HardDeleteFileOutput, error)
	UpdateFileContent(ctx context.Context, input *presenters.UpdateFileContentInput) (*presenters.UpdateFileContentOutput, error)
}
