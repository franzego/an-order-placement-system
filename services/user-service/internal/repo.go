package internal

import (
	"context"
	"fmt"
	"log"

	db "github.com/franzego/user-service/db/sqlc"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
)

type RepoServicer interface {
	CreateUser(ctx context.Context, username string, email string, passwordhash string,
		firstname pgtype.Text, lastname pgtype.Text) (db.CreateUserRow, error)
	GetUserByEmail(ctx context.Context, email string) (db.GetUserByEmailRow, error)
	DeleteUser(ctx context.Context, id int64) error
	UpdateUser(ctx context.Context, args db.UpdateUserParams) (db.UpdateUserRow, error)
}

type Repo struct {
	dbq db.Queries
	db  *pgxpool.Pool
}

func NewRepoService(dbconn *pgxpool.Pool) *Repo {

	if dbconn == nil {
		return nil
	}
	return &Repo{
		dbq: *db.New(dbconn),
		db:  dbconn,
	}
}

func (r *Repo) CreateUser(ctx context.Context, username string, email string,
	passwordhash string, firstname pgtype.Text, lastname pgtype.Text) (db.CreateUserRow, error) {

	user, err := r.dbq.CreateUser(ctx, db.CreateUserParams{
		Username:     username,
		Email:        email,
		PasswordHash: passwordhash,
		FirstName:    firstname,
		LastName:     lastname,
	})
	if err != nil {
		log.Printf("error in creating user:%v", err)

	}
	return user, nil
}

func (r *Repo) GetUserByEmail(ctx context.Context, email string) (db.GetUserByEmailRow, error) {
	user, err := r.dbq.GetUserByEmail(ctx, email)
	if err != nil {
		return db.GetUserByEmailRow{}, fmt.Errorf("error in getting user from email: %v", err)
	}
	return user, nil
}

//if the user wants to delete their account

func (r *Repo) DeleteUser(ctx context.Context, id int64) error {
	err := r.dbq.DeleteUser(ctx, id)
	if err != nil {
		log.Printf("error in deleting user account, %v", err)
	}
	return nil
}

// update user info maybe.. not so sure for now..will test it
func (r *Repo) UpdateUser(ctx context.Context, args db.UpdateUserParams) (db.UpdateUserRow, error) {
	user, err := r.dbq.UpdateUser(ctx, args)
	if err != nil {
		log.Printf("error in updating user %v", err)
	}
	return user, nil
}
