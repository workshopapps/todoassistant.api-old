package mySqlRepo

import (
	"context"
	"database/sql"
	"fmt"
	"test-va/internals/Repository/taskRepo"
	"test-va/internals/entity/taskEntity"
)

type sqlRepo struct {
	conn *sql.DB
}

func (s *sqlRepo) SetTaskToExpired(id string) error {
	stmt := fmt.Sprintf(`UPDATE Tasks SET STATUS='EXPIRED' WHERE task_id ='%v'`, id)
	_, err := s.conn.Exec(stmt)
	if err != nil {
		return err
	}
	return nil
}

func NewSqlRepo(conn *sql.DB) taskRepo.TaskRepository {
	return &sqlRepo{conn: conn}
}

func (s *sqlRepo) GetPendingTasks(userId string, ctx context.Context) ([]*taskEntity.GetPendingTasksRes, error) {

	query := fmt.Sprintf(`
		SELECT task_id, user_id, title, description, start_time, end_time, status
		FROM Tasks
		WHERE user_id = '%s' AND status = 'PENDING'
	`, userId)

	rows, err := s.conn.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tasks []*taskEntity.GetPendingTasksRes

	for rows.Next() {
		var task taskEntity.GetPendingTasksRes
		err := rows.Scan(
			&task.TaskId,
			&task.UserId,
			&task.Title,
			&task.Description,
			&task.StartTime,
			&task.EndTime,
			&task.Status,
		)
		if err != nil {
			return nil, err
		}
		tasks = append(tasks, &task)
	}
	if rows.Err(); err != nil {
		return nil, err
	}
	return tasks, nil
}

func (s *sqlRepo) Persist(ctx context.Context, req *taskEntity.CreateTaskReq) error {
	tx, err := s.conn.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	defer func() {
		if err != nil {
			tx.Rollback()
		} else {
			tx.Commit()
		}
	}()

	stmt := fmt.Sprintf(`INSERT
		INTO Tasks(
				   task_id,
				   title,
				   description,
				   user_id,
				   start_time,
				   end_time
				   )
		VALUES ('%v','%v','%v','%v','%v','%v')`, req.TaskId, req.Title, req.Description, req.UserId, req.StartTime, req.EndTime)

	_, err = tx.ExecContext(ctx, stmt)
	if err != nil {
		return err
	}

	for _, file := range req.Files {
		stmt2 := fmt.Sprintf(`INSERT
		INTO Taskfiles(
		               task_id,
		               file_link,
		               file_type
		               )
		VALUES ('%v', '%v', '%v')`, req.TaskId, file.FileLink, file.FileType)
		_, err = tx.ExecContext(ctx, stmt2)
		if err != nil {
			return err
		}
	}

	return nil
}

// search by name
func (s *sqlRepo) SearchTasks(title *taskEntity.SearchTitleParams, ctx context.Context) ([]*taskEntity.SearchTaskRes, error) {

	//tx, err := s.conn.BeginTx(ctx, nil)
	db, err := s.conn.Begin()
	if err != nil {
		return nil, err
	}

	// defer func() {
	// 	if err != nil {
	// 		tx.Rollback()
	// 	} else {
	// 		tx.Commit()
	// 	}
	// }()

	stmt := fmt.Sprintf(`
		SELECT task_id, user_id, title, start_time
		FROM Tasks
		WHERE title LIKE '%s%%'
	`, title.SearchQuery)

	rows, err := db.QueryContext(ctx, stmt)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	Searchedtasks := []*taskEntity.SearchTaskRes{}

	for rows.Next() {
		var singleTask taskEntity.SearchTaskRes

		err := rows.Scan(
			&singleTask.TaskId,
			&singleTask.UserId,
			&singleTask.Title,
			&singleTask.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		Searchedtasks = append(Searchedtasks, &singleTask)
	}
	return Searchedtasks, nil
}

// get task by ID
func (s *sqlRepo) GetTaskByID(taskId string, ctx context.Context) (*taskEntity.GetTasksByIdRes, error) {

	var res taskEntity.GetTasksByIdRes
	tx, err := s.conn.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}

	defer func() {
		if err != nil {
			tx.Rollback()
		} else {
			tx.Commit()
		}
	}()

	stmt := fmt.Sprintf(`
		SELECT T.task_id, T.user_id, T.title, T.description, T.status, T.start_time, T.end_time
		FROM Tasks T
		WHERE task_id = '%s'
	`, taskId)

	stmt2 := fmt.Sprintf(`
		SELECT F.file_link, F.file_type
		FROM Tasks AS T
		JOIN Taskfiles as F
		ON T.task_id = F.task_id
		WHERE F.task_id = '%s'
	`, taskId)

	row := tx.QueryRow(stmt)
	if err := row.Scan(
		&res.TaskId,
		&res.UserId,
		&res.Title,
		&res.Description,
		&res.Status,
		&res.StartTime,
		&res.EndTime,
	); err != nil {
		return nil, err
	}

	rows, err := tx.QueryContext(ctx, stmt2)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var taskFile taskEntity.TaskFile

		err := rows.Scan(
			&taskFile.FileLink,
			&taskFile.FileType,
		)
		if err != nil {
			return nil, err
		}
		res.Files = append(res.Files, taskFile)
	}

	return &res, nil
}



func (s *sqlRepo) GetListOfExpiredTasks(ctx context.Context) ([]*taskEntity.GetAllExpiredRes, error) {

	//tx, err := s.conn.BeginTx(ctx, nil)
	db, err := s.conn.Begin()
	if err != nil {
		return nil, err
	}

	// defer func() {
	// 	if err != nil {
	// 		tx.Rollback()
	// 	} else {
	// 		tx.Commit()
	// 	}
	// }()

	stmt := fmt.Sprintf(`
		SELECT task_id, user_id, title, start_time
		FROM Tasks
		WHERE status = 'EXPIRED'`)

	rows, err := db.QueryContext(ctx, stmt)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	Searchedtasks := []*taskEntity.GetAllExpiredRes{}

	for rows.Next() {
		var singleTask taskEntity.GetAllExpiredRes

		err := rows.Scan(
			&singleTask.TaskId,
			&singleTask.UserId,
			&singleTask.Title,
			&singleTask.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		Searchedtasks = append(Searchedtasks, &singleTask)
	}
	return Searchedtasks, nil
}