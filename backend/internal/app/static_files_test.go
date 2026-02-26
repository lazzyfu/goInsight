package app

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
)

func TestSetupStaticFiles_AllowsMissingAssets(t *testing.T) {
	gin.SetMode(gin.TestMode)

	wd, err := os.Getwd()
	require.NoError(t, err)
	t.Cleanup(func() {
		require.NoError(t, os.Chdir(wd))
	})

	tmp := t.TempDir()
	require.NoError(t, os.Chdir(tmp))

	r := gin.New()
	err = setupStaticFiles(r)
	require.NoError(t, err)

	_, err = os.Stat(filepath.Join(tmp, "media"))
	require.NoError(t, err)
}
