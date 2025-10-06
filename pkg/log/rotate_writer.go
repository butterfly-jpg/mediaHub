package log

import (
	"io"
	"sync"
	
	"gopkg.in/natefinch/lumberjack.v2"
)

// fileRotateWriter 管理多个日志文件的写入器
type fileRotateWriter struct {
	data         map[string]io.Writer // 存储不同日志路径对应的写入器
	sync.RWMutex                      // 读写锁，保证并发安全
}

// getWriter 根据指定日志路径获取写入器
func (frw *fileRotateWriter) getWriter(logPath string) io.Writer {
	// 读锁
	frw.RLock()
	defer frw.RUnlock()
	w, ok := frw.data[logPath]
	if !ok {
		return nil
	}
	return w
}

// setWriter 设置指定日志路径的写入器
func (frw *fileRotateWriter) setWriter(logPath string, w io.Writer) io.Writer {
	// 写锁
	frw.Lock()
	defer frw.Unlock()
	frw.data[logPath] = w
	return w
}

// 定义全局单例对象
var _fileRotateWriter *fileRotateWriter

// init 包初始化时创建全局的fileRotateWriter对象
func init() {
	_fileRotateWriter = &fileRotateWriter{
		data: map[string]io.Writer{},
	}
}

// GetRotateWriter 获取或创建指定日志路径的轮转写入器
func GetRotateWriter(logPath string) io.Writer {
	if logPath == "" {
		panic("日志文件路径不能为空")
	}
	writer := _fileRotateWriter.getWriter(logPath)
	if writer != nil {
		return writer
	}
	writer = &lumberjack.Logger{
		//文件名
		Filename: logPath,
		//单个文件大小单位MB
		MaxSize: 1,
		//最大保留时间（天）
		MaxAge: 7,
		//最多保留文件数
		MaxBackups: 15,
		LocalTime:  true,
	}
	return _fileRotateWriter.setWriter(logPath, writer)
}
