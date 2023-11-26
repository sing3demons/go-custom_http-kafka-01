package category

import (
	"github.com/sing3demons/go-product-service/microservice"
)

type ICategoryHandler interface {
	FindCategories(c microservice.IContext)
	FindOne(c microservice.IContext)
	InsertProduct(c microservice.IContext)
	// DeleteProduct(c microservice.IContext)
}

type categoryHandler struct {
	svc ICategoryService
}

func NewCategoryHandler(svc ICategoryService) ICategoryHandler {
	return &categoryHandler{svc}
}

func (h *categoryHandler) FindOne(c microservice.IContext) {
	category, err := h.svc.FindOne(c)
	if err != nil {
		if err.Error() == "mongo: no documents in result" {
			c.Error(404, "Not Found", err)
			return
		}
		c.Error(500, "Internal Server Error", err)
		return
	}

	c.JSON(200, category)
}
func (h *categoryHandler) FindCategories(c microservice.IContext) {
	categories, err := h.svc.FindAll(c)
	if err != nil {
		c.Error(500, "Internal Server Error", err)
		return
	}

	c.JSON(200, categories)
}

func (h *categoryHandler) InsertProduct(c microservice.IContext) {
	var req CreateCategoryReq
	if err := c.Body(&req); err != nil {
		c.Error(400, "Bad Request", err)
		return
	}
	id, err := h.svc.CreateCategory(c, req)
	if err != nil {
		c.Error(500, "Internal Server Error", err)
		return
	}
	c.JSON(200, map[string]string{
		"message": "success",
		"id":      id,
	})
}
