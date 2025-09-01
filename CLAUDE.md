# Go REST API 設計書（User API / Clean Architecture + DDD + DI / fx + gin / ALB→Lambda / Docker ローカル）

**目的**: テナントごとのユーザ情報を扱う REST API を、変更に強く・テスト容易な設計（Clean Architecture + DDD + DI）で実装する。実行基盤は **AWS ALB → Lambda**（Go）を本番想定。ローカルは **Docker** で動作確認。

---

## 1. スコープ / 要件

### 1.1 機能要件

* エンティティ: `User`
* **ユーザ取得（検索条件付き）のみ**: `user_name` / `email` の部分一致、ページネーション（`limit`/`offset`）
* **テナント別テーブル**: `users_{TENANT_ID}` に格納（テナント ID をパスから取得）

### 1.2 非機能要件

* 変更耐性: レイヤ分離 + インターフェイス駆動
* テスト容易性: モック差し替え可能
* セキュリティ: パスワードはハッシュ化、OTP キー保持
* デプロイ: 本番は ALB→Lambda、ローカルは Docker

---

## 2. データベース定義

```sql
CREATE TABLE "users_{TENANT_ID}" (
    id                    SERIAL                       PRIMARY KEY,
    user_name             TEXT                         UNIQUE NOT NULL,
    email                 TEXT                         UNIQUE NOT NULL,
    password_hash         TEXT                         NOT NULL,
    otp_secret_key        TEXT,
    type                  INTEGER                      NOT NULL,
    authority_data        TEXT,
    failed_count          INTEGER                      NOT NULL DEFAULT 0,
    unlock_at             TIMESTAMP WITHOUT TIME ZONE,
    is_reset_password     BOOLEAN                      DEFAULT FALSE,
    last_update_password  TIMESTAMP WITHOUT TIME ZONE  NOT NULL DEFAULT now()
);
```

---

## 3. 全体アーキテクチャ

```text
[Client]
   │
HTTPS
   │
[Amazon ALB] ──(Lambda Target Group)──> [AWS Lambda (Go)]
                                         │
                                         └─> [Amazon RDS/Aurora PostgreSQL]
```

---

## 4. ディレクトリ構成

````text
project/
└─ src/
   ├─ domain/
   │   ├─ user.go               # Entity / validation
   │   └─ repository.go         # UserRepository インターフェイス
   ├─ usecase/
   │   └─ user_search.go        # 検索ユースケースのみ
   ├─ interface/web/
   │   ├─ controller/user_controller.go
   │   ├─ middleware/logging.go
   │   └─ router.go
   ├─ infrastructure/
   │   ├─ datasource/pg.go
   │   ├─ repository/user_pg.go
   │   ├─ http/response.go
   │   └─ log/logger.go
   ├─ runtime/
   │   ├─ config.go
   │   ├─ dependency.go
   │   ├─ router.go
   │   ├─ start_server.go
   │   └─ start_lambda.go
   ├─ main.go
   └─ internal/testutil/
```text
project/
└─ src/
   ├─ domain/
   │   ├─ user.go               # Entity / validation
   │   └─ repository.go         # UserRepository インターフェイス
   ├─ usecase/
   │   ├─ user_create.go
   │   ├─ user_get.go
   │   ├─ user_update.go
   │   ├─ user_delete.go
   │   └─ user_search.go
   ├─ interface/web/
   │   ├─ controller/user_controller.go
   │   ├─ middleware/logging.go
   │   └─ router.go
   ├─ infrastructure/
   │   ├─ datasource/pg.go
   │   ├─ repository/user_pg.go
   │   ├─ http/response.go
   │   └─ log/logger.go
   ├─ runtime/
   │   ├─ config.go
   │   ├─ dependency.go
   │   ├─ router.go
   │   ├─ start_server.go
   │   └─ start_lambda.go
   ├─ main.go
   └─ internal/testutil/
````

---

## 5. ドメイン層

### user.go

```go
package domain

type UserID int

type User struct {
    ID                 UserID
    UserName           string
    Email              string
    PasswordHash       string
    OtpSecretKey       *string
    Type               int
    AuthorityData      *string
    FailedCount        int
    UnlockAt           *time.Time
    IsResetPassword    bool
    LastUpdatePassword time.Time
}
```

### repository.go

```go
package domain

type UserRepository interface {
    Create(ctx context.Context, tenantID string, u *User) error
    FindByID(ctx context.Context, tenantID string, id UserID) (*User, error)
    Update(ctx context.Context, tenantID string, u *User) error
    Delete(ctx context.Context, tenantID string, id UserID) error
    Search(ctx context.Context, tenantID string, userName, email string, limit, offset int) ([]*User, int, error)
}
```

---

## 6. ユースケース層

### user\_search.go

```go
package usecase

type UserSearchIn struct {
    TenantID string
    UserName string
    Email    string
    Limit    int
    Offset   int
}

type UserSearchUsecase struct { repo domain.UserRepository }

func NewUserSearchUsecase(r domain.UserRepository) *UserSearchUsecase {
    return &UserSearchUsecase{repo: r}
}

func (uc *UserSearchUsecase) Do(ctx context.Context, in UserSearchIn) ([]*domain.User, int, error) {
    return uc.repo.Search(ctx, in.TenantID, in.UserName, in.Email, in.Limit, in.Offset)
}
```

---

## 7. プレゼンテーション層（gin）

### controller/user\_controller.go

```go
package controller

type UserController struct {
    searchUC *usecase.UserSearchUsecase
}

func NewUserController(su *usecase.UserSearchUsecase) *UserController {
    return &UserController{searchUC: su}
}

