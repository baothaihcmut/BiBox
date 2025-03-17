package impl

import (
	"context"

	"github.com/baothaihcmut/Bibox/storage-app/internal/common/constant"
	"github.com/baothaihcmut/Bibox/storage-app/internal/common/enums"
	commonModel "github.com/baothaihcmut/Bibox/storage-app/internal/common/models"
	"github.com/baothaihcmut/Bibox/storage-app/internal/modules/files/presenters"
	tagPresenter "github.com/baothaihcmut/Bibox/storage-app/internal/modules/tags/presenters"

	"github.com/baothaihcmut/Bibox/storage-app/internal/modules/tags/models"
	"github.com/samber/lo"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// GetFileTags implements FileInteractor.
func (f *FileInteractorImpl) GetFileTags(ctx context.Context, input *presenters.GetFileTagsInput) (*presenters.GetFileTagsOutput, error) {
	//check tag exist and check user have permission with file
	userCtx := ctx.Value(constant.UserContext).(*commonModel.UserContext)
	userId, _ := primitive.ObjectIDFromHex(userCtx.Id)
	fileId, _ := primitive.ObjectIDFromHex(input.Id)

	//check permission
	file, err := f.checkFilePermission(ctx, fileId, userId, enums.ViewPermission)
	if err != nil {
		return nil, err
	}
	tags, err := f.tagRepo.FindAllTagInList(ctx, file.TagIDs)
	if err != nil {
		return nil, err
	}
	//map to ouput

	return &presenters.GetFileTagsOutput{
		Tags: lo.Map(tags, func(item *models.Tag, _ int) *tagPresenter.TagOutput {
			return tagPresenter.MaptoOuput(item)
		}),
	}, nil
}
