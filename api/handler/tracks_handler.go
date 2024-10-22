package handler

import (
	"fmt"
	"github.com/go-playground/validator/v10"
	_ "github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"math"
	"music-echo/api/domain/dao"
	"music-echo/api/domain/dto"
	"music-echo/api/repository"
	"music-echo/utils"
	"net/http"
	"sync"
)

var lock sync.Mutex

type TracksHandler interface {
	GetTracksByID(e echo.Context) error
	CreateTracks(e echo.Context) error
	UpdateTracks(e echo.Context) error
	DeleteTracks(e echo.Context) error
	GetAllTracks(e echo.Context) error
	LikeTracks(e echo.Context) error
}

type TracksHandlerImpl struct {
	TracksRepository repository.TracksRepository
	ArtistRepository repository.ArtistRepository
	LikesRepository  repository.LikesRepository
	Validators       *validator.Validate
}

func NewTracksHandlerImpl(tracksRepository repository.TracksRepository, artistRepository repository.ArtistRepository, likesRepository repository.LikesRepository, validators *validator.Validate) TracksHandler {
	return &TracksHandlerImpl{
		TracksRepository: tracksRepository,
		ArtistRepository: artistRepository,
		LikesRepository:  likesRepository,
		Validators:       validators,
	}
}

func (t *TracksHandlerImpl) GetTracksByID(e echo.Context) error {
	// Prevent Race Conditions
	lock.Lock()
	defer lock.Unlock()

	// Get Tracks ID
	var id int64
	var err error

	id, err = utils.ReadIdParam(e)
	if err != nil {
		echo.NewHTTPError(http.StatusBadRequest, "invalid id parameter")
	}

	// Get Tracks
	var tracksGet *dao.Tracks
	var artistGet *dao.Artists
	var likeGet *int64

	tracksGet, artistGet, likeGet, err = t.TracksRepository.GetId(e.Request().Context(), id)
	if err != nil {
		e.JSON(http.StatusNotFound, "theres no track that match an id")
	}

	// Response
	var trackResponse = dto.TrackGetResponse{
		Track:  tracksGet,
		Artist: artistGet,
		Likes:  likeGet,
	}
	var webResponse = dto.WebResponse{
		Message: fmt.Sprintf("get tracks %d", trackResponse.Track.Id),
		Data:    trackResponse,
	}
	return e.JSON(http.StatusOK, webResponse)
}

func (t *TracksHandlerImpl) CreateTracks(e echo.Context) error {
	// Prevent Race Conditions
	lock.Lock()
	defer lock.Unlock()

	// Decode from JSON Request Body
	tracksRequest := new(dto.TrackPostRequest)
	err := utils.ReadJSON(e, tracksRequest)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}

	// Validate
	err = t.Validators.Struct(tracksRequest)
	if err != nil {
		return echo.NewHTTPError(http.StatusNotAcceptable, err)
	}

	// DAO
	var tracks = &dao.Tracks{
		Title:    tracksRequest.Title,
		Duration: tracksRequest.Duration,
		Year:     tracksRequest.Year,
		Genre:    tracksRequest.Genre,
	}

	// Get Artist ID by Name
	var artistGet *dao.Artists
	artistGet, err = t.ArtistRepository.GetByName(e.Request().Context(), tracksRequest.Artist.Name)
	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound, "artist not found!")
	}
	tracks.IdArtist = artistGet.Id

	// Create Track
	err = t.TracksRepository.Insert(e.Request().Context(), tracks)
	if err != nil {
		return echo.NewHTTPError(http.StatusConflict, "conflicting database")
	}

	// Encode into JSON Response Body
	var trackResponse = dto.TrackInsertResponse{
		Id:        tracks.Id,
		CreatedAt: tracks.CreatedAt,
	}
	var webResponse = dto.WebResponse{
		Message: "post new tracks",
		Data:    trackResponse,
	}
	return e.JSON(http.StatusOK, webResponse)
}

