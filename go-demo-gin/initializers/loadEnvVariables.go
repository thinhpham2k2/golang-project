package initializers

import (
	"fmt"
	"os"
	"strings"

	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
)

func LoadEnvVariables(paths ...string) {
	// Dev: cố gắng nạp .env (không có thì cảnh báo, không fatal)
	if len(paths) == 0 {
		paths = []string{".env"} // bạn có thể thêm ".env.local" tùy ý
	}
	loadedAny := false
	for _, p := range paths {
		if err := godotenv.Load(p); err == nil {
			logrus.WithField("source", "system").Infof("Loaded env file: %s", p)
			loadedAny = true
		}
	}
	if !loadedAny {
		logrus.WithField("source", "system").
			Warn("No .env file loaded; using OS environment only")
	}
}

// Nếu thiếu biến bắt buộc → trả lỗi để main quyết định dừng
func RequireEnv(keys ...string) error {
	var missing []string
	for _, k := range keys {
		if os.Getenv(k) == "" {
			missing = append(missing, k)
		}
	}
	if len(missing) > 0 {
		return fmt.Errorf("missing required env: %s", strings.Join(missing, ", "))
	}
	return nil
}
