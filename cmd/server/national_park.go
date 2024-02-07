package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"net/http"
	"tn-rest/internal/db"

	_ "github.com/mattn/go-sqlite3"
)

// nationParkId, err := qtx.CreateNationalPark(np.Ctx, db.CreateNationalParkParams{
// 	Name:                sql.NullString{String: "Taman Nasional Komodo", Valid: true},
// 	Description:         sql.NullString{String: "Taman Nasional Komodo terletak di daerah administrasi Provinsi Nusa Tenggara Timur", Valid: true},
// 	Region:              sql.NullString{String: "Bali dan Nusa Tenggara", Valid: true},
// 	Link:                sql.NullString{String: "https://id.wikipedia.org/wiki/Taman_Nasional_Gunung_Rinjani", Valid: true},
// 	Year:                sql.NullInt64{Int64: 1990, Valid: true},
// 	TotalAreaInKm:       sql.NullInt64{Int64: 413, Valid: true},
// 	TotalAreaInMiles:    sql.NullInt64{Int64: 159, Valid: true},
// 	WaterPercentages:    sql.NullString{},
// 	CoordinateLatitude:  sql.NullFloat64{Float64: -8.408055555555556, Valid: true},
// 	CoordinateLongitude: sql.NullFloat64{Float64: 116.44944444444445, Valid: true},
// 	MapUrl:              sql.NullString{String: "https://id.wikipedia.org/wiki/Berkas:Lombok_Locator_Topography.png", Valid: true},
// 	Location:            sql.NullString{String: "Lombok, Nusa Tenggara Barat, Indonesia", Valid: true},
// 	Established:         sql.NullInt64{Int64: 1990, Valid: true},
// 	Visitors:            sql.NullString{String: "117.715 (tahun 2007[1])", Valid: true},
// 	Management:          sql.NullString{String: "Kementerian Lingkungan Hidup dan Kehutanan", Valid: true},
// })

// if err != nil {
// 	return err
// }

// imageId, err := qtx.CreateImage(np.Ctx, db.CreateImageParams{
// 	Link:   sql.NullString{String: "https://id.wikipedia.org/wiki/Berkas:Rinjani_Caldera.jpg", Valid: true},
// 	Title:  sql.NullString{String: "Rinjani Caldera.jpg", Valid: true},
// 	Date:   sql.NullString{String: "2006-08-31T17:00:00.000Z", Valid: true},
// 	Source: sql.NullString{String: "Karya sendiri", Valid: true},
// 	Author: sql.NullString{String: "Thorsten Peters", Valid: true},
// 	Src:    sql.NullString{String: "https://upload.wikimedia.org/wikipedia/commons/c/cb/Rinjani_Caldera.jpg", Valid: true},
// })

// if err != nil {
// 	return err
// }

// licenseId, err := qtx.CreateLicense(np.Ctx, db.CreateLicenseParams{
// 	Type: sql.NullString{String: "GFDL", Valid: true},
// 	Name: sql.NullString{String: "GNU Free Documentation License", Valid: true},
// 	Link: sql.NullString{String: "http://www.gnu.org/copyleft/fdl.html", Valid: true},
// })

// if err != nil {
// 	return err
// }

// intlStatusId, err := qtx.CreateIntlStatus(np.Ctx, db.CreateIntlStatusParams{
// 	Name: sql.NullString{String: "Situs Warisan Dunia", Valid: true},
// 	Link: sql.NullString{String: "http://whc.unesco.org/en/list/609", Valid: true},
// })

type NewResponse struct {
	Message string `json:"message"`
	Data    any    `json:"data ,omitempty"`
	Error   any    `json:"error ,omitempty"`
}

type NationalPark struct {
	Ctx     context.Context
	Queries *db.Queries
	DB      *sql.DB
}

