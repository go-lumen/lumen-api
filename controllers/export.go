package controllers

import (
	"github.com/tealeg/xlsx/v3"
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

/*
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

	device := &models.Device{}
	if len(c.Param("deviceId")) > 20 {
		device, err = models.GetDevice(ctx, bson.M{"_id": c.Param("deviceId")})
	} else {
		device, err = models.GetDevice(ctx, bson.M{"sigfox_id": strings.TrimLeft(c.Param("deviceId"), "0")})
	}
	if err != nil {
		ec.AbortWithError(c, helpers.ErrorResourceNotFound(err))
		return
	}

	order := store.SortDescending
	if int64(queryParams.Order) == 1 {
		order = store.SortAscending
	}
	deviceMessages, err := models.GetDeviceMessages(ctx, bson.M{"device_id": device.ID, "timestamp": bson.M{"$gte": queryParams.StartTime, "$lte": queryParams.EndTime}}, order, int64(queryParams.Limit))

	layout := "-messages-export-20060102-150405.csv"
	fileName := device.SigfoxID + time.Now().UTC().Format(layout)
	b := &bytes.Buffer{}   // creates IO Writer
	wr := csv.NewWriter(b) // creates a csv writer that uses the io buffer.

	var header []string
	header = append(header, "Sequence number", "Date", "Real Date", "Loc source", "Loc status", "Latitude", "Longitude", "Radius", "Event delay", "Event type", "Temperature", "Min temp", "Max temp", "Avg temp", "Humidity", "Min hum", "Max hum", "Avg hum")
	wr.Write(header)
	for _, deviceMess := range deviceMessages {
		var line []string
		line = append(line, strconv.Itoa(int(deviceMess.SequenceNumber)))
		line = append(line, time.Unix(deviceMess.Timestamp, 0).Format("02-01-2006 15:04:05"))
		if deviceMess.Timestamp != 0 {
			line = append(line, time.Unix(deviceMess.Timestamp, 0).Format("02-01-2006 15:04:05"))
		} else {
			line = append(line, "")
		}
		if deviceMess.Location != nil {
			line = append(line, deviceMess.Location.Source)
			line = append(line, deviceMess.Location.Status)
			line = append(line, fmt.Sprintf("%f", deviceMess.Location.Coordinates[1]))
			line = append(line, fmt.Sprintf("%f", deviceMess.Location.Coordinates[0]))
			line = append(line, strconv.Itoa(int(deviceMess.Location.Radius)))
		} else {
			line = append(line, "", "", "", "", "")
		}
		line = append(line, deviceMess.Data.EventType)
		line = append(line, strconv.Itoa(int(deviceMess.Data.Temperature)))
		wr.Write(line)
		utils.CheckErr(err)
	}
	wr.Flush() // writes the csv writer data to  the buffered data io writer(b(bytes.buffer))

	c.Header("Content-Description", "File Transfer")
	c.Header("Content-Type", "text/csv")
	c.Header("Content-Disposition", "attachment; filename="+fileName)
	c.Data(http.StatusOK, "application/octet-stream", b.Bytes())
}
*/
