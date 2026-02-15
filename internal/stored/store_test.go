package stored

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"alielgamal.com/myservice/internal/config"
	testDB "alielgamal.com/myservice/internal/db/test"
)

func TestStore(t *testing.T) {
	appConfig, _ := config.ReadConfig()

	tableName := "stored"

	type content struct {
		I int    `json:"i"`
		B bool   `json:"b"`
		S string `json:"s"`
	}
	admin := "admin@example.com"

	prepareMockDB := func(t *testing.T) (func(), Store[content]) {
		db, tearDown, err := testDB.SetupTestDB(t.Name(), 0, appConfig)
		require.Nil(t, err)

		_, err = db.ExecContext(context.Background(), fmt.Sprintf(`
		CREATE TABLE %v (
			id VARCHAR(36) NOT NULL PRIMARY KEY CHECK(length(id) > 0),
			content JSONB NOT NULL,
			created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
			modified_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
			created_by VARCHAR(50) NOT NULL CHECK(length(created_by) > 0),
			modified_by VARCHAR(50) NOT NULL CHECK(length(modified_by) > 0)
			)`,
			tableName))
		require.Nil(t, err)

		_, err = db.ExecContext(context.Background(), fmt.Sprintf(
			"CREATE INDEX %v_content_idx ON stored USING GIN(content jsonb_path_ops)",
			tableName))
		require.Nil(t, err)

		s := NewStore[content](db, tableName)

		return tearDown, s
	}

	ctx := context.Background()
	id := "id1"
	fixture := content{
		I: 5,
		B: true,
		S: "Some Text",
	}

	t.Run("Add", func(t *testing.T) {

		t.Run("Successfully inserts and retrieves the stored content after", func(t *testing.T) {
			tearDown, s := prepareMockDB(t)
			defer tearDown()

			added, err := s.Add(ctx, admin, id, fixture)
			require.NoError(t, err)
			assert.NotNil(t, added.CreatedAt)
			assert.Equal(t, added.CreatedAt, added.ModifiedAt)
			assert.Equal(t, admin, added.CreatedBy)
			assert.Equal(t, admin, added.ModifiedBy)
			assert.Equal(t, fixture, added.Content)
			assert.Equal(t, id, added.ID)
			assert.Equal(t, fixture, added.Content)

			fetched, err := s.Get(context.Background(), id)
			assert.NoError(t, err)
			assert.Equal(t, *added, *fetched)
		})

		t.Run("Fails when", func(t *testing.T) {
			tests := []struct {
				name    string
				id      string
				creator string
				value   content
			}{
				{"missing id", "", admin, fixture},
				{"missing creator", "abc", "", fixture},
			}
			for _, tt := range tests {
				t.Run(tt.name, func(t *testing.T) {
					tearDown, s := prepareMockDB(t)
					defer tearDown()
					added, err := s.Add(ctx, tt.creator, tt.id, tt.value)
					assert.Error(t, err)
					assert.Nil(t, added)
				})
			}
		})

	})

	t.Run("Patch", func(t *testing.T) {

		t.Run("Successfully patches and retrieves the stored content after", func(t *testing.T) {
			tearDown, s := prepareMockDB(t)
			defer tearDown()

			id := "id1"
			_, err := s.Add(ctx, admin, id, fixture)
			require.NoError(t, err)
			newAdmin := "admin2@example.com"
			newContent := content{
				I: 2,
				B: false,
				S: "New Text",
			}
			patched, err := s.Patch(ctx, newAdmin, id, map[string]any{"i": newContent.I, "b": newContent.B, "s": newContent.S})
			require.NoError(t, err)
			assert.Equal(t, newContent, patched.Content)
			assert.Equal(t, newAdmin, patched.ModifiedBy)
			assert.True(t, patched.CreatedAt.Before(patched.ModifiedAt))

			fetched, err := s.Get(ctx, id)
			assert.NoError(t, err)
			assert.Equal(t, *patched, *fetched)
		})

		t.Run("Allows patching unmodeleted attribte", func(t *testing.T) {
			tearDown, s := prepareMockDB(t)
			defer tearDown()

			id := "id1"
			_, err := s.Add(ctx, admin, id, fixture)
			require.NoError(t, err)
			newAdmin := "admin2@example.com"
			patched, err := s.Patch(ctx, newAdmin, id, map[string]any{"Z": 3})
			require.NoError(t, err)

			fetched, err := s.Get(ctx, id)
			assert.NoError(t, err)
			assert.Equal(t, *patched, *fetched)
		})

		t.Run("Fails if patching breaks modeled attribute", func(t *testing.T) {
			tearDown, s := prepareMockDB(t)
			defer tearDown()

			id := "id1"
			added, err := s.Add(ctx, admin, id, fixture)
			require.NoError(t, err)
			newAdmin := "admin2@example.com"
			patched, err := s.Patch(ctx, newAdmin, id, map[string]any{"i": "Not Int!"})
			assert.Error(t, err)
			assert.Nil(t, patched)

			// No changes applied
			fetched, err := s.Get(ctx, id)
			assert.NoError(t, err)
			assert.Equal(t, *added, *fetched)
		})

	})

	t.Run("List", func(t *testing.T) {
		fixtures := map[string]content{
			"1": {I: 0, B: true, S: "a"},
			"2": {I: 10, B: false, S: "b"},
			"3": {I: 20, B: true, S: "c"},
			"4": {I: 30, B: false, S: "d"},
		}
		tests := []struct {
			name             string
			cond             []Condition
			expectedFixtures []string
		}{
			{"empty result", []Condition{{Attribute: "i", Op: EqualOperator, Value: -1}}, []string{}},
			{"all results", nil, []string{"1", "2", "3", "4"}},
			{"equal int", []Condition{{Attribute: "i", Op: EqualOperator, Value: fixtures["1"].I}}, []string{"1"}},
			{"not equal int", []Condition{{Attribute: "i", Op: NotEqualOperator, Value: fixtures["1"].I}}, []string{"2", "3", "4"}},
			{"greater than int", []Condition{{Attribute: "i", Op: GreaterThanOperator, Value: fixtures["2"].I}}, []string{"3", "4"}},
			{"greater than equal int", []Condition{{Attribute: "i", Op: GreaterThanOrEqualOpertor, Value: fixtures["2"].I}}, []string{"2", "3", "4"}},
			{"less than int", []Condition{{Attribute: "i", Op: LessThanOperator, Value: fixtures["2"].I}}, []string{"1"}},
			{"less than equal int", []Condition{{Attribute: "i", Op: LessThanOrEqualOperator, Value: fixtures["2"].I}}, []string{"1", "2"}},
			{"equal bool", []Condition{{Attribute: "b", Op: EqualOperator, Value: true}}, []string{"1", "3"}},
			{"not equal bool", []Condition{{Attribute: "b", Op: NotEqualOperator, Value: true}}, []string{"2", "4"}},
			{"equal string", []Condition{{Attribute: "s", Op: EqualOperator, Value: fixtures["1"].S}}, []string{"1"}},
			{"not equal string", []Condition{{Attribute: "s", Op: NotEqualOperator, Value: fixtures["1"].S}}, []string{"2", "3", "4"}},
			{"greater than string", []Condition{{Attribute: "s", Op: GreaterThanOperator, Value: fixtures["2"].S}}, []string{"3", "4"}},
			{"greater than equal string", []Condition{{Attribute: "s", Op: GreaterThanOrEqualOpertor, Value: fixtures["2"].S}}, []string{"2", "3", "4"}},
			{"less than string", []Condition{{Attribute: "s", Op: LessThanOperator, Value: fixtures["2"].S}}, []string{"1"}},
			{"less than equal string", []Condition{{Attribute: "s", Op: LessThanOrEqualOperator, Value: fixtures["2"].S}}, []string{"1", "2"}},
			{"Ands Multiple Conditions", []Condition{{Attribute: "b", Op: EqualOperator, Value: true}, {Attribute: "i", Op: GreaterThanOperator, Value: fixtures["2"].I}}, []string{"3"}},
		}
		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				tearDown, s := prepareMockDB(t)
				defer tearDown()

				// add fixtures
				addedFixtures := map[string]Stored[content]{}
				for id, f := range fixtures {
					added, err := s.Add(ctx, admin, id, f)
					require.NoError(t, err)
					require.NotNil(t, added)
					addedFixtures[id] = *added
				}

				expected := []Stored[content]{}
				for _, i := range tt.expectedFixtures {
					expected = append(expected, addedFixtures[i])
				}

				result, err := s.List(ctx, tt.cond...)
				require.NoError(t, err)
				assert.ElementsMatch(t, expected, result)

			})
		}
	})
}
