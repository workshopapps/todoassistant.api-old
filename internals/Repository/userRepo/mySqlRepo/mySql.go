package mySqlRepo

import (
	"context"
	"database/sql"
	"fmt"
	"test-va/internals/Repository/userRepo"
	"test-va/internals/entity/userEntity"
	"time"
)

type mySql struct {
	conn *sql.DB
}

func NewMySqlUserRepo(conn *sql.DB) userRepo.UserRepository {
	return &mySql{conn: conn}
}

func (m *mySql) AssignVAToUser(user_id, va_id string) error {
	query := fmt.Sprintf(`
		SELECT user_id, virtual_assistant_id
		FROM Users
		WHERE user_id = '%s'
	`, user_id)
	var userId string
	var vaId string
	err := m.conn.QueryRowContext(context.Background(), query).Scan(&userId, &vaId)
	if err != nil {
		switch {
		case err.Error() == "sql: Scan error on column index 1, name \"virtual_assistant_id\": converting NULL to string is unsupported":
			vaId = ""
		default:
			return err
		}
	}
	if vaId != "" {
		return fmt.Errorf("user already has a VA")
	}

	query = fmt.Sprintf(`
		UPDATE Users SET virtual_assistant_id = '%s' WHERE user_id = '%s'
	`, va_id, user_id)

	_, err = m.conn.ExecContext(context.Background(), query)
	if err != nil {
		fmt.Println(err)
		return err
	}
	return nil
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
		SELECT user_id, email, password, first_name, last_name, phone, gender, avatar
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
		&user.Avatar,
	)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	return &user, nil
}

func (m *mySql) GetById(user_id string) (*userEntity.GetByIdRes, error) {
	query := fmt.Sprintf(`
		SELECT user_id, password, email, first_name, last_name, phone, gender, avatar
		FROM Users
		WHERE user_id = '%s'
	`, user_id)

	var user userEntity.GetByIdRes
	ctx := context.Background()
	err := m.conn.QueryRowContext(ctx, query).Scan(
		&user.UserId,
		&user.Password,
		&user.Email,
		&user.FirstName,
		&user.LastName,
		&user.Phone,
		&user.Gender,
		&user.Avatar,
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

func (m *mySql) UpdateUser(req *userEntity.UpdateUserReq, userId string) error {
	ctx, cancelFunc := context.WithTimeout(context.TODO(), time.Second*60)
	defer cancelFunc()

	stmt := fmt.Sprintf(`UPDATE Users SET 
                 first_name ='%s',
                 last_name='%s',
                 email ='%s',
                 phone='%s',
                 gender='%s',
                 date_of_birth='%s' WHERE user_id ='%s'
                 `, req.FirstName, req.LastName, req.Email, req.Phone, req.Gender, req.DateOfBirth, userId)

	_, err := m.conn.ExecContext(ctx, stmt)
	if err != nil {
		return err
	}
	return nil
}

func (m *mySql) UpdateImage(userId, fileName string) error {
	ctx, cancelFunc := context.WithTimeout(context.TODO(), time.Second*60)
	defer cancelFunc()

	stmt := fmt.Sprintf(`UPDATE Users SET 
                 avatar = '%s'
                 WHERE user_id ='%s'
                 `, fileName, userId)

	_, err := m.conn.ExecContext(ctx, stmt)
	if err != nil {
		return err
	}
	return nil
}

// Auxillary function to update user
func updateField(tx *sql.Tx, userId string, field string, val interface{}) (sql.Result, error) {
	return tx.Exec(fmt.Sprintf(`UPDATE Users SET %s = '%v' WHERE user_id = '%v'`, field, val, userId))
}

// Auxillary function to update user
func updateFieldIfSet(tx *sql.Tx, userId string, field string, val interface{}) (sql.Result, error) {
	v, ok := val.(string)
	if ok && v != "" {
		return updateField(tx, userId, field, v)
	}
	return nil, nil
}

func (m *mySql) ChangePassword(user_id, newPassword string) error {
	query := fmt.Sprintf(`UPDATE Users SET password = '%v' WHERE user_id = '%v'`, newPassword, user_id)
	_, err := m.conn.Exec(query)
	if err != nil {
		fmt.Println(err)
		return err
	}

	return nil
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

func (m *mySql) AddToken(req *userEntity.ResetPasswordRes) error {
	stmt := fmt.Sprintf(` INSERT INTO Reset_Token(
                   token_id,
                   user_id,
                   token,
                   expiry
                   ) VALUES ('%v', '%v', '%v', '%v')`,
		req.TokenId, req.UserId, req.Token, req.Expiry)

	_, err := m.conn.Exec(stmt)
	if err != nil {
		return err
	}
	return nil
}

func (m *mySql) GetTokenById(token, userId string) (*userEntity.ResetPasswordWithTokenRes, error) {
	query := fmt.Sprintf(`
		SELECT token_id, user_id, token, expiry
		FROM Reset_Token
		WHERE token = '%s'
		AND user_id = '%s'
	`, token, userId)

	var tokenRes userEntity.ResetPasswordWithTokenRes
	ctx := context.Background()
	err := m.conn.QueryRowContext(ctx, query).Scan(
		&tokenRes.TokenId,
		&tokenRes.UserId,
		&tokenRes.Token,
		&tokenRes.Expiry,
	)

	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	return &tokenRes, nil
}

func (m *mySql) DeleteToken(userId string) error {
	query := fmt.Sprintf(`DELETE FROM Reset_Token WHERE user_id = "%s"`, userId)
	_, err := m.conn.Exec(query)
	if err != nil {
		fmt.Println(err)
		return err
	}

	return nil
}
