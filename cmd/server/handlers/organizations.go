package handlers

import (
	"context"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"github.com/cvcio/mediawatch/models/deprecated/account"
	"github.com/cvcio/mediawatch/models/organization"
	"github.com/cvcio/mediawatch/pkg/auth"
	"github.com/cvcio/mediawatch/pkg/db"
	mailer "github.com/cvcio/mediawatch/pkg/mailer/v1"
	"github.com/cvcio/mediawatch/pkg/web"
	"go.opencensus.io/trace"
)

// OrganizationsHandler is the handler struct for organizations' related endpoints
type OrganizationsHandler struct {
	log         *zap.Logger
	mongoDB     *db.MongoDB
	mailService *mailer.Mailer
}

// NewOrganizationsHandler creates a new organization handler
func NewOrganizationsHandler(log *zap.Logger, mongoDB *db.MongoDB, mailService *mailer.Mailer) (*OrganizationsHandler, error) {
	ctx, span := trace.StartSpan(context.Background(), "organization.EnsureIndex")
	defer span.End()

	if err := organization.EnsureIndex(ctx, mongoDB); err != nil {
		return nil, err
	}

	return &OrganizationsHandler{
		log:         log,
		mongoDB:     mongoDB,
		mailService: mailService,
	}, nil
}

// List returns all the existing user in the system.
func (u *OrganizationsHandler) List(c *gin.Context) {
	ctx, span := trace.StartSpan(c.Request.Context(), "handlers.OrganizationsHandler.List")
	defer span.End()

	// Get claims from request
	claims, ok := ctx.Value(auth.Key).(auth.Claims)
	if !ok {
		c.JSON(http.StatusUnauthorized, web.ErrUnauthorized)
		return
	}
	// If you are not an org or admin and looking to retrieve someone else then you are rejected.
	if !claims.HasRole(auth.RoleAdmin) && !claims.HasRole(auth.RolePowerUser) {
		c.JSON(http.StatusUnauthorized, web.ErrUnauthorized)
		return
	}

	// Create Options for the list query
	opts := db.NewListOpts()
	l, err := strconv.Atoi(c.Query("limit"))
	if err == nil {
		opts = append(opts, db.Limit(l))
	}

	o, err := strconv.Atoi(c.Query("offset"))
	if err == nil {
		opts = append(opts, db.Offset(o))
	}

	data, err := organization.List(ctx, u.mongoDB, opts...)
	if err != nil {
		u.log.Sugar().Debug(err)
		c.JSON(http.StatusInternalServerError, web.ErrInternalError)
		return
	}

	c.JSON(http.StatusOK, web.NewOrganizationListResponse(data))
}

// Get an existing organization by an id.
func (u *OrganizationsHandler) Get(c *gin.Context) {
	ctx, span := trace.StartSpan(c.Request.Context(), "handlers.OrganizationsHandler.Get")
	defer span.End()

	id := c.Param("id")
	if !db.Valid(id) {
		c.JSON(http.StatusBadRequest, web.ErrInvalidID)
		return
	}

	claims, ok := ctx.Value(auth.Key).(auth.Claims)
	if !ok {
		c.JSON(http.StatusUnauthorized, web.ErrUnauthorized)
		return
	}

	// If you are not an org or admin and looking to retrieve someone else then you are rejected.
	if !claims.HasRole(auth.RoleAdmin) && claims.Organization != id {
		c.JSON(http.StatusUnauthorized, web.ErrUnauthorized)
		return
	}

	data, err := organization.Get(ctx, u.mongoDB, id)
	if err != nil {
		u.log.Sugar().Debug(err)
		c.JSON(http.StatusInternalServerError, web.ErrInternalError)
		return
	}

	c.JSON(http.StatusOK, web.NewOrganizationResponse(data))
}

