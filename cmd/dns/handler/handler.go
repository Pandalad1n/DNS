package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/Pandalad1n/DNS/internal/drone"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/rs/zerolog/log"
	"io"
	"math"
	"net/http"
	"runtime/debug"
	"strconv"
)

const bodySizeLimit = 1000

type Handler struct {
	router http.Handler
}

func NewHandler(sectorID float64) *Handler {
	r := http.NewServeMux()

	r.HandleFunc("/health", health)
	r.Handle("/metrics", promhttp.Handler())
	r.HandleFunc("/v1/locate", locateDrone(sectorID))
	return &Handler{router: r}
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	defer func() {
		if err := recover(); err != nil {
			log.Ctx(ctx).Error().Str("error", fmt.Sprint(err)).Bytes("stack", debug.Stack()).Msg("ServeHTTP panic.")
			newAPIError(http.StatusInternalServerError, "Internal server error").Write(w)
		}
	}()
	wr := newResponseWriter(w)
	h.router.ServeHTTP(wr, r.WithContext(ctx))
	log.Ctx(ctx).Info().Str("path", r.URL.Path).Int("code", wr.code).Msg("HTTP request served.")
	httpRequests.With(prometheus.Labels{"code": fmt.Sprint(wr.code), "path": r.URL.Path}).Inc()
}

func health(w http.ResponseWriter, _ *http.Request) {
	_, _ = w.Write([]byte("ok"))
}

func locateDrone(sectorID float64) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var drData droneData
		err := json.NewDecoder(io.LimitReader(r.Body, bodySizeLimit)).Decode(&drData)
		if err != nil {
			log.Ctx(r.Context()).Err(err).Msg("Invalid drone data.")
			newAPIError(http.StatusBadRequest, "Invalid drone data").Write(w)
			return
		}
		err = drData.validate()
		if err != nil {
			log.Ctx(r.Context()).Err(err).Msg("Invalid drone data.")
			newAPIError(http.StatusBadRequest, "Invalid drone data").Write(w)
			return
		}
		log.Ctx(r.Context()).Debug().Interface("droneData", drData).Float64("sectorId", sectorID).Msg("Locating Drone.")
		dr, err := newDrone(drData)
		if err != nil {
			log.Ctx(r.Context()).Err(err).Msg("Invalid drone data.")
			newAPIError(http.StatusBadRequest, "Invalid drone data").Write(w)
			return
		}

		type locationData struct {
			Loc float64 `json:"loc"`
		}
		loc := dr.Locate(sectorID)
		loc = math.Round(loc*100) / 100
		apiResponse{Code: http.StatusOK, Body: locationData{Loc: loc}}.Write(w)
	}
}

type responseWriter struct {
	http.ResponseWriter
	code int
}

func newResponseWriter(w http.ResponseWriter) *responseWriter {
	return &responseWriter{ResponseWriter: w, code: http.StatusOK}
}

func (w *responseWriter) WriteHeader(code int) {
	w.code = code
	w.ResponseWriter.WriteHeader(code)
}

var (
	httpRequests = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests",
		},
		[]string{"code", "path"},
	)
)

type apiResponse struct {
	Code int
	Body interface{}
}

func (r apiResponse) Write(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(r.Code)
	_ = json.NewEncoder(w).Encode(r.Body)
}

type errBody struct {
	Msg string `json:"msg"`
}

func newAPIError(code int, msg string) apiResponse {
	return apiResponse{Code: code, Body: errBody{Msg: msg}}
}

func newDrone(dd droneData) (drone.Drone, error) {
	var dr drone.Drone
	parsedX, err := strconv.ParseFloat(dd.X, 64)
	if err != nil {
		return drone.Drone{}, err
	}
	dr.X = parsedX
	parsedY, err := strconv.ParseFloat(dd.Y, 64)
	if err != nil {
		return drone.Drone{}, err
	}
	dr.Y = parsedY
	parsedZ, err := strconv.ParseFloat(dd.Z, 64)
	if err != nil {
		return drone.Drone{}, err
	}
	dr.Z = parsedZ
	parsedVel, err := strconv.ParseFloat(dd.Vel, 64)
	if err != nil {
		return drone.Drone{}, err
	}
	dr.Vel = parsedVel
	return dr, nil
}

type droneData struct {
	X   string `json:"x"`
	Y   string `json:"y"`
	Z   string `json:"z"`
	Vel string `json:"vel"`
}

func (dd droneData) validate() error {
	if dd.X == "" {
		return errors.New("x cord not set")
	}
	if dd.Y == "" {
		return errors.New("y cord not set")
	}
	if dd.Z == "" {
		return errors.New("z cord not set")
	}
	if dd.Vel == "" {
		return errors.New("vel param not set")
	}
	return nil
}
