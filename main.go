package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"time"

	"github.com/kelseyhightower/envconfig"
	minio "github.com/minio/minio-go"
	"k8s.io/klog"
)

type Config struct {
	Host         string `default:"localhost"`
	Port         string `default:"5432"`
	Database     string `required:"true"`
	User         string `default:"root"`
	Password     string `default:"root"`
	PgDumpBinary string `default:"/usr/bin/pg_dump"`
	AccessKey    string `required:"true"`
	SecretKey    string `required:"true"`
	Endpoint     string `required:"true"`
	Bucket       string `required:"true"`
	Path         string `default:""`
	Prefix       string `default:""`
}

func dumpToFile(cfg *Config, target *os.File) error {
	cmd := exec.Command(cfg.PgDumpBinary, "-Fc", "-h", cfg.Host, "-p", cfg.Port, "-U", cfg.User, cfg.Database)
	cmd.Env = []string{"PGPASSWORD=" + cfg.Password}
	cmd.Stdout = target
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return err
	}
	return nil
}

const uploadDateFormat = "2006-01-02-15-04"

func uploadFileToS3(cfg *Config, source *os.File) error {
	// Reset pointer
	source.Seek(0, 0)
	stat, err := source.Stat()
	if err != nil {
		return err
	}
	// Setup S3 client
	client, err := minio.New(cfg.Endpoint, cfg.AccessKey, cfg.SecretKey, true)
	if err != nil {
		return err
	}
	// Upload file
	objectName := fmt.Sprintf("%s%s_%s", cfg.Prefix, cfg.Database, time.Now().Format(uploadDateFormat))
	_, err = client.PutObject(cfg.Bucket, filepath.Join(cfg.Path, objectName), source, stat.Size(), minio.PutObjectOptions{})
	if err != nil {
		return err
	}
	return nil
}

func main() {
	var cfg Config
	err := envconfig.Process("", &cfg)
	if err != nil {
		klog.Fatalf("config init failed: %v", err)
	}
	klog.Info("parsed environment configuration")
	// Create temporary on-disk backup
	tmpfile, err := ioutil.TempFile("", cfg.Database)
	if err != nil {
		klog.Fatalf("file init failed: %v", err)
	}
	klog.Infof("created tmpfile in %s", tmpfile.Name())
	defer func() {
		tmpname := tmpfile.Name()
		tmpfile.Close()
		os.Remove(tmpname)
		klog.Infof("deleted tmpfile %s", tmpname)
	}()
	// Execute backup
	if err := dumpToFile(&cfg, tmpfile); err != nil {
		klog.Fatalf("pgsql dump failed: %v", err)
	}
	klog.Info("dumped database on disk")
	// Upload backup to S3
	if err := uploadFileToS3(&cfg, tmpfile); err != nil {
		klog.Fatalf("s3 upload failed: %v", err)
	}
	klog.Info("backup successful")
}
