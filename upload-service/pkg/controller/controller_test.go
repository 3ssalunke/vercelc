package controller

import (
	"os"
	"testing"

	"github.com/3ssalunke/vercelc/shared/config"

	"github.com/3ssalunke/vercelc/upload-service/pkg/services"
)

var (
	c *services.Container
)

func TestMain(m *testing.M) {
	// Set the environment to test
	config.SwitchEnvironment(config.EnvTest)

	// Create a new container
	c = services.NewContainer()

	// Run tests
	exitVal := m.Run()

	// Shutdown the container
	if err := c.Shutdown(); err != nil {
		panic(err)
	}

	os.Exit(exitVal)
}

// func TestController_Redirect(t *testing.T) {
// 	c.Web.GET("/path/:first/and/:second", func(c echo.Context) error {
// 		return nil
// 	}).Name = "redirect-test"

// 	ctx, _ := tests.NewContext(c.Web, "/abc")
// 	ctr := NewController(c)
// 	err := ctr.Redirect(ctx, "redirect-test", "one", "two")
// 	require.NoError(t, err)
// 	assert.Equal(t, "/path/one/and/two", ctx.Response().Header().Get(echo.HeaderLocation))
// 	assert.Equal(t, http.StatusFound, ctx.Response().Status)
// }

// func TestController_RenderPage(t *testing.T) {
// 	setup := func() (echo.Context, *httptest.ResponseRecorder, Controller, Page) {
// 		ctx, rec := tests.NewContext(c.Web, "/test/TestController_RenderPage")
// 		tests.InitSession(ctx)
// 		ctr := NewController(c)

// 		p := NewPage(ctx)
// 		p.Name = "home"
// 		p.Layout = "main"
// 		p.Headers["A"] = "b"
// 		p.Headers["C"] = "d"
// 		p.StatusCode = http.StatusCreated
// 		return ctx, rec, ctr, p
// 	}

// 	t.Run("missing name", func(t *testing.T) {
// 		// Rendering should fail if the Page has no name
// 		ctx, _, ctr, p := setup()
// 		p.Name = ""
// 		err := ctr.RenderPage(ctx, p)
// 		assert.Error(t, err)
// 	})

// 	t.Run("htmx rendering", func(t *testing.T) {
// 		ctx, _, ctr, p := setup()
// 		p.HTMX.Request.Enabled = true
// 		p.HTMX.Response = &htmx.Response{
// 			Trigger: "trigger",
// 		}
// 		err := ctr.RenderPage(ctx, p)
// 		require.NoError(t, err)

// 		// Check HTMX header
// 		assert.Equal(t, "trigger", ctx.Response().Header().Get(htmx.HeaderTrigger))

// 		// Check the template cache
// 		parsed, err := c.TemplateRenderer.Load("page:htmx", string(p.Name))
// 		require.NoError(t, err)

// 		// Check that all expected templates were parsed.
// 		// This includes the name, htmx and all components
// 		expectedTemplates := make(map[string]bool)
// 		expectedTemplates[fmt.Sprintf("%s%s", p.Name, config.TemplateExt)] = true
// 		expectedTemplates["htmx"+config.TemplateExt] = true
// 		components, err := templates.Get().ReadDir("components")
// 		require.NoError(t, err)
// 		for _, f := range components {
// 			expectedTemplates[f.Name()] = true
// 		}

// 		for _, v := range parsed.Template.Templates() {
// 			delete(expectedTemplates, v.Name())
// 		}
// 		assert.Empty(t, expectedTemplates)
// 	})

// 	t.Run("page cache", func(t *testing.T) {
// 		ctx, rec, ctr, p := setup()
// 		p.Cache.Enabled = true
// 		p.Cache.Tags = []string{"tag1"}
// 		err := ctr.RenderPage(ctx, p)
// 		require.NoError(t, err)

// 		// Fetch from the cache
// 		res, err := c.Cache.
// 			Get().
// 			Group(middleware.CachedPageGroup).
// 			Key(p.URL).
// 			Type(new(middleware.CachedPage)).
// 			Fetch(context.Background())
// 		require.NoError(t, err)

// 		// Compare the cached page
// 		cp, ok := res.(*middleware.CachedPage)
// 		require.True(t, ok)
// 		assert.Equal(t, p.URL, cp.URL)
// 		assert.Equal(t, p.Headers, cp.Headers)
// 		assert.Equal(t, p.StatusCode, cp.StatusCode)
// 		assert.Equal(t, rec.Body.Bytes(), cp.HTML)

// 		// Clear the tag
// 		err = c.Cache.
// 			Flush().
// 			Tags(p.Cache.Tags[0]).
// 			Execute(context.Background())
// 		require.NoError(t, err)

// 		// Refetch from the cache and expect no results
// 		_, err = c.Cache.
// 			Get().
// 			Group(middleware.CachedPageGroup).
// 			Key(p.URL).
// 			Type(new(middleware.CachedPage)).
// 			Fetch(context.Background())
// 		assert.Error(t, err)
// 	})
// }
