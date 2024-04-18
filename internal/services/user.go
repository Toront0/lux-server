package services

import (
	"guthub.com/Toront0/lux-server/internal/types"
	
	"github.com/jackc/pgx/v5/pgxpool"
	
	"context"
	"fmt"
)

type UserStorer interface {
	GetUserDetail(userID, requesterID int) (*types.User, error)
	GetUserPosts(userID, page int) ([]*types.Post, error)
	GetUserVideos(userID, page int) ([]*types.VideoPreview, error)
	GetUserFriends(userID, page int, search string) ([]*types.UserFriend, error)

	GetUserMusic(userID, page int) ([]*types.Song, error)
	GetUserPlaylists(userID int) ([]*types.MusicPlaylistPreview, error)

	GetUserFollowers(userID, page int, search string) ([]*types.UserPreview, error)
	GetUserFollowings(userID, page int, search string) ([]*types.UserPreview, error)

	GetDialogs(receiverID int) ([]*types.UserDialog, error)
	GetDialogMessages(receiverID int, senderID int) ([]*types.UserMessage ,error)
	InsertMessage(senderID, receiverID int, message string) error

	AddFollower(followerID, followeeID int) error
	DeleteFollow(followerID, followeeID int) error

	AddFriend(userFId, userSId int) error
	DeleteFriendship(userFId, userSId int) error

	GetSettingsData(userID int) (*types.UserSettingsData, error)
	UpdateUser(query string) error
}

type userStore struct {
	conn *pgxpool.Pool
}

func NewUserStore(conn *pgxpool.Pool) *userStore {

	return &userStore{
		conn: conn,
	}
}

func (s *userStore) GetUserDetail(userID, requesterID int) (*types.User, error) {
	res := &types.User{}

	err := s.conn.QueryRow(context.Background(), `select id, created_at, first_name, last_name, profile_img, banner_img, status,
	(select count(*) from followers where followee_id = $1),
	(select count(*) from followers where follower_id = $1),
	(
		case when (
			select count(*) from friends where user_f_id = $1 and user_s_id = $2 or user_f_id = $2 and user_s_id = $1
		) > 0 then 'friends'
		when (
			select id from followers where follower_id = $1 and followee_id = $2
		) > 0 then 'follower'
		when (
			select id from followers where follower_id = $2 and followee_id = $1
		) > 0 then 'followee'
		else ''
		end
	)
	from users where id = $1`, userID, requesterID).Scan(&res.ID, &res.CreatedAt, &res.FirstName, &res.LastName, &res.ProfileImg, &res.BannerImg, &res.Status, &res.FollowersAmount, &res.FolloweesAmount, &res.UsersRelations)



	return res, err
}

func (s *userStore) GetUserPosts(userID, page int) ([]*types.Post, error) {
	ps := []*types.Post{}

	rows, err := s.conn.Query(context.Background(), `select t1.id, t1.created_at, t1.content, t2.id, concat(t2.first_name, ' ', t2.last_name), t2.profile_img, (select count(*) from post_likes where post_id = t1.id), (select count(*) from post_likes where post_id = t1.id and user_id = $1), (select count(*) from post_comments where post_id = t1.id), (select json_agg(url) from post_media where post_id = t1.id) from posts t1 join users t2 on t1.user_id = t2.id where t1.user_id = $1 limit 20 offset $2`, userID, page * 20)
	// rows, err := s.conn.Query(context.Background(), `select t1.id, t1.created_at, t1.content, t2.id, concat(t2.first_name, ' ', t2.last_name), t2.profile_img, (select count(*) from post_likes where post_id = t1.id), (select id from post_likes where user_id = $1 and post_id = t1.id), (select json_agg(img_url) from post_images ) from posts t1 join users t2 on t1.user_id = t2.id where t1.user_id = $1`, userID)

	if err != nil {
		return ps, err
	}

	fmt.Println("$@!")

	for rows.Next() {
		p := &types.Post{}

		rows.Scan(&p.ID, &p.CreatedAt, &p.Content, &p.AuthorID, &p.AuthorName, &p.AuthorPImg, &p.LikesAmount, &p.IsRequesterLiked, &p.LikesAmount, &p.PostMedia)

		ps = append(ps, p)
	}

	return ps, nil
}

