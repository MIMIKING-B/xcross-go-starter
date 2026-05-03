package file

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/gogf/gf/v2/container/gvar"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/os/gcache"
	"github.com/gogf/gf/v2/os/gfile"
	"github.com/gogf/gf/v2/util/gconv"
)

type (
	// AdapterFile 基于文件系统实现的gcache缓存适配器
	// 实现了gcache.Adapter接口，将缓存数据持久化到文件中
	AdapterFile struct {
		dir string // 缓存文件存储的根目录
	}

	// fileContent 缓存文件内容结构体
	// 用于序列化存储到文件中的缓存数据
	fileContent struct {
		Duration int64       `json:"duration"`       // 缓存过期时间戳（秒），0表示永久有效
		Data     interface{} `json:"data,omitempty"` // 缓存的实际数据
	}
)

const perm = 0o666 // 文件权限：读写权限

// CacheExpiredErr 缓存过期错误变量
var (
	CacheExpiredErr = errors.New("cache expired")
)

// NewAdapterFile 创建并返回一个文件缓存适配器实例
// 参数：dir - 缓存文件存储的目录路径
func NewAdapterFile(dir string) gcache.Adapter {
	return &AdapterFile{
		dir: dir,
	}
}

// Set 设置缓存
// key：缓存键
// value：缓存值
// lifeTime：缓存有效期，小于0表示删除缓存，等于0表示永久有效
func (c *AdapterFile) Set(ctx context.Context, key interface{}, value interface{}, lifeTime time.Duration) (err error) {
	fileKey := gconv.String(key)
	// 如果值为nil或者有效期小于0，直接删除缓存
	if value == nil || lifeTime < 0 {
		return c.Delete(fileKey)
	}
	// 保存缓存数据到文件
	return c.Save(fileKey, gconv.String(value), lifeTime)
}

// SetMap 批量设置缓存（未实现）
func (c *AdapterFile) SetMap(ctx context.Context, data map[interface{}]interface{}, duration time.Duration) (err error) {
	return gerror.New("implement me")
}

// SetIfNotExist 键不存在时才设置缓存（未实现）
func (c *AdapterFile) SetIfNotExist(ctx context.Context, key interface{}, value interface{}, duration time.Duration) (ok bool, err error) {
	return false, gerror.New("implement me")
}

// SetIfNotExistFunc 键不存在时执行函数并设置缓存（未实现）
func (c *AdapterFile) SetIfNotExistFunc(ctx context.Context, key interface{}, f gcache.Func, duration time.Duration) (ok bool, err error) {
	return false, gerror.New("implement me")
}

// SetIfNotExistFuncLock 带锁的键不存在时执行函数并设置缓存（未实现）
func (c *AdapterFile) SetIfNotExistFuncLock(ctx context.Context, key interface{}, f gcache.Func, duration time.Duration) (ok bool, err error) {
	return false, gerror.New("implement me")
}

// Get 获取缓存数据
// 返回gvar.Var类型包装的值，方便类型转换
func (c *AdapterFile) Get(ctx context.Context, key interface{}) (*gvar.Var, error) {
	// 从文件中读取缓存数据
	fetch, err := c.Fetch(gconv.String(key))
	if err != nil {
		return nil, err
	}
	// 包装为gvar.Var返回
	return gvar.New(fetch), nil
}

// GetOrSet 获取缓存，不存在则设置默认值
func (c *AdapterFile) GetOrSet(ctx context.Context, key interface{}, value interface{}, duration time.Duration) (result *gvar.Var, err error) {
	// 先尝试获取缓存
	result, err = c.Get(ctx, key)
	// 非过期错误直接返回
	if err != nil && !errors.Is(err, CacheExpiredErr) {
		return nil, err
	}
	// 缓存不存在，设置新值并返回
	if result.IsNil() {
		return gvar.New(value), c.Set(ctx, key, value, duration)
	}
	return
}

