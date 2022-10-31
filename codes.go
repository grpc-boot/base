package base

import "google.golang.org/grpc/codes"

const (
	// OK 成功返回
	OK = codes.OK

	// ErrCanceled 表示操作已取消（通常由调用者取消）
	ErrCanceled = codes.Canceled

	// ErrUnknown 未知错误
	ErrUnknown = codes.Unknown

	// ErrInvalidArgument 客户端传递了一个无效的参数
	ErrInvalidArgument = codes.InvalidArgument

	// ErrDeadlineExceeded 表示操作在完成前过期。对于改变系统状态的操作，即使操作成功完成，也可能会返回此错误。
	// 例如，来自服务器的成功响应可能已延迟足够长的时间以使截止日期到期
	ErrDeadlineExceeded = codes.DeadlineExceeded

	// ErrNotFound 表示未找到某些请求的实体（例如文件或目录）
	ErrNotFound = codes.NotFound

	// ErrAlreadyExists 表示创建实体的尝试失败，因为实体已经存在
	ErrAlreadyExists = codes.AlreadyExists

	// ErrPermissionDenied 表示调用者没有执行指定操作的权限
	ErrPermissionDenied = codes.PermissionDenied

	// ErrResourceExhausted 表示某些资源已用尽，可能是每个用户的配额，或者可能是整个文件系统空间不足
	ErrResourceExhausted = codes.ResourceExhausted

	// ErrFailedPrecondition 操作被拒绝，因为系统未处于操作执行所需的状态。
	// 例如，要删除的目录可能是非空的，rmdir 操作应用于非目录等
	ErrFailedPrecondition = codes.FailedPrecondition

	// ErrAborted 操作已取消，通常是由于并发问题，如排序器检查失败、事务中止
	ErrAborted = codes.Aborted

	// ErrOutRange 超出范围的操作
	ErrOutRange = codes.OutOfRange

	// ErrUnimplemented 未实现或不支持的操作
	ErrUnimplemented = codes.Unimplemented

	// ErrInternal 系统内部错误
	ErrInternal = codes.Internal

	// ErrUnavailable 表示该服务当前不可用。这很可能是一种瞬态情况，可以通过回退重试来纠正。
	// 请注意，重试非幂等操作是不安全的
	ErrUnavailable = codes.Unavailable

	// ErrDataLoss 表示不可恢复的数据丢失或损坏
	ErrDataLoss = codes.DataLoss

	// ErrUnauthenticated 表示请求没有用于操作的有效身份验证凭据
	ErrUnauthenticated = codes.Unauthenticated
)

var defaultCodeMsg = map[codes.Code]string{
	OK:                    "OK",
	ErrCanceled:           "CANCELLED",
	ErrUnknown:            "UNKNOWN",
	ErrInvalidArgument:    "INVALID_ARGUMENT",
	ErrDeadlineExceeded:   "DEADLINE_EXCEEDED",
	ErrNotFound:           "NOT_FOUND",
	ErrAlreadyExists:      "ALREADY_EXISTS",
	ErrPermissionDenied:   "PERMISSION_DENIED",
	ErrResourceExhausted:  "RESOURCE_EXHAUSTED",
	ErrFailedPrecondition: "FAILED_PRECONDITION",
	ErrAborted:            "ABORTED",
	ErrOutRange:           "OUT_OF_RANGE",
	ErrUnimplemented:      "UNIMPLEMENTED",
	ErrInternal:           "INTERNAL",
	ErrUnavailable:        "UNAVAILABLE",
	ErrDataLoss:           "DATA_LOSS",
	ErrUnauthenticated:    "UNAUTHENTICATED",
}
