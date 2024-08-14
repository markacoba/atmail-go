package middleware

import (
	"atmail/backend/helpers/response"
	"atmail/backend/model"
	"fmt"
	"net/http"

	"github.com/euroteltr/rbac"
	"github.com/gin-gonic/gin"
)

func getRBACAction(httpMethod string) rbac.Action {
	switch httpMethod {
	case "POST":
		return rbac.Create
	case "PUT":
		return rbac.Update
	case "DELETE":
		return rbac.Delete
	}
	return rbac.Read
}

func RBACMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// RBAC Sample
		// TODO: Roles and Permissions should be saved in the DB
		R := rbac.New(nil)

		serverError := false
		// Users Resource Perms
		usersPerm, err := R.RegisterPermission("users", "User resource", rbac.Create, rbac.Read, rbac.Update, rbac.Delete)
		if err != nil {
			serverError = true
		}

		// Admin Role has access to Users Actions
		adminRole, err := R.RegisterRole("admin", "Admin role")
		if err != nil {
			serverError = true
		}
		if err = R.Permit(adminRole.ID, usersPerm, rbac.Create, rbac.Read, rbac.Update, rbac.Delete); err != nil {
			serverError = true
		}

		// Manager Role can only access Users Read Action
		managerRole, err := R.RegisterRole("manager", "Manager role")
		if err != nil {
			serverError = true
		}
		if err = R.Permit(managerRole.ID, usersPerm, rbac.Read); err != nil {
			serverError = true
		}

		if serverError {
			response.Error(c, http.StatusFailedDependency, fmt.Errorf("server error"))
			c.Abort()
			return
		}

		unauthorized := true
		// TODO: Users should be saved in the DB along with respective roles
		currentUserRole := "admin"

		action := getRBACAction(c.Request.Method)

		switch currentUserRole {
		case model.ROLE_ADMIN:
			if R.IsGranted(adminRole.ID, usersPerm, action) {
				unauthorized = false
			} else {
				break
			}
		case model.ROLE_MANAGER:
			if R.IsGranted(managerRole.ID, usersPerm, action) {
				unauthorized = false
			} else {
				break
			}
		}

		if unauthorized {
			response.Error(c, http.StatusUnauthorized, fmt.Errorf("you don't have access to this resource"))
			c.Abort()
			return
		}

		c.Next()
	}
}