// Create a new organization.
func (u *OrganizationsHandler) Create(c *gin.Context) {
	// setup tracing TODO
	ctx, span := trace.StartSpan(c.Request.Context(), "handlers.organizations.Create")
	defer span.End()

	claims, ok := ctx.Value(auth.Key).(auth.Claims)
	if !ok {
		c.JSON(http.StatusInternalServerError, web.ErrInternalError)
		return
	}

	oldUser, err := account.Get(ctx, u.mongoDB, claims.Subject)
	if err != nil {
		u.log.Sugar().Debug(err)
		c.JSON(http.StatusInternalServerError, web.ErrInternalError)
		return
	}

	if oldUser.Organization != "" {
		u.log.Debug("Error creating org, user already have one")
		c.JSON(http.StatusUnauthorized, web.ErrUnauthorized)
		return
	}

	var newOrganization organization.Organization
	if err := c.ShouldBindJSON(&newOrganization); err != nil {
		u.log.Sugar().Debug(err)
		c.JSON(http.StatusBadRequest, web.ErrInvalidRequest(err))
		return
	}

	created, err := organization.Create(ctx, u.mongoDB, &newOrganization, time.Now())
	if err != nil {
		u.log.Sugar().Debug(err)
		c.JSON(http.StatusBadRequest, web.ErrInvalidRequest(err))
		return
	}

	orgID := created.ID.Hex()
	roles := append(oldUser.Roles, auth.RoleOrgAdmin)
	updUser := account.UpdateAccount{
		Roles:        &roles,
		Organization: &orgID,
	}

	err = account.Update(ctx, u.mongoDB, oldUser.ID.Hex(), &updUser, time.Now())
	if err != nil {
		u.log.Sugar().Debug(err)
		c.JSON(http.StatusBadRequest, web.ErrInvalidRequest(err))
		return
	}

	c.JSON(http.StatusOK, web.NewOrganizationResponse(created))
}

// Delete an existing organization from an id.
func (u *OrganizationsHandler) Delete(c *gin.Context) {
	// setup tracing TODO
	ctx, span := trace.StartSpan(c.Request.Context(), "handlers.OrganizationsHandler.Delete")
	defer span.End()

	// Get id
	id := c.Param("id")
	if !db.Valid(id) {
		c.JSON(http.StatusBadRequest, web.ErrInvalidID)
		return
	}

	claims, ok := ctx.Value(auth.Key).(auth.Claims)
	if !ok {
		c.JSON(http.StatusUnauthorized, web.ErrUnauthorized)
		return
	}
	// If you are not an admin and looking to retrieve someone else then you are rejected.
	if !claims.HasRole(auth.RoleAdmin) {
		c.JSON(http.StatusUnauthorized, web.ErrUnauthorized)
		return
	}
	// Delete id

	err := organization.Delete(ctx, u.mongoDB, id)
	if err != nil {
		u.log.Sugar().Debug(err)
		c.JSON(http.StatusInternalServerError, web.ErrInternalError)
		return
	}

	c.JSON(http.StatusOK, web.Deleted)
}

// Update an existing organization by an id.
func (u *OrganizationsHandler) Update(c *gin.Context) {
	ctx, span := trace.StartSpan(c.Request.Context(), "handlers.OrganizationsHandler.Update")
	defer span.End()

	id := c.Param("id")
	if !db.Valid(id) {
		c.JSON(http.StatusBadRequest, web.ErrInvalidID)
		return
	}

	claims, ok := ctx.Value(auth.Key).(auth.Claims)
	if !ok {
		c.JSON(http.StatusUnauthorized, web.ErrUnauthorized)
		return
	}

	// If you are not an admin and looking to retrieve someone else then you are rejected.
	if !claims.HasRole(auth.RoleAdmin) && claims.Organization != id {
		c.JSON(http.StatusUnauthorized, web.ErrUnauthorized)
		return
	}
	// if not org admin, reject
	if !claims.HasRole(auth.RoleOrgAdmin) {
		c.JSON(http.StatusUnauthorized, web.ErrUnauthorized)
		return
	}

	var updOrg organization.UpdOrg
	if err := c.ShouldBindJSON(&updOrg); err != nil {
		u.log.Sugar().Debug(err)
		c.JSON(http.StatusBadRequest, web.ErrInvalidRequest(err))
		return
	}

	err := organization.Update(ctx, u.mongoDB, id, &updOrg, time.Now())
	if err != nil {
		u.log.Sugar().Debug(err)
		c.JSON(http.StatusBadRequest, web.ErrInvalidRequest(err))
		return
	}

	c.JSON(http.StatusOK, web.Updated)
}

