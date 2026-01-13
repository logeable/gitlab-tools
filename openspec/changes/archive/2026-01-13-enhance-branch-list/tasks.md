## 1. 实现 --hide-empty 参数
- [x] 1.1 在 `main.go` 中定义 `branchListHideEmpty` 布尔变量
- [x] 1.2 在 `init()` 函数中为 `branchListCmd` 添加 `--hide-empty` 标志
- [x] 1.3 在单个项目模式下，如果启用 `--hide-empty` 且没有分支，则直接返回
- [x] 1.4 在多项目模式下，如果启用 `--hide-empty` 且项目没有分支，则跳过该项目

## 2. 实现 --quiet 参数
- [x] 2.1 在 `main.go` 中定义 `branchListQuiet` 布尔变量
- [x] 2.2 在 `init()` 函数中为 `branchListCmd` 添加 `--quiet` 标志
- [x] 2.3 修改 `printBranchesList` 函数，添加 `quiet` 参数
- [x] 2.4 在 quiet 模式下，只显示项目名（每行一个），仅当项目有分支时显示
- [x] 2.5 在多项目模式下，quiet 时不在项目之间添加空行

## 3. 添加提交者信息显示
- [x] 3.1 在 `printBranchesList` 函数中，当分支有提交信息时，显示提交者姓名（如果存在）
- [x] 3.2 在 `printBranchesList` 函数中，显示提交者邮箱（如果存在）

## 4. 优化输出格式
- [x] 4.1 调整多项目模式下空行分隔的位置，在检查 `--hide-empty` 之后添加

## 5. 验证
- [x] 5.1 测试 `--hide-empty` 参数：`gitlab-tools branch list --hide-empty`（应隐藏没有分支的项目）
- [x] 5.2 测试 `--quiet` 参数：`gitlab-tools branch list --quiet`（应只显示项目名）
- [x] 5.3 测试组合使用：`gitlab-tools branch list --quiet --hide-empty`（应只显示有分支的项目名）
- [x] 5.4 验证提交者信息正确显示
- [x] 5.5 验证输出格式清晰易读

