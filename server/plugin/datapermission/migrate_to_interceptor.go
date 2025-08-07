package datapermission

import (
	"fmt"
	"go/ast"
	"go/format"
	"go/parser"
	"go/token"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

// MigrationHelper 迁移助手
type MigrationHelper struct {
	fileSet *token.FileSet
}

// NewMigrationHelper 创建迁移助手
func NewMigrationHelper() *MigrationHelper {
	return &MigrationHelper{
		fileSet: token.NewFileSet(),
	}
}

// MigrateDirectory 迁移整个目录
func (m *MigrationHelper) MigrateDirectory(dirPath string) error {
	return filepath.Walk(dirPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !strings.HasSuffix(path, ".go") {
			return nil
		}

		return m.MigrateFile(path)
	})
}

// MigrateFile 迁移单个文件
func (m *MigrationHelper) MigrateFile(filePath string) error {
	content, err := ioutil.ReadFile(filePath)
	if err != nil {
		return err
	}

	// 解析Go文件
	node, err := parser.ParseFile(m.fileSet, filePath, content, parser.ParseComments)
	if err != nil {
		return err
	}

	// 检查是否需要迁移
	if !m.needsMigration(node) {
		return nil
	}

	// 执行迁移
	m.migrateAST(node)

	// 生成新的代码
	var buf strings.Builder
	if err := format.Node(&buf, m.fileSet, node); err != nil {
		return err
	}

	// 写回文件
	return ioutil.WriteFile(filePath, []byte(buf.String()), 0644)
}

// needsMigration 检查是否需要迁移
func (m *MigrationHelper) needsMigration(node *ast.File) bool {
	for _, decl := range node.Decls {
		if funcDecl, ok := decl.(*ast.FuncDecl); ok {
			if m.containsOldPermissionCalls(funcDecl) {
				return true
			}
		}
	}
	return false
}

// containsOldPermissionCalls 检查是否包含旧的权限调用
func (m *MigrationHelper) containsOldPermissionCalls(funcDecl *ast.FuncDecl) bool {
	if funcDecl.Body == nil {
		return false
	}

	for _, stmt := range funcDecl.Body.List {
		if m.isOldPermissionCall(stmt) {
			return true
		}
	}
	return false
}

// isOldPermissionCall 检查是否是旧的权限调用
func (m *MigrationHelper) isOldPermissionCall(stmt ast.Stmt) bool {
	// 这里可以实现具体的AST分析逻辑
	// 检查是否包含 ApplyDataPermissionToQuery 等调用
	return false
}

// migrateAST 迁移AST
func (m *MigrationHelper) migrateAST(node *ast.File) {
	// 这里可以实现具体的AST转换逻辑
	// 将旧的权限调用替换为新的拦截器调用
}

