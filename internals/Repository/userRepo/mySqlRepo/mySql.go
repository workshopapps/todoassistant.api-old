package mySqlRepo

import (
	"context"
	"database/sql"
	"fmt"
	"test-va/internals/Repository/userRepo"
	"test-va/internals/entity/userEntity"
)

type mySql struct {
	conn *sql.DB
}

func NewMySqlUserRepo(conn *sql.DB) userRepo.UserRepository {
	return &mySql{conn: conn}
}

func (m *mySql) GetUsers(page int) ([]*userEntity.UsersRes, error) {
	var allUsers []*userEntity.UsersRes
	limit := 20
	offset := limit * (page - 1)
	query := fmt.Sprintf(`SELECT user_id, email, first_name, last_name, phone, date_of_birth, date_created
							FROM Users
							ORDER BY user_id
							LIMIT %d
							OFFSET %d`, limit, offset)

	ctx := context.Background()
	rows, err := m.conn.QueryContext(ctx, query)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	defer rows.Close()

	fmt.Println(rows.Next())
	// for rows.NextResultSet() {
	for rows.Next() {
		var user userEntity.UsersRes
		err := rows.Scan(
			&user.UserId,
			&user.Email,
			&user.FirstName,
			&user.LastName,
			&user.Phone,
			&user.DateOfBirth,
			&user.DateCreated,
		)

		if err != nil {
			return allUsers, err
		}

		allUsers = append(allUsers, &user)
	}
	// }

	if err = rows.Err(); err != nil {
		return allUsers, err
	}
	return allUsers, nil
}

func (m *mySql) GetByEmail(email string) (*userEntity.GetByEmailRes, error) {
	query := fmt.Sprintf(`
		SELECT user_id, email, password, first_name, last_name, phone, gender 
		FROM Users
		WHERE email = '%s'
	`, email)
	var user userEntity.GetByEmailRes
	ctx := context.Background()
	err := m.conn.QueryRowContext(ctx, query).Scan(
		&user.UserId,
		&user.Email,
		&user.Password,
		&user.FirstName,
		&user.LastName,
		&user.Phone,
		&user.Gender,
	)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	return &user, nil
}

func (m *mySql) GetById(user_id string) (*userEntity.GetByIdRes, error) {
	query := fmt.Sprintf(`
		SELECT user_id, email, first_name, last_name, phone, gender 
		FROM Users
		WHERE user_id = '%s'
	`, user_id)

	var user userEntity.GetByIdRes
	ctx := context.Background()
	err := m.conn.QueryRowContext(ctx, query).Scan(
		&user.UserId,
		&user.Email,
		&user.FirstName,
		&user.LastName,
		&user.Phone,
		&user.Gender,
	)

	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	return &user, nil
}

func (m *mySql) Persist(req *userEntity.CreateUserReq) error {
	stmt := fmt.Sprintf(` INSERT INTO Users(
                   user_id,
                   first_name,
                   last_name,
                   email,
                   phone,
                   password,
                   gender,
                   date_of_birth,
                   account_status,
                   payment_status,
                   date_created
                   ) VALUES ('%v', '%v', '%v', '%v', '%v', '%v', '%v', '%v', '%v', '%v', '%v')`,
		req.UserId, req.FirstName, req.LastName, req.Email, req.Phone, req.Password, req.Gender, req.DateOfBirth, req.AccountStatus, req.PaymentStatus, req.DateCreated)

	_, err := m.conn.Exec(stmt)
	if err != nil {
		return err
	}
	return nil
}

// Create function to update user in database
func (m *mySql) UpdateUser(req *userEntity.UpdateUserReq, userId string) (*userEntity.GetByIdRes, error) {
	// stmt := fmt.Sprintf(`UDPATE Users
	// 	SET first_name = '%v',
	// 	last_name = '%v',
	// 	email = '%v',
	// 	phone = '%v',
	// 	gender = '%v',
	// 	date_of_birth = '%v',
	// 	account_status = '%v',
	// 	payment_status = '%v',
	// 	WHERE user_id='%v'
	// 	`,
	// 	req.FirstName, req.LastName, req.Email, req.Phone, req.Gender, req.DateOfBirth, req.AccountStatus, req.PaymentStatus, userId)

	// _, err := m.conn.Exec(stmt)
	// if err != nil {
	// 	return err
	// }

	tx, err := m.conn.BeginTx(context.Background(), nil)
	if err != nil {
		return nil, err
	}

	_, err = updateFieldIfSet(tx, "first_name", req.FirstName)
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	_, err = updateFieldIfSet(tx, "last_name", req.LastName)
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	err = tx.Commit()
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	return m.GetById(userId)
}

func updateField(tx *sql.Tx, field string, val interface{}) (sql.Result, error) {
	return tx.Exec(fmt.Sprintf(`UPDATE Users SET %s = '%v'`, field, val))
}

func updateFieldIfSet(tx *sql.Tx, field string, val interface{}) (sql.Result, error) {
	v, ok := val.(string)
	if ok && v != "" {
		return updateField(tx, field, v)
	}
	return nil, nil
}

func (m *mySql) DeleteUser(user_id string) error {
	query := fmt.Sprintf(`DELETE FROM Users WHERE user_id = "%s"`, user_id)
	_, err := m.conn.Exec(query)
	if err != nil {
		fmt.Println(err)
		return err
	}

	return nil
}
