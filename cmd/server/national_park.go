package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
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

type GetNationalParks struct {
	ID          int64  `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Region      string `json:"region"`
	Images      []struct {
		Id     int64  `json:"id"`
		Link   string `json:"link"`
		Title  string `json:"title"`
		Date   string `json:"date"`
		Source string `json:"source"`
		Author string `json:"author"`
		Src    string `json:"src"`
	} `json:"images"`
	IntlStatuses []struct {
		Id   int     `json:"id"`
		Name string  `json:"name"`
		Link *string `json:"link"`
	} `json:"intlStatuses"`
	TotalArea struct {
		Km    int `json:"km"`
		Miles int `json:"miles"`
	} `json:"totalrea"`
	Coordinate struct {
		Lat  float64 `json:"lat"`
		Long float64 `json:"long"`
	} `json:"coordinate"`
	WaterPercentages *string `json:"waterPercentages"`
	MapUrl           *string `json:"mapUrl"`
	Location         string  `json:"location"`
	EstablishedYear  int     `json:"establishedYear"`
	Visitors         *string `json:"visitors"`
	Management       *string `json:"management"`
}

func (np NationalPark) GetNationalParks() ([]GetNationalParks, error) {
	queries := db.New(np.DB)

	rows, err := queries.GetNationalParks(np.Ctx)

	if err != nil {
		return nil, err
	}

	for i := range rows {
		if err := json.Unmarshal([]byte(fmt.Sprint(rows[i].Images)), &rows[i].Images); err != nil {
			return nil, err
		}
		if err := json.Unmarshal([]byte(fmt.Sprint(rows[i].TotalArea)), &rows[i].TotalArea); err != nil {
			fmt.Println("error", err)
			return nil, err
		}
		if err := json.Unmarshal([]byte(fmt.Sprint(rows[i].IntlStatuses)), &rows[i].IntlStatuses); err != nil {
			fmt.Println("error", err)
			return nil, err
		}
		if err := json.Unmarshal([]byte(fmt.Sprint(rows[i].Coordinate)), &rows[i].Coordinate); err != nil {
			fmt.Println("error", err)
			return nil, err
		}
	}

	rowsByte, err := json.Marshal(rows)
	var list []GetNationalParks
	if err := json.Unmarshal(rowsByte, &list); err != nil {
		return nil, err
	}

	return list, nil
}

type CreateNationalParkInput struct {
	Name             string  `json:"name"`
	Description      string  `json:"description"`
	Region           string  `json:"region"`
	Link             *string `json:"link"`
	WaterPercentages *string `json:"water_percentages"`
	MapUrl           *string `json:"map_url"`
	Location         string  `json:"location"`
	EstablishedYear  int64   `json:"established_year"`
	Visitors         *string `json:"visitors"`
	Management       *string `json:"management"`
	TotalArea        struct {
		Km    *int64 `json:"km"`
		Miles *int64 `json:"miles"`
	} `json:"total_area"`
	Coordinate struct {
		Lat  float64 `json:"latitude"`
		Long float64 `json:"longitude"`
	} `json:"coordinate"`
	License struct {
		Type string  `json:"type"`
		Name string  `json:"name"`
		Link *string `json:"link"`
	} `json:"license"`
	Image struct {
		Link   string  `json:"link"`
		Title  string  `json:"title"`
		Date   *string `json:"date"`
		Source *string `json:"source"`
		Author *string `json:"author"`
		Src    *string `json:"src"`
	} `json:"image"`
	IntlStatus struct {
		Name string  `json:"name"`
		Link *string `json:"link"`
	} `json:"intl_status"`
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
		Message string             `json:"message"`
		Data    []GetNationalParks `json:"data"`
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
