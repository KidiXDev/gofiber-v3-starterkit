package controllers

import (
	"math"
	"strconv"

	"gofiber-starterkit/app/api/services"
	"gofiber-starterkit/app/api/types"
	"gofiber-starterkit/app/shared"
	"gofiber-starterkit/pkg/middlewares"

	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
)

type UserController struct {
	service *services.UserService
}

func NewUserController(service *services.UserService) *UserController {
	return &UserController{service: service}
}

func (c *UserController) Register(ctx fiber.Ctx) error {
	var req types.RegisterRequest
	if err := middlewares.ValidateBody(ctx, &req); err != nil {
		return err
	}

	user, authResp, err := c.service.Register(ctx.Context(), &req)
	if err != nil {
		return err
	}

	return shared.RespondSuccess(ctx, "Registration successful", fiber.Map{
		"user": user,
		"auth": authResp,
	})
}

func (c *UserController) Login(ctx fiber.Ctx) error {
	var req types.LoginRequest
	if err := middlewares.ValidateBody(ctx, &req); err != nil {
		return err
	}

	user, authResp, err := c.service.Login(ctx.Context(), &req)
	if err != nil {
		return err
	}

	return shared.RespondSuccess(ctx, "Login successful", fiber.Map{
		"user": user,
		"auth": authResp,
	})
}

func (c *UserController) Refresh(ctx fiber.Ctx) error {
	var req types.RefreshTokenRequest
	if err := middlewares.ValidateBody(ctx, &req); err != nil {
		return err
	}

	authResp, err := c.service.RefreshToken(ctx.Context(), req.RefreshToken)
	if err != nil {
		return err
	}

	return shared.RespondSuccess(ctx, "Token refreshed", authResp)
}

func (c *UserController) Me(ctx fiber.Ctx) error {
	userID := ctx.Locals("userID").(uuid.UUID)
	user, err := c.service.GetByID(ctx.Context(), userID)
	if err != nil {
		return err
	}
	return shared.RespondSuccess(ctx, "User retrieved", user)
}

func (c *UserController) UpdateProfile(ctx fiber.Ctx) error {
	userID := ctx.Locals("userID").(uuid.UUID)

	var req types.UpdateProfileRequest
	if err := middlewares.ValidateBody(ctx, &req); err != nil {
		return err
	}

	user, err := c.service.Update(ctx.Context(), userID, &req)
	if err != nil {
		return err
	}

	return shared.RespondSuccess(ctx, "Profile updated", user)
}

func (c *UserController) Logout(ctx fiber.Ctx) error {
	userID := ctx.Locals("userID").(uuid.UUID)
	tokenID := ctx.Locals("tokenID").(string)

	if err := c.service.Logout(ctx.Context(), tokenID, userID.String()); err != nil {
		return err
	}

	return shared.RespondSuccess(ctx, "Logout successful", nil)
}

func (c *UserController) LogoutAll(ctx fiber.Ctx) error {
	userID := ctx.Locals("userID").(uuid.UUID)

	if err := c.service.LogoutAll(ctx.Context(), userID.String()); err != nil {
		return err
	}

	return shared.RespondSuccess(ctx, "Logged out from all devices", nil)
}

func (c *UserController) List(ctx fiber.Ctx) error {
	page, _ := strconv.Atoi(ctx.Query("page", "1"))
	perPage, _ := strconv.Atoi(ctx.Query("per_page", "10"))

	if page < 1 {
		page = 1
	}
	if perPage < 1 || perPage > 100 {
		perPage = 10
	}

	users, total, err := c.service.List(ctx.Context(), page, perPage)
	if err != nil {
		return err
	}

	pages := int(math.Ceil(float64(total) / float64(perPage)))

	return shared.RespondSuccessWithMeta(ctx, "Users retrieved", users, &shared.Metadata{
		Total:   total,
		Page:    page,
		PerPage: perPage,
		Pages:   pages,
	})
}

func (c *UserController) Get(ctx fiber.Ctx) error {
	idParam := ctx.Params("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		return shared.ErrBadRequest("Invalid user ID")
	}

	user, err := c.service.GetByID(ctx.Context(), id)
	if err != nil {
		return err
	}

	return shared.RespondSuccess(ctx, "User retrieved", user)
}

func (c *UserController) Update(ctx fiber.Ctx) error {
	idParam := ctx.Params("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		return shared.ErrBadRequest("Invalid user ID")
	}

	var req types.UpdateProfileRequest
	if err := middlewares.ValidateBody(ctx, &req); err != nil {
		return err
	}

	user, err := c.service.Update(ctx.Context(), id, &req)
	if err != nil {
		return err
	}

	return shared.RespondSuccess(ctx, "User updated", user)
}

func (c *UserController) Delete(ctx fiber.Ctx) error {
	idParam := ctx.Params("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		return shared.ErrBadRequest("Invalid user ID")
	}

	if err := c.service.Delete(ctx.Context(), id); err != nil {
		return err
	}

	return shared.RespondSuccess(ctx, "User deleted", nil)
}
