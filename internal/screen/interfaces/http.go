package interfaces

import (
	"encoding/json"
	"net/http"

	pb "github.com/zaeemnadeem/golang-example-repo/api/proto/v1/content"
	"github.com/zaeemnadeem/golang-example-repo/internal/screen/app"
	"go.uber.org/zap"
)

type HttpHandler struct {
	app         *app.ScreenService
	contentGrpc pb.ContentServiceClient
	logger      *zap.Logger
}

func NewHttpHandler(app *app.ScreenService, contentGrpc pb.ContentServiceClient, logger *zap.Logger) *HttpHandler {
	return &HttpHandler{
		app:         app,
		contentGrpc: contentGrpc,
		logger:      logger,
	}
}

func (h *HttpHandler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /health", h.handleHealth)
	mux.HandleFunc("POST /screens", h.handleCreateScreen)
	mux.HandleFunc("GET /screens/{id}", h.handleGetScreen)
	mux.HandleFunc("PUT /screens/{id}/status", h.handleUpdateStatus)
	mux.HandleFunc("GET /screens/{id}/content", h.handleGetScreenContent)
	mux.HandleFunc("POST /screens/{id}/content", h.handleCreateAndAssignContent)
}

func (h *HttpHandler) handleHealth(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}

type CreateScreenRequest struct {
	Name     string `json:"name" example:"Lobby Display"`
	Location string `json:"location" example:"Main Entrance"`
}

// @Summary      Create a new screen
// @Description  Provisions a new digital signage screen in the system.
// @Tags         Screens
// @Accept       json
// @Produce      json
// @Param        request body CreateScreenRequest true "Screen creation info"
// @Success      201  {object}  domain.Screen
// @Failure      400  {string}  string "invalid request"
// @Failure      500  {string}  string "internal server error"
// @Router       /screens [post]
func (h *HttpHandler) handleCreateScreen(w http.ResponseWriter, r *http.Request) {
	var req CreateScreenRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}

	screen, err := h.app.CreateScreen(r.Context(), req.Name, req.Location)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(screen)
}

// @Summary      Get screen by ID
// @Description  Retrieves the metadata of a specific screen.
// @Tags         Screens
// @Produce      json
// @Param        id   path      string  true  "UUID of the screen"
// @Success      200  {object}  domain.Screen
// @Failure      404  {string}  string "Screen not found"
// @Router       /screens/{id} [get]
func (h *HttpHandler) handleGetScreen(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	screen, err := h.app.GetScreen(r.Context(), id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(screen)
}

type UpdateStatusRequest struct {
	Status string `json:"status" example:"ONLINE"`
}

// @Summary      Update screen status
// @Description  Updates the operational status of a screen.
// @Tags         Screens
// @Accept       json
// @Produce      json
// @Param        id       path      string               true  "UUID of the screen"
// @Param        request  body      UpdateStatusRequest  true  "New status"
// @Success      200      {object}  domain.Screen
// @Failure      400      {string}  string "invalid request"
// @Failure      500      {string}  string "internal server error"
// @Router       /screens/{id}/status [put]
func (h *HttpHandler) handleUpdateStatus(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	var req UpdateStatusRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}

	screen, err := h.app.UpdateScreenStatus(r.Context(), id, req.Status)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(screen)
}

// @Summary      Get content assigned to screen
// @Description  Retrieves the currently assigned media contents for a given screen from the internal microservice.
// @Tags         Content
// @Produce      json
// @Param        id   path      string  true  "UUID of the screen"
// @Success      200  {array}   domain.Content
// @Failure      502  {string}  string "failed to fetch content internally"
// @Router       /screens/{id}/content [get]
func (h *HttpHandler) handleGetScreenContent(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	// Call the Content Service internally via gRPC!
	resp, err := h.contentGrpc.GetScreenContent(r.Context(), &pb.GetScreenContentRequest{
		ScreenId: id,
	})
	if err != nil {
		h.logger.Error("Failed to call content-service", zap.Error(err))
		http.Error(w, "failed to fetch content internally", http.StatusBadGateway)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp.Contents)
}

type CreateAndAssignContentRequest struct {
	Title           string `json:"title" example:"Welcome Video"`
	Url             string `json:"url" example:"https://example.com/video.mp4"`
	Type            string `json:"type" example:"VIDEO"` // IMAGE, VIDEO, HTML
	DurationSeconds int    `json:"duration_seconds" example:"30"`
}

// @Summary      Create and assign content to screen
// @Description  Creates a new content item and immediately assigns it to the specified screen.
// @Tags         Content
// @Accept       json
// @Produce      json
// @Param        id       path      string                        true  "UUID of the screen"
// @Param        request  body      CreateAndAssignContentRequest true  "Content details"
// @Success      201      {object}  domain.Content
// @Failure      400      {string}  string "invalid request"
// @Failure      500      {string}  string "internal server error"
// @Router       /screens/{id}/content [post]
func (h *HttpHandler) handleCreateAndAssignContent(w http.ResponseWriter, r *http.Request) {
	screenID := r.PathValue("id")
	var req CreateAndAssignContentRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}

	// 1. Map type to enum
	cType := pb.ContentType_CONTENT_TYPE_UNSPECIFIED
	switch req.Type {
	case "IMAGE":
		cType = pb.ContentType_CONTENT_TYPE_IMAGE
	case "VIDEO":
		cType = pb.ContentType_CONTENT_TYPE_VIDEO
	case "HTML":
		cType = pb.ContentType_CONTENT_TYPE_HTML
	}

	// 2. Create Content via gRPC
	createResp, err := h.contentGrpc.CreateContent(r.Context(), &pb.CreateContentRequest{
		Title:           req.Title,
		Url:             req.Url,
		Type:            cType,
		DurationSeconds: int32(req.DurationSeconds),
	})
	if err != nil {
		h.logger.Error("Failed to create content", zap.Error(err))
		http.Error(w, "failed to create content internally", http.StatusBadGateway)
		return
	}

	// 3. Assign Content to Screen via gRPC
	_, err = h.contentGrpc.AssignContentToScreen(r.Context(), &pb.AssignContentToScreenRequest{
		ContentId: createResp.Content.Id,
		ScreenId:  screenID,
	})
	if err != nil {
		h.logger.Error("Failed to assign content", zap.Error(err))
		http.Error(w, "failed to assign content internally", http.StatusBadGateway)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(createResp.Content)
}