// UpsertMember adds a new member to an organization.
func (u *OrganizationsHandler) UpsertMember(c *gin.Context) {
	ctx, span := trace.StartSpan(c.Request.Context(), "handlers.OrganizationsHandler.UpsertMember")
	defer span.End()

	id := c.Param("id")
	if !db.Valid(id) {
		c.JSON(http.StatusBadRequest, web.ErrInvalidID)
		return
	}

	claims, ok := ctx.Value(auth.Key).(auth.Claims)
	if !ok {
		c.JSON(http.StatusUnauthorized, web.ErrUnauthorized)
		return
	}

	// If you are not an admin and looking to retrieve someone else then you are rejected.
	if !claims.HasRole(auth.RoleAdmin) && claims.Organization != id {
		c.JSON(http.StatusUnauthorized, web.ErrUnauthorized)
		return
	}
	// if not org admin, reject
	if !claims.HasRole(auth.RoleOrgAdmin) {
		c.JSON(http.StatusUnauthorized, web.ErrUnauthorized)
		return
	}
	// parse request body
	var member organization.Member
	if err := c.ShouldBindJSON(&member); err != nil {
		u.log.Sugar().Debug(err)
		c.JSON(http.StatusBadRequest, web.ErrInvalidRequest(err))
		return
	}
	member.Email = strings.ToLower(member.Email)

	// get the organization
	org, err := organization.Get(ctx, u.mongoDB, id)
	if err != nil {
		c.JSON(http.StatusUnauthorized, web.ErrUnauthorized)
		return
	}

	// get the inviter
	inviter, err := account.ByEmail(ctx, u.mongoDB, claims.User)
	if err != nil {
		c.JSON(http.StatusUnauthorized, web.ErrUnauthorized)
		return
	}

	member.Nonce = account.OTP(4)
	member.Status = "pending"

	// check whether member exists
	exists, err := account.ByEmail(ctx, u.mongoDB, member.Email)
	if err != nil {
		u.log.Sugar().Debug(err)
		c.JSON(http.StatusNotFound, web.ErrNotFound)
		return
	}

	// notify user and create accept request
	member.ID = exists.ID.Hex()
	err = mailer.SendInviteExistingUser(ctx, u.mailService, member.Email, inviter.FirstName, inviter.LastName, inviter.Email, member.Nonce, org.Name, id, member.ID)
	if err != nil {
		u.log.Sugar().Debug(err)
		return
	}

	if err := organization.InsertMember(ctx, u.mongoDB, id, &member); err != nil {
		u.log.Sugar().Debug(err)
		c.JSON(http.StatusBadRequest, web.ErrInvalidRequest(err))
		return
	}

	c.JSON(http.StatusOK, member)
}

