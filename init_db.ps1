# PowerShell script to initialize the database
# You may need to adjust the MySQL path based on your installation

$mysqlPath = "C:\Program Files\MySQL\MySQL Server 8.0\bin\mysql.exe"
$sqlFile = "d:\Workspace\project\app\power-admin\power-admin-server\db\init.sql"

# Check if mysql exists
if (Test-Path $mysqlPath) {
    Write-Host "Initializing database..."
    & $mysqlPath -u root @($sqlFile)
    Write-Host "Database initialization completed!"
} else {
    Write-Host "MySQL not found at $mysqlPath"
    Write-Host "Please ensure MySQL is installed and adjust the path accordingly"
}
