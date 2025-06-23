-- +goose Up
-- +goose StatementBegin
-- Create a new schema (if not exists)
CREATE SCHEMA IF NOT EXISTS task_manager;

CREATE TABLE IF NOT EXISTS task_manager.task (
                                   id VARCHAR(200) NOT NULL UNIQUE,
                                   title VARCHAR(100) NOT NULL UNIQUE,
                                   description VARCHAR(1000) NOT NULL,
                                   due_date TIMESTAMP NOT NULL,
                                   status INTEGER NOT NULL,
                                   created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
) ;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS task_manager.task;
DROP SCHEMA IF EXISTS task_manager;
-- +goose StatementEnd
