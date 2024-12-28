package utils

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	errs "gopher_tix/packages/common/errors"
)

func GetUUIDParam(ctx *fiber.Ctx, param string) (uuid.UUID, error) {
	id, err := uuid.Parse(ctx.Params(param))
	if err != nil {
		return uuid.Nil, errs.NewIncorrectParameter(fmt.Sprintf("The provided %s is not valid", param))
	}
	return id, nil
}

func CurrentUserID(ctx *fiber.Ctx) uuid.UUID {
	return ctx.Locals("user_id").(uuid.UUID)
}