func (uc *UserController) Search(c *gin.Context) {
    tenantID := c.Param("tenant_id")
    in := usecase.UserSearchIn{
        TenantID: tenantID,
        UserName: c.Query("user_name"),
        Email:    c.Query("email"),
        Limit:    parseInt(c.Query("limit"), 20),
        Offset:   parseInt(c.Query("offset"), 0),
    }
    users, total, err := uc.searchUC.Do(c, in)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    c.JSON(http.StatusOK, gin.H{"items": users, "total": total})
}
```

### router.go

````go
func NewRouter(uc *controller.UserController) *gin.Engine {
  r := gin.New()
  r.Use(gin.Recovery())

  v1 := r.Group("/:tenant_id")
  {
    // 仕様: {tenant_id}/Users ユーザ取得（検索条件付き）
    v1.GET("/Users", uc.Search)
  }
  return r
}
```go
func NewRouter(uc *controller.UserController) *gin.Engine {
  r := gin.New()
  r.Use(gin.Recovery())

  v1 := r.Group(":tenant_id")
  {
    v1.GET("/users", uc.Search)
    // POST / PUT / DELETE も同様に追加予定
  }
  return r
}
````

---

## 8. インフラストラクチャ層

### datasource/pg.go

```go
func NewPgPool(cfg runtime.DBConfig) (*pgxpool.Pool, error) {
  ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
  defer cancel()

  pool, err := pgxpool.New(ctx, cfg.DSN())
  if err != nil { return nil, err }
  if err := pool.Ping(ctx); err != nil { return nil, err }
  return pool, nil
}
```

### repository/user\_pg.go

```go
type UserPg struct { db *pgxpool.Pool }

func NewUserPg(db *pgxpool.Pool) *UserPg { return &UserPg{db: db} }

func (r *UserPg) Search(ctx context.Context, tenantID, userName, email string, limit, offset int) ([]*domain.User, int, error) {
  table := fmt.Sprintf("users_%s", tenantID)
  where := []string{"1=1"}
  args := []any{}

  if userName != "" {
    where = append(where, fmt.Sprintf("user_name ILIKE '%%' || $%d || '%%'", len(args)+1))
    args = append(args, userName)
  }
  if email != "" {
    where = append(where, fmt.Sprintf("email ILIKE '%%' || $%d || '%%'", len(args)+1))
    args = append(args, email)
  }

  countSQL := fmt.Sprintf("SELECT COUNT(1) FROM %s WHERE %s", table, strings.Join(where, " AND "))
  var total int
  if err := r.db.QueryRow(ctx, countSQL, args...).Scan(&total); err != nil { return nil, 0, err }

  listSQL := fmt.Sprintf("SELECT id,user_name,email,type,authority_data,failed_count,unlock_at,is_reset_password,last_update_password FROM %s WHERE %s LIMIT $%d OFFSET $%d", table, strings.Join(where, " AND "), len(args)+1, len(args)+2)
  args = append(args, limit, offset)

  rows, err := r.db.Query(ctx, listSQL, args...)
  if err != nil { return nil, 0, err }
  defer rows.Close()

  var items []*domain.User
  for rows.Next() {
    var u domain.User
    if err := rows.Scan(&u.ID, &u.UserName, &u.Email, &u.Type, &u.AuthorityData, &u.FailedCount, &u.UnlockAt, &u.IsResetPassword, &u.LastUpdatePassword); err != nil {
      return nil, 0, err
    }
    items = append(items, &u)
  }
  return items, total, nil
}
```

---

## 9. runtime / fx 配線

### dependency.go

```go
var Module = fx.Options(
  fx.Provide(
    runtime.LoadDBConfig,
    datasource.NewPgPool,
    repository.NewUserPg,
    usecase.NewUserSearchUsecase,
    controller.NewUserController,
    web.NewRouter,
  ),
  fx.Invoke(runtime.Start),
)
```

### main.go

```go
func main() { fx.New(Module).Run() }
```

---

## 10. API 仕様（ユーザ取得のみ）

**パス形式**: `/{tenant_id}/Users`

**クエリ**: `user_name`（部分一致）, `email`（部分一致）, `limit`（既定 20）, `offset`（既定 0）

| Method | Path                | 説明            |
| -----: | ------------------- | ------------- |
|    GET | /{tenant\_id}/Users | ユーザ取得（検索条件付き） |

\-------:|--------------------------|------|
\| GET    | /{tenant\_id}/users       | ユーザ検索（条件付き）|
\| POST   | /{tenant\_id}/users       | ユーザ作成 |
\| GET    | /{tenant\_id}/users/\:id   | ユーザ取得 |
\| PUT    | /{tenant\_id}/users/\:id   | ユーザ更新 |
\| DELETE | /{tenant\_id}/users/\:id   | ユーザ削除 |

---

## 11. 設定/環境変数（必須のみ）

| 変数            | 例          | 用途      |
| ------------- | ---------- | ------- |
| `DB_HOST`     | `db`       | DB ホスト名 |
| `DB_PORT`     | `5432`     | ポート番号   |
| `DB_USER`     | `postgres` | 接続ユーザー  |
| `DB_PASSWORD` | `postgres` | 接続パスワード |
| `DB_NAME`     | `app`      | データベース名 |
| `DB_SSLMODE`  | `disable`  | SSL モード |

---

### 本設計のキーメッセージ

* **スコープはユーザ機能（User）のみ**に限定
* **テナントごとに users\_{TENANT\_ID} テーブルを利用**
* **インターフェイス駆動**で疎結合、**fx** で依存解決
* **gin** はプレゼンテーションに限定
* 本番は **ALB→Lambda**、ローカルは **Docker** で同一コードを実行
