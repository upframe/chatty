package boltdb

import (
	"fmt"

	"github.com/boltdb/bolt"
	"github.com/upframe/chatty"
)

// TeamService allows the management of the team object.
type TeamService struct{}

// Loads loads all the teams from the database.
func (s *TeamService) Loads() ([]*chatty.Team, error) {
	var teams []*chatty.Team
	err := db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("teams"))
		return b.ForEach(func(k, v []byte) error {
			var tm chatty.Team
			err := jsonDecode(v, &tm)
			if err != nil {
				return err
			}
			teams = append(teams, &tm)
			return nil
		})
	})

	return teams, err
}

// Load loads a team from the database.
func (s *TeamService) Load(id string) (*chatty.Team, error) {
	var team *chatty.Team
	err := loadFromDB(team, "teams", id)
	if err != nil {
		return team, err
	}
	if team.ID != id {
		return team, fmt.Errorf("no team with ID '%s'", id)
	}
	return team, nil
}

// Save saves a team to the database.
func (s *TeamService) Save(t *chatty.Team) error {
	return saveToDB("teams", t.ID, t)
}

// Delete removes a team from the database.
func (s *TeamService) Delete(t *chatty.Team) error {
	return db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("teams"))
		return b.Delete([]byte(t.ID))
	})
}
