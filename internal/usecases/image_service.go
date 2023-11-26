package usecases

import "image-processing/internal/domain"

type ImageService struct {
	Processor domain.ImageProcessor
}

func NewImageService(processor domain.ImageProcessor) *ImageService {
	return &ImageService{Processor: processor}
}

func (s *ImageService) ResizeImage(filename string, width, height int) ([]byte, error) {
	return s.Processor.ResizeImage(filename, width, height)
}
