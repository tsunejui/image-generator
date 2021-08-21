package cmd

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"path/filepath"

	"generator/pkg"

	"github.com/spf13/cobra"
)

var (
	MergeCommand = &cobra.Command{
		Use:   "merge",
		Short: "Merge the images from the specified folders",
		Long:  `Merge the images from the specified folders`,
		RunE:  runMergeCommand,
	}
	supportExtensions = []string{
		".png",
	}
	folders = []string{}
	output  string
)

func init() {
	MergeCommand.Flags().StringSliceVarP(&folders, "folder", "f", []string{}, "specify a folder")
	MergeCommand.Flags().StringVarP(&output, "out", "o", "", "export the images")
}

type MergeImages struct {
	FolderName string
	Files      []string
}

func runMergeCommand(cmd *cobra.Command, args []string) error {

	// specify the output path
	if output == "" {
		return fmt.Errorf("the output path is not specified")
	}

	isOutputFolder, err := pkg.IsDirectory(output)
	if err != nil {
		return fmt.Errorf("failed to check the output directory: %v", err)
	}
	if !isOutputFolder {
		return fmt.Errorf("failed to find the output directory: %s", output)
	}

	// check directories
	for _, f := range folders {
		isFolder, err := pkg.IsDirectory(f)
		if err != nil {
			return fmt.Errorf("failed to check the directory: %v", err)
		}
		if !isFolder {
			return fmt.Errorf("%s is not a directory", f)
		}
	}

	var images []MergeImages

	// check files
	for _, f := range folders {
		img := MergeImages{}
		files, err := ioutil.ReadDir(f)
		if err != nil {
			return fmt.Errorf("failed to get files: %v", err)
		}

		img.FolderName = filepath.Base(f)
		for _, file := range files {
			if file.IsDir() {
				continue
			}
			fileName := file.Name()
			extension := filepath.Ext(fileName)

			var isSupport bool
			for _, ext := range supportExtensions {
				if ext == extension {
					isSupport = true
				}
			}
			if isSupport {
				img.Files = append(img.Files, filepath.Join(f, fileName))
			}
		}
		images = append(images, img)
	}

	// continue to check or exit
	if len(images) == 0 {
		return fmt.Errorf("failed to find images")
	}

	var number int = 1
	fmt.Printf("\nThe analytical result:\n\n")
	for _, m := range images {
		count := len(m.Files)
		fmt.Printf("%s -> %d \n", m.FolderName, count)
		number *= count
	}
	fmt.Printf("\nNumber of combinations: %d", number)

	fmt.Printf("\n\nDo you want to continue? (Y/N):")
	var isContinue string
	fmt.Scanf("%s", &isContinue)
	if !pkg.CheckYes(isContinue) {
		fmt.Println("bye")
		return nil
	}

	if err := mergeImages(images); err != nil {
		return fmt.Errorf("failed to merge images: %v", err)
	}
	return nil
}

type MergeError struct {
	Files []string
	Err   error
}

type GeneratorErrors struct {
	AddErrors   []error
	MergeErrors []MergeError
}

type MergeResult struct {
	Files []string `json:"files"`
	Path  string   `json:"path"`
}

const resultFile = "result.json"

func mergeImages(images []MergeImages) error {

	// combinations from multiple arrays
	emptyCom := []string{}
	combinations := [][]string{emptyCom}
	for i := 0; i < len(images); i++ {
		var tmp [][]string
		for _, v1 := range combinations {
			for _, v2 := range images[i].Files {
				newGroup := []string{}
				newGroup = append(newGroup, v1...)
				newGroup = append(newGroup, v2)
				tmp = append(tmp, newGroup)
			}
		}
		combinations = tmp
	}

	// merge images
	var mergeResults []MergeResult
	var generatorErr GeneratorErrors
	for k, com := range combinations {
		var generator pkg.ImageGenerator
		var isContinue bool

		fileName := filepath.Join(output, fmt.Sprintf("%d.png", k))
		generator.SetPath(fileName)
		for _, f := range com {
			if err := generator.AddFile(f); err != nil {
				generatorErr.AddErrors = append(generatorErr.AddErrors, fmt.Errorf("failed to add file into the generator: %v", err))
				isContinue = true
				break
			}
		}
		if err := generator.Merge(); err != nil {
			mergeErr := MergeError{
				Files: generator.GetFiles(),
				Err:   fmt.Errorf("failed to merge the images: %v", err),
			}
			generatorErr.MergeErrors = append(generatorErr.MergeErrors, mergeErr)
			isContinue = true
		}

		if isContinue {
			continue
		}

		var baseFilesName []string
		for _, f := range generator.GetFiles() {
			baseFilesName = append(baseFilesName, filepath.Base(f))
		}
		mergeResults = append(mergeResults, MergeResult{
			Files: baseFilesName,
			Path:  fileName,
		})
	}

	// error handling
	fmt.Printf("\nShowing all errors below:\n")
	fmt.Printf("\n=== adding errors ===\n")
	for _, e := range generatorErr.AddErrors {
		fmt.Println(e)
	}
	fmt.Printf("\n=== merging errors ===\n")
	for _, e := range generatorErr.MergeErrors {
		fmt.Printf("\nerror message: %v\n", e.Err)
		fmt.Println("files:")
		for _, f := range e.Files {
			fmt.Println(f)
		}
	}

	fmt.Printf("\nit's done, sucess count: %d\n", len(mergeResults))

	// export result json
	data, err := json.Marshal(mergeResults)
	if err != nil {
		return fmt.Errorf("failed to marshal data: %v", err)
	}
	if err := pkg.ExportFile(filepath.Join(output, resultFile), string(data)); err != nil {
		return fmt.Errorf("failed to export data: %v", err)
	}

	return nil
}
