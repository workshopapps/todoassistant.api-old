package mySqlRepo

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"test-va/internals/Repository/taskRepo"
	"test-va/internals/entity/taskEntity"
	"test-va/internals/entity/vaEntity"
)

type sqlRepo struct {
	conn *sql.DB
}

func (s *sqlRepo) AssignTaskToVa(ctx context.Context, vaId, taskId string) error {
	log.Println(vaId)
	log.Println(taskId)
	stmt := fmt.Sprintf(`UPDATE Tasks SET va_id ='%v' WHERE task_id ='%v'`, vaId, taskId)
	_, err := s.conn.ExecContext(ctx, stmt)
	if err != nil {
		return err
	}
	return nil
}

func (s *sqlRepo) GetVADetails(ctx context.Context, userId string) (string, error) {
	var vaId *string
	stmt := fmt.Sprintf(`
SELECT
	virtual_Assistant_id from Users
WHERE user_id = '%v'
`, userId)
	row := s.conn.QueryRowContext(ctx, stmt)
	err := row.Scan(&vaId)
	if err != nil {
		return "", err
	}
	return *vaId, nil
}

func (s *sqlRepo) GetAllTaskAssignedToVA(ctx context.Context, vaId string) ([]*vaEntity.VATask, error) {
	stmt := fmt.Sprintf(`SELECT
    T.task_id,
    T.title,
    T.end_time,
    T.status,
    T.description,
    concat(U.first_name, ' ', U.last_name) AS 'name',
    T.user_id,
    U.phone
FROM Tasks T
         join Users U on T.va_id = U.virtual_assistant_id
WHERE T.va_id = '%s'
;`, vaId)

	queryRow, err := s.conn.QueryContext(ctx, stmt)
	if err != nil {
		return nil, err
	}

	var Results []*vaEntity.VATask

	for queryRow.Next() {
		var res vaEntity.VATask
		err := queryRow.Scan(&res.TaskId, &res.Title, &res.EndTime, &res.Status, &res.Description, &res.User.Name, &res.User.UserId, &res.User.Phone)
		if err != nil {
			return nil, err
		}
		Results = append(Results, &res)
	}

	return Results, nil
}