func (np NationalPark) CreateNationalPark(input CreateNationalParkInput) error {
	tx, err := np.DB.Begin()

	if err != nil {
		return err
	}

	defer tx.Rollback()
	qtx := np.Queries.WithTx(tx)

	nationParkId, err := qtx.CreateNationalPark(np.Ctx, db.CreateNationalParkParams{
		Name:                sql.NullString{String: input.Name, Valid: true},
		Description:         sql.NullString{String: input.Description, Valid: true},
		Region:              sql.NullString{String: input.Region, Valid: true},
		Link:                sql.NullString{String: input.Link, Valid: true},
		Year:                sql.NullInt64{Int64: input.Year, Valid: true},
		TotalAreaInKm:       sql.NullInt64{Int64: input.TotalArea.Km, Valid: true},
		TotalAreaInMiles:    sql.NullInt64{Int64: input.TotalArea.Miles, Valid: true},
		WaterPercentages:    sql.NullString{String: input.WaterPercentages, Valid: len(input.WaterPercentages) != 0},
		CoordinateLatitude:  sql.NullFloat64{Float64: input.Coordinate.Lat, Valid: true},
		CoordinateLongitude: sql.NullFloat64{Float64: input.Coordinate.Long, Valid: true},
		MapUrl:              sql.NullString{String: input.MapUrl, Valid: true},
		Location:            sql.NullString{String: input.Location, Valid: true},
		Established:         sql.NullInt64{Int64: input.EstablishedYear, Valid: true},
		Visitors:            sql.NullString{String: input.Visitors, Valid: true},
		Management:          sql.NullString{String: input.Management, Valid: true},
	})

	if err != nil {
		return err
	}

	imageId, err := qtx.CreateImage(np.Ctx, db.CreateImageParams{
		Link:   sql.NullString{String: input.Link, Valid: true},
		Title:  sql.NullString{String: input.Image.Title, Valid: true},
		Date:   sql.NullString{String: input.Image.Date, Valid: true},
		Source: sql.NullString{String: input.Image.Source, Valid: true},
		Author: sql.NullString{String: input.Image.Author, Valid: true},
		Src:    sql.NullString{String: input.Image.Src, Valid: true},
	})

	if err != nil {
		return err
	}

	licenseId, err := qtx.CreateLicense(np.Ctx, db.CreateLicenseParams{
		Type: sql.NullString{String: input.License.Type, Valid: true},
		Name: sql.NullString{String: input.License.Name, Valid: true},
		Link: sql.NullString{String: input.License.Link, Valid: true},
	})

	if err != nil {
		return err
	}

	intlStatusId, err := qtx.CreateIntlStatus(np.Ctx, db.CreateIntlStatusParams{
		Name: sql.NullString{String: input.IntlStatus.Name, Valid: true},
		Link: sql.NullString{String: input.IntlStatus.Link, Valid: true},
	})

	if err != nil {
		return err
	}

	err = qtx.CreateImageLicense(np.Ctx, db.CreateImageLicenseParams{
		ImageID:   sql.NullInt64{Int64: imageId, Valid: true},
		LicenseID: sql.NullInt64{Int64: licenseId, Valid: true},
	})

	if err != nil {
		return err
	}

	err = qtx.CreateNationalParkImage(np.Ctx, db.CreateNationalParkImageParams{
		NationalParkID: sql.NullInt64{Int64: nationParkId, Valid: true},
		ImageID:        sql.NullInt64{Int64: imageId, Valid: true},
	})

	if err != nil {
		return err
	}

	err = qtx.CreateNationalParkIntlStatus(np.Ctx, db.CreateNationalParkIntlStatusParams{
		NationalParkID: sql.NullInt64{Int64: nationParkId, Valid: true},
		IntlStatusID:   sql.NullInt64{Int64: intlStatusId, Valid: true},
	})

	if err != nil {
		return err
	}

	return tx.Commit()
}

type CreateTotalAreaInput struct {
	Km    int64 `json:"km"`
	Miles int64 `json:"miles"`
}

type CreateCoordinateInput struct {
	Lat  float64 `json:"latitude"`
	Long float64 `json:"longitude"`
}

type CreateIntlStatusInput struct {
	Name string `json:"name"`
	Link string `json:"link"`
}

type CreateLicenseInput struct {
	Type string `json:"type"`
	Name string `json:"name"`
	Link string `json:"link"`
}

type CreateImageInput struct {
	Link   string `json:"link"`
	Title  string `json:"title"`
	Date   string `json:"date"`
	Source string `json:"source"`
	Author string `json:"author"`
	Src    string `json:"src"`
}

type CreateNationalParkInput struct {
	Name             string                `json:"name"`
	Description      string                `json:"description"`
	Region           string                `json:"region"`
	Link             string                `json:"link"`
	Year             int64                 `json:"year"`
	WaterPercentages string                `json:"water_percentages"`
	MapUrl           string                `json:"map_url"`
	Location         string                `json:"location"`
	EstablishedYear  int64                 `json:"established_year"`
	Visitors         string                `json:"visitors"`
	Management       string                `json:"management"`
	TotalArea        CreateTotalAreaInput  `json:"total_area"`
	Coordinate       CreateCoordinateInput `json:"coordinate"`
	License          CreateLicenseInput    `json:"license"`
	Image            CreateImageInput      `json:"image"`
	IntlStatus       CreateIntlStatusInput `json:"intl_status"`
}

type NationalParkHandler struct {
	Ctx     context.Context
	Queries *db.Queries
	DB      *sql.DB
}

func (h NationalParkHandler) Create(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	payload := CreateNationalParkInput{}

	if err := decoder.Decode(&payload); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	np := NationalPark{
		Ctx:     h.Ctx,
		Queries: h.Queries,
		DB:      h.DB,
	}

	if err := np.CreateNationalPark(payload); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	res, err := json.Marshal(NewResponse{Message: "successfully create a national park"})

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(res)
}

func (h NationalParkHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	rows, err := h.Queries.GetNationalParks(h.Ctx)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	data := struct {
		Message string `json:"message"`
		Data    any    `json:"data"`
	}{
		Message: "successfully get national parks",
		Data:    rows,
	}

	res, err := json.Marshal(data)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(res)
}