// RemoveMember removes a member from an organization.
func (u *OrganizationsHandler) RemoveMember(c *gin.Context) {
	// setup tracing TODO
	ctx, span := trace.StartSpan(c.Request.Context(), "handlers.OrganizationsHandler.RemoveMember")
	defer span.End()

	id := c.Param("id")
	if !db.Valid(id) {
		c.JSON(http.StatusBadRequest, web.ErrInvalidID)
		return
	}

	memberId := c.Param("memberId")
	if !db.Valid(memberId) {
		c.JSON(http.StatusBadRequest, web.ErrInvalidID)
		return
	}

	claims, ok := ctx.Value(auth.Key).(auth.Claims)
	if !ok {
		c.JSON(http.StatusUnauthorized, web.ErrUnauthorized)
		return
	}

	// If you are not an admin and looking to retrieve someone else then you are rejected.
	if !claims.HasRole(auth.RoleAdmin) && claims.Organization != id {
		c.JSON(http.StatusUnauthorized, web.ErrUnauthorized)
		return
	}
	// if not org admin, reject
	if !claims.HasRole(auth.RoleOrgAdmin) {
		c.JSON(http.StatusUnauthorized, web.ErrUnauthorized)
		return
	}

	if err := organization.RemoveMember(ctx, u.mongoDB, id, memberId); err != nil {
		u.log.Sugar().Debug(err)
		c.JSON(http.StatusBadRequest, web.ErrInvalidRequest(err))
		return
	}

	c.JSON(http.StatusOK, web.Deleted)
}

// ProcessInvitation processes an invitation to an organization.
func (u *OrganizationsHandler) ProcessInvitation(c *gin.Context) {
	ctx, span := trace.StartSpan(c.Request.Context(), "handlers.OrganizationsHandler.ProcessInvitation")
	defer span.End()

	id := c.Param("id")
	if !db.Valid(id) {
		c.JSON(http.StatusBadRequest, web.ErrInvalidID)
		return
	}

	memberId := c.Param("memberId")
	if !db.Valid(memberId) {
		c.JSON(http.StatusBadRequest, web.ErrInvalidID)
		return
	}

	method := c.Param("method")

	// parse request body
	var member organization.Member
	if err := c.ShouldBindJSON(&member); err != nil {
		u.log.Sugar().Debug(err)
		c.JSON(http.StatusBadRequest, web.ErrInvalidRequest(err))
		return
	}

	// get the organization
	org, err := organization.Get(ctx, u.mongoDB, id)
	if err != nil {
		c.JSON(http.StatusUnauthorized, web.ErrUnauthorized)
		return
	}

	var found *organization.Member
	for _, v := range org.Members {
		if v.ID == memberId {
			found = v
		}
	}

	if found == nil {
		c.JSON(http.StatusNotFound, web.ErrNotFound)
		return
	}

	if member.Nonce != found.Nonce {
		c.JSON(http.StatusBadRequest, web.ErrInvalidID)
		return
	}

	if found.Status == "accepted" {
		c.JSON(http.StatusOK, web.Updated)
		return
	}

	// check whether member exists
	exists, err := account.Get(ctx, u.mongoDB, memberId)
	if err != nil {
		u.log.Sugar().Debug(err)
		c.JSON(http.StatusNotFound, web.ErrNotFound)
		return
	}

	if method == "accept" {
		exists.Organization = id
		exists.Roles = append(exists.Roles, found.Role)
		if err := account.Update(ctx, u.mongoDB, exists.ID.Hex(), &account.UpdateAccount{
			Organization: &id,
			Roles:        &exists.Roles,
		}, time.Now().Truncate(time.Millisecond)); err != nil {
			u.log.Sugar().Debug(err)
			c.JSON(http.StatusInternalServerError, web.ErrRender(err))
			return
		}

		found.Nonce = ""
		found.Status = "accepted"
		if err := organization.UpdateMember(ctx, u.mongoDB, id, found); err != nil {
			u.log.Sugar().Debug(err)
			c.JSON(http.StatusInternalServerError, web.ErrRender(err))
			return
		}
	} else if method == "reject" {
		found.Status = "rejected"
		if err := organization.UpdateMember(ctx, u.mongoDB, id, found); err != nil {
			u.log.Sugar().Debug(err)
			c.JSON(http.StatusInternalServerError, web.ErrRender(err))
			return
		}
	}

	c.JSON(http.StatusOK, web.Updated)
}
