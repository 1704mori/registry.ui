package handlers

import (
	"fmt"
	"log"
	"net/http"
	"sort"

	"github.com/labstack/echo/v4"

	"github.com/1704mori/registry.ui/internal/api"
	"github.com/1704mori/registry.ui/internal/config"
	"github.com/1704mori/registry.ui/internal/templates/components"
	"github.com/1704mori/registry.ui/internal/templates/layouts"
	"github.com/1704mori/registry.ui/internal/templates/pages"
)

// Handlers struct holds the dependencies for HTTP handlers
type Handlers struct {
	registryClient *api.Client
	config         *config.Config
}

// NewHandlers creates a new Handlers instance
func NewHandlers(registryClient *api.Client, config *config.Config) *Handlers {
	return &Handlers{
		registryClient: registryClient,
		config:         config,
	}
}

// Dashboard renders the dashboard page
func (h *Handlers) Dashboard(c echo.Context) error {
	// Get repositories
	repositories, err := h.registryClient.ListRepositories()
	if err != nil {
		return renderErrorPage(c, err)
	}

	// Sort repositories by name
	sort.Strings(repositories)

	// Render dashboard page
	component := pages.Dashboard(repositories, h.config.DefaultTheme)
	return layouts.Base(h.config.DefaultTheme, component).Render(c.Request().Context(), c.Response().Writer)
}

// ListImages renders the images list page
func (h *Handlers) ListImages(c echo.Context) error {
	// Get repositories
	repositories, err := h.registryClient.ListRepositories()
	if err != nil {
		return renderErrorPage(c, err)
	}

	// Sort repositories by name
	sort.Strings(repositories)

	// Render images page
	component := pages.Images(repositories, h.config.DefaultTheme)
	return layouts.Base(h.config.DefaultTheme, component).Render(c.Request().Context(), c.Response().Writer)
}

// GetImage renders the image detail page
func (h *Handlers) GetImage(c echo.Context) error {
	// Get image name from URL parameter
	name := c.Param("*")
	log.Println(name)

	if name == "" {
		return c.Redirect(http.StatusFound, "/images")
	}

	// Get image info
	imageInfo, err := h.registryClient.GetImageInfo(name)
	if err != nil {
		return renderErrorPage(c, err)
	}

	// Render image detail page
	component := pages.ImageDetail(imageInfo, h.config.DefaultTheme)
	return layouts.Base(h.config.DefaultTheme, component).Render(c.Request().Context(), c.Response().Writer)
}

// ListTags renders the tags list page
func (h *Handlers) ListTags(c echo.Context) error {
	// Get repository name from URL parameter
	name := c.Param("*")
	log.Printf("ListTags name: %s", name)
	if name == "" {
		return c.Redirect(http.StatusFound, "/images")
	}

	// Get tags
	tags, err := h.registryClient.ListTags(name)
	if err != nil {
		return renderErrorPage(c, err)
	}

	// Sort tags
	sort.Strings(tags)

	// Render tags page
	component := components.TagsList(name, tags)
	return layouts.Base(h.config.DefaultTheme, component).Render(c.Request().Context(), c.Response().Writer)
}

// GetTag renders the tag detail page
func (h *Handlers) GetTag(c echo.Context) error {
	// Get repository and tag names from URL parameters
	tag := c.Param("tag")
	name := c.Param("*")
	if name == "" || tag == "" {
		return c.Redirect(http.StatusFound, "/images")
	}

	// Get tag info
	tagInfo, err := h.registryClient.GetTagInfo(name, tag)
	if err != nil {
		return renderErrorPage(c, err)
	}

	// Render tag detail page
	component := pages.TagDetail(tagInfo, h.config.DefaultTheme)
	return layouts.Base(h.config.DefaultTheme, component).Render(c.Request().Context(), c.Response().Writer)
}

// DeleteTag handles tag deletion
func (h *Handlers) DeleteTag(c echo.Context) error {
	// Get repository and tag names from URL parameters
	name := c.Param("*")
	tag := c.Param("tag")
	if name == "" || tag == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Missing parameters"})
	}

	// Get tag info to get the digest
	tagInfo, err := h.registryClient.GetTagInfo(name, tag)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	// Delete tag
	err = h.registryClient.DeleteTag(name, tagInfo.Digest)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	// Redirect to the image page
	return c.Redirect(http.StatusFound, fmt.Sprintf("/images/%s", name))
}

// Settings renders the settings page
func (h *Handlers) Settings(c echo.Context) error {
	// Render settings page
	component := pages.Settings(h.config, h.config.DefaultTheme)
	return layouts.Base(h.config.DefaultTheme, component).Render(c.Request().Context(), c.Response().Writer)
}

// UpdateSettings handles settings update
func (h *Handlers) UpdateSettings(c echo.Context) error {
	// Parse form
	registryURL := c.FormValue("registry_url")
	registryUsername := c.FormValue("registry_username")
	registryPassword := c.FormValue("registry_password")
	registryInsecure := c.FormValue("registry_insecure") == "on"
	defaultTheme := c.FormValue("default_theme")

	// Update config
	err := h.config.UpdateConfig(registryURL, registryUsername, registryPassword, registryInsecure, defaultTheme)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	// Recreate registry client
	registryClient, err := api.NewClient(h.config.RegistryURL, h.config.RegistryUsername, h.config.RegistryPassword, h.config.RegistryInsecure)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	h.registryClient = registryClient

	// Redirect to the settings page
	return c.Redirect(http.StatusFound, "/settings")
}

// HTMX specific handlers

// HtmxListImages handles HTMX request for images list
func (h *Handlers) HtmxListImages(c echo.Context) error {
	// Get repositories
	repositories, err := h.registryClient.ListRepositories()
	if err != nil {
		return c.HTML(http.StatusInternalServerError, fmt.Sprintf("<div class='text-red-500'>Error: %s</div>", err.Error()))
	}

	// Sort repositories by name
	sort.Strings(repositories)

	// Render images list component
	component := components.ImagesList(repositories)
	return component.Render(c.Request().Context(), c.Response().Writer)
}

// HtmxListTags handles HTMX request for tags list
func (h *Handlers) HtmxListTags(c echo.Context) error {
	// Get repository name from URL parameter
	name := c.Param("*")
	if name == "" {
		return c.HTML(http.StatusBadRequest, "<div class='text-red-500'>Error: Missing repository name</div>")
	}

	// Get tags
	tags, err := h.registryClient.ListTags(name)
	if err != nil {
		return c.HTML(http.StatusInternalServerError, fmt.Sprintf("<div class='text-red-500'>Error: %s</div>", err.Error()))
	}

	// Sort tags
	sort.Strings(tags)

	// Render tags list component
	component := components.TagsList(name, tags)
	return component.Render(c.Request().Context(), c.Response().Writer)
}

// HtmxThemeToggle handles HTMX request for theme toggling
func (h *Handlers) HtmxThemeToggle(c echo.Context) error {
	currentTheme := c.QueryParam("theme")
	newTheme := "light"
	if currentTheme == "light" {
		newTheme = "dark"
	}

	component := components.ThemeSwitcher(newTheme)
	return component.Render(c.Request().Context(), c.Response().Writer)
}

// Helper function to render error page
func renderErrorPage(c echo.Context, err error) error {
	return c.HTML(http.StatusInternalServerError, fmt.Sprintf("<h1>Error</h1><p>%s</p>", err.Error()))
}
