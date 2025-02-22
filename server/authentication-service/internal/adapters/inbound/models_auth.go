package inbound

import (
	"github.com/LeonLow97/internal/core/domain"
	pb "github.com/LeonLow97/proto"
)

func ToLoginInput(req *pb.LoginRequest) domain.LoginInput {
	return domain.LoginInput{
		Username: req.Username,
		Password: req.Password,
	}
}

func ToSignUpInput(req *pb.SignUpRequest) domain.SignUpInput {
	return domain.SignUpInput{
		FirstName: req.FirstName,
		LastName:  req.LastName,
		Username:  req.Username,
		Password:  req.Password,
		Email:     req.Email,
	}
}
