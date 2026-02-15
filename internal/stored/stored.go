// Package stored provides the ability to store arbitrary content as JSON in a JSON-capable SQL store (like Postgres)
package stored

import "time"

// Stored struct represents a stored item that has a specific type of content.
// All stored items use a string to identify them. These IDs may or may not
// have business implications; for example just a random UUID if there is no
// need for specifying a unique identifier for an object.
// All Stored Objects are required to have a consistent schema the follows:
//
// CREATE TABLE <stored_name> (
// id VARCHAR(36) NOT NULL PRIMARY KEY CHECK(length(name) > 0),
// content JSONB NOT NULL,
// created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
// modified_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
// created_by VARCHAR(50) NOT NULL CHECK(length(created_by) > 0),
// modified_by VARCHAR(50) NOT NULL CHECK(length(modified_by) > 0)'
// );
//
// Additionally, you must add an index on the content to ensure that querying for
// content is efficient:
// CREATE INDEX <stored_named>_content_idx ON app USING GIN(content jsonb_path_ops);
//
// You can potentially add constraints and unique indexes on the content if needed.
// The Content Struct the fields of the type of the content must be exported and have
// JSON tags associated with them.
// The table that you use can be partitioned if you wish and you can use the content
// JSON values attributes to control how your data is partitioned across.
// It is best to avoid nested structures and keep your content struct flat.
type Stored[T any] struct {

	// ID The unique ID of the storable object
	ID string `json:"id" uri:"id" binding:"required,min=1,max=36"`

	// CreatedAt the time at which the stored item was created
	CreatedAt time.Time `json:"createdAt,omitempty" binding:"isdefault"`

	// ModifiedAt the time at which the stored item was modified
	ModifiedAt time.Time `json:"modifiedAt,omitempty" binding:"isdefault"`

	// CreatedBy the identification of the user who created this stored item
	CreatedBy string `json:"createdBy,omitempty" binding:"isdefault"`

	// ModifiedBy the identification of the user who last modified this stored item
	ModifiedBy string `json:"modifiedBy,omitempty" binding:"isdefault"`

	// The content of the storable
	Content T `json:"content"`
}

// Operator Defines an operator that is used for creating conditions
type Operator string

// The possible operators to use with Conditions
const (
	EqualOperator             Operator = "="
	NotEqualOperator          Operator = "<>"
	GreaterThanOperator       Operator = ">"
	LessThanOperator          Operator = "<"
	GreaterThanOrEqualOpertor Operator = ">="
	LessThanOrEqualOperator   Operator = "<="
)

// Condition models a condition on an attribute. These conditions can then be passed to Store operations like List to control the result returns.
type Condition struct {
	Attribute string
	Op        Operator
	Value     any
}
