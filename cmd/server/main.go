package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"tn-rest/internal/db"

	_ "github.com/mattn/go-sqlite3"
)

func CreateNationalPark(queries *db.Queries, ctx context.Context, dbx *sql.DB) error {
	tx, err := dbx.Begin()
	
	if err != nil {
		return err
	}

	defer tx.Rollback()
	qtx := queries.WithTx(tx)

	nationParkId, err := qtx.CreateNationalPark(ctx, db.CreateNationalParkParams{
		Name: sql.NullString{String: "Taman Nasional Komodo", Valid: true},
		Description: sql.NullString{String: "Taman Nasional Komodo terletak di daerah administrasi Provinsi Nusa Tenggara Timur", Valid: true},
		Region: sql.NullString{String: "Bali dan Nusa Tenggara", Valid: true},
		Link: sql.NullString{String: "https://id.wikipedia.org/wiki/Taman_Nasional_Gunung_Rinjani", Valid: true},
		Year: sql.NullInt64{Int64: 1990, Valid: true},
		TotalAreaInKm: sql.NullInt64{Int64: 413, Valid: true},
		TotalAreaInMiles: sql.NullInt64{Int64: 159, Valid: true},
		WaterPercentages: sql.NullString{},
		CoordinateLatitude: sql.NullFloat64{Float64: -8.408055555555556, Valid: true},
		CoordinateLongitude: sql.NullFloat64{Float64: 116.44944444444445, Valid: true},
		MapUrl: sql.NullString{String: "https://id.wikipedia.org/wiki/Berkas:Lombok_Locator_Topography.png", Valid: true},
		Location: sql.NullString{String: "Lombok, Nusa Tenggara Barat, Indonesia", Valid: true},
		Established: sql.NullInt64{Int64: 1990, Valid: true},
		Visitors: sql.NullString{String: "117.715 (tahun 2007[1])", Valid: true},
		Management: sql.NullString{String: "Kementerian Lingkungan Hidup dan Kehutanan", Valid: true},
	})

	if err != nil {
		return err
	}

	imageId, err := qtx.CreateImage(ctx, db.CreateImageParams{
		Link: sql.NullString{String: "https://id.wikipedia.org/wiki/Berkas:Rinjani_Caldera.jpg", Valid: true},
		Title: sql.NullString{String: "Rinjani Caldera.jpg", Valid: true},
		Date: sql.NullString{String: "2006-08-31T17:00:00.000Z", Valid: true},
		Source: sql.NullString{String: "Karya sendiri", Valid: true},
		Author: sql.NullString{String: "Thorsten Peters", Valid: true},
		Src: sql.NullString{String: "https://upload.wikimedia.org/wikipedia/commons/c/cb/Rinjani_Caldera.jpg", Valid: true},
	})

	if err != nil {
		return err
	}
	
	licenseId, err := qtx.CreateLicense(ctx, db.CreateLicenseParams{
		Type: sql.NullString{String: "GFDL", Valid: true},
		Name: sql.NullString{String: "GNU Free Documentation License", Valid: true},
		Link: sql.NullString{String:  "http://www.gnu.org/copyleft/fdl.html", Valid: true},
	})
	
	if err != nil {
		return err
	}
	
	intlStatusId, err := qtx.CreateIntlStatus(ctx, db.CreateIntlStatusParams{
		Name: sql.NullString{String: "Situs Warisan Dunia", Valid: true},
		Link: sql.NullString{String: "http://whc.unesco.org/en/list/609", Valid: true},
	})

	if err != nil {
		return err
	}

	err = qtx.CreateImageLicense(ctx, db.CreateImageLicenseParams{
		ImageID: sql.NullInt64{Int64: imageId, Valid: true},
		LicenseID: sql.NullInt64{Int64: licenseId, Valid: true},
	})

	if err != nil {
		return err
	}
	
	err = qtx.CreateNationalParkImage(ctx, db.CreateNationalParkImageParams{
		NationalParkID: sql.NullInt64{Int64: nationParkId, Valid: true},
		ImageID: sql.NullInt64{Int64: imageId, Valid: true},
	})
	
	if err != nil {
		return err
	}
	
	err = qtx.CreateNationalParkIntlStatus(ctx, db.CreateNationalParkIntlStatusParams{
		NationalParkID: sql.NullInt64{Int64: nationParkId, Valid: true},
		IntlStatusID: sql.NullInt64{Int64: intlStatusId, Valid: true},
	})

	if err != nil {
		return err
	}

	return tx.Commit()
}

func dump(data interface{}){
	b,_:=json.MarshalIndent(data, "", "  ")
	fmt.Print(string(b))
}

func main () {
	ctx := context.Background()
	dbx, err := sql.Open("sqlite3", "national_park.db")
	
	if err != nil {
		panic(err)
	}

	queries := db.New(dbx)

	// err = CreateNationalPark(queries, ctx, dbx)

	// if err != nil {
	// 	log.Fatal(err)
	// }

	rows, err := queries.GetNationalParks(ctx)

	if err != nil {
		log.Fatal(err)
	}
	
	dump(rows)
}