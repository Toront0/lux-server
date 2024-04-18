package services

import (
	"guthub.com/Toront0/lux-server/internal/types"
	
	"github.com/jackc/pgx/v5/pgxpool"
	
	"context"
	// "reflect"
	// "fmt"
)

type PostStorer interface {
	GetPosts(page, requesterID int) ([]*types.Post, error)
	GetPost(postID, requesterID int) (*types.Post, error)
	LikePost(postID, userID int) error
	DeleteLike(postID, userID int) error
	CreatePost(content string, userID int, imgs []string) error

	GetComments(postID, page, requesterID int) ([]*types.Comment, error)
	InsertComment(postID, userID int, content string) error
	LikeComment(commentID, userID int) error
	DeleteLikeComment(commentID, userID int) error

	GetCommentReplies(commentID, limit, page, requsterID int) ([]*types.CommentReply, error)
	InsertCommentReply(userID, commentID, postID int, content string) error
	LikeCommentReply(commentID, userID int) error
	DeleteCommentLikeReply(commentID, userID int) error
}

type postStore struct {
	conn *pgxpool.Pool
}

func NewPostStore(conn *pgxpool.Pool) *postStore {
	return &postStore{
		conn: conn,
	}
}

func (s *postStore) GetPosts(page, requesterID int) ([]*types.Post, error) {
	ps := []*types.Post{}
	
	rows, err := s.conn.Query(context.Background(), `select t1.id, t1.created_at, t1.content, t2.id, concat(t2.first_name, ' ', t2.last_name), t2.profile_img, (select count(*) from post_likes where post_id = t1.id), (select count(*) from post_likes where post_id = t1.id and user_id = $1), (select count(*) from post_comments where post_id = t1.id), (select json_agg(url) from post_media where post_id = t1.id) from posts t1 join users t2 on t1.user_id = t2.id limit 20 offset $2`, requesterID, page * 20)

	if err != nil {
		return ps, err
	}


	for rows.Next() {
		p := &types.Post{}

		rows.Scan(&p.ID, &p.CreatedAt, &p.Content, &p.AuthorID, &p.AuthorName, &p.AuthorPImg, &p.LikesAmount, &p.IsRequesterLiked, &p.CommentsAmount, &p.PostMedia)

		ps = append(ps, p)
	}

	return ps, nil
}

func (s *postStore) GetPost(postID, requesterID int) (*types.Post, error) {
	res := &types.Post{}

	err := s.conn.QueryRow(context.Background(), `select t1.id, t1.created_at, t1.content, t2.id, concat(t2.first_name, ' ', t2.last_name), t2.profile_img, (select count(*) from post_likes where post_id = t1.id), (select count(*) from post_likes where post_id = t1.id and user_id = $1), (SELECT (select count(*) from post_comments where post_id = t1.id) + (select count(*) from post_comment_replies where post_id = t1.id)), (select json_agg(url) from post_media where post_id = t1.id) from posts t1 join users t2 on t1.user_id = t2.id where t1.id = $2`, requesterID, postID).Scan(&res.ID, &res.CreatedAt, &res.Content, &res.AuthorID, &res.AuthorName, &res.AuthorPImg, &res.LikesAmount, &res.IsRequesterLiked, &res.CommentsAmount, &res.PostMedia)


	// err := s.conn.QueryRow(context.Background(), `select t1.id, t1.created_at, t1.content, t2.id, t2.first_name, t2.last_name, t2.profile_img, (select count(*) from post_likes where post_id = t1.id), (select count(*) from post_likes where post_id = t1.id and user_id = $1), (select count(*) from post_comments where post_id = t1.id) from posts t1 join users t2 on t1.author_id = t2.id where t1.id = $2`, 7, postID).Scan(&res.ID, &res.CreatedAt, &res.Content, &res.AuthorID, &res.AuthorFName, &res.AuthorLName, &res.AuthorPImg, &res.LikesAmount, &res.IsRequesterLiked, &res.CommentsAmount)


	return res, err
}

