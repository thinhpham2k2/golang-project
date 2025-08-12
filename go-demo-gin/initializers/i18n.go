package initializers

import (
	"embed"
	"io/fs"

	"github.com/nicksnyder/go-i18n/v2/i18n"
	"github.com/pelletier/go-toml/v2"
	"github.com/sirupsen/logrus"
	"golang.org/x/text/language"
)

// LocaleFS is the embedded filesystem containing locale files
// It contains all the .toml files in the locales directory
//
//go:embed locales/*.toml
var LocaleFS embed.FS

var Bundle *i18n.Bundle

func LoadI18n() error {
	Bundle = i18n.NewBundle(language.English)
	Bundle.RegisterUnmarshalFunc("toml", toml.Unmarshal)

	// Tìm tất cả file locale từ embed FS
	files, err := fs.Glob(LocaleFS, "locales/*.toml")
	if err != nil {
		logrus.WithField("source", "system").WithError(err).Error("Failed to glob locale files")
		return err
	}
	if len(files) == 0 {
		logrus.WithField("source", "system").Warn("No locale files found in embed FS")
	}

	for _, f := range files {
		data, err := LocaleFS.ReadFile(f)
		if err != nil {
			logrus.WithField("source", "system").WithError(err).
				Errorf("Failed to read language bundle: %s", f)
			return err
		}
		if _, err := Bundle.ParseMessageFileBytes(data, f); err != nil {
			logrus.WithField("source", "system").WithError(err).
				Errorf("Failed to parse language bundle: %s", f)
			return err
		}
		logrus.WithField("source", "system").
			Infof("Loaded language bundle %s successfully", f)
	}
	return nil
}
