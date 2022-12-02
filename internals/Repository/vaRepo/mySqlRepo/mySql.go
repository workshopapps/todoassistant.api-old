package mySqlRepo

import (
	"context"
	"database/sql"
	"fmt"
	"test-va/internals/Repository/vaRepo"
	"test-va/internals/entity/vaEntity"
)

type mySql struct {
	conn *sql.DB
}

func (m *mySql) GetUserAssignedToVa(ctx context.Context, vaId string) ([]*vaEntity.VAStruct, error) {
	stmt := fmt.Sprintf(`
SELECT first_name, last_name, user_id, email, phone, account_status 
FROM Users
WHERE virtual_assistant_id = '%v'`, vaId)

	queryRow, err := m.conn.QueryContext(ctx, stmt)
	if err != nil {
		return nil, err
	}
	var Results []*vaEntity.VAStruct
	for queryRow.Next() {
		var res vaEntity.VAStruct
		err := queryRow.Scan(&res.FirstName, &res.LastName, &res.UserId, &res.Email, &res.Phone, &res.Status)
		if err != nil {
			return nil, err
		}
		Results = append(Results, &res)
	}

	return Results, nil
}

func (m *mySql) Persist(ctx context.Context, req *vaEntity.CreateVAReq) error {
	stmt := fmt.Sprintf(`
	INSERT INTO va_table(va_id,
	                     first_name,
	                     last_name,
	                     email,
	                     phone,
	                     password,
	                     profile_picture,
	                     account_type,
	                     created_at) VALUES ('%v', '%v', '%v', '%v', '%v', '%v', '%v', '%v', '%v')
`, req.VaId, req.FirstName, req.LastName, req.Email, req.Phone, req.Password, req.ProfilePicture, req.AccountType, req.CreatedAt)
	_, err := m.conn.ExecContext(ctx, stmt)
	if err != nil {
		return err
	}
	return nil
}

func (m *mySql) FindById(ctx context.Context, id string) (*vaEntity.FindByIdRes, error) {
	stmt := fmt.Sprintf(` SELECT 
                   va_id,
                        first_name,
                        last_name,
                        email,
                        phone,
                        password,
                        profile_picture,
                        account_type,
                        created_at
FROM va_table where va_id = '%v'`, id)
	var res vaEntity.FindByIdRes
	row := m.conn.QueryRowContext(ctx, stmt)
	err := row.Scan(&res.VaId, &res.FirstName, &res.LastName, &res.Email, &res.Phone,
		&res.Password, &res.ProfilePicture, &res.AccountType, &res.CreatedAt)
	if err != nil {
		return nil, err
	}

	return &res, nil
}

func (m *mySql) FindByEmail(ctx context.Context, email string) (*vaEntity.FindByIdRes, error) {
	stmt := fmt.Sprintf(` SELECT 
						 va_id,
	                     first_name,
	                     last_name,
	                     email,
	                     phone,
	                     password,
	                     profile_picture,
	                     account_type,
	                     created_at
FROM va_table where email = '%v'`, email)
	var res vaEntity.FindByIdRes
	row := m.conn.QueryRowContext(ctx, stmt)
	err := row.Scan(&res.VaId, &res.FirstName, &res.LastName, &res.Email, &res.Phone,
		&res.Password, &res.ProfilePicture, &res.AccountType, &res.CreatedAt)
	if err != nil {
		return nil, err
	}

	return &res, nil
}

func (m *mySql) DeleteUser(ctx context.Context, id string) error {
	stmt := fmt.Sprintf(`DELETE FROM va_table WHERE va_id = '%v'`, id)
	_, err := m.conn.ExecContext(ctx, stmt)
	if err != nil {
		return err
	}
	return nil
}

func (m *mySql) UpdateUser(ctx context.Context, req *vaEntity.EditVaReq, id string) error {
	stmt := fmt.Sprintf(`UPDATE va_table 
SET
			first_name='%v',
			last_name='%v',
			email='%v',
			phone='%v',
			profile_picture='%v'
			WHERE va_id ='%v'
`, req.FirstName, req.LastName, req.Email, req.Phone, req.ProfilePicture, id)
	_, err := m.conn.ExecContext(ctx, stmt)
	if err != nil {
		return err
	}
	return nil
}

func (m *mySql) ChangePassword(ctx context.Context, req *vaEntity.ChangeVAPassword) error {
	stmt := fmt.Sprintf(`UPDATE va_table SET password ='%v' WHERE va_id='%v'`, req.NewPassword, req.VaId)
	_, err := m.conn.ExecContext(ctx, stmt)
	if err != nil {
		return err
	}
	return nil
}

func NewVASqlRepo(conn *sql.DB) vaRepo.VARepo {
	return &mySql{conn: conn}
}
