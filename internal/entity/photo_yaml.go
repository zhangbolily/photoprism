package entity

import (
	"errors"
	"os"
	"path/filepath"
	"sync"
	"time"

	"gopkg.in/yaml.v2"

	"github.com/photoprism/photoprism/pkg/fs"
)

var photoYamlMutex = sync.Mutex{}
var photoYamlLockMap = make(map[string]bool, 0)

// Yaml returns photo data as YAML string.
func (m *Photo) Yaml() ([]byte, error) {
	// Load details if not done yet.
	m.GetDetails()

	m.CreatedAt = m.CreatedAt.UTC().Truncate(time.Second)
	m.UpdatedAt = m.UpdatedAt.UTC().Truncate(time.Second)

	out, err := yaml.Marshal(m)

	if err != nil {
		return []byte{}, err
	}

	return out, err
}

// SaveAsYaml saves photo data as YAML file.
func (m *Photo) SaveAsYaml(fileName string) error {
	data, err := m.Yaml()

	if err != nil {
		return err
	}

	// Make sure directory exists.
	if err := os.MkdirAll(filepath.Dir(fileName), fs.ModeDir); err != nil {
		return err
	}

	if _, ok := photoYamlLockMap[m.PhotoUID]; ok {
		return errors.New("photo: already locked")
	}

	photoYamlMutex.Lock()
	// Lock photo to prevent concurrent writes.
	photoYamlLockMap[m.PhotoUID] = true
	photoYamlMutex.Unlock()

	defer func() {
		photoYamlMutex.Lock()
		// Unlock photo.
		delete(photoYamlLockMap, m.PhotoUID)
		photoYamlMutex.Unlock()
	}()

	// Write YAML data to file.
	if err := os.WriteFile(fileName, data, fs.ModeFile); err != nil {
		return err
	}

	return nil
}

// LoadFromYaml photo data from a YAML file.
func (m *Photo) LoadFromYaml(fileName string) error {
	data, err := os.ReadFile(fileName)

	if err != nil {
		return err
	}

	if err := yaml.Unmarshal(data, m); err != nil {
		return err
	}

	return nil
}

// YamlFileName returns the YAML file name.
func (m *Photo) YamlFileName(originalsPath, sidecarPath string) string {
	return fs.FileName(filepath.Join(originalsPath, m.PhotoPath, m.PhotoName), sidecarPath, originalsPath, fs.ExtYAML)
}
