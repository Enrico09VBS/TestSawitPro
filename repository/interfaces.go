// This file contains the interfaces for the repository layer.
// The repository layer is responsible for interacting with the database.
// For testing purpose we will generate mock implementations of these
// interfaces using mockgen. See the Makefile for more information.
package repository

import (
	"context"
)

type RepositoryInterface interface {
	Register(ctx context.Context, input RegistrationParam) (insertedID int64, err error)
	Login(ctx context.Context, input LoginParam) (id int64, err error)
	IncreaseLoginCount(ctx context.Context, userID int64) (err error)
	GetMyProfile(ctx context.Context, id int64) (user *UserData, err error)
	UpdateMyProfile(ctx context.Context, params UpdateProfileParam) (err error)
}
