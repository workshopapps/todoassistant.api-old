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
