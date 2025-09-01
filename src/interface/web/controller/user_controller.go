package controller

import (
	"net/http"
	"strconv"
	"users_api/src/usecase"

	"github.com/gin-gonic/gin"
)

type UserController struct {
	searchUC *usecase.UserSearchUsecase
}

func NewUserController(su *usecase.UserSearchUsecase) *UserController {
	return &UserController{searchUC: su}
}

func (uc *UserController) Search(c *gin.Context) {
	tenantID := c.Param("tenant_id")
	
	in := usecase.UserSearchIn{
		TenantID: tenantID,
		UserName: c.Query("user_name"),
		Email:    c.Query("email"),
		Limit:    parseInt(c.Query("limit"), 20),
		Offset:   parseInt(c.Query("offset"), 0),
	}
	
	users, total, err := uc.searchUC.Do(c, in)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{"items": users, "total": total})
}

func parseInt(s string, defaultVal int) int {
	if s == "" {
		return defaultVal
	}
	if val, err := strconv.Atoi(s); err == nil {
		return val
	}
	return defaultVal
}