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
func (s *sqlRepo) SearchTasks(title *taskEntity.SearchTitleParams, ctx context.Context) ([]*taskEntity.SearchTaskRes ,error) {

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
		WHERE title LIKE "%%%s%%"
	`, title.SearchQuery)


	rows, err := db.QueryContext(ctx, stmt)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	Searchedtasks := []*taskEntity.SearchTaskRes{}

	for rows.Next(){
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
