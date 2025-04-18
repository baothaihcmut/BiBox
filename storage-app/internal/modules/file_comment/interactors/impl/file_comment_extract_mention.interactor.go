package impl

import (
	"context"
	"regexp"

	"github.com/baothaihcmut/Bibox/storage-app/internal/common/exception"
	"github.com/baothaihcmut/Bibox/storage-app/internal/modules/users/models"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var mentionRegex = regexp.MustCompile(`@(\w+)`)

func (f *FileCommentInteractorImpl) extractMention(ctx context.Context, content string) ([]primitive.ObjectID, error) {
	matches := mentionRegex.FindAllStringSubmatch(content, -1)
	var emails []string
	for _, match := range matches {
		emails = append(emails, match[1])
	}
	if len(emails) == 0 {
		return []primitive.ObjectID{}, nil
	}
	users, err := f.userRepo.FindUsersByEmailList(ctx, emails)
	if err != nil {
		return nil, err
	}
	mapUser := make(map[string]*models.User)
	for _, user := range users {
		mapUser[user.Email] = user
	}
	mentions := make([]primitive.ObjectID, 0, len(emails))
	for _, email := range emails {
		user, exist := mapUser[email]
		if !exist {
			return nil, exception.ErrUserNotFound
		}
		mentions = append(mentions, user.ID)
	}
	return mentions, nil
}
