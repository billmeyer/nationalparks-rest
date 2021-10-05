package http

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"go.opentelemetry.io/otel"
	"nationalparks-rest/pkg"
	"nationalparks-rest/pkg/db"
	"net/http"
	"strconv"
)

var dtb *sql.DB

func SetDB(d *sql.DB) {
	dtb = d
}

func RouteHealthCheck(w http.ResponseWriter, r *http.Request) {
	// Create a child span.
	tracer := otel.GetTracerProvider().Tracer(pkg.TRACER_NAME)
	ctx, span := tracer.Start(r.Context(), "RouteHealthCheck")
	defer span.End()

	respondWithSuccess(ctx, "API is up and running", w)
}

func RouteGetNationalParkById(w http.ResponseWriter, r *http.Request) {
	// Create a child span.
	tracer := otel.GetTracerProvider().Tracer(pkg.TRACER_NAME)
	ctx, span := tracer.Start(r.Context(), "RouteGetNationalParkById")
	defer span.End()

	vars := mux.Vars(r)
	id, _ := strconv.Atoi(vars["id"])

	np, err := db.DBGetNationalParkById(ctx, dtb, id)
	if err != nil {
		respondWithError(ctx, err, w)
	} else {
		respondWithSuccess(ctx, np, w)
	}
}

func RouteGetNationalParkByName(w http.ResponseWriter, r *http.Request) {
	// Create a child span.
	tracer := otel.GetTracerProvider().Tracer(pkg.TRACER_NAME)
	ctx, span := tracer.Start(r.Context(), "RouteGetNationalParkByName")
	defer span.End()

	vars := mux.Vars(r)
	var name = vars["name"]

	np, err := db.DBGetNationalParkByName(ctx, dtb, name)
	if err != nil {
		respondWithError(ctx, err, w)
	} else {
		respondWithSuccess(ctx, np, w)
	}
}

func RouteGetNationalParks(w http.ResponseWriter, r *http.Request) {
	// Create a child span.
	tracer := otel.GetTracerProvider().Tracer(pkg.TRACER_NAME)
	ctx, span := tracer.Start(r.Context(), "RouteGetNationalParks")
	defer span.End()

	var err error
	var nps []db.NationalPark
	var start, count int

	var city = r.FormValue("city")
	var state = r.FormValue("state")
	var zipcode = r.FormValue("zipcode")

	start, err = strconv.Atoi(r.FormValue("start"))
	if err != nil {
		start = 0
	}
	count, err = strconv.Atoi(r.FormValue("count"))
	if err != nil {
		count = 5
	}

	nps, err = db.DBGetNationalParks(ctx, dtb, city, state, zipcode, start, count)
	if err != nil {
		respondWithError(ctx, err, w)
	} else {
		respondWithSuccess(ctx, nps, w)
	}
}

func RouteGetNationalParksByCity(w http.ResponseWriter, r *http.Request) {
	// Create a child span.
	tracer := otel.GetTracerProvider().Tracer(pkg.TRACER_NAME)
	ctx, span := tracer.Start(r.Context(), "RouteGetNationalParksByCity")
	defer span.End()

	var err error
	var nps []db.NationalPark
	var start, count int

	vars := mux.Vars(r)
	var city = vars["city"]
	start, err = strconv.Atoi(vars["start"])
	if err != nil {
		start = 0
	}
	count, err = strconv.Atoi(vars["count"])
	if err != nil {
		count = 5
	}

	nps, err = db.DBGetNationalParksByCity(ctx, dtb, city, start, count)
	if err != nil {
		respondWithError(ctx, err, w)
	} else {
		respondWithSuccess(ctx, nps, w)
	}
}

func RouteGetNationalParksByState(w http.ResponseWriter, r *http.Request) {
	// Create a child span.
	tracer := otel.GetTracerProvider().Tracer(pkg.TRACER_NAME)
	ctx, span := tracer.Start(r.Context(), "RouteGetNationalParksByState")
	defer span.End()

	var err error
	var nps []db.NationalPark
	var start, count int

	vars := mux.Vars(r)
	var state = vars["stateabbr"]
	start, err = strconv.Atoi(vars["start"])
	if err != nil {
		start = 0
	}
	count, err = strconv.Atoi(vars["count"])
	if err != nil {
		count = 5
	}

	nps, err = db.DBGetNationalParksByState(ctx, dtb, state, start, count)
	if err != nil {
		respondWithError(ctx,err, w)
	} else {
		respondWithSuccess(ctx, nps, w)
	}
}

func RouteGetNationalParksByZipCode(w http.ResponseWriter, r *http.Request) {
	// Create a child span.
	tracer := otel.GetTracerProvider().Tracer(pkg.TRACER_NAME)
	ctx, span := tracer.Start(r.Context(), "RouteGetNationalParksByZipCode")
	defer span.End()

	var err error
	var nps []db.NationalPark
	var zipCode, start, count int

	vars := mux.Vars(r)
	zipCode, err = strconv.Atoi(vars["zipcode"])
	if err != nil {
		respondWithError(ctx, fmt.Errorf("bad ZipCode: %s", vars["zipCode"]), w)
	} else {
		start, err = strconv.Atoi(vars["start"])
		if err != nil {
			start = 0
		}
		count, err = strconv.Atoi(vars["count"])
		if err != nil {
			count = 5
		}

		nps, err = db.DBGetNationalParksByZipCode(ctx, dtb, zipCode, start, count)
		if err != nil {
			respondWithError(ctx, err, w)
		} else {
			respondWithSuccess(ctx, nps, w)
		}
	}
}

// Helper functions for respond with 200 or 500 code
func respondWithError(ctx context.Context, err error, w http.ResponseWriter) {
	// Create a child span.
	tracer := otel.GetTracerProvider().Tracer(pkg.TRACER_NAME)
	_, span := tracer.Start(ctx, "respondWithError")
	defer span.End()

	w.WriteHeader(http.StatusInternalServerError)
	json.NewEncoder(w).Encode(err.Error())
}

func respondWithSuccess(ctx context.Context, data interface{}, w http.ResponseWriter) {
	// Create a child span.
	tracer := otel.GetTracerProvider().Tracer(pkg.TRACER_NAME)
	_, span := tracer.Start(ctx, "respondWithSuccess")
	defer span.End()

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(data)
}
