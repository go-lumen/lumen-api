package controllers

import (
	"bytes"
	"encoding/csv"
	"github.com/gin-gonic/gin"
	"github.com/go-lumen/lumen-api/config"
	"github.com/go-lumen/lumen-api/helpers"
	"github.com/go-lumen/lumen-api/models"
	"github.com/go-lumen/lumen-api/store"
	"github.com/go-lumen/lumen-api/utils"
	"github.com/tealeg/xlsx/v3"
	"go.mongodb.org/mongo-driver/bson"
	"net/http"
	"time"
)

// ExportController holds all controller functions related to the export
type ExportController struct {
	BaseController
}

// NewExportController instantiates the controller
func NewExportController() ExportController {
	return ExportController{}
}

func addExcelHeader(cellsContent []string, r *xlsx.Row) {
	var hStyle = xlsx.NewStyle()
	font := *xlsx.NewFont(10, "Arial")
	hStyle.Font.Bold = true
	hStyle.Font = font
	fill := *xlsx.NewFill("solid", "00C8D9EE", "FF000000")
	hStyle.Fill = fill
	border := *xlsx.NewBorder("thin", "thin", "thin", "thin")
	hStyle.Border = border
	hStyle.ApplyBorder = true
	hStyle.ApplyFill = true
	for _, item := range cellsContent {
		cell := r.AddCell()
		cell.SetStyle(hStyle)
		cell.Value = item
	}
	r.SetHeight(20)
}

// ExportDeviceMessages allows to export device messages
func (ec ExportController) ExportDeviceMessages(c *gin.Context) {
	ctx := store.AuthContext(c)
	var queryParams models.QueryParams
	if err := c.ShouldBind(&queryParams); err == nil {
		if queryParams.Limit == 0 {
			queryParams.Limit = 2000
		}
		if queryParams.Order == 0 {
			queryParams.Order = -1
		}
	} else {
		ec.AbortWithError(c, helpers.ErrorInvalidInput(err))
		return
	}
	if queryParams.StartTime < 0 {
		queryParams.StartTime = 0
	}
	if queryParams.EndTime < 100000 {
		queryParams.EndTime = 2000000000
	}

	encodedKey := []byte(config.GetString(c, "rsa_private"))
	claims, err := helpers.ValidateJwtToken(c.Param("token"), encodedKey, "access")
	if err != nil {
		ec.AbortWithError(c, helpers.ErrorInvalidToken(err))
		return
	}

	user, _ := models.GetUser(ctx, bson.M{"_id": claims["sub"].(string)})

	group, err := models.GetGroup(ctx, bson.M{"_id": user.GroupID})
	if err != nil {
		ec.AbortWithError(c, helpers.ErrorResourceNotFound(err))
		return
	}
	if group.Role == store.RoleCustomer {
		ec.AbortWithError(c, helpers.ErrorUserUnauthorized)
		return
	}

	layout := "-messages-export-20060102-150405.csv"
	fileName := time.Now().UTC().Format(layout)
	b := &bytes.Buffer{}   // creates IO Writer
	wr := csv.NewWriter(b) // creates a csv writer that uses the io buffer.

	var header []string
	header = append(header, "User name", "User Email", "Group role", "Group name")
	wr.Write(header)
	//for _, user := range dbUsers {
	var line []string
	line = append(line, user.FirstName+" "+user.LastName)
	line = append(line, user.Email)
	line = append(line, group.Role)
	line = append(line, group.Name)
	wr.Write(line)
	utils.CheckErr(err)
	//}
	wr.Flush() // writes the csv writer data to  the buffered data io writer(b(bytes.buffer))

	c.Header("Content-Description", "File Transfer")
	c.Header("Content-Type", "text/csv")
	c.Header("Content-Disposition", "attachment; filename="+fileName)
	c.Data(http.StatusOK, "application/octet-stream", b.Bytes())
}
