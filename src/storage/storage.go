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

type Storage interface {
	CreateCommand(ctx context.Context, command_uuid, description, script, status, output string) error
	GetCommand(ctx context.Context, command_uuid string) (Command, error)
	GetCommands(ctx context.Context, status string, limit int, offset int) ([]Command, error)
	UpdateCommandStatus(ctx context.Context, command_uuid string, status string) error
	UpdateCommandOutput(ctx context.Context, command_uuid string, output string) error
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

	var command Command

	rows, err := pg.db.Query(ctx, query)
	if err != nil {
		log.Printf("unable to query: %v", err)
		return command, err
	}
	defer rows.Close()

	command, err = pgx.CollectOneRow(rows, pgx.RowToStructByName[Command])

	if errors.Is(err, pgx.ErrNoRows) {
		return command, errors.New("command not found")
	}

	if err != nil {
		log.Printf("CollectRows error: %v", err)
		return command, err
	}

	return command, nil
}

func (pg *postgres) GetCommands(ctx context.Context, status string, limit int, offset int) ([]Command, error) {

	condition := ""
	if status != "" {
		condition += fmt.Sprintf(`WHERE status = '%s'`, status)
	}

	query := fmt.Sprintf(`SELECT * FROM command %s ORDER BY id DESC `, condition)

	if limit > 0 {
		query += fmt.Sprintf(`LIMIT %d `, limit)
	}
	if offset > 0 {
		query += fmt.Sprintf(`OFFSET %d `, offset)
	}
	query += `;`

	var commands []Command

	rows, err := pg.db.Query(ctx, query)
	if err != nil {
		return commands, fmt.Errorf("unable to query: %w", err)
	}
	defer rows.Close()

	commands, err = pgx.CollectRows(rows, pgx.RowToStructByName[Command])
	if err != nil {
		log.Printf("CollectRows error: %v", err)
		return commands, err
	}

	return commands, nil
}

func (pg *postgres) UpdateCommandStatus(ctx context.Context, command_uuid string, status string) error {
	query := fmt.Sprintf(`UPDATE command SET status = '%s' WHERE command_uuid = '%s'`, status, command_uuid)

	_, err := pg.db.Exec(ctx, query)
	if err != nil {
		return fmt.Errorf("unable to insert row: %w", err)
	}

	return nil
}

func (pg *postgres) UpdateCommandOutput(ctx context.Context, command_uuid string, output string) error {
	query := fmt.Sprintf(`UPDATE command SET output = '%s' WHERE command_uuid = '%s'`, output, command_uuid)

	_, err := pg.db.Exec(ctx, query)
	if err != nil {
		return fmt.Errorf("unable to insert row: %w", err)
	}

	return nil
}
