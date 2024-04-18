package services

import (
	"guthub.com/Toront0/lux-server/internal/types"
	
	"github.com/jackc/pgx/v5/pgxpool"
	
	"context"
)

type CommunityStorer interface {
	GetCommunities(search string) ([]*types.CommunityPreview, error)
}

type communityStore struct {
	conn *pgxpool.Pool
}

func NewCommunityStore(conn *pgxpool.Pool) *communityStore {
	
	return &communityStore{
		conn: conn,
	}
}

func (s *communityStore) GetCommunities(search string) ([]*types.CommunityPreview, error) {
	cs := []*types.CommunityPreview{}

	arg := "%" + search + "%"

	rows, err := s.conn.Query(context.Background(), `select id, title, category, profile_img from communities where title ilike $1`, arg)

	if err != nil {
		return cs, err
	}

	for rows.Next() {
		c := &types.CommunityPreview{}

		rows.Scan(&c.ID, &c.Title, &c.Category, &c.ProfileImg)

		cs = append(cs, c)
	}

	return cs, nil
}