// GenerateMigrationReport 生成迁移报告
func (m *MigrationHelper) GenerateMigrationReport(dirPath string) (string, error) {
	report := &strings.Builder{}
	report.WriteString("# 数据权限迁移报告\n\n")

	var totalFiles, migratedFiles int
	var issues []string

	err := filepath.Walk(dirPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !strings.HasSuffix(path, ".go") {
			return nil
		}

		totalFiles++

		content, err := ioutil.ReadFile(path)
		if err != nil {
			issues = append(issues, fmt.Sprintf("无法读取文件 %s: %v", path, err))
			return nil
		}

		contentStr := string(content)

		// 检查是否包含旧的权限调用
		if strings.Contains(contentStr, "ApplyDataPermissionToQuery") ||
			strings.Contains(contentStr, "DataPermissionGinMiddleware") ||
			strings.Contains(contentStr, "FilterFieldsByPermission") {

			migratedFiles++
			report.WriteString(fmt.Sprintf("## 需要迁移的文件: %s\n\n", path))

			// 分析具体的调用
			if strings.Contains(contentStr, "ApplyDataPermissionToQuery") {
				report.WriteString("- 包含 `ApplyDataPermissionToQuery` 调用\n")
			}
			if strings.Contains(contentStr, "DataPermissionGinMiddleware") {
				report.WriteString("- 包含 `DataPermissionGinMiddleware` 中间件\n")
			}
			if strings.Contains(contentStr, "FilterFieldsByPermission") {
				report.WriteString("- 包含 `FilterFieldsByPermission` 调用\n")
			}

			report.WriteString("\n### 建议的迁移步骤:\n\n")
			report.WriteString("1. 将 `DataPermissionGinMiddleware()` 替换为 `DataPermissionInterceptorMiddleware()`\n")
			report.WriteString("2. 将 `utils.PermissionHelperApp.ApplyDataPermissionToQuery(db, c, \"table\")` 替换为 `utils.EnhancedPermissionHelperApp.GetDBWithTable(c, \"table\")`\n")
			report.WriteString("3. 将 `utils.PermissionHelperApp.FilterFieldsByPermission(c, data, \"table\", \"view\")` 替换为 `utils.EnhancedPermissionHelperApp.FilterResponseFields(c, data, \"table\")`\n")
			report.WriteString("\n")
		}

		return nil
	})

	if err != nil {
		return "", err
	}

	report.WriteString(fmt.Sprintf("## 统计信息\n\n"))
	report.WriteString(fmt.Sprintf("- 总文件数: %d\n", totalFiles))
	report.WriteString(fmt.Sprintf("- 需要迁移的文件数: %d\n", migratedFiles))
	report.WriteString(fmt.Sprintf("- 迁移进度: %.2f%%\n\n", float64(migratedFiles)/float64(totalFiles)*100))

	if len(issues) > 0 {
		report.WriteString("## 问题列表\n\n")
		for _, issue := range issues {
			report.WriteString(fmt.Sprintf("- %s\n", issue))
		}
		report.WriteString("\n")
	}

	report.WriteString("## 迁移建议\n\n")
	report.WriteString("1. **备份代码**: 在开始迁移前，请确保代码已提交到版本控制系统\n")
	report.WriteString("2. **逐步迁移**: 建议按模块逐步迁移，而不是一次性迁移所有文件\n")
	report.WriteString("3. **测试验证**: 每次迁移后都要进行充分的测试\n")
	report.WriteString("4. **保持兼容**: 新旧API可以并存，可以逐步替换\n")
	report.WriteString("5. **参考文档**: 详细使用方法请参考 [INTERCEPTOR_GUIDE.md](./INTERCEPTOR_GUIDE.md)\n")

	return report.String(), nil
}

// CreateMigrationScript 创建迁移脚本
func (m *MigrationHelper) CreateMigrationScript(outputPath string) error {
	script := `#!/bin/bash

# 数据权限迁移脚本
# 此脚本帮助您从传统模式迁移到拦截器模式

echo "开始数据权限迁移..."

# 1. 备份当前代码
echo "1. 备份当前代码..."
git add .
git commit -m "迁移前备份: 数据权限传统模式"

# 2. 替换中间件调用
echo "2. 替换中间件调用..."
find . -name "*.go" -type f -exec sed -i 's/DataPermissionGinMiddleware()/DataPermissionInterceptorMiddleware()/g' {} +
find . -name "*.go" -type f -exec sed -i 's/DataPermissionGinMiddleware/DataPermissionAutoMiddleware/g' {} +

# 3. 替换权限助手调用
echo "3. 替换权限助手调用..."
find . -name "*.go" -type f -exec sed -i 's/utils\.PermissionHelperApp/utils.EnhancedPermissionHelperApp/g' {} +

# 4. 替换常用方法调用
echo "4. 替换常用方法调用..."
# 这里可以添加更多的替换规则

# 5. 添加导入
echo "5. 检查导入..."
echo "请手动检查并添加必要的导入语句"

# 6. 提交更改
echo "6. 提交更改..."
git add .
git commit -m "迁移到数据权限拦截器模式"

echo "迁移完成！请运行测试确保功能正常。"
echo "详细使用指南请查看: server/plugin/datapermission/INTERCEPTOR_GUIDE.md"
`

	return ioutil.WriteFile(outputPath, []byte(script), 0755)
}

// 使用示例
func ExampleUsage() {
	migrationHelper := NewMigrationHelper()

	// 生成迁移报告
	report, err := migrationHelper.GenerateMigrationReport("./server")
	if err != nil {
		fmt.Printf("生成迁移报告失败: %v\n", err)
		return
	}

	// 保存报告
	ioutil.WriteFile("migration_report.md", []byte(report), 0644)
	fmt.Println("迁移报告已生成: migration_report.md")

	// 创建迁移脚本
	err = migrationHelper.CreateMigrationScript("migrate.sh")
	if err != nil {
		fmt.Printf("创建迁移脚本失败: %v\n", err)
		return
	}

	fmt.Println("迁移脚本已创建: migrate.sh")
	fmt.Println("运行 ./migrate.sh 开始自动迁移")
}
