package pkg

import (
	"fmt"
	"image"
	"image/draw"
	"image/png"
	"os"
)

type ImageGenerator struct {
	directory string
	path      string
	files     []Image
}

type Image struct {
	Directory string
	Data      *os.File
}

func (i *ImageGenerator) SetDirectory(directory string) *ImageGenerator {
	i.directory = directory
	return i
}

func (i *ImageGenerator) SetPath(path string) *ImageGenerator {
	i.path = path
	return i
}

func (i *ImageGenerator) AddFile(directory string, path string) error {
	f, err := os.Open(path)
	if err != nil {
		return fmt.Errorf("failed to open file (%s): %v", path, err)
	}
	i.files = append(i.files, Image{
		Directory: directory,
		Data:      f,
	})
	return nil
}

func (i *ImageGenerator) GetFiles() []Image {
	return i.files
}

func (i *ImageGenerator) GetFilesName() []string {
	var files []string
	for _, f := range i.files {
		files = append(files, f.Data.Name())
	}
	return files
}

func (i *ImageGenerator) Merge() error {
	if len(i.files) == 0 {
		return fmt.Errorf("failed to find the images")
	}

	if i.path == "" {
		return fmt.Errorf("failed to find image path")
	}

	defer i.close()
	var rectangle image.Rectangle
	var outputImage *image.RGBA
	for k, f := range i.files {
		img, err := png.Decode(f.Data)
		if err != nil {
			return fmt.Errorf("failed to decode file: %v", err)
		}
		if k == 0 {
			rectangle = img.Bounds()
			outputImage = image.NewRGBA(rectangle)
		}

		draw.Draw(outputImage, rectangle, img, image.Point{}, draw.Over)
	}

	outputFile, err := os.Create(i.path)
	if err != nil {
		return fmt.Errorf("failed to create file: %v", err)
	}
	defer outputFile.Close()

	if err := png.Encode(outputFile, outputImage); err != nil {
		return fmt.Errorf("failed to encode png: %v", err)
	}
	return nil
}

func (i *ImageGenerator) close() error {
	for _, f := range i.files {
		if err := f.Data.Close(); err != nil {
			fmt.Printf("failed to close file: %v", err)
		}
	}
	return nil
}
