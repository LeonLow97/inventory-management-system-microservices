package inbound

import (
	"github.com/LeonLow97/internal/core/domain"
	"github.com/LeonLow97/internal/pkg/utils"
	pb "github.com/LeonLow97/proto"
)

func ToUpdateUserInput(req *pb.UpdateUserRequest) domain.UpdateUserInput {
	return domain.UpdateUserInput{
		FirstName: utils.SanitizePointer(&req.FirstName),
		LastName:  utils.SanitizePointer(&req.LastName),
		Password:  utils.SanitizePointer(&req.Password),
	}
}
