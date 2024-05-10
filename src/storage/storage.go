package storage

import (
	"context"
	"errors"
	"fmt"
	"log"
	"sync"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Command struct {
	Id           int    `json:"id"`
	Command_uuid string `json:"command_uuid"`
	Description  string `json:"description"`
	Script       string `json:"script"`
	Status       string `json:"status"`
	Output       string `json:"output"`
}

type Storage interface {
	CreateCommand(ctx context.Context, command_uuid, description, script, status, output string) error
	GetCommand(ctx context.Context, command_uuid string) (Command, error)
}

type postgres struct {
	db *pgxpool.Pool
}

func NewPgStorage(ctx context.Context, connString string) (*postgres, error) {
	var pgInstance *postgres
	var pgOnce sync.Once
	pgOnce.Do(func() {
		db, err := pgxpool.New(ctx, connString)
		if err != nil {
			log.Printf("unable to create connection pool: %v\n", err)
			return
		}

		pgInstance = &postgres{db}
	})

	return pgInstance, nil
}

func (pg *postgres) Ping(ctx context.Context) error {
	return pg.db.Ping(ctx)
}

func (pg *postgres) Close() {
	pg.db.Close()
}

func (pg *postgres) CreateCommand(ctx context.Context, command_uuid, description, script, status, output string) error {
	query := `INSERT INTO command (command_uuid, description, script, status, output)
	VALUES (@command_uuid, @description, @script, @status, @output)`
	args := pgx.NamedArgs{
		"command_uuid": command_uuid,
		"description":  description,
		"script":       script,
		"status":       status,
		"output":       output,
	}

	_, err := pg.db.Exec(ctx, query, args)
	if err != nil {
		return fmt.Errorf("unable to insert row: %w", err)
	}

	return nil
}

func (pg *postgres) GetCommand(ctx context.Context, command_uuid string) (Command, error) {
	query := fmt.Sprintf(`SELECT * FROM command WHERE command_uuid = '%s'`, command_uuid)

	rows, err := pg.db.Query(ctx, query)

	var command Command

	if err != nil {
		log.Printf("unable to query: %v", err)
		return command, err
	}
	defer rows.Close()

	command, err = pgx.CollectOneRow(rows, pgx.RowToStructByName[Command])

	if errors.Is(err, pgx.ErrNoRows) {
		return command, errors.New("banner not found")
	}

	if err != nil {
		log.Printf("CollectRows error: %v", err)
		return command, err
	}

	return command, nil
}