// GetOrSetFunc 获取缓存，不存在则执行函数获取值并设置
func (c *AdapterFile) GetOrSetFunc(ctx context.Context, key interface{}, f gcache.Func, duration time.Duration) (result *gvar.Var, err error) {
	// 先尝试获取缓存
	v, err := c.Get(ctx, key)
	// 非过期错误直接返回
	if err != nil && !errors.Is(err, CacheExpiredErr) {
		return nil, err
	}
	// 缓存不存在，执行函数获取值
	if v.IsNil() {
		value, err := f(ctx)
		if err != nil {
			return nil, err
		}
		if value == nil {
			return nil, nil
		}
		// 设置新缓存并返回值
		return gvar.New(value), c.Set(ctx, key, value, duration)
	} else {
		return v, nil
	}
}

// GetOrSetFuncLock 带锁的获取/设置缓存，直接复用GetOrSetFunc
func (c *AdapterFile) GetOrSetFuncLock(ctx context.Context, key interface{}, f gcache.Func, duration time.Duration) (result *gvar.Var, err error) {
	return c.GetOrSetFunc(ctx, key, f, duration)
}

// Contains 检查缓存键是否存在
func (c *AdapterFile) Contains(ctx context.Context, key interface{}) (bool, error) {
	return c.Has(gconv.String(key)), nil
}

// Size 获取缓存数量（未实现）
func (c *AdapterFile) Size(ctx context.Context) (size int, err error) {
	return 0, nil
}

// Data 获取所有缓存数据（未实现）
func (c *AdapterFile) Data(ctx context.Context) (data map[interface{}]interface{}, err error) {
	return nil, gerror.New("implement me")
}

// Keys 获取所有缓存键（未实现）
func (c *AdapterFile) Keys(ctx context.Context) (keys []interface{}, err error) {
	return nil, gerror.New("implement me")
}

// Values 获取所有缓存值（未实现）
func (c *AdapterFile) Values(ctx context.Context) (values []interface{}, err error) {
	return nil, gerror.New("implement me")
}

// Update 更新缓存值（未实现）
func (c *AdapterFile) Update(ctx context.Context, key interface{}, value interface{}) (oldValue *gvar.Var, exist bool, err error) {
	return nil, false, gerror.New("implement me")
}

// UpdateExpire 更新缓存有效期
// 返回旧的有效期，duration<0表示删除缓存
func (c *AdapterFile) UpdateExpire(ctx context.Context, key interface{}, duration time.Duration) (oldDuration time.Duration, err error) {
	var (
		v       *gvar.Var
		oldTTL  int64
		fileKey = gconv.String(key)
	)
	// 获取原有效期
	expire, err := c.GetExpire(ctx, fileKey)
	if err != nil {
		return
	}
	oldTTL = int64(expire)
	// 缓存不存在，直接返回
	if oldTTL == -2 {
		return
	}
	oldDuration = time.Duration(oldTTL) * time.Second
	// 有效期小于0，删除缓存
	if duration < 0 {
		err = c.Delete(fileKey)
		return
	}
	// 获取原缓存值
	v, err = c.Get(ctx, fileKey)
	if err != nil {
		return
	}
	// 重新设置缓存（更新有效期）
	err = c.Set(ctx, fileKey, v.Val(), duration)
	return
}

// GetExpire 获取缓存剩余有效期
// 返回：剩余时间，不存在返回-1
func (c *AdapterFile) GetExpire(ctx context.Context, key interface{}) (time.Duration, error) {
	// 读取缓存文件
	content, err := c.read(gconv.String(key))
	if err != nil {
		return -1, nil
	}
	// 已过期
	if content.Duration <= time.Now().Unix() {
		return -1, nil
	}
	// 计算剩余有效期
	return time.Duration(content.Duration-time.Now().Unix()) * time.Second, nil
}

// Remove 删除指定的多个缓存键
// 返回最后一个被删除键的旧值
func (c *AdapterFile) Remove(ctx context.Context, keys ...interface{}) (lastValue *gvar.Var, err error) {
	if len(keys) == 0 {
		return nil, nil
	}
	// 获取最后一个键的旧值
	if lastValue, err = c.Get(ctx, gconv.String(keys[len(keys)-1])); err != nil {
		return nil, err
	}
	// 批量删除所有键
	err = c.DeleteMulti(gconv.Strings(keys)...)
	return
}

