package controller

import (
	"cine/entity/model"
	"cine/entity/schemas"
	"cine/pkg/fault"
	"cine/server/middleware"
	"cine/service"
	"github.com/MarcusSanchez/go-parse"
	"github.com/MarcusSanchez/go-z"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"net/http"
)

type UserController struct {
	user service.UserService
}

func NewUserController(userService service.UserService) *UserController {
	return &UserController{user: userService}
}

func (uc *UserController) Routes(router fiber.Router, mw *middleware.Middleware) {
	users := router.Group("/users")

	users.Get("/me", mw.SignedIn, uc.GetMe)
	users.Get("/:userID", mw.SignedIn, mw.ParseUUID("userID"), uc.GetUser)

	users.Get("/detailed/me", mw.SignedIn, uc.GetDetailedMe)
	users.Get("/detailed/:userID", mw.SignedIn, mw.ParseUUID("userID"), uc.GetDetailedUser)

	users.Put("/", mw.SignedIn, mw.CSRF, uc.UpdateUser)
	users.Delete("/", mw.SignedIn, mw.CSRF, uc.DeleteUser)

	users.Post("/:userID/follow", mw.SignedIn, mw.CSRF, mw.ParseUUID("userID"), uc.FollowUser)
	users.Delete("/:userID/unfollow", mw.SignedIn, mw.CSRF, mw.ParseUUID("userID"), uc.UnfollowUser)
}

// GetMe [GET] /api/users/me
func (uc *UserController) GetMe(c *fiber.Ctx) error {
	session := c.Locals("session").(*model.Session)

	user, err := uc.user.GetUser(c.Context(), session.UserID)
	if err != nil {
		return err
	}

	return c.Status(http.StatusOK).JSON(fiber.Map{"user": user})
}

// GetUser [GET] /api/users/:userID
func (uc *UserController) GetUser(c *fiber.Ctx) error {
	userID := c.Locals("userID").(uuid.UUID)

	user, err := uc.user.GetUser(c.Context(), userID)
	if err != nil {
		return err
	}

	return c.Status(http.StatusOK).JSON(fiber.Map{"user": user})
}

// GetDetailedMe [GET] /api/users/detailed/me
func (uc *UserController) GetDetailedMe(c *fiber.Ctx) error {
	session := c.Locals("session").(*model.Session)

	user, err := uc.user.GetDetailedUser(c.Context(), session.UserID, uuid.Nil)
	if err != nil {
		return err
	}

	return c.Status(http.StatusOK).JSON(fiber.Map{"detailed_user": user})
}

// GetDetailedUser [GET] /api/users/detailed/:userID
func (uc *UserController) GetDetailedUser(c *fiber.Ctx) error {
	session := c.Locals("session").(*model.Session)
	userID := c.Locals("userID").(uuid.UUID)

	user, err := uc.user.GetDetailedUser(c.Context(), userID, session.UserID)
	if err != nil {
		return err
	}

	return c.Status(http.StatusOK).JSON(fiber.Map{"detailed_user": user})
}

// UpdateUser [PUT] /api/users
func (uc *UserController) UpdateUser(c *fiber.Ctx) error {

	type Payload struct {
		DisplayName    *string `json:"display_name,optional"    z:"display_name" `
		Email          *string `json:"email,optional"           z:"email"`
		Username       *string `json:"username,optional"        z:"username"`
		Password       *string `json:"password,optional"        z:"password"`
		ProfilePicture *string `json:"profile_picture,optional" z:"profile_picture"`
	}

	p, err := parse.JSON[Payload](c.Body())
	if err != nil {
		return fault.BadRequest(err.Error())
	}

	schema := z.Struct{
		"display_name":    schemas.DisplayNameSchema.Optional(),
		"email":           schemas.EmailSchema.Optional(),
		"username":        schemas.UsernameSchema.Optional(),
		"password":        schemas.PasswordSchema.Optional(),
		"profile_picture": schemas.ProfilePictureSchema.Optional(),
	}
	if errs := schema.Validate(p); errs != nil {
		return fault.Validation(errs.One())
	}

	session := c.Locals("session").(*model.Session)

	user, err := uc.user.UpdateUser(c.Context(),
		session.UserID, &model.UserU{
			DisplayName:    p.DisplayName,
			Username:       p.Username,
			Email:          p.Email,
			Password:       p.Password,
			ProfilePicture: p.ProfilePicture,
		},
	)
	if err != nil {
		return err
	}

	return c.Status(http.StatusOK).JSON(fiber.Map{"user": user})
}

// DeleteUser [DELETE] /api/users
func (uc *UserController) DeleteUser(c *fiber.Ctx) error {
	session := c.Locals("session").(*model.Session)

	err := uc.user.DeleteUser(c.Context(), session.UserID)
	if err != nil {
		return err
	}

	return c.SendStatus(http.StatusNoContent)
}

// FollowUser [POST] /api/users/:userID/follow
func (uc *UserController) FollowUser(c *fiber.Ctx) error {
	session := c.Locals("session").(*model.Session)
	userID := c.Locals("userID").(uuid.UUID)

	err := uc.user.FollowUser(c.Context(), session.UserID, userID)
	if err != nil {
		return err
	}

	return c.SendStatus(http.StatusNoContent)
}

// UnfollowUser [POST] /api/users/:userID/unfollow
func (uc *UserController) UnfollowUser(c *fiber.Ctx) error {
	session := c.Locals("session").(*model.Session)
	userID := c.Locals("userID").(uuid.UUID)

	err := uc.user.UnfollowUser(c.Context(), session.UserID, userID)
	if err != nil {
		return err
	}

	return c.SendStatus(http.StatusNoContent)
}
