package storage

import (
	"context"
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/kingxl111/cakes-database-app/internal/models"
)

var _ AdminAuthorization = (*AdminAuthPostgres)(nil)

type AdminAuthPostgres struct {
	db *DB
}

func NewAdminAuthPostgres(db *DB) *AdminAuthPostgres {
	return &AdminAuthPostgres{db: db}
}

func (a *AdminAuthPostgres) GetAdmin(username, passwordHash string) (int, error) {
	const op = "pgsql.GetAdmin"
	var id int

	builderSelect := sq.Select(idColumn).
		From(AdminTable).
		PlaceholderFormat(sq.Dollar).
		Where(sq.Eq{usernameColumn: username}).
		Where(sq.Eq{passwordColumn: passwordHash})

	query, args, err := builderSelect.ToSql()
	if err != nil {
		return id, fmt.Errorf("op: %s, error building select query: %w", op, err)
	}

	err = a.db.pool.QueryRow(context.Background(), query, args...).Scan(&id)
	if err != nil {
		return id, fmt.Errorf("op: %s, error executing select query: %w", op, err)
	}
	return id, nil
}

var _ Admin = (*AdminPostgres)(nil)

type AdminPostgres struct {
	db *DB
}

func NewAdminPostgres(db *DB) *AdminPostgres {
	return &AdminPostgres{db: db}
}

func (a *AdminPostgres) GetUsers() ([]models.User, error) {
	const op = "pgsql.GetUsers"
	users := make([]models.User, 0, 10)

	builderSelect := sq.Select(idColumn, fullnameColumn, usernameColumn, emailColumn, phoneColumn).
		From(UserTable).
		PlaceholderFormat(sq.Dollar)

	query, args, err := builderSelect.ToSql()
	if err != nil {
		return users, fmt.Errorf("op: %s, error building select query: %w", op, err)
	}

	rows, err := a.db.pool.Query(context.Background(), query, args...)
	if err != nil {
		return users, fmt.Errorf("op: %s, error executing select query: %w", op, err)
	}
	for rows.Next() {
		var usr models.User
		err = rows.Scan(&usr.ID, &usr.FullName, &usr.Username, &usr.Email, &usr.PhoneNumber)
		if err != nil {
			return users, fmt.Errorf("op: %s, error scanning row: %w", op, err)
		}
		users = append(users, usr)
	}
	if err := rows.Err(); err != nil {
		return users, fmt.Errorf("op: %s, error scanning rows: %w", op, err)
	}

	return users, nil
}

func (db *DB) backupTable(tableName string) error {
	const op = "pgsql.BackupTable"

	query, args, err := sq.Select("*").
		From(tableName).
		ToSql()
	if err != nil {
		return fmt.Errorf("op: %s, error building query: %w", op, err)
	}

	rows, err := db.pool.Query(context.Background(), query, args...)
	if err != nil {
		return fmt.Errorf("op: %s, error executing query: %w", op, err)
	}
	defer rows.Close()

	fieldDescriptions := rows.FieldDescriptions()
	columns := make([]string, len(fieldDescriptions))
	for i, fd := range fieldDescriptions {
		columns[i] = string(fd.Name) // Convert pgconn.FieldDescription.Name to string
	}

	file, err := os.Create(fmt.Sprintf("%s_backup.csv", tableName))
	if err != nil {
		return fmt.Errorf("op: %s, error opening file: %w", op, err)
	}
	defer func() {
		if err := file.Close(); err != nil {
			log.Printf("error closing file: %v", err) // Handle error
		}
	}()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	if err := writer.Write(columns); err != nil {
		return fmt.Errorf("op: %s, error writing csv: %w", op, err)
	}

	// using interface{} for universal data retrieve
	values := make([]interface{}, len(columns))
	for i := range values {
		values[i] = new(interface{})
	}

	for rows.Next() {
		if err := rows.Scan(values...); err != nil {
			return fmt.Errorf("op: %s, error scanning row: %w", op, err)
		}

		record := make([]string, len(columns))
		for i, value := range values {
			switch v := (*value.(*interface{})).(type) {
			case nil:
				record[i] = ""
			case []byte:
				record[i] = string(v)
			case time.Time:
				record[i] = v.Format(time.RFC3339)
			default:
				record[i] = fmt.Sprintf("%v", v)
			}
		}

		if err := writer.Write(record); err != nil {
			return fmt.Errorf("op: %s, error writing csv: %w", op, err)
		}
	}

	if err := rows.Err(); err != nil {
		return fmt.Errorf("op: %s, error scanning rows: %w", op, err)
	}
	return nil
}

