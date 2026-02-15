package stored

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"strings"

	"go.opentelemetry.io/otel"

	internalDB "alielgamal.com/myservice/internal/db"
)

var tracer = otel.Tracer("myservice/store")

const (
	addTemplateStmt      = "INSERT INTO %v(id, content, created_by, modified_by) VALUES ($1, $2, $3, $3) RETURNING created_at, modified_at"
	getTemplateStmt      = "SELECT content, created_by, created_at, modified_by, modified_at FROM %v WHERE id=$1"
	patchTemplateStmt    = "UPDATE %v SET modified_by=$1, modified_at=CURRENT_TIMESTAMP, %v  WHERE id=$2 RETURNING content, created_by, created_at, modified_by, modified_at"
	listTemplateStmt     = "SELECT id, content, created_by, created_at, modified_by, modified_at FROM %v %v"
	patchArgStartIndex   = 3
	setAttributeTemplate = "content['%v'] = $%v"
	conditionTemplate    = "content['%v'] %v $%v"
)

// Store An interface that provides Storage facility for any object that can be represents in JSON format
type Store[T any] interface {
	// Add a new Stored item with a specific id and content
	Add(ctx context.Context, creator string, id string, content T) (*Stored[T], error)

	// Updates a single attribute in the content. This method doesn't check the attribute existence but guarantees
	//  that content stored is still a valid. If it is not, the patch operation will fail without impacting storage.
	Patch(ctx context.Context, updater string, id string, attributes map[string]any) (*Stored[T], error)

	// Get finds a storable by its id
	Get(ctx context.Context, id string) (*Stored[T], error)

	// List returns all items that fill certain all conditions (AND operator between the conditions).
	// if no conditions are passed, all stored items are returned.
	List(ctx context.Context, conds ...Condition) ([]Stored[T], error)
}

// NewStore creates a new store for a specific Stored T in a specific table
func NewStore[T any](db internalDB.DB, table string) Store[T] {
	return sqlStore[T]{
		db:      db,
		table:   table,
		addStmt: fmt.Sprintf(addTemplateStmt, table),
		getStmt: fmt.Sprintf(getTemplateStmt, table),
	}
}

type sqlStore[T any] struct {
	db      internalDB.DB
	table   string
	addStmt string
	getStmt string
}

func (s sqlStore[T]) Add(ctx context.Context, creator string, id string, content T) (*Stored[T], error) {
	ctx, span := tracer.Start(ctx, "store.add")
	defer span.End()

	contentJSON, err := json.Marshal(content)
	if err != nil {
		return nil, err
	}

	row := s.db.QueryRowContext(ctx, s.addStmt, id, contentJSON, creator)

	result := &Stored[T]{
		ID:         id,
		Content:    content,
		CreatedBy:  creator,
		ModifiedBy: creator,
	}

	err = row.Scan(&result.CreatedAt, &result.ModifiedAt)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (s sqlStore[T]) Patch(ctx context.Context, updater string, id string, attributes map[string]any) (*Stored[T], error) {
	ctx, span := tracer.Start(ctx, "store.patch")
	defer span.End()

	attributeStmts := []string{}
	nextValueIndex := patchArgStartIndex
	queryParams := []any{updater, id}
	for k, v := range attributes {
		jsonValue, err := json.Marshal(v)
		if err != nil {
			return nil, err
		}
		attributeStmts = append(attributeStmts, fmt.Sprintf(setAttributeTemplate, k, nextValueIndex))
		queryParams = append(queryParams, jsonValue)
		nextValueIndex++
	}

	patchStmt := fmt.Sprintf(patchTemplateStmt, s.table, strings.Join(attributeStmts, ","))

	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	row := tx.QueryRowContext(ctx, patchStmt, queryParams...)

	result := &Stored[T]{
		ID: id,
	}

	err = s.scanStored(result, row)
	if err != nil {
		return nil, err
	}

	err = tx.Commit()
	if err != nil {
		return nil, err
	}
	return result, err
}

func (s sqlStore[T]) Get(ctx context.Context, id string) (*Stored[T], error) {
	ctx, span := tracer.Start(ctx, "store.get")
	defer span.End()

	row := s.db.QueryRowContext(ctx, s.getStmt, id)

	result := &Stored[T]{
		ID: id,
	}

	if err := s.scanStored(result, row); err != nil {
		return nil, err
	}

	return result, nil
}
func (s sqlStore[T]) List(ctx context.Context, conds ...Condition) ([]Stored[T], error) {
	ctx, span := tracer.Start(ctx, "store.list")
	defer span.End()

	conditionStmts := []string{}
	nextValueIndex := 1
	queryParams := []any{}
	for _, c := range conds {
		conditionStmts = append(conditionStmts, fmt.Sprintf(conditionTemplate, c.Attribute, c.Op, nextValueIndex))
		nextValueIndex++
		queryValue, err := json.Marshal(c.Value)
		if err != nil {
			return nil, err
		}
		queryParams = append(queryParams, queryValue)
	}

	listStmtCondition := ""
	if len(conds) > 0 {
		listStmtCondition = fmt.Sprintf("WHERE %v", strings.Join(conditionStmts, " AND "))
	}
	listStmt := fmt.Sprintf(listTemplateStmt, s.table, listStmtCondition)
	var rows *sql.Rows
	var err error
	rows, err = s.db.QueryContext(ctx, listStmt, queryParams...)

	if err != nil {
		return nil, err
	}

	result := []Stored[T]{}
	for rows.Next() {
		r := Stored[T]{}
		var contentJSON []byte
		err := rows.Scan(&r.ID, &contentJSON, &r.CreatedBy, &r.CreatedAt, &r.ModifiedBy, &r.ModifiedAt)
		if err != nil {
			return nil, err
		}
		err = json.Unmarshal(contentJSON, &r.Content)
		if err != nil {
			return nil, err
		}
		result = append(result, r)
	}

	return result, nil
}

func (s sqlStore[T]) scanStored(result *Stored[T], row *sql.Row) error {
	var contentJSON []byte
	err := row.Scan(&contentJSON, &result.CreatedBy, &result.CreatedAt, &result.ModifiedBy, &result.ModifiedAt)
	if err != nil {
		return err
	}
	err = json.Unmarshal(contentJSON, &result.Content)
	if err != nil {
		return err
	}
	return nil
}
