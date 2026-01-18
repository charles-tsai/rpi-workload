package api

import (
	"context"
	"errors"
	"testing"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

// MockDB is a mock implementation of DBInterface
type MockDB struct {
	ExecFunc  func(ctx context.Context, sql string, arguments ...any) (pgconn.CommandTag, error)
	QueryFunc func(ctx context.Context, sql string, arguments ...any) (pgx.Rows, error)
}

func (m *MockDB) Exec(ctx context.Context, sql string, arguments ...any) (pgconn.CommandTag, error) {
	if m.ExecFunc != nil {
		return m.ExecFunc(ctx, sql, arguments...)
	}
	return pgconn.NewCommandTag(""), nil
}

func (m *MockDB) Query(ctx context.Context, sql string, arguments ...any) (pgx.Rows, error) {
	if m.QueryFunc != nil {
		return m.QueryFunc(ctx, sql, arguments...)
	}
	return &MockRows{}, nil
}

// MockRows implements pgx.Rows
type MockRows struct {
	data [][]interface{}
	idx  int
	err  error
}

func (m *MockRows) Close()                                         {}
func (m *MockRows) Err() error                                     { return m.err }
func (m *MockRows) CommandTag() pgconn.CommandTag                  { return pgconn.NewCommandTag("") }
func (m *MockRows) FieldDescriptions() []pgconn.FieldDescription   { return nil }
func (m *MockRows) Next() bool                                     { m.idx++; return m.idx <= len(m.data) }
func (m *MockRows) Scan(dest ...any) error {
	if m.idx < 1 || m.idx > len(m.data) {
		return errors.New("invalid row index")
	}
	row := m.data[m.idx-1]
	for i, v := range row {
		if i >= len(dest) {
			break
		}
		// In a real implementation we would handle type conversion, but for now we just assign if types match or pointer
		// Here we keep it simple for the test case which knows the types.
		switch d := dest[i].(type) {
		case *string:
			*d = v.(string)
		}
	}
	return nil
}
func (m *MockRows) Values() ([]any, error) { return nil, nil }
func (m *MockRows) RawValues() [][]byte    { return nil }
func (m *MockRows) Conn() *pgx.Conn        { return nil }

func TestCreateApp(t *testing.T) {
	var execCalled bool
	mockDB := &MockDB{
		ExecFunc: func(ctx context.Context, sql string, arguments ...any) (pgconn.CommandTag, error) {
			execCalled = true
			if sql != "INSERT INTO apps (id, name) VALUES ($1, $2)" {
				t.Errorf("expected sql 'INSERT INTO apps (id, name) VALUES ($1, $2)', got '%s'", sql)
			}
			if len(arguments) != 2 {
				t.Errorf("expected 2 arguments, got %d", len(arguments))
			}
			if arguments[0] != "test-id" {
				t.Errorf("expected id 'test-id', got '%v'", arguments[0])
			}
			if arguments[1] != "test-app" {
				t.Errorf("expected name 'test-app', got '%v'", arguments[1])
			}
			return pgconn.NewCommandTag("INSERT 0 1"), nil
		},
	}

	server := NewServer(mockDB)
	req := CreateAppRequestObject{
		Body: &App{
			Id:   "test-id",
			Name: "test-app",
		},
	}

	resp, err := server.CreateApp(context.Background(), req)
	if err != nil {
		t.Fatalf("CreateApp failed: %v", err)
	}

	if !execCalled {
		t.Error("Exec was not called")
	}

	jsonResp, ok := resp.(CreateApp201JSONResponse)
	if !ok {
		t.Errorf("expected CreateApp201JSONResponse, got %T", resp)
	}
	if jsonResp.Id != "test-id" || jsonResp.Name != "test-app" {
		t.Errorf("unexpected response body: %+v", jsonResp)
	}
}

func TestGetApps(t *testing.T) {
	mockDB := &MockDB{
		QueryFunc: func(ctx context.Context, sql string, arguments ...any) (pgx.Rows, error) {
			if sql != "SELECT id, name FROM apps" {
				t.Errorf("expected sql 'SELECT id, name FROM apps', got '%s'", sql)
			}
			return &MockRows{
				data: [][]interface{}{
					{"1", "App 1"},
					{"2", "App 2"},
				},
			}, nil
		},
	}

	server := NewServer(mockDB)
	resp, err := server.GetApps(context.Background(), GetAppsRequestObject{})
	if err != nil {
		t.Fatalf("GetApps failed: %v", err)
	}

	jsonResp, ok := resp.(GetApps200JSONResponse)
	if !ok {
		t.Errorf("expected GetApps200JSONResponse, got %T", resp)
	}

	if len(jsonResp) != 2 {
		t.Errorf("expected 2 apps, got %d", len(jsonResp))
	}
	if jsonResp[0].Id != "1" || jsonResp[0].Name != "App 1" {
		t.Errorf("unexpected app 1: %+v", jsonResp[0])
	}
}