func (t *TracksHandlerImpl) UpdateTracks(e echo.Context) error {
	var err error
	var id int64
	var tracksRequest *dto.TrackUpdateRequest
	var trackGet *dao.Tracks
	var artistGet *dao.Artists
	var likeGet *int64
	var tracksResponse dto.TrackUpdateResponse
	var response dto.WebResponse

	// Prevent Race Conditions
	lock.Lock()
	defer lock.Unlock()

	// Read ID of Tracks
	id, err = utils.ReadIdParam(e)
	if err != nil {
		echo.NewHTTPError(http.StatusBadRequest, "invalid id parameter")
	}

	// Read Request Body
	tracksRequest = new(dto.TrackUpdateRequest)
	err = utils.ReadJSON(e, tracksRequest)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}

	// Validate
	err = t.Validators.Struct(tracksRequest)
	if err != nil {
		return echo.NewHTTPError(http.StatusNotAcceptable, err)
	}

	// Get Copy of Current Track Version
	trackGet, artistGet, likeGet, err = t.TracksRepository.GetId(e.Request().Context(), id)

	// Update Track
	if tracksRequest.Artist != nil && tracksRequest.Artist.Name != nil {
		artistGet, err = t.ArtistRepository.GetByName(e.Request().Context(), *tracksRequest.Artist.Name)
		if err != nil {
			return echo.NewHTTPError(http.StatusNotFound, "artist doesnt exist")
		}
		trackGet.IdArtist = artistGet.Id
	}
	if tracksRequest.Title != nil {
		trackGet.Title = *tracksRequest.Title
	}
	if tracksRequest.Duration != nil {
		trackGet.Duration = *tracksRequest.Duration
	}
	if tracksRequest.Year != nil {
		trackGet.Year = *tracksRequest.Year
	}
	if tracksRequest.Genre != nil {
		trackGet.Genre = *tracksRequest.Genre
	}

	err = t.TracksRepository.Update(e.Request().Context(), trackGet)
	if err != nil {
		if err.Error() == "trackGet doesnt exist" {
			return echo.NewHTTPError(http.StatusNotFound, err.Error())
		}
		return echo.NewHTTPError(http.StatusConflict, "conflicting database")
	}

	// Response
	tracksResponse = dto.TrackUpdateResponse{
		Track:  trackGet,
		Artist: artistGet,
		Likes:  likeGet,
	}
	response = dto.WebResponse{
		Message: fmt.Sprintf("Success update tracks %d", id),
		Data:    tracksResponse,
	}

	return e.JSON(http.StatusOK, response)
}

func (t *TracksHandlerImpl) DeleteTracks(e echo.Context) error {
	var id int64
	var err error
	var response dto.WebResponse

	// Prevent Race Conditions
	lock.Lock()
	defer lock.Unlock()

	// Read ID
	id, err = utils.ReadIdParam(e)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid id parameter")
	}

	// Delete
	err = t.TracksRepository.Delete(e.Request().Context(), id)
	if err != nil {
		if err.Error() == "track doesnt exist" {
			return echo.NewHTTPError(http.StatusNotFound, err.Error())
		}
		return echo.NewHTTPError(http.StatusConflict, "conflicting database")
	}

	// Response
	response = dto.WebResponse{
		Message: fmt.Sprintf("succesfully delete task %d", id),
	}

	return e.JSON(http.StatusOK, response)
}

func (t *TracksHandlerImpl) GetAllTracks(e echo.Context) error {
	var err error
	var title, artist string
	var genre []string
	var sorting utils.Sortings
	var paginating utils.Paginatings
	var tracksGetAll []*dao.Tracks
	var artistGetAll []*dao.Artists
	var totalRecord int64
	var likes []int64
	var metadata dto.MetadataResponse
	var tracksResponse []dto.TrackGetAllResponse
	var response dto.WebResponse

	// Prevent race condition
	lock.Lock()
	defer lock.Unlock()

	// Query Parameter
	title = utils.ReadStrQuery(e, "title", "")

	artist = utils.ReadStrQuery(e, "artist", "")

	genre = utils.ReadCSVQuery(e, "genres", []string{})

	paginating.Page = utils.ReadIntQuery(e, "page", 1)
	paginating.PageSize = utils.ReadIntQuery(e, "page_size", 5)
	err = paginating.Validate(t.Validators)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}

	sorting.Sorts = utils.ReadStrQuery(e, "sort", "id")
	sorting.SafeSortLists = []string{
		"id", "title", "year",
		"-id", "-title", "-year",
	}

	// Get all tracks
	tracksGetAll, artistGetAll, likes, totalRecord, err = t.TracksRepository.GetAll(e.Request().Context(), title, artist, sorting, paginating, genre)
	if err != nil {
		return echo.NewHTTPError(http.StatusConflict, "conflicting database")
	}

	// Response
	metadata.CurrentPage = paginating.Page
	metadata.PageSize = paginating.PageSize
	metadata.FirstPage = 1
	metadata.LastPage = int64(math.Ceil(float64(totalRecord) / float64(paginating.PageSize)))
	metadata.TotalRecord = totalRecord

	tracksResponse = make([]dto.TrackGetAllResponse, len(tracksGetAll))
	for i := 0; i < len(tracksGetAll); i++ {
		tracksResponse[i].Track = tracksGetAll[i]
		tracksResponse[i].Artist = artistGetAll[i]
		tracksResponse[i].Likes = likes[i]
	}

	response = dto.WebResponse{
		Message:  fmt.Sprintf("Title:%s Artist:%s Genre:%s Page:%d PageSize:%d Sort:%s", title, artist, genre, paginating.Page, paginating.PageSize, sorting.Sorts),
		Metadata: metadata,
		Data:     tracksResponse,
	}
	return e.JSON(http.StatusOK, response)
}

func (t *TracksHandlerImpl) LikeTracks(e echo.Context) error {
	//TODO implement me
	panic("implement me")
}
