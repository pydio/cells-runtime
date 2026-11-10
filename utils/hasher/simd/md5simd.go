package simd

import (
	"crypto/md5"
	"hash"
	"os"
	"sync"

	md5simd "github.com/minio/md5-simd"

	"github.com/pydio/cells-runtime/runtime"
)

var (
	mServer     md5simd.Server
	mServerOnce sync.Once
)

func init() {
	runtime.RegisterApplicationEnvVariable("ENABLE_SIMDMD5", "", "Empty is false by default, if set this will switch the md5 hasher to simd implementation", true)
}

func MD5() hash.Hash {
	if os.Getenv(runtime.EnvironmentName("ENABLE_SIMDMD5")) != "" {
		mServerOnce.Do(func() {
			mServer = md5simd.NewServer()
		})
		return mServer.NewHash()
	}
	return md5.New()
}
