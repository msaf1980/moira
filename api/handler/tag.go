package handler

import (
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/render"
	"github.com/moira-alert/moira/api"
	"github.com/moira-alert/moira/api/controller"
	"github.com/moira-alert/moira/api/dto"
	"github.com/moira-alert/moira/api/middleware"
)

func tag(router chi.Router) {
	router.Post("/", createTags)
	router.Get("/", getAllTags)
	router.With(middleware.AdminOnlyMiddleware()).Get("/stats", getAllTagsAndSubscriptions)
	router.Route("/{tag}", func(router chi.Router) {
		router.Use(middleware.TagContext)
		router.Use(middleware.AdminOnlyMiddleware())
		router.Delete("/", removeTag)
	})
}

// nolint: gofmt,goimports
//
//	@summary	Get all tags
//	@id			get-all-tags
//	@tags		tag
//	@produce	json
//	@success	200	{object}	dto.TagsData					"Tags fetched successfully"
//	@failure	422	{object}	api.ErrorRenderExample			"Render error"
//	@failure	500	{object}	api.ErrorInternalServerExample	"Internal server error"
//	@router		/tag [get]
func getAllTags(writer http.ResponseWriter, request *http.Request) {
	tagData, err := controller.GetAllTags(database)
	if err != nil {
		render.Render(writer, request, err) //nolint
		return
	}

	if err := render.Render(writer, request, tagData); err != nil {
		render.Render(writer, request, api.ErrorRender(err)) //nolint
		return
	}
}

// nolint: gofmt,goimports
//
//	@summary	Create new tags
//	@id			create-tags
//	@tags		tag
//	@accept		json
//	@produce	json
//	@param		tags	body	dto.TagsData	true	"Tags data"
//	@success	200		"Create tags successfully"
//	@failure	400		{object}	api.ErrorInvalidRequestExample	"Bad request from client"
//	@failure	422		{object}	api.ErrorRenderExample			"Render error"
//	@failure	500		{object}	api.ErrorInternalServerExample	"Internal server error"
//	@router		/tag [post]
func createTags(writer http.ResponseWriter, request *http.Request) {
	tags := dto.TagsData{}
	if err := render.Bind(request, &tags); err != nil {
		render.Render(writer, request, api.ErrorInvalidRequest(err)) //nolint:errcheck
		return
	}

	if err := controller.CreateTags(database, &tags); err != nil {
		render.Render(writer, request, err) //nolint
	}
}

// nolint: gofmt,goimports
//
//	@summary	Get all tags and their subscriptions
//	@id			get-all-tags-and-subscriptions
//	@tags		tag
//	@produce	json
//	@success	200	{object}	dto.TagsStatistics				"Successful"
//	@failure	422	{object}	api.ErrorRenderExample			"Render error"
//	@failure	500	{object}	api.ErrorInternalServerExample	"Internal server error"
//	@router		/tag/stats [get]
func getAllTagsAndSubscriptions(writer http.ResponseWriter, request *http.Request) {
	logger := middleware.GetLoggerEntry(request)
	data, err := controller.GetAllTagsAndSubscriptions(database, logger)
	if err != nil {
		render.Render(writer, request, err) //nolint
		return
	}
	if err := render.Render(writer, request, data); err != nil {
		render.Render(writer, request, api.ErrorRender(err)) //nolint
		return
	}
}

// nolint: gofmt,goimports
//
//	@summary	Remove a tag
//	@id			remove-tag
//	@tags		tag
//	@produce	json
//	@param		tag	path		string							true	"Name of the tag to remove"	default(cpu)
//	@success	200	{object}	dto.MessageResponse				"Tag removed successfully"
//	@failure	400	{object}	api.ErrorInvalidRequestExample	"Bad request from client"
//	@failure	422	{object}	api.ErrorRenderExample			"Render error"
//	@failure	500	{object}	api.ErrorInternalServerExample	"Internal server error"
//	@router		/tag/{tag} [delete]
func removeTag(writer http.ResponseWriter, request *http.Request) {
	tagName := middleware.GetTag(request)
	response, err := controller.RemoveTag(database, tagName)
	if err != nil {
		render.Render(writer, request, err) //nolint
		return
	}
	if err := render.Render(writer, request, response); err != nil {
		render.Render(writer, request, api.ErrorRender(err)) //nolint
		return
	}
}
