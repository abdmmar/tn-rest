package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"tn-rest/internal/db"
	"tn-rest/internal/sqlc"
)

func CreateNationalPark(queries *db.Queries, ctx context.Context, dbx *sql.DB) error {
	tx, err := dbx.Begin()
	
	if err != nil {
		return err
	}

	defer tx.Rollback()
	qtx := queries.WithTx(tx)

	nationParkId, err := qtx.CreateNationalPark(ctx, db.CreateNationalParkParams{
		Name: sql.NullString{String: "Taman Nasional Komodo"},
		Description: sql.NullString{String: "Taman Nasional Komodo terletak di daerah administrasi Provinsi Nusa Tenggara Timur"},
		Region: sql.NullString{String: "Bali dan Nusa Tenggara"},
		Link: sql.NullString{String: "https://id.wikipedia.org/wiki/Taman_Nasional_Gunung_Rinjani"},
		Year: sql.NullInt64{Int64: 1990},
		TotalAreaInKm: sql.NullInt64{Int64: 413},
		TotalAreaInMiles: sql.NullInt64{Int64: 159},
		WaterPercentages: sql.NullString{},
		CoordinateLatitude: sql.NullFloat64{Float64: -8.408055555555556},
		CoordinateLongitude: sql.NullFloat64{Float64: 116.44944444444445},
		MapUrl: sql.NullString{String: "https://id.wikipedia.org/wiki/Berkas:Lombok_Locator_Topography.png"},
		Location: sql.NullString{String: "Lombok, Nusa Tenggara Barat, Indonesia"},
		Established: sql.NullInt64{Int64: 1990},
		Visitors: sql.NullString{String: "117.715 (tahun 2007[1])"},
		Management: sql.NullString{String: "Kementerian Lingkungan Hidup dan Kehutanan"},
	})

	if err != nil {
		return err
	}

	imageId, err := qtx.CreateImage(ctx, db.CreateImageParams{
		Link: sql.NullString{String: "https://id.wikipedia.org/wiki/Berkas:Rinjani_Caldera.jpg"},
		Title: sql.NullString{String: "Rinjani Caldera.jpg"},
		Date: sql.NullString{String: "2006-08-31T17:00:00.000Z"},
		Source: sql.NullString{String: "Karya sendiri"},
		Author: sql.NullString{String: "Thorsten Peters"},
		Src: sql.NullString{String: "https://upload.wikimedia.org/wikipedia/commons/c/cb/Rinjani_Caldera.jpg"},
	})

	if err != nil {
		return err
	}
	
	licenseId, err := qtx.CreateLicense(ctx, db.CreateLicenseParams{
		Type: sql.NullString{String: "GFDL"},
		Name: sql.NullString{String: "GNU Free Documentation License"},
		Link: sql.NullString{String:  "http://www.gnu.org/copyleft/fdl.html"},
	})
	
	if err != nil {
		return err
	}
	
	intlStatusId, err := qtx.CreateIntlStatus(ctx, db.CreateIntlStatusParams{
		Name: sql.NullString{String: "Situs Warisan Dunia"},
		Link: sql.NullString{String: "http://whc.unesco.org/en/list/609"},
	})

	if err != nil {
		return err
	}

	err = qtx.CreateImageLicense(ctx, db.CreateImageLicenseParams{
		ImageID: sql.NullInt64{Int64: imageId},
		LicenseID: sql.NullInt64{Int64: licenseId},
	})

	if err != nil {
		return err
	}
	
	err = qtx.CreateNationalParkImage(ctx, db.CreateNationalParkImageParams{
		NationalParkID: sql.NullInt64{Int64: nationParkId},
		ImageID: sql.NullInt64{Int64: imageId},
	})

	if err != nil {
		return err
	}
	
	err = qtx.CreateNationalParkIntlStatus(ctx, db.CreateNationalParkIntlStatusParams{
		NationalParkID: sql.NullInt64{Int64: nationParkId},
		IntlStatusID: sql.NullInt64{Int64: intlStatusId},
	})

	if err != nil {
		return err
	}

	return tx.Commit()
}

func main () {
	ctx := context.Background()
	dbx, err := sql.Open("sqlite3", "db.sqlite3")
	if err != nil {
		fmt.Println(err)
		panic(err)
	}

	if _, err = dbx.ExecContext(ctx, sqlc.DDL); err != nil {
		panic(err)
	}

	queries := db.New(dbx)

	err = CreateNationalPark(queries, ctx, dbx)

	if err != nil {
		log.Fatal(err)
	}

	queries.GetNationalParks(ctx)
}