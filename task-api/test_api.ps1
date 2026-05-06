$ErrorActionPreference = "Stop"

Write-Host "1. Testing /health endpoint..."
$health = Invoke-RestMethod -Method Get -Uri "http://localhost:8080/health"
Write-Host "Health Response: $($health | ConvertTo-Json -Depth 10 -Compress)`n"

Write-Host "2. Registering a user..."
$registerBody = @{
    name = "Test User"
    email = "testuser@example.com"
    password = "password123"
} | ConvertTo-Json -Compress

try {
    $register = Invoke-RestMethod -Method Post -Uri "http://localhost:8080/api/v1/auth/register" -ContentType "application/json" -Body $registerBody
    Write-Host "Register Response: $($register | ConvertTo-Json -Depth 10 -Compress)`n"
} catch {
    Write-Host "Register Failed: $_"
    if ($_.Exception.Response) {
        $reader = New-Object System.IO.StreamReader($_.Exception.Response.GetResponseStream())
        $reader.BaseStream.Position = 0
        $reader.DiscardBufferedData()
        $responseBody = $reader.ReadToEnd()
        Write-Host "Response Body: $responseBody`n"
    }
}

Write-Host "3. Logging in..."
$loginBody = @{
    email = "testuser@example.com"
    password = "password123"
} | ConvertTo-Json -Compress

try {
    $login = Invoke-RestMethod -Method Post -Uri "http://localhost:8080/api/v1/auth/login" -ContentType "application/json" -Body $loginBody
    Write-Host "Login Response: $($login | ConvertTo-Json -Depth 10 -Compress)`n"
    
    $token = $login.data.access_token
    
    Write-Host "4. Creating a project..."
    $projectBody = @{
        name = "My Test Project"
        description = "This is a test project"
    } | ConvertTo-Json -Compress
    
    $headers = @{
        Authorization = "Bearer $token"
    }
    
    $project = Invoke-RestMethod -Method Post -Uri "http://localhost:8080/api/v1/projects" -Headers $headers -ContentType "application/json" -Body $projectBody
    Write-Host "Project Created: $($project | ConvertTo-Json -Depth 10 -Compress)`n"
    
    $projectId = $project.data.id
    
    Write-Host "5. Creating a task in the project..."
    $taskBody = @{
        title = "First Task"
        description = "This is the first task"
        status = "todo"
        priority = "high"
    } | ConvertTo-Json -Compress
    
    $task = Invoke-RestMethod -Method Post -Uri "http://localhost:8080/api/v1/projects/$projectId/tasks" -Headers $headers -ContentType "application/json" -Body $taskBody
    Write-Host "Task Created: $($task | ConvertTo-Json -Depth 10 -Compress)`n"
    
    Write-Host "6. Getting projects..."
    $projects = Invoke-RestMethod -Method Get -Uri "http://localhost:8080/api/v1/projects" -Headers $headers
    Write-Host "Projects: $($projects | ConvertTo-Json -Depth 10 -Compress)`n"
} catch {
    Write-Host "Operation Failed: $_"
    if ($_.Exception.Response) {
        $reader = New-Object System.IO.StreamReader($_.Exception.Response.GetResponseStream())
        $reader.BaseStream.Position = 0
        $reader.DiscardBufferedData()
        $responseBody = $reader.ReadToEnd()
        Write-Host "Response Body: $responseBody"
    }
}
