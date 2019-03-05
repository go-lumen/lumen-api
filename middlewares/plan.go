package middlewares

import (
	"net/http"

	"github.com/adrien3d/base-api/config"
	"github.com/adrien3d/base-api/helpers"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"time"
)

func PlanMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		expired := config.GetBool(c, "plan_expired")
		expiration := config.GetInt(c, "plan_expiration")

		if expiration <= int(time.Now().Unix()) {
			c.AbortWithError(http.StatusUnauthorized, helpers.ErrorWithCode("plan_expired", "The plan is expired, please renew your plan", errors.New("The plan is expired")))
			return
		}

		if expired == true {
			c.AbortWithError(http.StatusUnauthorized, helpers.ErrorWithCode("plan_expired", "The plan is expired, please renew your plan", errors.New("The plan is expired")))
			return
		}

		/*if mailCredit <= 0 {
			c.AbortWithError(http.StatusUnauthorized, helpers.ErrorWithCode("credit_expired", "The mail credit is empty, please buy a new pack", errors.New("Mail credit empty")))
			return
		}
		if textCredit <= 0 {
			c.AbortWithError(http.StatusUnauthorized, helpers.ErrorWithCode("credit_expired", "The text credit is empty, please buy a new pack", errors.New("Text credit empty")))
			return
		}
		if wifiCredit <= 0 {
			c.AbortWithError(http.StatusUnauthorized, helpers.ErrorWithCode("credit_expired", "The wifi resolving credit is empty, please buy a new pack", errors.New("Wifi credit empty")))
			return
		}*/

		c.Next()
	}
}