// Clear 清空所有缓存
func (c *AdapterFile) Clear(ctx context.Context) error {
	return c.Flush()
}

// Close 关闭缓存适配器（文件适配器无需关闭，空实现）
func (c *AdapterFile) Close(ctx context.Context) error {
	return nil
}

// createName 根据缓存键生成唯一的缓存文件名
// 使用sha256哈希处理键，避免特殊字符和文件名长度问题
func (c *AdapterFile) createName(key string) string {
	h := sha256.New()
	_, _ = h.Write([]byte(key))
	hash := hex.EncodeToString(h.Sum(nil))
	// 拼接缓存文件路径：目录/哈希值.cache
	return filepath.Join(c.dir, fmt.Sprintf("%s.cache", hash))
}

// read 读取缓存文件内容并解析，自动处理过期逻辑
// 过期会自动删除文件并返回CacheExpiredErr
func (c *AdapterFile) read(key string) (*fileContent, error) {
	// 获取文件真实路径
	rp := gfile.RealPath(c.createName(key))
	if rp == "" {
		return nil, nil
	}

	// 读取文件内容
	value, err := os.ReadFile(rp)
	if err != nil {
		return nil, err
	}

	// 反序列化JSON数据
	content := &fileContent{}
	if err := json.Unmarshal(value, content); err != nil {
		return nil, err
	}

	// 永久有效缓存，直接返回
	if content.Duration == 0 {
		return content, nil
	}

	// 检查是否过期
	if content.Duration <= time.Now().Unix() {
		_ = c.Delete(key)
		return nil, CacheExpiredErr
	}
	return content, nil
}

// Has 检查缓存键是否存在且有效
func (c *AdapterFile) Has(key string) bool {
	fc, err := c.read(key)
	return err == nil && fc != nil
}

// Delete 删除单个缓存文件
func (c *AdapterFile) Delete(key string) error {
	// 检查文件是否存在
	_, err := os.Stat(c.createName(key))
	if err != nil && os.IsNotExist(err) {
		return nil
	}
	// 删除文件
	return os.Remove(c.createName(key))
}

// DeleteMulti 批量删除多个缓存
func (c *AdapterFile) DeleteMulti(keys ...string) (err error) {
	for _, key := range keys {
		if err = c.Delete(key); err != nil {
			return
		}
	}
	return
}

// Fetch 根据键获取缓存数据
func (c *AdapterFile) Fetch(key string) (interface{}, error) {
	content, err := c.read(key)
	if err != nil {
		return nil, err
	}

	if content == nil {
		return nil, nil
	}
	return content.Data, nil
}

// FetchMulti 批量获取多个缓存
func (c *AdapterFile) FetchMulti(keys []string) map[string]interface{} {
	result := make(map[string]interface{})
	for _, key := range keys {
		if value, err := c.Fetch(key); err == nil {
			result[key] = value
		}
	}
	return result
}

// Flush 清空所有缓存文件
func (c *AdapterFile) Flush() error {
	// 打开缓存目录
	dir, err := os.Open(c.dir)
	if err != nil {
		return err
	}

	defer func() {
		_ = dir.Close()
	}()

	// 读取目录下所有文件名
	names, _ := dir.Readdirnames(-1)

	// 逐个删除文件
	for _, name := range names {
		_ = os.Remove(filepath.Join(c.dir, name))
	}
	return nil
}

// Save 将缓存数据写入文件
// key：缓存键
// value：缓存值
// lifeTime：缓存有效期
func (c *AdapterFile) Save(key string, value string, lifeTime time.Duration) error {
	duration := int64(0)

	// 计算过期时间戳
	if lifeTime > 0 {
		duration = time.Now().Unix() + int64(lifeTime.Seconds())
	}

	// 构造缓存内容结构体
	content := &fileContent{duration, value}

	// 序列化为JSON
	data, err := json.Marshal(content)
	if err != nil {
		return err
	}

	// 写入文件
	err = os.WriteFile(c.createName(key), data, perm)
	return err
}
