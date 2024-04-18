package services

import (
	"guthub.com/Toront0/lux-server/internal/types"
	
	"github.com/jackc/pgx/v5/pgxpool"
	
	"context"
)

type VideoStorer interface {
	GetVideos(page int) ([]*types.VideoPreview, error)
	GetVideo(ID int) (*types.Video, error)
	LikeVideo(videoID, userID int) error
	DeleteLike(videoID, userID int) error

	GetComments(videoID, page, requesterID int) ([]*types.Comment, error)
	InsertComment(videoID, userID int, content string) error
	LikeComment(commentID, userID int) error
	DeleteLikeComment(commentID, userID int) error

	GetCommentReplies(commentID, limit, page, requsterID int) ([]*types.CommentReply, error)
	InsertCommentReply(userID, commentID, videoID int, content string) error
	LikeCommentReply(commentID, userID int) error
	DeleteCommentLikeReply(commentID, userID int) error
}

type videoStore struct {
	conn *pgxpool.Pool
}

func NewVideoStore(conn *pgxpool.Pool) *videoStore {

	return &videoStore{
		conn: conn,
	}

}

func (s *videoStore) GetVideos(page int) ([]*types.VideoPreview, error) {
	vs := []*types.VideoPreview{}

	rows, err := s.conn.Query(context.Background(), `select t1.id, t1.created_at, t1.title, t1.thumbnail, t1.author_id, t1.url, t1.viewsamount, t2.first_name, t2.last_name, t2.profile_img from videos t1 join users t2 on t1.author_id = t2.id limit 50 offset $1`, page * 50)

	if err != nil {
		return vs, err
	}

	for rows.Next() {
		v := &types.VideoPreview{}

		rows.Scan(&v.ID, &v.CreatedAt, &v.Title, &v.Thumbnail, &v.AuthorID, &v.Url, &v.ViewsAmount, &v.AuthorFName, &v.AuthorLName, &v.AuthorPImg)

		vs = append(vs, v)
	}

	return vs, nil
}

func (s *videoStore) GetVideo(ID int) (*types.Video, error) {
	v := &types.Video{}

	err := s.conn.QueryRow(context.Background(), `select t1.id, t1.created_at, t1.title, t1.author_id, t1.url, t1.description, t1.viewsamount, t2.first_name, t2.last_name, t2.profile_img from videos t1 join users t2 on t1.author_id = t2.id where t1.id = $1`, ID).Scan(&v.ID, &v.CreatedAt, &v.Title, &v.AuthorID, &v.Url, &v.Description, &v.ViewsAmount, &v.AuthorFName, &v.AuthorLName, &v.AuthorPImg)

	return v, err
}

func (s *videoStore) LikeVideo(videoID, userID int) error {

	_, err := s.conn.Exec(context.Background(), `insert into video_likes (video_id, user_id) values($1, $2)`, videoID, userID)

	return err
}

func (s *videoStore) DeleteLike(videoID, userID int) error {

	_, err := s.conn.Exec(context.Background(), `delete from video_likes where video_id = $1 and user_id = $2`, videoID, userID)

	return err
}

func (s *videoStore) GetComments(videoID, page, requesterID int) ([]*types.Comment, error) {
	res := []*types.Comment{}

	rows, err := s.conn.Query(context.Background(), `select t1.id, t1.created_at, t1.content, t1.user_id, t3.first_name, t3.last_name, t3.profile_img, (select count(*) from video_comment_likes where comment_id = t1.id), (select count(*) from video_comment_replies where comment_id = t1.id), (select count(*) from video_comment_likes where comment_id = t1.id and user_id = $3) from video_comments t1 join users t3 on t3.id = t1.user_id where t1.video_id = $1 order by created_at ASC limit 20 offset $2`, videoID, page * 20, requesterID)
	
	// rows, err := s.conn.Query(context.Background(), `select t1.id, t1.created_at, t1.content, t1.user_id, t3.first_name, t3.last_name, t3.profile_img, (select count(*) from post_comment_likes where comment_id = t1.id), (select json_agg(json_build_object('id', t2.id, 'createdAt', t2.created_at, 'content', t2.content, 'authorId', t2.user_id, 'authorFName', t3.first_name, 'authorLName', t3.last_name, 'authorPImg', t3.profile_img, 'likesAmount', (select count(*) from post_comment_reply_likes where comment_id = t2.id))) from post_comment_replies t2 join users t3 on t2.user_id = t3.id where t2.comment_id = t1.id) from post_comments t1 join users t3 on t3.id = t1.user_id where t1.post_id = $1 order by created_at ASC limit 20 offset $2 `, postID, page * 20)



	if err != nil {
		return res, err
	}

	for rows.Next() {
		r := &types.Comment{}

		rows.Scan(&r.ID, &r.CreatedAt, &r.Content, &r.AuthorID, &r.AuthorFName, &r.AuthorLName, &r.AuthorPImg, &r.LikesAmount, &r.RepliesAmount, &r.IsRequesterLiked)

		res = append(res, r)
	}


	return res, nil
}

func (s *videoStore) InsertComment(videoID, userID int, content string) error {

	_, err := s.conn.Exec(context.Background(), `insert into video_comments (video_id, user_id, content) values ($1, $2, $3)`, videoID, userID, content)

	return err
}

func (s *videoStore) LikeComment(commentID, userID int) error {

	_, err := s.conn.Exec(context.Background(), `insert into video_comment_likes (comment_id, user_id) values ($1, $2)`, commentID, userID)

	return err
}

func (s *videoStore) DeleteLikeComment(commentID, userID int) error {


	_, err := s.conn.Exec(context.Background(), `delete from video_comment_likes where comment_id = $1 and user_id = $2`, commentID, userID)

	return err
}

func (s *videoStore) GetCommentReplies(commentID, limit, page, requesterID int) ([]*types.CommentReply, error) {
	res := []*types.CommentReply{}

	rows, err := s.conn.Query(context.Background(), `select t1.id, t1.created_at, t1.content, t1.user_id, t3.first_name, t3.last_name, t3.profile_img, (select count(*) from video_comment_reply_likes where comment_id = t1.id), (select count(*) from video_comment_reply_likes where comment_id = t1.id and user_id = $4) from video_comment_replies t1 join users t3 on t3.id = t1.user_id where t1.comment_id = $1 order by id ASC limit $2 offset $3`, commentID, limit, page, requesterID)
	
	

	if err != nil {
		return res, err
	}

	for rows.Next() {
		r := &types.CommentReply{}

		rows.Scan(&r.ID, &r.CreatedAt, &r.Content, &r.AuthorID, &r.AuthorFName, &r.AuthorLName, &r.AuthorPImg, &r.LikesAmount, &r.IsRequesterLiked)

		res = append(res, r)
	}


	return res, nil


}

func (s *videoStore) InsertCommentReply(userID, commentID, videoID int, content string) error {


	_, err := s.conn.Exec(context.Background(), `insert into video_comment_replies (user_id, comment_id, content, video_id) values ($1, $2, $3, $4)`, userID, commentID, content, videoID)

	return err
}

func (s *videoStore) LikeCommentReply(commentID, userID int) error {

	_, err := s.conn.Exec(context.Background(), `insert into video_comment_reply_likes (comment_id, user_id) values ($1, $2)`, commentID, userID)
	
	return err
}

func (s *videoStore) DeleteCommentLikeReply(commentID, userID int) error {

	_, err := s.conn.Exec(context.Background(), `delete from video_comment_reply_likes where comment_id = $1 and user_id = $2`, commentID, userID)

	return err
}