func (s *userStore) GetUserVideos(userID, page int) ([]*types.VideoPreview, error) {
	vs := []*types.VideoPreview{}

	rows, err := s.conn.Query(context.Background(), `select t1.id, t1.created_at, t1.title, t1.thumbnail, t1.author_id, t1.url, t1.viewsamount, t2.first_name, t2.last_name, t2.profile_img from videos t1 join users t2 on t1.author_id = t2.id where t2.id = $1 limit 18 offset $2`, userID, page * 18)

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

func (s *userStore) GetUserFriends(userID, page int, search string) ([]*types.UserFriend, error) {
	fs := []*types.UserFriend{}

	searchQuery := "%" + search + "%"

	rows, err := s.conn.Query(context.Background(), `select t1.id, t2.id, t2.first_name, t2.last_name, t2.profile_img from friends t1 join users t2 on t2.id = case when t1.user_f_id = $1 then t1.user_s_id else t1.user_f_id end where t1.user_f_id = $1 and t2.first_name || t2.last_name ilike $3 or t1.user_s_id = $1 and t2.first_name || t2.last_name ilike $3 limit 30 offset $2`, userID, page * 30, searchQuery)

	if err != nil {	
		return fs, err
	}

	for rows.Next() {
		f := &types.UserFriend{}

		rows.Scan(&f.ID, &f.FriendID, &f.FirstName, &f.LastName, &f.ProfileImg)

		fs = append(fs, f)
	}

	return fs, nil
}

func (s *userStore) GetUserMusic(userID, page int) ([]*types.Song, error) {
	res := []*types.Song{}


	rows, err := s.conn.Query(context.Background(), `select t1.id, t1.title, t1.cover, t1.performer, t1.url, (select count(*) from user_music where user_id = $1 and song_id = t1.id) from music t1 join user_music t2 on t2.song_id = t1.id where t2.user_id = $1 and playlist_id is null limit 30 offset $2`, userID, page * 30)

	if err != nil {
		return res, err
	}

	

	for rows.Next() {
		m := &types.Song{}

		rows.Scan(&m.ID, &m.Title, &m.Cover, &m.Performer, &m.Url, &m.IsInMyList)

		res = append(res, m)

	}

	return res, nil


}

func (s *userStore) GetUserPlaylists(userID int) ([]*types.MusicPlaylistPreview, error) {
	res := []*types.MusicPlaylistPreview{}


	rows, err := s.conn.Query(context.Background(), `select id, title, cover_img, user_id from user_music_playlists where user_id = $1`, userID)

	if err != nil {
		return res, err
	}

	for rows.Next() {
		p := &types.MusicPlaylistPreview{}


		rows.Scan(&p.ID, &p.Title, &p.CoverImg, &p.CreatorID)

		res = append(res, p)
	}

	return res, nil

}


func (s *userStore) GetUserFollowers(userID, page int, search string) ([]*types.UserPreview, error) {
	res := []*types.UserPreview{}

	searchQuery := "%" + search + "%"

	rows, err := s.conn.Query(context.Background(), `select t1.id, t2.id, t2.first_name, t2.last_name, t2.profile_img from followers t1 join users t2 on t2.id = t1.follower_id where t1.followee_id = $1 and t2.first_name || t2.last_name ilike $3 limit 30 offset $2`, userID, page * 30, searchQuery)

	if err != nil {	
		return res, err
	}

	for rows.Next() {
		f := &types.UserPreview{}

		rows.Scan(&f.ID, &f.UserID, &f.FirstName, &f.LastName, &f.ProfileImg)

		res = append(res, f)
	}

	return res, nil
}

func (s *userStore) GetUserFollowings(userID, page int, search string) ([]*types.UserPreview, error) {
	res := []*types.UserPreview{}

	searchQuery := "%" + search + "%"

	rows, err := s.conn.Query(context.Background(), `select t1.id, t2.id, t2.first_name, t2.last_name, t2.profile_img from followers t1 join users t2 on t2.id = t1.followee_id where t1.follower_id = $1 and t2.first_name || t2.last_name ilike $3 limit 30 offset $2`, userID, page * 30, searchQuery)

	if err != nil {	
		return res, err
	}

	for rows.Next() {
		f := &types.UserPreview{}

		rows.Scan(&f.ID, &f.UserID, &f.FirstName, &f.LastName, &f.ProfileImg)

		res = append(res, f)
	}

	return res, nil
}

func (s *userStore) GetDialogs(receiverID int) ([]*types.UserDialog, error) {
	ms := []*types.UserDialog{}


	rows, err := s.conn.Query(context.Background(), `
	SELECT
 	 subquery.sender_id as sender_id,
     users2.first_name, 
     users2.last_name, 
     users2.profile_img, 
     subquery.message as message,
     subquery.id as id,
	 subquery.created_at as created_at,
	 subquery.receiver_id as receiver_id
 FROM
     users
     JOIN
     (
     SELECT
         message,
         row_number() OVER ( PARTITION BY  sender_id + receiver_id ORDER BY id DESC) AS row_num,
         receiver_id,
         sender_id,
         id,
		 created_at
     FROM
         personal_messages
     GROUP BY
         id,
         sender_id,
         receiver_id,
         message,
		 created_at
     ) AS subquery ON ( ( subquery.sender_id = users.id OR subquery.receiver_id = users.id)  AND row_num = 1 )
	 JOIN users as users2 ON ( users2.id = CASE WHEN users.id = subquery.sender_id THEN subquery.receiver_id ELSE subquery.sender_id END )
 WHERE
     users.id = $1
 ORDER BY
     subquery.id DESC  
	
	`, receiverID)

	if err != nil {
		return ms, err
	}

	for rows.Next() {
		m := &types.UserDialog{}

		rows.Scan(&m.SenderID, &m.ContactFName, &m.ContactLName, &m.ContactPImg, &m.Message, &m.ID, &m.CreatedAt, &m.ReceiverID)

		ms = append(ms, m)
	}


	return ms, nil
}

func (s *userStore) GetDialogMessages(receiverID int, senderID int) ([]*types.UserMessage, error) {
	ms := []*types.UserMessage{}

	rows, err := s.conn.Query(context.Background(), `select t1.*, t2.profile_img from personal_messages t1 join users t2 on t1.sender_id = t2.id where receiver_id = $1 and sender_id = $2 or sender_id = $1 and receiver_id = $2 order by created_at ASC`, receiverID, senderID)
	// rows, err := s.conn.Query(context.Background(), `select * from personal_messages where receiver_id = $1 and sender_id = $2 or sender_id = $1 and receiver_id = $2 order by created_at desc limit 2`, receiverID, senderID)

	if err != nil {
		return ms, err
	}

	for rows.Next() {
		m := &types.UserMessage{}

		rows.Scan(&m.ID, &m.CreatedAt, &m.ReceiverID, &m.SenderID, &m.Message, &m.ProfileImg)

		ms = append(ms, m)
	}

	return ms, nil



}

func (s *userStore) InsertMessage(senderID, receiverID int, message string) error {


	_, err := s.conn.Exec(context.Background(), `insert into personal_messages (sender_id, receiver_id, message) values ($1, $2, $3)`, senderID, receiverID, message)

	return err

}

func (s *userStore) AddFollower(followerID, followeeID int) error {

	_, err := s.conn.Exec(context.Background(), `insert into followers (follower_id, followee_id) values ($1, $2)`, followerID, followeeID)

	return err
}

func (s *userStore) DeleteFollow(followerID, followeeID int) error {

	_, err := s.conn.Exec(context.Background(), `delete from followers where follower_id = $1 and followee_id = $2`, followerID, followeeID)

	return err
}

func (s *userStore) AddFriend(userFId, userSId int) error {

	_, err := s.conn.Exec(context.Background(), `insert into friends (user_f_id, user_s_id) values ($1, $2)`, userFId, userSId)

	return err

}

func (s *userStore) DeleteFriendship(userFId, userSId int) error {

	_, err := s.conn.Exec(context.Background(), `delete from friends where user_f_id = $1 and user_s_id = $2 or user_f_id = $2 and user_s_id = $1`, userFId, userSId)

	return err

}

func (s *userStore) GetSettingsData(userID int) (*types.UserSettingsData, error) {
	res := &types.UserSettingsData{}

	err := s.conn.QueryRow(context.Background(), `select id, first_name, last_name, status, profile_img, banner_img, email  from users where id = $1`, userID).Scan(&res.ID, &res.FirstName, &res.LastName, &res.Status, &res.ProfileImg, &res.BannerImg, &res.Email)

	return res, err



}

func (s *userStore) UpdateUser(query string) error {

	_, err := s.conn.Exec(context.Background(), query)

	return err

}