package shorturl

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"

	"github.com/alexyslozada/shorturl/domain/shorturl"
	"github.com/alexyslozada/shorturl/model"
)

type handler struct {
	useCase shorturl.UseCase
	logger  *zap.SugaredLogger
}

func newHandler(uc shorturl.UseCase, l *zap.SugaredLogger) handler {
	return handler{useCase: uc, logger: l}
}

func (h handler) Create(c echo.Context) error {
	s := model.ShortURL{}
	err := c.Bind(&s)
	if err != nil {
		h.logger.Infow("can't bind short url", "func", "Create", "internal", err)
		return c.JSON(http.StatusBadRequest, "Please verify the short_url structure")
	}

	err = h.useCase.Create(&s)
	if err != nil {
		h.logger.Errorw("can't create short url", "func", "Create", "short_url", s, "internal", err)
		return c.JSON(http.StatusInternalServerError, "Ups!!! can't create short url")
	}

	return c.JSON(http.StatusCreated, nil)
}

func (h handler) Update(c echo.Context) error {
	s := model.ShortURL{}
	err := c.Bind(&s)
	if err != nil {
		h.logger.Infow("can't bind short url", "func", "Update", "internal", err)
		return c.JSON(http.StatusBadRequest, "Please verify the short_url structure")
	}

	err = h.useCase.Update(&s)
	if err != nil {
		h.logger.Errorw("can't update short url", "func", "Update", "short_url", s, "internal", err)
		return c.JSON(http.StatusInternalServerError, "Ups! can't update short url")
	}

	return c.JSON(http.StatusOK, nil)
}

func (h handler) Delete(c echo.Context) error {
	ID := c.Param("id")
	uuidID, err := uuid.Parse(ID)
	if err != nil {
		h.logger.Infow("ID is not a valid uuid type", "func", "Delete", "id", ID, "internal", err)
		return c.JSON(http.StatusBadRequest, "Please verify the ID is a valid uuid type")
	}

	err = h.useCase.Delete(uuidID)
	if err != nil {
		h.logger.Errorw("can't delete short url", "func", "Delete", "id", uuidID, "internal", err)
		return c.JSON(http.StatusInternalServerError, "Ups! can't delete short url")
	}

	return c.JSON(http.StatusOK, nil)
}

func (h handler) ByShort(c echo.Context) error {
	s := c.Param("short")
	shortURL, err := h.useCase.ByShort(s)
	if err != nil {
		h.logger.Errorw("can't find short url", "func", "ByShort", "short", s, "internal", err)
		return c.JSON(http.StatusInternalServerError, "Ups! can't find short url")
	}

	return c.JSON(http.StatusOK, map[string]model.ShortURL{"data": shortURL})
}

func (h handler) All(c echo.Context) error {
	ss, err := h.useCase.All()
	if err != nil {
		h.logger.Errorw("can't get short urls", "func", "All", "internal", err)
		return c.JSON(http.StatusInternalServerError, "Ups! can't get short urls")
	}

	return c.JSON(http.StatusOK, map[string]model.ShortURLs{"data": ss})
}
