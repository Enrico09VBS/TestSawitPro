package handler

import (
	"errors"
	"net/http"
	"testing"

	"github.com/SawitProRecruitment/UserService/repository"
	"github.com/golang/mock/gomock"
	"github.com/labstack/echo/v4"
)

func TestRegister(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockRepo := repository.NewMockRepositoryInterface(ctrl)
	mockCtx := echo.New().NewContext(&http.Request{}, nil)

	type args struct {
		ctx    echo.Context
		params repository.RegistrationParam
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
		patch   func()
	}{
		{
			name: "invalid params",
			args: args{
				ctx: mockCtx,
				params: repository.RegistrationParam{
					PhoneNumber: "1",
					FullName:    "1",
					Password:    "1",
				},
			},
			wantErr: true,
		},
		{
			name: "success register",
			args: args{
				ctx: mockCtx,
				params: repository.RegistrationParam{
					PhoneNumber: "+6200000000",
					FullName:    "asdfghjkl",
					Password:    "1A*vdhjas",
				},
			},
			wantErr: false,
			patch: func() {
				mockRepo.EXPECT().Register(mockCtx, repository.RegistrationParam{
					PhoneNumber: "+6200000000",
					FullName:    "asdfghjkl",
					Password:    "1A*vdhjas",
				}).Return(int64(1), nil)
			},
		},
		{
			name: "error register",
			args: args{
				ctx: mockCtx,
				params: repository.RegistrationParam{
					PhoneNumber: "+6200000000",
					FullName:    "asdfghjkl",
					Password:    "1A*vdhjas",
				},
			},
			wantErr: false,
			patch: func() {
				mockRepo.EXPECT().Register(mockCtx, repository.RegistrationParam{
					PhoneNumber: "+6200000000",
					FullName:    "asdfghjkl",
					Password:    "1A*vdhjas",
				}).Return(int64(0), errors.New("expected error"))
			},
		},
	}
	for _, tt := range tests {
		tt.patch()
		s := Server{
			Repository: mockRepo,
		}
		t.Run(tt.name, func(t *testing.T) {
			if gotError := s.Register(tt.args.ctx, tt.args.params); (gotError != nil) != tt.wantErr {
				t.Errorf("validateLineItem() = %v, want %v", gotError, tt.wantErr)
			}
		})
	}
}

func TestLogin(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockRepo := repository.NewMockRepositoryInterface(ctrl)
	mockCtx := echo.New().NewContext(&http.Request{}, nil)

	type args struct {
		ctx    echo.Context
		params repository.LoginParam
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
		patch   func()
	}{
		{
			name: "success register",
			args: args{
				ctx: mockCtx,
				params: repository.LoginParam{
					PhoneNumber: "+6200000000",
					Password:    "1A*vdhjas",
				},
			},
			wantErr: false,
			patch: func() {
				mockRepo.EXPECT().Login(mockCtx, repository.LoginParam{
					PhoneNumber: "+6200000000",
					Password:    "1A*vdhjas",
				}).Return(int64(1), nil)

				mockRepo.EXPECT().IncreaseLoginCount(mockCtx, int64(1)).Return(nil)
			},
		},
		{
			name: "error login",
			args: args{
				ctx: mockCtx,
				params: repository.LoginParam{
					PhoneNumber: "+6200000000",
					Password:    "1A*vdhjas",
				},
			},
			wantErr: true,
			patch: func() {
				mockRepo.EXPECT().Login(mockCtx, repository.LoginParam{
					PhoneNumber: "+6200000000",
					Password:    "1A*vdhjas",
				}).Return(int64(0), errors.New("expected error"))
			},
		},
		{
			name: "error increase login count",
			args: args{
				ctx: mockCtx,
				params: repository.LoginParam{
					PhoneNumber: "+6200000000",
					Password:    "1A*vdhjas",
				},
			},
			wantErr: true,
			patch: func() {
				mockRepo.EXPECT().Login(mockCtx, repository.LoginParam{
					PhoneNumber: "+6200000000",
					Password:    "1A*vdhjas",
				}).Return(int64(1), nil)

				mockRepo.EXPECT().IncreaseLoginCount(mockCtx, int64(1)).Return(errors.New("expected error"))
			},
		},
	}
	for _, tt := range tests {
		tt.patch()
		s := Server{
			Repository: mockRepo,
		}
		t.Run(tt.name, func(t *testing.T) {
			if gotError := s.Login(tt.args.ctx, tt.args.params); (gotError != nil) != tt.wantErr {
				t.Errorf("validateLineItem() = %v, want %v", gotError, tt.wantErr)
			}
		})
	}
}
