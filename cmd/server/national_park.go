package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"net/http"
	"tn-rest/internal/db"

	_ "github.com/mattn/go-sqlite3"
)

type NewResponse struct {
	Message string `json:"message"`
	Data    any    `json:"data ,omitempty"`
	Error   any    `json:"error ,omitempty"`
}

type NationalPark struct {
	Ctx context.Context
	DB  *sql.DB
}

func (np NationalPark) GetNationalParks() ([]db.GetNationalParksRow, error) {
	queries := db.New(np.DB)

	rows, err := queries.GetNationalParks(np.Ctx)

	if err != nil {
		return nil, err
	}

	return rows, nil
}

type CreateTotalAreaInput struct {
	Km    *int64 `json:"km"`
	Miles *int64 `json:"miles"`
}

type CreateCoordinateInput struct {
	Lat  float64 `json:"latitude"`
	Long float64 `json:"longitude"`
}

type CreateIntlStatusInput struct {
	Name string  `json:"name"`
	Link *string `json:"link"`
}

type CreateLicenseInput struct {
	Type string  `json:"type"`
	Name string  `json:"name"`
	Link *string `json:"link"`
}

type CreateImageInput struct {
	Link   string  `json:"link"`
	Title  string  `json:"title"`
	Date   *string `json:"date"`
	Source *string `json:"source"`
	Author *string `json:"author"`
	Src    *string `json:"src"`
}

type CreateNationalParkInput struct {
	Name             string                `json:"name"`
	Description      string                `json:"description"`
	Region           string                `json:"region"`
	Link             *string               `json:"link"`
	WaterPercentages *string               `json:"water_percentages"`
	MapUrl           *string               `json:"map_url"`
	Location         string                `json:"location"`
	EstablishedYear  int64                 `json:"established_year"`
	Visitors         *string               `json:"visitors"`
	Management       *string               `json:"management"`
	TotalArea        CreateTotalAreaInput  `json:"total_area"`
	Coordinate       CreateCoordinateInput `json:"coordinate"`
	License          CreateLicenseInput    `json:"license"`
	Image            CreateImageInput      `json:"image"`
	IntlStatus       CreateIntlStatusInput `json:"intl_status"`
}

func (np NationalPark) CreateNationalPark(input CreateNationalParkInput) error {
	queries := db.New(np.DB)

	tx, err := np.DB.Begin()

	if err != nil {
		return err
	}

	defer tx.Rollback()
	qtx := queries.WithTx(tx)

	nationParkId, err := qtx.CreateNationalPark(np.Ctx, db.CreateNationalParkParams{
		Name:                input.Name,
		Description:         input.Description,
		Region:              input.Region,
		Link:                input.Link,
		TotalAreaInKm:       input.TotalArea.Km,
		TotalAreaInMiles:    input.TotalArea.Miles,
		WaterPercentages:    input.WaterPercentages,
		CoordinateLatitude:  input.Coordinate.Lat,
		CoordinateLongitude: input.Coordinate.Long,
		MapUrl:              input.MapUrl,
		Location:            input.Location,
		EstablishedYear:     input.EstablishedYear,
		Visitors:            input.Visitors,
		Management:          input.Management,
	})

	if err != nil {
		return err
	}

	imageId, err := qtx.CreateImage(np.Ctx, db.CreateImageParams{
		Link:   input.Image.Link,
		Title:  input.Image.Title,
		Date:   input.Image.Date,
		Source: input.Image.Source,
		Author: input.Image.Author,
		Src:    input.Image.Src,
	})

	if err != nil {
		return err
	}

	licenseId, err := qtx.CreateLicense(np.Ctx, db.CreateLicenseParams{
		Type: input.License.Type,
		Name: input.License.Name,
		Link: input.License.Link,
	})

	if err != nil {
		return err
	}

	intlStatusId, err := qtx.CreateIntlStatus(np.Ctx, db.CreateIntlStatusParams{
		Name: input.IntlStatus.Name,
		Link: input.IntlStatus.Link,
	})

	if err != nil {
		return err
	}

	err = qtx.CreateImageLicense(np.Ctx, db.CreateImageLicenseParams{
		ImageID:   imageId,
		LicenseID: licenseId,
	})

	if err != nil {
		return err
	}

	err = qtx.CreateNationalParkImage(np.Ctx, db.CreateNationalParkImageParams{
		NationalParkID: nationParkId,
		ImageID:        imageId,
	})

	if err != nil {
		return err
	}

	err = qtx.CreateNationalParkIntlStatus(np.Ctx, db.CreateNationalParkIntlStatusParams{
		NationalParkID: nationParkId,
		IntlStatusID:   intlStatusId,
	})

	if err != nil {
		return err
	}

	return tx.Commit()
}

type NationalParkHandler struct {
	Service *NationalPark
}

func (h NationalParkHandler) Create(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	payload := CreateNationalParkInput{}

	if err := decoder.Decode(&payload); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := h.Service.CreateNationalPark(payload); err != nil {
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
	rows, err := h.Service.GetNationalParks()

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
