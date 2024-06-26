package repository

import (
	"cine/entity/model"
	"context"
	"github.com/google/uuid"
)

type Repository[E, F, U any] interface {
	One(ctx context.Context, filters ...F) (E, error)
	All(ctx context.Context, filters ...F) ([]E, error)

	Exists(ctx context.Context, filters ...F) (bool, error)
	Count(ctx context.Context, filters ...F) (int, error)

	Insert(ctx context.Context, entity E) (E, error)
	InsertBulk(ctx context.Context, entities []E) ([]E, error)

	Update(ctx context.Context, id uuid.UUID, updater U) (E, error)
	UpdateExec(ctx context.Context, updater U, filters ...F) (int, error)

	Delete(ctx context.Context, id uuid.UUID) error
	DeleteExec(ctx context.Context, filters ...F) (int, error)
}

type (
	SessionRepository Repository[*model.Session, *model.SessionF, *model.SessionU]
	LikeRepository    Repository[*model.Like, *model.LikeF, *model.LikeU]
	MediaRepository   Repository[*model.Media, *model.MediaF, *model.MediaU]
)

type UserRepository interface {
	Repository[*model.User, *model.UserF, *model.UserU]

	OneDetailed(ctx context.Context, id, userID uuid.UUID) (*model.DetailedUser, error)

	OneFollowed(ctx context.Context, user *model.User, followedID uuid.UUID) (*model.User, error)
	AllFollowed(ctx context.Context, user *model.User) ([]*model.User, error)
	FollowUser(ctx context.Context, user *model.User, userToFollowID uuid.UUID) error
	UnfollowUser(ctx context.Context, user *model.User, followedID uuid.UUID) error

	OneFollower(ctx context.Context, user *model.User, followerID uuid.UUID) (*model.User, error)
	AllFollowers(ctx context.Context, user *model.User) ([]*model.User, error)
}

type ListRepository interface {
	Repository[*model.List, *model.ListF, *model.ListU]

	AllWithMedia(ctx context.Context, filters ...*model.ListF) ([]*model.ListWithMedia, error)
	OneWithMedia(ctx context.Context, filters ...*model.ListF) (*model.ListWithMedia, error)

	AllMembers(ctx context.Context, list *model.List) ([]*model.User, error)
	AddMember(ctx context.Context, list *model.List, userID uuid.UUID) error
	RemoveMember(ctx context.Context, list *model.List, userID uuid.UUID) error

	AddMedia(ctx context.Context, list *model.List, mediaID uuid.UUID) error
	RemoveMedia(ctx context.Context, list *model.List, mediaID uuid.UUID) error
	AllMedia(ctx context.Context, list *model.List) ([]*model.Media, error)
}

type CommentRepository interface {
	Repository[*model.Comment, *model.CommentF, *model.CommentU]

	AllAsDetailed(ctx context.Context, mediaID uuid.UUID, userID uuid.UUID) ([]*model.DetailedComment, error)
	AllRepliesAsDetailed(ctx context.Context, comment *model.Comment, userID uuid.UUID) ([]*model.DetailedComment, error)
}

type ReviewRepository interface {
	Repository[*model.Review, *model.ReviewF, *model.ReviewU]

	AllWithUser(ctx context.Context, reviewFs ...*model.ReviewF) ([]*model.DetailedReview, error)
}
