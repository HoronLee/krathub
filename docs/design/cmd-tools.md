# Servora 命令工具设计分析（更新版：`svr gen dao` Ent-only）

## 背景与本次决策

基于当前仓库实现与团队维护成本考量，本次将命令策略收敛为：

1. 服务级 GORM 生成命令已统一为 `make gen.gorm`（替代旧命名）。
2. `svr gen dao` 只支持 **Ent** 生成能力。
3. 中心化 CLI **不再支持 gorm-gen**；需要 gorm-gen 的服务继续在服务内自行维护（例如 `cmd/genDao`）。

---

## 结论

这个方向是正确的，原因很直接：

- Ent 生成是代码驱动、仓库内可静态判定，适合中心化 CLI。
- gorm-gen 依赖数据库连接与服务本地配置，中心化后会放大环境差异与失败面。
- 统一入口应优先处理“可确定、可重复、低外部依赖”的生成流程。

因此：

- `svr gen dao` = Ent 入口（统一规范）
- `make gen.gorm` / `cmd/genDao` = gorm-gen 本地入口（服务自治）

---

## 仓库证据（当前状态）

### 1) Make 目标关系

- `app.mk:129`：`gen: wire api openapi gen.ent`
- `app.mk:132`：`gen.gorm:`
- `app.mk:138`：`gen.ent:`

说明：当前 `make gen` 默认只串 Ent，不串 gorm-gen。

### 2) Ent 可检测标记（本仓）

- `app/servora/service/internal/data/schema/user.go`
- `app/servora/service/internal/data/generate.go:3`（包含 `entgo.io/ent/cmd/ent` 的 `go:generate`）

### 3) gorm-gen 本地入口（本仓）

- `app/servora/service/cmd/genDao/genDao.go:78`（`gen.NewGenerator(...)`）
- `app/servora/service/internal/data/gorm/dao/gen.go`

---

## `svr gen dao` 新设计（Ent-only）

## 命令契约

```bash
svr gen dao [--service <name>|--all] [--dry-run] [--fail-fast]
```

> 不再提供 `--orm` 选项。

## 识别逻辑（推荐）

采用单一入口条件即可：

1. 存在 `internal/data/generate.go` 且包含 `entgo.io/ent/cmd/ent`

满足该条件即判定该服务支持 `svr gen dao`。

## 执行映射

- Ent 服务：执行等价于服务目录下 `make gen.ent`
- 非 Ent 服务：输出明确提示并标记 `unsupported`

提示文案应包含：

- 检测失败原因（缺 generate.go 或 generate.go 不含 ent 生成入口）
- 如果该服务使用 gorm-gen，应该执行 `make gen.gorm` 或服务本地 `cmd/genDao`

---

## 与 gorm-gen 的职责边界

## 中心化（`svr`）

- 只负责 Ent 生成与批量调度。

## 服务本地（`make gen.gorm` / `cmd/genDao`）

- 负责数据库连接相关的 gorm-gen 生成。
- 由服务自己维护连接配置、模型生成策略与失败处理。

这样做的好处：

- 降低 CLI 复杂度与跨环境故障。
- 避免把数据库可用性耦合到中心化命令。
- 允许不同服务按自身数据库场景演进 gorm-gen 逻辑。

---

## 迁移策略

## Phase 1（已完成）

- [x] `gen.gorm` 目标命名统一
- [x] 文档命令引用同步至 `gen.gorm`

## Phase 2（进行中）

- [ ] `svr gen dao` 明确为 Ent-only
- [ ] 识别逻辑采用 `generate.go(ent)` 单条件
- [ ] 输出针对 gorm-gen 的“本地执行指引”

## Phase 3（稳定后）

- [ ] 评估是否将 `make gen.ent` 的中心化入口逐步切到 `svr gen dao`
- [ ] 保持 `make gen.gorm` 作为服务自治能力长期存在

---

## 风险与防护

1. **误把非 Ent 目录结构当成 Ent 服务**
   - 防护：仅以 `generate.go` 中的 ent 生成入口作为判定依据。

2. **开发者误以为 `svr gen dao` 还能处理 gorm-gen**
   - 防护：帮助文档和错误提示明确“Ent-only”。

3. **混合 ORM 服务认知混乱**
   - 防护：统一原则：Ent 走 `svr`，gorm-gen 走服务本地命令。

---

## 外部参考（用于支撑本决策）

