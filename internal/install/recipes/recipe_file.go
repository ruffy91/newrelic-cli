package recipes

import (
	"io/ioutil"
	"net/http"
	"net/url"

	"gopkg.in/yaml.v2"

	"github.com/newrelic/newrelic-cli/internal/install/types"
)

// RecipeFile represents a recipe file as defined in the Open Installation Library.
type RecipeFile struct {
	Description    string                 `yaml:"description"`
	InputVars      []VariableConfig       `yaml:"inputVars"`
	Install        map[string]interface{} `yaml:"install"`
	InstallTargets []RecipeInstallTarget  `yaml:"installTargets"`
	Keywords       []string               `yaml:"keywords"`
	LogMatch       []types.LogMatch       `yaml:"logMatch"`
	Name           string                 `yaml:"name"`
	DisplayName    string                 `yaml:"displayName"`
	PreInstall     RecipePreInstall       `yaml:"preInstall"`
	PostInstall    RecipePostInstall      `yaml:"postInstall"`
	ProcessMatch   []string               `yaml:"processMatch"`
	Repository     string                 `yaml:"repository"`
	ValidationNRQL string                 `yaml:"validationNrql"`
}

type RecipePreInstall struct {
	Info   string `yaml:"info"`
	Prompt string `yaml:"prompt"`
}

type RecipePostInstall struct {
	Info   string `yaml:"info"`
	Prompt string `yaml:"prompt"`
}

type VariableConfig struct {
	Name    string `yaml:"name"`
	Prompt  string `yaml:"prompt"`
	Secret  bool   `secret:"prompt"`
	Default string `yaml:"default"`
}

type RecipeInstallTarget struct {
	Type            string `yaml:"type"`
	OS              string `yaml:"os"`
	Platform        string `yaml:"platform"`
	PlatformFamily  string `yaml:"platformFamily"`
	PlatformVersion string `yaml:"platformVersion"`
	KernelVersion   string `yaml:"kernelVersion"`
	KernelArch      string `yaml:"kernelArch"`
}

type RecipeFileFetcherImpl struct {
	HTTPGetFunc  func(string) (*http.Response, error)
	readFileFunc func(string) ([]byte, error)
}

func NewRecipeFileFetcher() RecipeFileFetcher {
	f := RecipeFileFetcherImpl{}
	f.HTTPGetFunc = defaultHTTPGetFunc
	f.readFileFunc = defaultReadFileFunc
	return &f
}

func defaultHTTPGetFunc(recipeURL string) (*http.Response, error) {
	return http.Get(recipeURL)
}

func defaultReadFileFunc(filename string) ([]byte, error) {
	return ioutil.ReadFile(filename)
}

func (f *RecipeFileFetcherImpl) FetchRecipeFile(recipeURL *url.URL) (*RecipeFile, error) {
	response, err := f.HTTPGetFunc(recipeURL.String())
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}
	return StringToRecipeFile(string(body))
}

func (f *RecipeFileFetcherImpl) LoadRecipeFile(filename string) (*RecipeFile, error) {
	out, err := f.readFileFunc(filename)
	if err != nil {
		return nil, err
	}
	return StringToRecipeFile(string(out))
}

func StringToRecipeFile(content string) (*RecipeFile, error) {
	f, err := NewRecipeFile(content)
	if err != nil {
		return nil, err
	}
	return f, nil
}

func NewRecipeFile(recipeFileString string) (*RecipeFile, error) {
	var f RecipeFile
	err := yaml.Unmarshal([]byte(recipeFileString), &f)
	if err != nil {
		return nil, err
	}

	return &f, nil
}

func (f *RecipeFile) String() (string, error) {
	out, err := yaml.Marshal(f)
	if err != nil {
		return "", err
	}

	return string(out), nil
}

func (f *RecipeFile) ToRecipe() (*types.Recipe, error) {
	fileStr, err := f.String()
	if err != nil {
		return nil, err
	}

	r := types.Recipe{
		File:        fileStr,
		Name:        f.Name,
		DisplayName: f.DisplayName,
		Description: f.Description,
		Repository:  f.Repository,
		Keywords:    f.Keywords,
		PreInstall: types.RecipePreInstall{
			Info:   f.PreInstall.Info,
			Prompt: f.PreInstall.Prompt,
		},
		PostInstall: types.RecipePostInstall{
			Info:   f.PostInstall.Info,
			Prompt: f.PostInstall.Prompt,
		},
		ProcessMatch:   f.ProcessMatch,
		LogMatch:       f.LogMatch,
		ValidationNRQL: f.ValidationNRQL,
	}

	return &r, nil
}

func RecipeToRecipeFile(r types.Recipe) (*RecipeFile, error) {
	var f RecipeFile
	err := yaml.Unmarshal([]byte(r.File), &f)
	if err != nil {
		return nil, err
	}

	return &f, nil
}
