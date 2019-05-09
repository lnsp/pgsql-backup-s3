package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"time"

	"github.com/kelseyhightower/envconfig"
	minio "github.com/minio/minio-go"
)

type Config struct {
	Container string `required:"true"`
	Database  string `required:"true"`
	User      string `default:"root"`
	AccessKey string `required:"true"`
	SecretKey string `required:"true"`
	Endpoint  string `required:"true"`
	Bucket    string `required:"true"`
	Path      string `default:""`
}

func main() {
	var cfg Config
	err := envconfig.Process("", &cfg)
	if err != nil {
		log.Fatalf("init failed: %v", err)
	}
	// Create temporary on-disk backup
	tmpfile, err := ioutil.TempFile("", cfg.Container)
	if err != nil {
		log.Fatalf("file init failed: %v", err)
	}
	defer func() {
		tmpname := tmpfile.Name()
		tmpfile.Close()
		os.Remove(tmpname)
	}()
	// Execute backup
	cmd := exec.Command("docker", "exec", cfg.Container, "pg_dump", "-U", cfg.User, cfg.Database)
	cmd.Stdout = tmpfile
	if err := cmd.Run(); err != nil {
		log.Fatalf("backup failed: %v", err)
	}
	// Reset read/write offset
	tmpfile.Seek(0, 0)
	fileInfo, err := tmpfile.Stat()
	if err != nil {
		log.Fatalf("failed to get fileinfo: %v", err)
	}
	// Setup S3 client
	client, err := minio.New(cfg.Endpoint, cfg.AccessKey, cfg.SecretKey, true)
	if err != nil {
		log.Fatalf("failed to connect: %v", err)
	}
	// Upload file
	objectName := fmt.Sprintf("%s_%s.sql", cfg.Container, time.Now().Format("2006-01-02-15-04"))
	_, err = client.PutObject(cfg.Bucket, filepath.Join(cfg.Path, objectName), tmpfile, fileInfo.Size(), minio.PutObjectOptions{})
	if err != nil {
		log.Fatalf("failed to upload: %v", err)
	}
}
