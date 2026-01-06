package record

import (
	"os"

	"github.com/shashankC011/replayr/internal/config"
)

type Handler struct {
	cfg  *config.Config
	file *os.File
}

func New(cfg *config.Config, file *os.File) *Handler {
	return &Handler{
		cfg:  cfg,
		file: file,
	}
}