- Ent 官方实践强调 `ent/schema` + `generate.go` 的生成入口模式。
- Kratos 风格 CLI 强调命令分组与工具职责清晰。
- GoFrame 的 `gen` 命令体系强调“生成命令统一入口”，但本项目按依赖特征将 DB-introspective 生成留在服务本地更稳妥。

---

## 最终建议

按你这次的决策执行即可：

- `svr gen dao` 收敛为 Ent-only（并明确文案）
- gorm-gen 坚持服务本地自治（`make gen.gorm` / `cmd/genDao`）

这会显著降低中心化 CLI 的长期维护复杂度，同时保留 gorm-gen 的灵活性。

---

## `svr` 命令行程序骨架设计（可扩展优先）

为支持后续持续新增命令（如 `svr new api`、`svr new svc`），`svr` 采用“命令分组 + 模块化注册 + 模板层解耦”的骨架。

### 设计目标

1. 新增子命令时避免修改大量已有代码。
2. 命令层只做参数校验和流程编排，业务逻辑下沉到独立模块。
3. 输出与错误码统一，方便接入 CI 与自动化脚本。

### 命令树（v1）

```bash
svr
├── gen
│   └── dao        # Ent-only：批量或单服务执行 ent 生成
├── new
│   ├── api        # 创建 API/proto 骨架
│   └── svc        # 创建微服务骨架（可选 with-ent）
├── version
└── help
```

### 目录骨架（建议）

```text
cmd/svr/
  main.go
  internal/
    root/
      root.go                # 根命令与全局 flags
    cmd/
      gen/
        gen.go               # 注册 gen 命令组
        dao.go               # Ent-only 生成入口
      new/
        new.go               # 注册 new 命令组
        api.go               # 创建 api/proto 骨架
        svc.go               # 创建服务骨架
    discovery/
      services.go            # 扫描 app/*/service
      ent.go                 # Ent 能力判定（generate.go 含 ent 入口）
    scaffold/
      renderer.go            # 模板渲染与文件输出
      templates/
        api/
        service/
    ux/
      output.go              # 统一输出（success/skip/fail）
      errors.go              # 统一错误码与提示文案
```

### 扩展机制（新增命令的标准方式）

每个命令组暴露统一注册函数：

```go
// 伪代码示意
func Register(parent *cobra.Command) {
    parent.AddCommand(NewCmd())
}
```

新增一个命令的流程固定为：

1. 在 `internal/cmd/<group>/` 新建命令文件。
2. 在对应组的 `Register(...)` 中注册。
3. 如涉及发现/脚手架能力，分别放入 `discovery/` 或 `scaffold/`，避免命令层膨胀。

### v1 命令契约（最小可用）

#### `svr gen dao`

- `svr gen dao [--service <name>|--all] [--dry-run] [--fail-fast]`
- 识别规则：仅 `internal/data/generate.go` 中包含 `entgo.io/ent/cmd/ent`
- 非 Ent 服务输出 `unsupported` 并给出 gorm 本地命令指引

#### `svr new api`

- `svr new api --name <name> [--version v1] [--http] [--grpc] [--dry-run]`
- 生成 `api/protos/<name>/service/<version>/` 基础 proto 骨架

#### `svr new svc`

- `svr new svc --name <name> [--with-ent] [--standalone] [--dry-run]`
- 生成 `app/<name>/service` 基础目录与 `Makefile(include ../../../app.mk)`
- `--with-ent` 时附带 `internal/data/generate.go` 与最小 schema 模板

### 与 Make 的协作边界（更新）

- Make 保持稳定流水线入口（`api/wire/openapi/ent`）。
- `svr` 负责高逻辑命令：服务发现、脚手架、策略执行。
- gorm-gen 继续服务本地自治（`make gen.gorm` / `cmd/genDao`）。

### 演进路线（补充）

#### v1

- [ ] 完成 `svr` 根命令骨架与 `gen/new` 分组
- [ ] 完成 `svr gen dao`（Ent-only）
- [ ] 完成 `svr new api` 与 `svr new svc` 的最小模板生成

#### v1.5

- [ ] 增加 `svr doctor`（工具链/目录规范检查）
- [ ] 增加统一 JSON 输出模式（便于 CI 解析）

#### v2

- [ ] 模板版本化与可配置模板源
- [ ] 可选插件机制（仅在确有需求时启用）
