package price

import "github.com/sing3demons/go-product-service/microservice"

type IProductPriceHandler interface {
	FindAll(c microservice.IContext)
	FindOne(c microservice.IContext)
	InsertProductPrice(c microservice.IContext)
	DeleteProductPrice(c microservice.IContext)
}

type productPriceHandler struct {
	svc IProductPriceService
}

func NewProductPriceHandler(svc IProductPriceService) IProductPriceHandler {
	return &productPriceHandler{svc}
}

func (h *productPriceHandler) InsertProductPrice(c microservice.IContext) {
	var req CreateProductPrice
	if err := c.Body(&req); err != nil {
		c.Error(400, "Bad Request", err)
		return
	}
	id, err := h.svc.CreateProductPrice(c, req)
	if err != nil {
		c.Error(500, "Internal Server Error", err)
		return
	}
	c.JSON(200, map[string]string{
		"message": "success",
		"id":      id,
	})
}

func (h *productPriceHandler) FindOne(c microservice.IContext) {
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

func (h *productPriceHandler) FindAll(c microservice.IContext) {
	result, err := h.svc.FindAll(c)
	if err != nil {
		c.Error(500, "Internal Server Error", err)
		return
	}

	c.JSON(200, result)
}

func (h *productPriceHandler) DeleteProductPrice(c microservice.IContext) {
	id, err := h.svc.DeleteProductPrice(c)
	if err != nil {
		c.Error(500, "Internal Server Error", err)
		return
	}
	c.JSON(200, map[string]string{
		"message": "success",
		"id":      id,
	})
}
