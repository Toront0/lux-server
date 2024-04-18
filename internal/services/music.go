package services

import (
	"guthub.com/Toront0/lux-server/internal/types"
	
	"github.com/jackc/pgx/v5/pgxpool"
	
	"context"
)

type MusicStorer interface {
	GetSongs(userID, page int) ([]*types.Song, error)
	GetPlaylists() ([]*types.MusicPlaylistPreview, error)
	GetPlaylistDetail(playlistID int) (*types.MusicPlaylist, error)
	GetPlaylistSongs(playlistID, userID, page int) ([]*types.Song, error)
	CreatePlaylist(query string) (int, error)

	UpdateSongsPlaylistID(query string) error
	DeletePlaylist(playlistID, userID int) error
	GetAvailableAndPlaylistSongs(playlistID, userID, page int) ([]*types.PlaylistSong, error)

	AddSongToUser(songID, userID int) error
	DeleteUserSong(songID, userID int) error
}

type musicStore struct {

	conn *pgxpool.Pool

}

func NewMusicStore(conn *pgxpool.Pool) *musicStore {
	return &musicStore{
		conn: conn,
	}
}

func (s *musicStore) GetSongs(userID, page int) ([]*types.Song, error) {
	ms := []*types.Song{}

	rows, err := s.conn.Query(context.Background(), `select t1.id, t1.title, t1.performer, t1.cover, t1.url, (select count(*) from user_music where user_id = $1 and song_id = t1.id) from music t1 limit 50 offset $2`, userID, page * 50)

	if err != nil {
		return ms, err
	}

	for rows.Next() {
		m := &types.Song{}
		
		rows.Scan(&m.ID, &m.Title, &m.Performer, &m.Cover, &m.Url, &m.IsInMyList)

		ms = append(ms, m)
	}

	return ms, nil
}

func (s *musicStore) GetPlaylists() ([]*types.MusicPlaylistPreview, error) {
	res := []*types.MusicPlaylistPreview{}

	rows, err := s.conn.Query(context.Background(), `select id, title, cover_img, user_id from user_music_playlists where is_private = false order by random() limit 5`)

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

func (s *musicStore) GetPlaylistDetail(playlistID int) (*types.MusicPlaylist, error) {
	res := &types.MusicPlaylist{}

	err := s.conn.QueryRow(context.Background(), `select t1.id, t1.title, t1.cover_img, t2.id, t2.first_name, t2.last_name, t2.profile_img from user_music_playlists t1 join users t2 on t1.user_id = t2.id where t1.id = $1`, playlistID).Scan(&res.ID, &res.Title, &res.CoverImg, &res.CreatorID, &res.CreatorFName, &res.CreatorLName, &res.CreatorPImg)

	return res, err
}

func (s *musicStore) GetPlaylistSongs(playlistID, userID, page int) ([]*types.Song, error) {
	ms := []*types.Song{}

	rows, err := s.conn.Query(context.Background(), `select t1.id, t1.title, t1.performer, t1.cover, t1.url, (select count(*) from user_music where user_id = $1 and song_id = t1.id) from music t1 join user_music t2 on t1.id = t2.song_id where t2.playlist_id = $1 and t2.user_id = $2 limit 30 offset $3`, playlistID, userID,  page * 30)

	if err != nil {
		return ms, err
	}

	for rows.Next() {
		m := &types.Song{}
		
		rows.Scan(&m.ID, &m.Title, &m.Performer, &m.Cover, &m.Url, &m.IsInMyList)

		ms = append(ms, m)
	}

	return ms, nil


}

func (s *musicStore) CreatePlaylist(query string) (int, error) {

	var id int

	err := s.conn.QueryRow(context.Background(), query).Scan(&id)

	return id, err
}

func (s *musicStore) UpdateSongsPlaylistID(query string) error {

	_, err := s.conn.Query(context.Background(), query)

	return err

}

func (s *musicStore) DeletePlaylist(playlistID, userID int) error {

	_, err := s.conn.Exec(context.Background(), `update user_music set playlist_id = null where playlist_id = $1 and user_id = $2`, playlistID, userID)

	if err != nil {
		return err
	}

	_, err = s.conn.Exec(context.Background(), `delete from user_music_playlists where id = $1`, playlistID)

	return err
}

func (s *musicStore) GetAvailableAndPlaylistSongs(playlistID, userID, page int)  ([]*types.PlaylistSong, error) {
	res := []*types.PlaylistSong{}

	rows, err := s.conn.Query(context.Background(), `select t1.id, t1.title, t1.performer, t1.cover, t1.url, t2.playlist_id from music t1 join user_music t2 on t1.id = t2.song_id where t2.playlist_id = $1 and t2.user_id = $2 or t2.playlist_id is null and t2.user_id = $2 limit 30 offset $3`, playlistID, userID, page * 30)

	if err != nil {
		return res, err
	}

	for rows.Next() {
		s := &types.PlaylistSong{}

		rows.Scan(&s.ID, &s.Title, &s.Performer, &s.Cover, &s.Url, &s.PlaylistID)

		res = append(res, s)

	}

	return res, nil

}

func (s *musicStore) AddSongToUser(songID, userID int) error {

	_, err := s.conn.Exec(context.Background(), `insert into user_music (song_id, user_id) values ($1, $2)`, songID, userID)

	return err
}

func (s *musicStore) DeleteUserSong(songID, userID int) error {

	_, err := s.conn.Exec(context.Background(), `delete from user_music where song_id = $1 and user_id = $2`, songID, userID)

	return err

}