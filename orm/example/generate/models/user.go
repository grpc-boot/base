package models

import (
	"context"
	"database/sql"

	"github.com/grpc-boot/base/v2/orm"
	"github.com/grpc-boot/base/v2/orm/basis"
	"github.com/grpc-boot/base/v2/orm/condition"
	"github.com/grpc-boot/base/v2/orm/sqlite"
)

var (
	DefaultUser = &User{}
	UserMapping *basis.Mapping
)

func init() {
	// 预加载mapping
	UserMapping = basis.ParseMapping(DefaultUser)
}

type User struct {
	Id        int64   `json:"id" field:"id" size:"0" label:""`
	UserName  string  `json:"userName" field:"user_name" size:"64" label:""`
	HeadImg   string  `json:"headImg" field:"head_img" size:"0" label:""`
	CreatedAt uint64  `json:"createdAt" field:"created_at" size:"0" label:""`
	Status    int8    `json:"status" field:"status" size:"0" label:""`
	Amount    float64 `json:"amount" field:"amount" size:"0" label:""`
	Remark    string  `json:"remark" field:"remark" size:"0" label:""`
}

func (u *User) TableName() string {
	return "user"
}

func (u *User) PrimaryField() string {
	return "id"
}

func (u *User) PrimaryValue() any {
	return u.Id
}

func (u *User) GetDefault(name string) string {
	return UserMapping.Default(name)
}

func (u *User) GetLabel(name string) string {
	return UserMapping.Label(name)
}

func (u *User) GetEnums(name string) []string {
	return UserMapping.Enums(name)
}

func (u *User) GetSize(name string) int {
	return UserMapping.Size(name)
}

func (u *User) Query() *sqlite.Query {
	return sqlite.AcquireQuery()
}

func (u *User) Deleter() basis.Delete {
	return sqlite.Delete
}

func (u *User) Inserter() basis.Insert {
	return sqlite.Insert
}

func (u *User) Updater() basis.Update {
	return sqlite.Update
}

func (u *User) Add(ctx context.Context, executor basis.Executor, user User) (lastInsertId int64, err error) {
	return orm.InsertRow(ctx, executor, user.TableName(), basis.Columns{
		UserMapping.Field("UserName"),
		UserMapping.Field("HeadImg"),
		UserMapping.Field("CreatedAt"),
		UserMapping.Field("Status"),
		UserMapping.Field("Amount"),
		UserMapping.Field("Remark"),
	}, basis.Row{
		user.UserName,
		user.HeadImg,
		user.CreatedAt,
		user.Status,
		user.Amount,
		user.Remark,
	}, false, u.Inserter())
}

func (u *User) FindOne(ctx context.Context, executor basis.Executor, user User) (User, error) {
	qObj := u.Query().
		From(user.TableName()).
		Where(condition.Equal{
			Field: "id",
			Value: user.Id,
		})
	defer qObj.Close()

	list, err := u.FindByQuery(ctx, executor, qObj)
	if err != nil || len(list) == 0 {
		return User{}, err
	}

	return list[0], nil
}

func (u *User) FindSome(ctx context.Context, executor basis.Executor, idList ...any) ([]User, error) {
	qObj := u.Query().
		Where(condition.In[any]{
			Field: "id",
			Value: idList,
		})
	defer qObj.Close()

	return u.FindByQuery(ctx, executor, qObj)
}

func (u *User) FindByQuery(ctx context.Context, executor basis.Executor, query basis.Query) ([]User, error) {
	if !query.HasFrom() {
		query.From(u.TableName())
	}

	rows, err := orm.QueryWithQuery(ctx, executor, query)
	if err != nil {
		return nil, err
	}

	list, err := u.scanRows(rows)
	if err != nil || len(list) == 0 {
		return nil, err
	}

	return list, err
}

func (u *User) scanRows(rows *sql.Rows) ([]User, error) {
	records, err := basis.ScanRecords(rows)
	if err != nil {
		return nil, err
	}

	list := make([]User, len(records))
	for index, row := range records {
		list[index] = User{
			Id:        row.Int64(UserMapping.Field("Id")),
			UserName:  row.String(UserMapping.Field("UserName")),
			HeadImg:   row.String(UserMapping.Field("HeadImg")),
			CreatedAt: row.Uint64(UserMapping.Field("CreatedAt")),
			Status:    row.Int8(UserMapping.Field("Status")),
			Amount:    row.Float64(UserMapping.Field("Amount")),
			Remark:    row.String(UserMapping.Field("Remark")),
		}
	}

	return list, nil
}
