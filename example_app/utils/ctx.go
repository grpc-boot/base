package utils

import (
	"strconv"

	"github.com/gofiber/fiber/v3"
	"github.com/grpc-boot/base/v2/utils"
)

func CtxSet(ctx fiber.Ctx, key string, value any) {
	ctx.Context().SetUserValue(key, value)
}

func CtxGet(ctx fiber.Ctx, key string) any {
	return ctx.Context().UserValue(key)
}

func ClientIp(ctx fiber.Ctx) []byte {
	userIp := ctx.Request().Header.Peek("X-Real-IP")
	return userIp
}

func ArgsBytes(ctx fiber.Ctx, key string) []byte {
	if ctx.Context().IsPost() && FormExists(ctx, key) {
		return FormBytes(ctx, key)
	}

	return QueryBytes(ctx, key)
}

func ArgsString(ctx fiber.Ctx, key string) string {
	if ctx.Context().IsPost() && FormExists(ctx, key) {
		return FormString(ctx, key)
	}

	return QueryString(ctx, key)
}

func ArgsInt64(ctx fiber.Ctx, key string) int64 {
	if ctx.Context().IsPost() && FormExists(ctx, key) {
		return FormInt64(ctx, key)
	}

	return QueryInt64(ctx, key)
}

func ArgsExists(ctx fiber.Ctx, key string) bool {
	if ctx.Context().IsPost() && FormExists(ctx, key) {
		return FormExists(ctx, key)
	}

	return QueryExists(ctx, key)
}

func ArgsUintOrZero(ctx fiber.Ctx, key string) int {
	if ctx.Context().IsPost() && FormExists(ctx, key) {
		return FormUintOrZero(ctx, key)
	}

	return QueryUintOrZero(ctx, key)
}

func FormBytes(ctx fiber.Ctx, key string) []byte {
	return ctx.Context().PostArgs().Peek(key)
}

func FormString(ctx fiber.Ctx, key string) string {
	return utils.Bytes2String(FormBytes(ctx, key))
}

func FormInt64(ctx fiber.Ctx, key string) (value int64) {
	val := FormString(ctx, key)
	value, _ = strconv.ParseInt(val, 10, 64)
	return
}

func FormUintOrZero(ctx fiber.Ctx, key string) (value int) {
	return ctx.Context().PostArgs().GetUintOrZero(key)
}

func FormExists(ctx fiber.Ctx, key string) bool {
	return ctx.Context().PostArgs().Has(key)
}

func QueryBytes(ctx fiber.Ctx, key string) []byte {
	return ctx.Context().QueryArgs().Peek(key)
}

func QueryString(ctx fiber.Ctx, key string) string {
	return utils.Bytes2String(ctx.Context().QueryArgs().Peek(key))
}

func QueryInt64(ctx fiber.Ctx, key string) (value int64) {
	val := QueryString(ctx, key)
	value, _ = strconv.ParseInt(val, 10, 64)
	return
}

func QueryUintOrZero(ctx fiber.Ctx, key string) (value int) {
	return ctx.Context().QueryArgs().GetUintOrZero(key)
}

func QueryExists(ctx fiber.Ctx, key string) bool {
	return ctx.Context().QueryArgs().Has(key)
}
