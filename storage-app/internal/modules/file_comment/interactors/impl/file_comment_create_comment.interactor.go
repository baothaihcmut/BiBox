package impl

import (
	"context"
	"sync"

	"github.com/baothaihcmut/Bibox/storage-app/internal/common/enums"
	"github.com/baothaihcmut/Bibox/storage-app/internal/common/exception"
	"github.com/baothaihcmut/Bibox/storage-app/internal/common/response"
	"github.com/baothaihcmut/Bibox/storage-app/internal/common/utils"
	"github.com/baothaihcmut/Bibox/storage-app/internal/modules/file_comment/models"
	"github.com/baothaihcmut/Bibox/storage-app/internal/modules/file_comment/presenters"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (f *FileCommentInteractorImpl) CreateFileComment(ctx context.Context, input *presenters.CreateFileCommentInput) (*presenters.CreateFileCommentOutput, error) {
	userContext := utils.GetUserContext(ctx)
	userId, err := primitive.ObjectIDFromHex(userContext.Id)
	if err != nil {
		return nil, exception.ErrInvalidObjectId
	}
	fileId, err := primitive.ObjectIDFromHex(input.FileId)
	if err != nil {
		return nil, exception.ErrInvalidObjectId
	}
	wg := &sync.WaitGroup{}
	errCh := make(chan error, 1)
	mentionsCh := make(chan []primitive.ObjectID, 1)
	wg.Add(1)
	go func() {
		defer wg.Done()
		file, err := f.fileRepo.FindFileById(ctx, fileId)
		if err != nil {
			errCh <- err
			return
		}
		if file == nil {
			errCh <- exception.ErrFileNotFound
			return
		}
	}()
	wg.Add(1)
	go func() {
		defer wg.Done()
		hasPermission, err := f.permissionService.CheckPermission(ctx, fileId, userId, enums.CommentPermission)
		if err != nil {
			errCh <- err
			return
		}
		if !hasPermission {
			errCh <- exception.ErrPermissionDenied
			return
		}
	}()
	wg.Add(1)
	go func() {
		defer wg.Done()
		mentions, err := f.extractMention(ctx, input.Content)
		if err != nil {
			errCh <- err
			return
		}
		mentionsCh <- mentions
	}()
	doneCh := make(chan struct{}, 1)
	go func() {
		wg.Wait()
		doneCh <- struct{}{}
	}()
	select {
	case err := <-errCh:
		return nil, err
	case <-doneCh:
	}
	comment := models.NewFileComment(
		fileId,
		userId,
		input.Content,
		<-mentionsCh,
	)
	session, err := f.mongoService.BeginTransaction(ctx)
	if err != nil {
		f.logger.Errorf(ctx, nil, "Error init traction: ", err)
		return nil, err
	}
	defer f.mongoService.EndTransansaction(ctx, session)
	if err := f.commentRepo.CreateComment(ctx, comment); err != nil {
		return nil, err
	}
	if err := f.mongoService.CommitTransaction(ctx, session); err != nil {
		f.logger.Errorf(ctx, nil, "Error commit transaction: ", err)
		return nil, err
	}
	return &presenters.CreateFileCommentOutput{
		FileCommentOutput: response.MapToFileCommentOutput(*comment),
	}, nil

}
