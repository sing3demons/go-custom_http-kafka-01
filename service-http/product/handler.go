package product

import (
	"github.com/sing3demons/go-product-service/microservice"
)

type IProductHandler interface {
	FindAll(c microservice.IContext)
	FindOne(c microservice.IContext)
	InsertProduct(c microservice.IContext)
	DeleteProduct(c microservice.IContext)
}

type ProductHandler struct {
	svc IProductService
}

func NewProductHandler(svc IProductService) *ProductHandler {
	return &ProductHandler{svc}
}

func (h *ProductHandler) FindAll(c microservice.IContext) {
	products, err := h.svc.FindAll(c)
	if err != nil {
		c.Error(500, "Internal Server Error", err)
		return
	}

	c.JSON(200, products)
}

func (h *ProductHandler) InsertProduct(c microservice.IContext) {
	var req CreateProductRequest
	if err := c.Body(&req); err != nil {
		c.Error(400, "Bad Request", err)
		return
	}
	id, err := h.svc.EventCreateProduct(c, req)
	if err != nil {
		c.Error(500, "Internal Server Error", err)
		return
	}
	c.JSON(200, map[string]string{
		"message": "success",
		"id":      id,
	})
}

func (h *ProductHandler) FindOne(c microservice.IContext) {
	product, err := h.svc.FindOne(c)
	if err != nil {
		if err.Error() == "mongo: no documents in result" {
			c.Error(404, "Not Found", err)
			return
		}
		c.Error(500, "Internal Server Error", err)
		return
	}

	c.JSON(200, product)
}

func (h *ProductHandler) DeleteProduct(c microservice.IContext) {
	id, err := h.svc.EventDeleteProduct(c)
	if err != nil {
		c.Error(500, "Internal Server Error", err)
		return
	}
	c.JSON(200, map[string]string{
		"message": "success",
		"id":      id,
	})
}
