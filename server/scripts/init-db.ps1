# 衣搭库数据库重新初始化脚本
# 用法: 在 server 目录下执行 .\scripts\init-db.ps1
# 或从项目根: cd server; .\scripts\init-db.ps1

$ErrorActionPreference = "Stop"
$scriptDir = Split-Path -Parent $MyInvocation.MyCommand.Path
$serverDir = Split-Path -Parent $scriptDir
$migrationsDir = Join-Path $serverDir "migrations"

# 从 config.yaml 读取数据库密码（若存在）
$configPath = Join-Path $serverDir "config.yaml"
$mysqlUser = "root"
$mysqlPass = ""

if (Test-Path $configPath) {
    $content = Get-Content $configPath -Raw
    if ($content -match "password:\s*(\S+)") {
        $mysqlPass = $Matches[1].Trim()
    }
}

if (-not $mysqlPass) {
    Write-Host "未找到 config.yaml 或其中无 database.password，将使用无密码连接。" -ForegroundColor Yellow
}

$mysqlArgs = @("-u", $mysqlUser, "--default-character-set=utf8mb4")
if ($mysqlPass) { $mysqlArgs += "-p$mysqlPass" }

function Run-Sql {
    param([string]$Name, [string]$Path)
    $sql = Get-Content $Path -Raw -Encoding UTF8
    $sql | mysql @mysqlArgs 2>&1 | Out-Null
    if ($LASTEXITCODE -ne 0) {
        Write-Error "执行 $Name 失败"
    }
    Write-Host "  OK $Name" -ForegroundColor Green
}

Push-Location $serverDir
try {
    Write-Host "1. 删除数据库..." -ForegroundColor Cyan
    Run-Sql "000_drop" (Join-Path $migrationsDir "000_drop.sql")
    Write-Host "2. 建库建表..." -ForegroundColor Cyan
    Run-Sql "001_init" (Join-Path $migrationsDir "001_init.sql")
    Write-Host "3. 种子数据..." -ForegroundColor Cyan
    Run-Sql "002_seed_data" (Join-Path $migrationsDir "002_seed_data.sql")
    Write-Host "数据库初始化完成。" -ForegroundColor Green
} finally {
    Pop-Location
}
