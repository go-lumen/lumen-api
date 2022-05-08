package controllers

import (
	"encoding/binary"
	"encoding/csv"
	"fmt"
	"github.com/gin-gonic/gin"
	mgobson "github.com/globalsign/mgo/bson"
	"github.com/go-lumen/lumen-api/models"
	"github.com/go-lumen/lumen-api/store"
	"github.com/go-lumen/lumen-api/utils"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

// ImportController holds all controller functions related to the imports
type ImportController struct {
	BaseController
}

// NewImportController instantiates the controller
func NewImportController() ImportController {
	return ImportController{}
}

//id;category;city;country;cover;currency;date_online;description;external_id;is_online;lat;lng;location;metafields;min_qty;price;price_unit;quantity;quantity_selector;sku;supplier_email;title;url;"_Produit couvert par les certifications de l'usine";"_Origine géographique de votre produit";"_Quel est votre lien avec la matière première vendue ? ";"_DLUO ou DLC";"_DLUO / DLC";"_Conventionnel ou BIO";"_Type de Conditionnement";"_Poids du conditionnement (KG)";"_Type de palette";"_Poids palette (KG)";"_Raison de la vente";"_Livraison possible par le vendeur";"_Autre origine géographique :";"_Type de découpe";"_Autre conditionnement";"_Autre certification : ";"_Autre type de palette à préciser";"_Matière du conditionnement";"_Autre matière de conditionnement"

// NewObjectIDWithTime allows to generate an ObjectId from time
func NewObjectIDWithTime(t time.Time) mgobson.ObjectId {
	var b [12]byte
	binary.BigEndian.PutUint32(b[:4], uint32(t.Unix()))
	return mgobson.ObjectId(string(b[:]))
}

// ImportAdsFile allows parsing and create things from an .csv file
func (ic ImportController) ImportAdsFile(c *gin.Context) {
	ctx := store.AuthContext(c)
	if !ic.ShouldBeLogged(ctx) {
		return
	}

	var ret []*models.Ad
	file, err := c.FormFile("file")
	utils.CheckErr(err)
	// Upload the file to specific dst.
	dst := "/var/www/uploads/" + time.Now().Format("2006-01-02-15-04-05") + "-" + file.Filename
	err = c.SaveUploadedFile(file, dst)
	utils.CheckErr(err)
	//fi, _ := file.Open()

	fi, _ := os.Open(dst)
	r := csv.NewReader(fi)
	r.Comma = ';'
	for {
		line, err := r.Read()
		if (err == io.EOF) || (line[0] == "id") {
			break
		}
		t, _ := time.Parse("01/02/2006 15:04:05", line[6])
		fmt.Println(t)
		pri, _ := strconv.ParseFloat(line[15], 64)
		qty, _ := strconv.ParseFloat(line[17], 64)
		expDate, _ := time.Parse("01/02/2006", line[27])
		packWei, _ := strconv.ParseFloat(line[30], 64)
		pallWei, _ := strconv.ParseFloat(line[32], 64)
		user, _ := models.GetOrCreateUser(ctx, &models.User{Email: line[20]})
		ad := &models.Ad{
			ID:                   primitive.NewObjectIDFromTimestamp(t), //fmt.Sprintf("ObjectId(%q)", primitive.NewObjectIDFromTimestamp(t).Hex())
			UserID:               user.ID,
			Name:                 line[21],
			Category:             line[1],
			Place:                line[2],
			Region:               line[3],
			Quantity:             qty,
			Price:                pri,
			Unit:                 line[16],
			Description:          line[7],
			PicturesURLs:         []string{line[4]},
			Certifications:       strings.Split(line[23], "\n"),
			OriginType:           line[24],
			Origin:               "",
			VendorLinkToMaterial: line[25],
			ExpirationDateType:   line[26],
			ExpirationDate:       expDate.Unix(),
			IsOrganic:            line[28] == "BIO",
			PackingType:          line[29],
			Packing:              line[37],
			PackingWeight:        packWei,
			PalletType:           line[31],
			PalletWeight:         pallWei,
			CuttingType:          line[36],
			SellingReason:        line[33],
			DeliveryAvailable:    false,
		}
		models.CreateAd(ctx, ad)
		ret = append(ret, ad)
	}

	/*defer func(fi *os.File) {
		err := fi.Close()
		utils.CheckErr(err)
	}(fi)*/
	c.String(http.StatusOK, fmt.Sprintf("'%s' successfully uploaded", file.Filename))
}
