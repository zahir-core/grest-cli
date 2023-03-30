package app

import (
	"context"
	"io"
	"os"

	"github.com/jinzhu/copier"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

func FS() FSInterface {
	if fsClient == nil {
		fsClient = fsUtil{}.configure()
	}
	return fsClient
}

type FSInterface interface {
	GetFileUrl(fileName string, path ...string) string
	Upload(fileName string, src io.Reader, fileSize int64, opts ...FileUploadOption) (FileUploadInfo, error)
	Delete(fileName string, opts ...FileDeleteOption) error
}

var fsClient *fsUtil

type fsUtil struct {
	Driver        string
	LocalDirPath  string
	PublicDirPath string
	EndPoint      string
	Port          int
	Region        string
	BucketName    string
	AccessKey     string
	SecretKey     string
	mClient       *minio.Client
	ctx           context.Context
	err           error
}

func (f fsUtil) configure() *fsUtil {
	if FS_DRIVER == "local" {
		return f.setLocalFS()
	}

	// Cloud Storage Like AWS S3, Google Cloud Storage, etc
	f.Driver = FS_DRIVER
	f.EndPoint = FS_END_POINT
	f.Port = FS_PORT
	f.Region = FS_REGION
	f.BucketName = FS_BUCKET_NAME
	f.AccessKey = FS_ACCESS_KEY
	f.SecretKey = FS_SECRET_KEY
	f.mClient, f.err = minio.New(f.EndPoint, &minio.Options{
		Creds:  credentials.NewStaticV4(f.AccessKey, f.SecretKey, ""),
		Secure: true,
	})
	if f.err != nil {
		return f.setLocalFS()
	}
	f.ctx = context.Background()
	isBucketExists := false
	isBucketExists, f.err = f.mClient.BucketExists(f.ctx, f.BucketName)
	if !isBucketExists {
		f.err = f.mClient.MakeBucket(f.ctx, f.BucketName, minio.MakeBucketOptions{Region: f.Region})
	}
	if f.err != nil {
		return f.setLocalFS()
	}
	return &f
}

func (f fsUtil) setLocalFS() *fsUtil {
	if f.err != nil {
		Logger().Error().Msg(f.err.Error() + ", local filesystem will be used.")
	}
	f.Driver = "local"
	f.LocalDirPath = FS_LOCAL_DIR_PATH
	if f.LocalDirPath == "" {
		f.LocalDirPath = "storages"
	}
	f.PublicDirPath = FS_PUBLIC_DIR_PATH
	if f.PublicDirPath == "" {
		f.PublicDirPath = "storages"
	}
	f.createLocalDirPath()
	return &f
}

func (f fsUtil) createLocalDirPath() {
	_, err := os.Stat(f.LocalDirPath)
	if os.IsNotExist(err) {
		err = os.Mkdir(f.LocalDirPath, 0755)
		if err != nil {
			Logger().Error().
				Err(err).
				Str("FS_DRIVER", FS_DRIVER).
				Str("LocalDirPath", f.LocalDirPath).
				Msg("Failed to create local dir path.")
		}
	}
}

func (f *fsUtil) GetFileUrl(fileName string, path ...string) string {
	res := APP_URL + "/" + f.PublicDirPath
	if f.EndPoint == "s3.amazonaws.com" {
		res = "https://" + f.BucketName + ".s3." + f.Region + ".amazonaws.com/"
	} else if f.Driver != "local" {
		res = "https://" + f.BucketName + "." + f.EndPoint + "/"
	}

	// cek apakah di cloudkilat support multiple path apa g, kalo g support berarti nda lewat ini
	for _, p := range path {
		res += p + "/"
	}
	res += fileName

	return res
}

func (f *fsUtil) Upload(fileName string, src io.Reader, fileSize int64, opts ...FileUploadOption) (FileUploadInfo, error) {
	if f.Driver == "local" {
		dst, err := os.Create(f.LocalDirPath + "/" + fileName)
		if err != nil {
			return FileUploadInfo{}, nil
		}
		defer dst.Close()
		_, err = io.Copy(dst, src)
		return FileUploadInfo{}, err
	}

	opt := minio.PutObjectOptions{}
	opt.UserMetadata = map[string]string{"x-amz-acl": "public-read"}
	if len(opts) > 0 {
		copier.Copy(opt, opts[0])
	}
	info, err := f.mClient.PutObject(f.ctx, f.BucketName, fileName, src, fileSize, opt)
	return FileUploadInfo{info}, err
}

func (f *fsUtil) Delete(fileName string, opts ...FileDeleteOption) error {
	if f.Driver == "local" {
		return os.Remove(f.LocalDirPath + "/" + fileName)
	}

	opt := minio.RemoveObjectOptions{}
	opt.GovernanceBypass = true
	if len(opts) > 0 {
		copier.Copy(opt, opts[0])
	}
	return f.mClient.RemoveObject(f.ctx, f.BucketName, fileName, opt)
}

type FileUploadOption struct {
	minio.PutObjectOptions
}

type FileUploadInfo struct {
	minio.UploadInfo
}

type FileDeleteOption struct {
	minio.RemoveObjectOptions
}