func (a *AdminPostgres) Backup() error {
	const op = "pgsql.Backup"

	ctx := context.Background()
	tx, err := a.db.pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("op: %s, failed to begin transaction: %w", op, err)
	}
	tables := []string{
		UserTable,
		DeliveryTable,
		OrderTable,
		OrdersCakesTable,
		CakesTable,
		DeliveryPointTable,
		AdminTable,
	}

	for _, table := range tables {
		if err := a.db.backupTable(table); err != nil {
			return fmt.Errorf("op: %s, failed to backup table: %w", op, err)
		}
	}

	if err = tx.Commit(ctx); err != nil {
		return fmt.Errorf("op: %s, failed to commit transaction: %w", op, err)
	}

	return nil
}

func (db *DB) restoreTable(tableName string) error {
	const op = "pgsql.RestoreTable"

	file, err := os.Open(fmt.Sprintf("%s_backup.csv", tableName))
	if err != nil {
		return fmt.Errorf("op: %s, error opening backup file: %w", op, err)
	}
	defer func() {
		if err := file.Close(); err != nil {
			log.Printf("error closing file: %v", err)
		}
	}()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		return fmt.Errorf("error reading CSV: %w", err)
	}

	if len(records) == 0 {
		return fmt.Errorf("op: %s, error no records found, tableName: %s", op, tableName)
	}

	columns := records[0]

	placeholders := make([]string, len(columns))
	for i := range placeholders {
		placeholders[i] = fmt.Sprintf("$%d", i+1)
	}

	query := fmt.Sprintf("INSERT INTO %s (%s) VALUES (%s)",
		tableName,
		strings.Join(columns, ", "),
		strings.Join(placeholders, ", "))

	for _, record := range records[1:] {
		if len(record) != len(columns) {
			return fmt.Errorf("op: %s, error row has %d columns, expected %d", op, len(columns), len(record))
		}

		args := make([]interface{}, len(record))
		for i := range record {
			args[i] = record[i]
		}

		_, err := db.pool.Exec(context.Background(), query, args...)
		if err != nil {
			return fmt.Errorf("op: %s, error executing row: %w", op, err)
		}
	}

	return nil
}

func (a *AdminPostgres) Restore() error {
	const op = "pgsql.Restore"
	ctx := context.Background()
	tx, err := a.db.pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("op: %s, failed to begin transaction: %w", op, err)
	}
	tables := []string{
		UserTable,
		AdminTable,
		CakesTable,
		DeliveryPointTable,
		OrderTable,
		DeliveryTable,
		OrdersCakesTable,
	}

	for _, table := range tables {
		if err := a.db.restoreTable(table); err != nil {
			return fmt.Errorf("op: %s, error restoring table: %w", op, err)
		}
	}
	if err = tx.Commit(ctx); err != nil {
		return fmt.Errorf("op: %s, failed to commit transaction: %w", op, err)
	}

	return nil
}

func (a *AdminPostgres) AddCake(ctx context.Context, cake models.Cake) (int, error) {
	const op = "pgsql.AddCake"
	var id int
	builder := sq.Insert(CakesTable).
		PlaceholderFormat(sq.Dollar).
		Columns(descriptionColumn, priceColumn, weightColumn, fullDescriptionColumn, activeColumn).
		Values(cake.Description, cake.Price, cake.Weight, cake.FullDescription, true).
		Suffix("RETURNING id")

	query, args, err := builder.ToSql()
	if err != nil {
		return id, fmt.Errorf("op: %s, error building query: %w", op, err)
	}
	err = a.db.pool.QueryRow(ctx, query, args...).Scan(&id)
	if err != nil {
		return id, fmt.Errorf("op: %s, error executing query: %w", op, err)
	}

	return id, nil
}

func (a *AdminPostgres) RemoveCake(ctx context.Context, id int) error {
	const op = "pgsql.RemoveCake"

	builder := sq.Update(CakesTable).
		Set(activeColumn, false).
		PlaceholderFormat(sq.Dollar).
		Where(sq.Eq{idColumn: id})
	query, args, err := builder.ToSql()
	if err != nil {
		return fmt.Errorf("op: %s, error building query: %w", op, err)
	}

	_, err = a.db.pool.Exec(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("op: %s, error executing query: %w", op, err)
	}

	return nil
}
