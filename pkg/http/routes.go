package http

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
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
	log.Debugf("RouteHealthCheck() called from %s at %s", r.Header.Get("User-Agent"), r.RemoteAddr)

	// Create a child span.
	tracer := otel.GetTracerProvider().Tracer(pkg.TRACER_NAME)
	ctx, span := tracer.Start(r.Context(), "RouteHealthCheck")
	pkg.AddTraceParentToResponse(span, w)
	defer span.End()

	respondWithSuccess(ctx, "API is up and running", w)
}

func RouteGetNationalParkById(w http.ResponseWriter, r *http.Request) {
	log.Debugf("RouteGetNationalParkById() called from %s at %s", r.Header.Get("User-Agent"), r.RemoteAddr)

	// Create a child span.
	tracer := otel.GetTracerProvider().Tracer(pkg.TRACER_NAME)
	ctx, span := tracer.Start(r.Context(), "RouteGetNationalParkById")
	pkg.AddTraceParentToResponse(span, w)
	defer span.End()

	vars := mux.Vars(r)
	id, _ := strconv.Atoi(vars["id"])
	span.SetAttributes(attribute.Int("id", id))

	np, err := db.DBGetNationalParkById(ctx, dtb, id)
	if err != nil {
		respondWithError(ctx, err, w)
	} else {
		respondWithSuccess(ctx, np, w)
	}
}

func RouteGetNationalParkByName(w http.ResponseWriter, r *http.Request) {
	log.Debugf("RouteGetNationalParkByName() called from %s at %s", r.Header.Get("User-Agent"), r.RemoteAddr)

	tracer := otel.GetTracerProvider().Tracer(pkg.TRACER_NAME)
	ctx, span := tracer.Start(r.Context(), "RouteGetNationalParkByName")
	pkg.AddTraceParentToResponse(span, w)
	defer span.End()

	vars := mux.Vars(r)
	var name = vars["name"]
	span.SetAttributes(attribute.String("park-name", name))

	np, err := db.DBGetNationalParkByName(ctx, dtb, name)
	if err != nil {
		respondWithError(ctx, err, w)
	} else {
		respondWithSuccess(ctx, np, w)
	}
}

func RouteGetNationalParks(w http.ResponseWriter, r *http.Request) {
	log.Debugf("RouteGetNationalParks() called from %s at %s", r.Header.Get("User-Agent"), r.RemoteAddr)

	// Create a child span.
	tracer := otel.GetTracerProvider().Tracer(pkg.TRACER_NAME)
	ctx, span := tracer.Start(r.Context(), "RouteGetNationalParks")
	pkg.AddTraceParentToResponse(span, w)
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

	span.SetAttributes(attribute.Int("start", start))
	span.SetAttributes(attribute.Int("count", count))

	nps, err = db.DBGetNationalParks(ctx, dtb, city, state, zipcode, start, count)
	if err != nil {
		respondWithError(ctx, err, w)
	} else {
		respondWithSuccess(ctx, nps, w)
	}
}

func RouteGetNationalParksByCity(w http.ResponseWriter, r *http.Request) {
	log.Debugf("RouteGetNationalParksByCity() called from %s at %s", r.Header.Get("User-Agent"), r.RemoteAddr)

	// Create a child span.
	tracer := otel.GetTracerProvider().Tracer(pkg.TRACER_NAME)
	ctx, span := tracer.Start(r.Context(), "RouteGetNationalParksByCity")
	pkg.AddTraceParentToResponse(span, w)
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

	span.SetAttributes(attribute.Int("start", start))
	span.SetAttributes(attribute.Int("count", count))

	nps, err = db.DBGetNationalParksByCity(ctx, dtb, city, start, count)
	if err != nil {
		respondWithError(ctx, err, w)
	} else {
		respondWithSuccess(ctx, nps, w)
	}
}

func RouteGetNationalParksByState(w http.ResponseWriter, r *http.Request) {
	log.Debugf("RouteGetNationalParksByState() called from %s at %s", r.Header.Get("User-Agent"), r.RemoteAddr)

	// Create a child span.
	tracer := otel.GetTracerProvider().Tracer(pkg.TRACER_NAME)
	ctx, span := tracer.Start(r.Context(), "RouteGetNationalParksByState")
	pkg.AddTraceParentToResponse(span, w)
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

	span.SetAttributes(attribute.Int("start", start))
	span.SetAttributes(attribute.Int("count", count))

	nps, err = db.DBGetNationalParksByState(ctx, dtb, state, start, count)
	if err != nil {
		respondWithError(ctx,err, w)
	} else {
		respondWithSuccess(ctx, nps, w)
	}
}

func RouteGetNationalParksByZipCode(w http.ResponseWriter, r *http.Request) {
	log.Debugf("RouteGetNationalParksByZipCode() called from %s at %s", r.Header.Get("User-Agent"), r.RemoteAddr)

	// Create a child span.
	tracer := otel.GetTracerProvider().Tracer(pkg.TRACER_NAME)
	ctx, span := tracer.Start(r.Context(), "RouteGetNationalParksByZipCode")
	pkg.AddTraceParentToResponse(span, w)
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

		span.SetAttributes(attribute.Int("start", start))
		span.SetAttributes(attribute.Int("count", count))

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