func (s *sqlRepo) SetNewEvent(req *taskEntity.CreateTaskReq) error {
	stmt := fmt.Sprintf(`INSERT INTO Tasks(
                  task_id,
                  user_id,
                  title,
                  description,
                  start_time,
                  end_time,
                  created_at,
                  va_option,
                  repeat_frequency
                  )
	VALUES ('%v','%v','%v','%v','%v','%v','%v','%v','%v')
	`, req.TaskId, req.UserId, req.Title, req.Description, req.StartTime, req.EndTime, req.CreatedAt, req.VAOption, req.Repeat)
	_, err := s.conn.Exec(stmt)
	if err != nil {
		log.Println(stmt)
		log.Println(err)
		return err
	}
	return nil
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
	log.Printf("#%v\n", req)
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
                  user_id,
                  title,
                  description,
                  start_time,
                  end_time,
                  created_at,
                  va_option,
                  repeat_frequency
				   )
		VALUES ('%v','%v','%v','%v','%v','%v','%v', '%v', '%v')`, req.TaskId, req.UserId, req.Title, req.Description,
		req.StartTime, req.EndTime, req.CreatedAt, req.VAOption, req.Repeat)

	_, err = tx.ExecContext(ctx, stmt)
	if err != nil {
		log.Println(err)
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
			log.Println("err", err)
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

	var Searchedtasks []*taskEntity.SearchTaskRes

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

func (s *sqlRepo) GetTaskByID(ctx context.Context, taskId string) (*taskEntity.GetTasksByIdRes, error) {

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
		SELECT T.task_id, T.user_id, T.title, T.description, T.status, T.start_time, T.end_time, T.created_at
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
		&res.CreatedAt,
	); err != nil {
		return nil, err
	}
	log.Println("Created AT", res)
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
	db, err := s.conn.Begin()
	if err != nil {
		return nil, err
	}

	stmt := fmt.Sprintf(`
		SELECT task_id, user_id, title, start_time
		FROM Tasks
		WHERE status = 'EXPIRED'`)

	rows, err := db.QueryContext(ctx, stmt)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var Searchedtasks []*taskEntity.GetAllExpiredRes

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

func (s *sqlRepo) GetListOfPendingTasks(ctx context.Context) ([]*taskEntity.GetAllPendingRes, error) {
	db, err := s.conn.Begin()
	if err != nil {
		return nil, err
	}

	stmt := fmt.Sprintf(`
		SELECT task_id, user_id, title, end_time
		FROM Tasks
		WHERE status = 'PENDING'`)

	rows, err := db.QueryContext(ctx, stmt)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var SearchedPendingtasks []*taskEntity.GetAllPendingRes

	for rows.Next() {
		var singleTask taskEntity.GetAllPendingRes

		err := rows.Scan(
			&singleTask.TaskId,
			&singleTask.UserId,
			&singleTask.Title,
			// &singleTask.VAOption,
			&singleTask.EndTime,
		)
		// fmt.Println(err)
		if err != nil {
			return nil, err
		}
		SearchedPendingtasks = append(SearchedPendingtasks, &singleTask)
	}
	return SearchedPendingtasks, nil
}

// Get All task
func (s *sqlRepo) GetAllTasks(ctx context.Context, userId string) ([]*taskEntity.GetAllTaskRes, error) {

	//tx, err := s.conn.BeginTx(ctx, nil)
	db, err := s.conn.Begin()
	if err != nil {
		return nil, err
	}

	stmt := fmt.Sprintf(`
		SELECT task_id, title, description, repeat_frequency, va_option, status, start_time, end_time
		FROM Tasks WHERE user_id = '%s'`, userId)

	rows, err := db.QueryContext(ctx, stmt)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var AllTasks []*taskEntity.GetAllTaskRes

	for rows.Next() {
		var singleTask taskEntity.GetAllTaskRes

		err := rows.Scan(
			&singleTask.TaskId,
			&singleTask.Title,
			&singleTask.Description,
			&singleTask.Repeat,
			&singleTask.VAOption,
			&singleTask.Status,
			&singleTask.StartTime,
			&singleTask.EndTime,
		)
		if err != nil {
			return nil, err
		}
		AllTasks = append(AllTasks, &singleTask)
	}
	return AllTasks, nil
}

// Delete task by id
func (s *sqlRepo) DeleteTaskByID(ctx context.Context, taskId string) error {

	_, err := s.conn.ExecContext(ctx, fmt.Sprintf(`Delete from Tasks  WHERE task_id = '%s'`, taskId))
	if err != nil {
		log.Fatal(err)
	}
	return nil
}

// Delete All
func (s *sqlRepo) DeleteAllTask(ctx context.Context, userId string) error {

	//var res taskEntity.GetTasksByIdRes
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
	_, err = tx.ExecContext(ctx, fmt.Sprintf(`Delete from Tasks where  user_id = '%s'`, userId))
	if err != nil {
		log.Fatal(err)
	}

	return nil
}

func (s *sqlRepo) EditTaskById(ctx context.Context, taskId string, req *taskEntity.EditTaskReq) error {

	_, err := s.conn.ExecContext(ctx, fmt.Sprintf(`UPDATE Tasks SET
                 title = '%s',
                 description = '%s',
                 end_time = '%s',
                 updated_at = '%s',
                 va_option ='%s',
                 repeat_frequency= '%s'
             WHERE task_id = '%s'
            `, req.Title, req.Description, req.EndTime, req.UpdatedAt, req.VAOption, req.Repeat, taskId))
	if err != nil {
		log.Fatal(err)
	}
	return nil
}

func (s *sqlRepo) UpdateTaskStatusByID(ctx context.Context, taskId string) error {
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
	_, err = tx.ExecContext(ctx, fmt.Sprintf(`UPDATE Tasks SET status = 'DONE' WHERE task_id = '%s'`, taskId))
	if err != nil {
		log.Fatal(err)
	}
	return nil
}

// comment
func (s *sqlRepo) PersistComment(ctx context.Context, req *taskEntity.CreateCommentReq) error {

	stmt := fmt.Sprintf(`INSERT INTO Comments(
                  user_id,
                  task_id,
                  comment,
				  created_at
                  )
	VALUES ('%v','%v','%v','%v')
	`, req.UserId, req.TaskId, req.Comment, req.CreatedAt)

	_, err := s.conn.Exec(stmt)
	if err != nil {
		log.Println(stmt)
		log.Println(err)
		return err
	}

	return nil
}

func (s *sqlRepo) GetAllComments(ctx context.Context, taskId string) ([]*taskEntity.GetCommentRes, error) {

	//tx, err := s.conn.BeginTx(ctx, nil)
	db, err := s.conn.Begin()
	if err != nil {
		return nil, err
	}

	stmt := fmt.Sprintf(`
		SELECT user_id, task_id, comment, created_at
		FROM Comments WHERE task_id = '%s'`, taskId)

	rows, err := db.QueryContext(ctx, stmt)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var AllComment []*taskEntity.GetCommentRes

	for rows.Next() {
		var singleTask taskEntity.GetCommentRes

		err := rows.Scan(
			&singleTask.UserId,
			&singleTask.TaskId,
			&singleTask.Comment,
			&singleTask.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		AllComment = append(AllComment, &singleTask)
	}
	return AllComment, nil
}
