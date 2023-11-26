package domain

type Image struct {
	Width  int
	Height int
	Data   []byte
}

type ImageProcessor interface {
	ResizeImage(filename string, newWidth, newHeight int) ([]byte, error)
}
