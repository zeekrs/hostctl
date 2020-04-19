package file

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/guumaster/hostctl/pkg/host"
	"github.com/guumaster/hostctl/pkg/host/errors"
)

func TestFile_AddProfile(t *testing.T) {
	mem := CreateBasicFS(t)
	f, err := mem.Open("/tmp/etc/hosts")
	assert.NoError(t, err)

	t.Run("Add new", func(t *testing.T) {
		m, err := NewWithFs(f.Name(), mem)
		assert.NoError(t, err)
		r := strings.NewReader(`127.0.0.1 added.loc`)

		p, err := host.NewProfileFromReader(r, true)
		assert.NoError(t, err)
		p.Name = "awesome"
		p.Status = host.Enabled

		err = m.AddProfile(p)
		assert.NoError(t, err)

		assert.Equal(t, []string{"profile1", "awesome"}, m.GetEnabled())

		added, err := m.GetProfile("awesome")
		assert.NoError(t, err)

		assert.Equal(t, added, p)
	})

	t.Run("Add existing", func(t *testing.T) {
		m, err := NewWithFs(f.Name(), mem)
		assert.NoError(t, err)
		r := strings.NewReader(`127.0.0.1 added.loc`)

		p, err := host.NewProfileFromReader(r, true)
		assert.NoError(t, err)
		p.Name = "profile1"

		err = m.AddProfile(p)
		assert.NoError(t, err)

		assert.Equal(t, []string{"profile1"}, m.GetEnabled())

		added, err := m.GetProfile("profile1")
		assert.NoError(t, err)
		hosts, err := added.GetHostNames(Localhost.String())
		assert.NoError(t, err)

		assert.Equal(t, hosts, []string{"first.loc", "second.loc", "added.loc"})
	})

	t.Run("Add default error", func(t *testing.T) {
		m, err := NewWithFs(f.Name(), mem)
		assert.NoError(t, err)
		r := strings.NewReader(`127.0.0.1 added.loc`)

		p, err := host.NewProfileFromReader(r, true)
		assert.NoError(t, err)
		p.Name = "default"

		err = m.AddProfile(p)
		assert.Error(t, err, errors.ErrDefaultProfile)
	})
}
