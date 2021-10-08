package db

import (
	"context"
	"database/sql"
	"go.opentelemetry.io/otel"
	"strings"
)

type NationalPark struct {
	Id           int     `json:"id"`
	LocationNum  string  `json:"location_num"`
	LocationName string  `json:"location_name"`
	Address      string  `json:"address"`
	City         string  `json:"city"`
	State        string  `json:"state"`
	ZipCode      int     `json:"zip_code"`
	PhoneNum     string  `json:"phone_num"`
	FaxNum       string  `json:"fax_num"`
	Latitude     float32 `json:"latitude"`
	Longitude    float32 `json:"longitude"`
}

func DBGetNationalParkById(ctx context.Context, db *sql.DB, id int) (NationalPark, error) {
	// Create a child span.
	tracer := otel.GetTracerProvider().Tracer("")
	_, span := tracer.Start(ctx, "DBGetNationalParkById")
	defer span.End()

	var np = NationalPark{}
	var row = db.QueryRow("SELECT ID, LOCATION_NUM, LOCATION_NAME, ADDRESS, CITY, STATE, ZIP_CODE, PHONE_NUM, FAX_NUM, LATITUDE, LONGITUDE FROM NATIONAL_PARKS WHERE ID=?", id)
	return np, row.Scan(&np.Id, &np.LocationNum, &np.LocationName, &np.Address, &np.City, &np.State, &np.ZipCode, &np.PhoneNum, &np.FaxNum, &np.Latitude, &np.Longitude)
}

func DBGetNationalParkByName(ctx context.Context, db *sql.DB, name string) (NationalPark, error) {
	// Create a child span.
	tracer := otel.GetTracerProvider().Tracer("")
	_, span := tracer.Start(ctx, "DBGetNationalParkByName")
	defer span.End()

	var np = NationalPark{}
	var row = db.QueryRow("SELECT ID, LOCATION_NUM, LOCATION_NAME, ADDRESS, CITY, STATE, ZIP_CODE, PHONE_NUM, FAX_NUM, LATITUDE, LONGITUDE FROM NATIONAL_PARKS WHERE LOCATION_NAME=?", name)
	return np, row.Scan(&np.Id, &np.LocationNum, &np.LocationName, &np.Address, &np.City, &np.State, &np.ZipCode, &np.PhoneNum, &np.FaxNum, &np.Latitude, &np.Longitude)
}

func DBGetNationalParks(ctx context.Context, db *sql.DB, city string, state string, zipcode string, start int, count int) ([]NationalPark, error) {
	// Create a child span.
	tracer := otel.GetTracerProvider().Tracer("")
	_, span := tracer.Start(ctx, "DBGetNationalParks")
	defer span.End()

	if city == "" {
		city = "%"
	}
	city = strings.ToLower(city)

	if state == "" {
		state = "%"
	}
	state = strings.ToLower(state)

	if zipcode == "" {
		zipcode = "%"
	}

	const query = "SELECT ID, LOCATION_NUM, LOCATION_NAME, ADDRESS, CITY, STATE, ZIP_CODE, PHONE_NUM, " +
		"FAX_NUM, LATITUDE, LONGITUDE FROM NATIONAL_PARKS " +
		"WHERE LOWER(CITY) LIKE ? AND LOWER(STATE) LIKE ? AND ZIP_CODE LIKE ? " +
		"LIMIT ? OFFSET ?"

	rows, err := db.Query(query, city, state, zipcode, count, start)

	return processRows(ctx, rows, err)
}

func DBGetNationalParksByCity(ctx context.Context, db *sql.DB, city string, start int, count int) ([]NationalPark, error) {
	// Create a child span.
	tracer := otel.GetTracerProvider().Tracer("")
	newctx, span := tracer.Start(ctx, "DBGetNationalParksByCity")
	defer span.End()

	rows, err := db.Query("SELECT ID, LOCATION_NUM, LOCATION_NAME, ADDRESS, CITY, STATE, ZIP_CODE, PHONE_NUM, FAX_NUM, LATITUDE, LONGITUDE "+
		"FROM NATIONAL_PARKS WHERE CITY = ? LIMIT ? OFFSET ?", city, count, start)

	return processRows(newctx, rows, err)
}

func DBGetNationalParksByState(ctx context.Context, db *sql.DB, state string, start int, count int) ([]NationalPark, error) {
	// Create a child span.
	tracer := otel.GetTracerProvider().Tracer("")
	newctx, span := tracer.Start(ctx, "DBGetNationalParksByState")
	defer span.End()

	rows, err := db.Query("SELECT ID, LOCATION_NUM, LOCATION_NAME, ADDRESS, CITY, STATE, ZIP_CODE, PHONE_NUM, FAX_NUM, LATITUDE, LONGITUDE "+
		"FROM NATIONAL_PARKS WHERE STATE = ? LIMIT ? OFFSET ?", state, count, start)

	return processRows(newctx, rows, err)
}

func DBGetNationalParksByZipCode(ctx context.Context, db *sql.DB, zipCode int, start int, count int) ([]NationalPark, error) {
	// Create a child span.
	tracer := otel.GetTracerProvider().Tracer("")
	newctx, span := tracer.Start(ctx, "DBGetNationalParksByZipCode")
	defer span.End()

	rows, err := db.Query("SELECT ID, LOCATION_NUM, LOCATION_NAME, ADDRESS, CITY, STATE, ZIP_CODE, PHONE_NUM, FAX_NUM, LATITUDE, LONGITUDE "+
		"FROM NATIONAL_PARKS WHERE ZIP_CODE = ? LIMIT ? OFFSET ?", zipCode, count, start)

	return processRows(newctx, rows, err)
}

func processRows(ctx context.Context, rows *sql.Rows, err error) ([]NationalPark, error) {
	// Create a child span.
	tracer := otel.GetTracerProvider().Tracer("")
	_, span := tracer.Start(ctx, "processRows")
	defer span.End()

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	nationalParks := []NationalPark{}
	for rows.Next() {
		var np NationalPark
		if err = rows.Scan(&np.Id, &np.LocationNum, &np.LocationName, &np.Address, &np.City, &np.State, &np.ZipCode, &np.PhoneNum, &np.FaxNum, &np.Latitude, &np.Longitude); err != nil {
			return nil, err
		}
		nationalParks = append(nationalParks, np)
	}

	return nationalParks, nil
}
