package tasks

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func getTestTask() Task {
	return Task{
		Id:          "Id123",
		Title:       "Test Task Title",
		Description: "Task description",
		DueDate:     time.Date(2025, 6, 4, 5, 50, 0, 0, time.UTC),
		Status:      StatusInProgress,
	}
}

func TestMemStore_Add(t *testing.T) {
	type fields struct {
		list map[string]Task
	}
	type args struct {
		id   string
		task Task
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
		wantLen int
	}{
		{
			name: "Add to empty map",
			fields: fields{
				map[string]Task{},
			},
			args: args{
				id:   "Id123",
				task: getTestTask(),
			},
			wantLen: 1,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := MemStore{
				list: tt.fields.list,
			}
			err := m.Add(tt.args.id, tt.args.task)
			if !tt.wantErr {
				assert.NoError(t, err)
			}

			assert.Len(t, tt.fields.list, tt.wantLen)
		})
	}
}

func TestMemStore_Get(t *testing.T) {
	type fields struct {
		list map[string]Task
	}
	type args struct {
		id string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    Task
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "Find test task",
			fields: fields{
				map[string]Task{
					"Id123": getTestTask(),
				},
			},
			args: args{
				id: "Id123",
			},
			want:    getTestTask(),
			wantErr: nil,
		},
		{
			name: "Not Found Ratatouille",
			fields: fields{
				map[string]Task{
					"Id123": getTestTask(),
				},
			},
			args: args{
				id: "Id234",
			},
			want: Task{},
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return err == NotFoundErr
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := MemStore{
				list: tt.fields.list,
			}
			got, err := m.Get(tt.args.id)
			if tt.wantErr != nil {
				if !tt.wantErr(t, err, fmt.Sprintf("Get(%v)", tt.args.id)) {
					require.Failf(t, "Invalid error message", "Got: %v", err.Error())
				}
			} else {
				assert.NoError(t, err)
			}

			assert.Equalf(t, tt.want, got, "Get(%v)", tt.args.id)
		})
	}
}

func TestMemStore_List(t *testing.T) {
	type fields struct {
		list map[string]Task
	}
	tests := []struct {
		name    string
		fields  fields
		want    map[string]Task
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "Simple list",
			fields: fields{
				map[string]Task{
					"Id123": getTestTask(),
				},
			},
			want: map[string]Task{
				"Id123": getTestTask(),
			},
			wantErr: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := MemStore{
				list: tt.fields.list,
			}
			got, err := m.List()
			if tt.wantErr != nil {
				if !tt.wantErr(t, err, fmt.Sprintf("List()")) {
					assert.Fail(t, "Invalid error")
				}
			} else {
				assert.NoError(t, err)
			}

			assert.Equalf(t, tt.want, got, "List()")
		})
	}
}

func TestMemStore_Remove(t *testing.T) {
	type fields struct {
		list map[string]Task
	}
	type args struct {
		id string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr assert.ErrorAssertionFunc
		wantLen int
	}{
		{
			name: "Empty list",
			fields: fields{
				map[string]Task{
					"Id123": getTestTask(),
				},
			},
			args: args{
				id: "Id123",
			},
			wantErr: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := MemStore{
				list: tt.fields.list,
			}

			err := m.Remove(tt.args.id)

			if tt.wantErr != nil {
				if !tt.wantErr(t, err, fmt.Sprintf("List()")) {
					assert.Fail(t, "Invalid error")
				}
			} else {
				assert.NoError(t, err)
			}

			assert.Len(t, m.list, tt.wantLen)
		})
	}
}

func TestMemStore_Update(t *testing.T) {
	type fields struct {
		list map[string]Task
	}
	type args struct {
		id   string
		task Task
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr assert.ErrorAssertionFunc
		wantLen int
	}{
		{
			name: "Update butter to Ham and cheese",
			fields: fields{
				map[string]Task{
					"Id123": getTestTask(),
				},
			},
			args: args{
				id: "Id123",
				task: Task{
					Id:          "Id123",
					Title:       "Updated title",
					Description: "Updated Description",
					DueDate:     time.Date(2025, 9, 23, 0, 0, 0, 0, time.UTC),
					Status:      StatusToDo,
				},
			},
			wantErr: nil,
			wantLen: 1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := MemStore{
				list: tt.fields.list,
			}

			err := m.Update(tt.args.id, tt.args.task)
			if tt.wantErr != nil {
				if !tt.wantErr(t, err, fmt.Sprintf("List()")) {
					assert.Fail(t, "Invalid error")
				}
			} else {
				assert.NoError(t, err)
			}

			assert.Len(t, m.list, tt.wantLen)
		})
	}
}
