package repository

import (
	"context"
)

func (r *Repository) Register(ctx context.Context, params RegistrationParam) (insertedID int64, err error) {
	res, err := r.Db.ExecContext(ctx, "INSERT INTO users (full_name, password, phone_number) VALUES ($1, HASHBYTES('SHA2_256', $2), $3)", params.FullName, params.Password, params.PhoneNumber)
	if err != nil {
		return 0, err
	}

	return res.LastInsertId()
}

func (r *Repository) Login(ctx context.Context, params LoginParam) (userID int64, err error) {
	err = r.Db.QueryRowContext(ctx, "SELECT id FROM users where phone_number = $1 AND password = HASHBYTES('SHA2_256', $2)", params.PhoneNumber, params.Password).Scan(&userID)
	if err != nil {
		return 0, err
	}

	return
}

func (r *Repository) IncreaseLoginCount(ctx context.Context, userID int64) (err error) {
	_, err = r.Db.ExecContext(ctx, "UPDATE users SET login_count = login_count + 1 WHERE id = $1", userID)

	return
}

func (r *Repository) GetMyProfile(ctx context.Context, id int64) (user *UserData, err error) {
	err = r.Db.QueryRowContext(ctx, "SELECT phone_number, full_name FROM users where id = $1 AND", id).Scan(&user)
	if err != nil {
		return nil, err
	}

	return
}

func (r *Repository) UpdateMyProfile(ctx context.Context, params UpdateProfileParam) (err error) {
	if params.FullName != "" && params.PhoneNumber != "" {
		_, err = r.Db.ExecContext(ctx, "UPDATE users SET full_name = $1, phone_number = $2 WHERE id = $3", params.FullName, params.PhoneNumber, params.ID)
	} else if params.FullName != "" {
		_, err = r.Db.ExecContext(ctx, "UPDATE users SET full_name = $1 WHERE id = $2", params.FullName, params.ID)
	} else if params.PhoneNumber != "" {
		_, err = r.Db.ExecContext(ctx, "UPDATE users SET phone_number = $1 WHERE id = $2", params.PhoneNumber, params.ID)
	}

	return
}