func (s *postStore) LikePost(postID, userID int) error {

	_, err := s.conn.Exec(context.Background(), `insert into post_likes (post_id, user_id) values($1, $2)`, postID, userID)

	return err
}

func (s *postStore) DeleteLike(postID, userID int) error {

	_, err := s.conn.Exec(context.Background(), `delete from post_likes where post_id = $1 and user_id = $2`, postID, userID)

	return err
}

func (s *postStore) CreatePost(content string, userID int, imgs []string) error {
	var postId int

	err := s.conn.QueryRow(context.Background(), `insert into posts (content, user_id) values ($1, $2) returning id`, content, userID).Scan(&postId)

	if err != nil {
		return err
	}

	for _, img := range imgs {

		_, err := s.conn.Exec(context.Background(), `insert into post_images (post_id, img_url) values ($1, $2)`, postId, img)

		if err != nil {
			return err
		}
	}

	return err

}

func (s *postStore) GetComments(postID, page, requesterID int) ([]*types.Comment, error) {
	res := []*types.Comment{}

	rows, err := s.conn.Query(context.Background(), `select t1.id, t1.created_at, t1.content, t1.user_id, t3.first_name, t3.last_name, t3.profile_img, (select count(*) from post_comment_likes where comment_id = t1.id), (select count(*) from post_comment_replies where comment_id = t1.id), (select count(*) from post_comment_likes where comment_id = t1.id and user_id = $3) from post_comments t1 join users t3 on t3.id = t1.user_id where t1.post_id = $1 order by created_at ASC limit 20 offset $2`, postID, page * 20, requesterID)
	
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

func (s *postStore) InsertComment(postID, userID int, content string) error {

	_, err := s.conn.Exec(context.Background(), `insert into post_comments (post_id, user_id, content) values ($1, $2, $3)`, postID, userID, content)

	return err
}

func (s *postStore) LikeComment(commentID, userID int) error {

	_, err := s.conn.Exec(context.Background(), `insert into post_comment_likes (comment_id, user_id) values ($1, $2)`, commentID, userID)

	return err
}

func (s *postStore) DeleteLikeComment(commentID, userID int) error {


	_, err := s.conn.Exec(context.Background(), `delete from post_comment_likes where comment_id = $1 and user_id = $2`, commentID, userID)

	return err
}

func (s *postStore) GetCommentReplies(commentID, limit, page, requesterID int) ([]*types.CommentReply, error) {
	res := []*types.CommentReply{}

	rows, err := s.conn.Query(context.Background(), `select t1.id, t1.created_at, t1.content, t1.user_id, t3.first_name, t3.last_name, t3.profile_img, (select count(*) from post_comment_reply_likes where comment_id = t1.id), (select count(*) from post_comment_reply_likes where comment_id = t1.id and user_id = $4) from post_comment_replies t1 join users t3 on t3.id = t1.user_id where t1.comment_id = $1 order by created_at ASC limit $2 offset $3`, commentID, limit, page, requesterID)
	
	

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

func (s *postStore) InsertCommentReply(userID, commentID, postID int, content string) error {


	_, err := s.conn.Exec(context.Background(), `insert into post_comment_replies (user_id, comment_id, content, post_id) values ($1, $2, $3, $4)`, userID, commentID, content, postID)

	return err
}

func (s *postStore) LikeCommentReply(commentID, userID int) error {

	_, err := s.conn.Exec(context.Background(), `insert into post_comment_reply_likes (comment_id, user_id) values ($1, $2)`, commentID, userID)
	
	return err
}

func (s *postStore) DeleteCommentLikeReply(commentID, userID int) error {

	_, err := s.conn.Exec(context.Background(), `delete from post_comment_reply_likes where comment_id = $1 and user_id = $2`, commentID, userID)

	return err
}

