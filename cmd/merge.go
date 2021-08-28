package cmd

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"path/filepath"
	"strings"
	"sync"

	"generator/pkg"

	"github.com/spf13/cobra"
)

var (
	mergeCommand = &cobra.Command{
		Use:   "merge",
		Short: "Merge the images from the specified folders",
		Long:  `Merge the images from the specified folders`,
		RunE:  runMergeCommand,
	}
	supportExtensions = []string{
		".png",
	}
	folders       = []string{}
	rootDirectory string
	output        string
)

func init() {
	rootCmd.AddCommand(mergeCommand)
	mergeCommand.Flags().StringSliceVarP(&folders, "directory", "d", []string{}, "specify a directory")
	mergeCommand.Flags().StringVarP(&rootDirectory, "root", "D", "", "specify the root directory")
	mergeCommand.Flags().StringVarP(&output, "out", "o", "", "export the images")
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

	// Adjust the directories
	var directories []string
	if rootDirectory != "" {
		for _, f := range folders {
			if rootDirectory != "" {
				directories = append(directories, filepath.Join(rootDirectory, f))
			} else {
				directories = append(directories, f)
			}
		}
	}

	// check directories
	for _, f := range directories {
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
	for _, f := range directories {
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

	// Ask question to user
	fmt.Printf("\n\nDo you want to continue? (Y/N): ")
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
	Type  string
}

const typeGeneratorErr = "generatorErr"
const typeOtherErr = "otherErr"
const typeAddingErr = "addingErr"

type MergeAttribute struct {
	TraitType string `json:"trait_type"`
	Value     string `json:"value"`
}

type MergeResult struct {
	ID         int              `json:"id"`
	Attributes []MergeAttribute `json:"attributes"`
	Path       string           `json:"path"`
}

type MergeFile struct {
	Directory string
	Path      string
}

const resultFile = "result.json"

func mergeImages(images []MergeImages) error {

	// combinations from multiple arrays
	emptyCom := []MergeFile{}
	combinations := [][]MergeFile{emptyCom}
	for i := 0; i < len(images); i++ {
		var tmp [][]MergeFile
		for _, v1 := range combinations {
			for _, v2 := range images[i].Files {
				newGroup := []MergeFile{}
				newGroup = append(newGroup, v1...)
				newGroup = append(newGroup, MergeFile{
					Directory: images[i].FolderName,
					Path:      v2,
				})
				tmp = append(tmp, newGroup)
			}
		}
		combinations = tmp
	}

	// merge images
	var mergeResults []*MergeResult
	var generatorErrors []*MergeError
	chanErrors := make(chan *MergeError, 1)
	results := make(chan *MergeResult, 1)
	var wg sync.WaitGroup
	wg.Add(len(combinations))
	for k, com := range combinations {
		go func(id int, com []MergeFile) {
			defer wg.Done()
			m, genErr := genImage(id, fmt.Sprintf("%d.png", id), com)
			if genErr != nil {
				chanErrors <- genErr
			} else {
				results <- m
			}
		}(k, com)
	}

	// add result
	go func() {
		defer close(results)
		for m := range results {
			mergeResults = append(mergeResults, m)
		}
	}()
	// add error
	go func() {
		defer close(chanErrors)
		for err := range chanErrors {
			generatorErrors = append(generatorErrors, err)
		}
	}()
	wg.Wait()

	// error handling
	var addErrors []*MergeError
	var mergeErrors []*MergeError
	for _, e := range generatorErrors {
		switch e.Type {
		case typeAddingErr:
			addErrors = append(addErrors, e)
		case typeGeneratorErr:
			mergeErrors = append(mergeErrors, e)
		}
	}

	fmt.Printf("\nShowing all errors below:\n")
	fmt.Printf("\n=== adding errors ===\n")
	for _, e := range addErrors {
		fmt.Println(e)
	}
	fmt.Printf("\n=== merging errors ===\n")
	for _, e := range mergeErrors {
		fmt.Printf("\nerror message: %v\n", e.Err)
		fmt.Println("files:")
		for _, f := range e.Files {
			fmt.Println(f)
		}
	}
	fmt.Printf("\nit's done, sucess count: %d\n", len(mergeResults))

	// Bubble Sort
	for i := 0; i < len(mergeResults)-1; i++ {
		for j := 0; j < len(mergeResults)-i-1; j++ {
			if mergeResults[j].ID > mergeResults[j+1].ID {
				temp := mergeResults[j]
				mergeResults[j] = mergeResults[j+1]
				mergeResults[j+1] = temp
			}
		}
	}

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

func genImage(id int, imageName string, com []MergeFile) (*MergeResult, *MergeError) {
	var generator pkg.ImageGenerator

	fileName := filepath.Join(output, imageName)
	generator.SetPath(fileName)
	for _, f := range com {
		if err := generator.AddFile(f.Directory, f.Path); err != nil {
			return nil, &MergeError{
				Type: typeAddingErr,
				Err:  fmt.Errorf("failed to add file into the generator: %v", err),
			}
		}
	}
	if err := generator.Merge(); err != nil {
		return nil, &MergeError{
			Files: generator.GetFilesName(),
			Type:  typeGeneratorErr,
			Err:   fmt.Errorf("failed to add file into the generator: %v", err),
		}
	}

	// generate the json file of result
	var attributes []MergeAttribute
	for _, f := range generator.GetFiles() {
		fName := filepath.Base(f.Data.Name())
		attributes = append(attributes, MergeAttribute{
			TraitType: f.Directory,
			Value:     strings.Replace(fName, fmt.Sprintf("%s", filepath.Ext(fName)), "", 1),
		})
	}
	return &MergeResult{
		ID:         id,
		Attributes: attributes,
		Path:       fileName,
	}, nil
}
