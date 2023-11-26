package service

import (
	"encoding/json"
	"time"

	"github.com/sing3demons/go-category-service/model"
	"github.com/sing3demons/go-category-service/repository"
	"github.com/sirupsen/logrus"
)

type EventHandler interface {
	Handle(topic string, eventBytes []byte)
}

type categoryEventHandler struct {
	categoryRepo repository.CategoryRepository
	logger       *logrus.Logger
}

func NewCategoryEventHandler(categoryRepo repository.CategoryRepository, logger *logrus.Logger) EventHandler {
	return &categoryEventHandler{categoryRepo, logger}
}

type Event struct {
	Header map[string]any `json:"header"`
	Body   any            `json:"body"`
}

func (obj *categoryEventHandler) Handle(topic string, eventBytes []byte) {
	switch topic {
	case "category.created":
		var event Event
		if err := json.Unmarshal(eventBytes, &event); err != nil {
			obj.logger.WithFields(logrus.Fields{
				"topic": topic,
				"error": err,
			}).Error("unmarshal event error")
			return
		}
		header := event.Header

		bytes, err := json.Marshal(event.Body)
		if err != nil {
			obj.logger.WithFields(logrus.Fields{
				"topic": topic,
				"error": err,
			}).Error("marshal event body error")
			return
		}

		var body model.CreateCategoryReq
		if err := json.Unmarshal(bytes, &body); err != nil {
			obj.logger.WithFields(logrus.Fields{
				"topic": topic,
				"error": err,
			}).Error("marshal event body error")
			return
		}

		var doc model.CreateCategoryReq

		doc.ID = body.ID
		doc.Name = body.Name
		doc.Type = "category"
		doc.Status = body.Status

		if body.LastUpdate.IsZero() {
			body.LastUpdate = time.Now().UTC()
		}
		doc.LastUpdate = body.LastUpdate

		if err := obj.categoryRepo.Save(doc); err != nil {
			obj.logger.WithFields(logrus.Fields{
				"topic": topic,
				"heder": header,
				"body":  doc,
				"error": err,
			}).Error("insert category error")
			return
		}
		obj.logger.WithFields(logrus.Fields{
			"topic": topic,
			"heder": header,
			"body":  doc,
		}).Info("insert category success")
	case "category.updated":
		var event Event
		if err := json.Unmarshal(eventBytes, &event); err != nil {
			obj.logger.WithFields(logrus.Fields{
				"topic": topic,
				"error": err,
			}).Error("unmarshal event error")
			return
		}
		header := event.Header

		bytes, err := json.Marshal(event.Body)
		if err != nil {
			obj.logger.WithFields(logrus.Fields{
				"topic": topic,
				"error": err,
			}).Error("marshal event body error")
			return
		}

		var body model.UpdateCategoryReq
		if err := json.Unmarshal(bytes, &body); err != nil {
			obj.logger.WithFields(logrus.Fields{
				"topic": topic,
				"error": err,
			}).Error("marshal event body error")
			return
		}

		var doc model.UpdateCategoryReq
		doc.ID = body.ID
		doc.Type = "category"
		if body.Name != "" {
			doc.Name = body.Name
		}
		if body.Status != "" {
			doc.Status = body.Status
		}
		if len(body.Products) > 0 {
			doc.Products = body.Products
		}
		if body.LastUpdate.IsZero() {
			body.LastUpdate = time.Now().UTC()
		}
		doc.LastUpdate = body.LastUpdate

		category, err := obj.categoryRepo.Update(doc)
		if err != nil {
			obj.logger.WithFields(logrus.Fields{
				"topic": topic,
				"heder": header,
				"body":  doc,
				"error": err,
			}).Error("insert category error")
			return
		}
		obj.logger.WithFields(logrus.Fields{
			"topic":  topic,
			"heder":  header,
			"body":   doc,
			"result": category,
		}).Info("insert category success")
	}
}
