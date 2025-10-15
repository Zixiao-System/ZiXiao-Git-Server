# API Test Script for Windows
$ErrorActionPreference = "Stop"

$BASE_URL = "http://localhost:8080/api/v1"
$TOKEN = ""

Write-Host "======================================"
Write-Host "ZiXiao Git Server - API Test"
Write-Host "======================================"
Write-Host ""

function Invoke-API {
    param(
        [string]$Method,
        [string]$Uri,
        [string]$Body = $null,
        [hashtable]$Headers = @{}
    )

    try {
        $params = @{
            Method = $Method
            Uri = $Uri
            Headers = $Headers
            ContentType = "application/json"
        }

        if ($Body) {
            $params.Body = $Body
        }

        $response = Invoke-RestMethod @params
        return $response
    } catch {
        Write-Host "Error: $_" -ForegroundColor Red
        throw
    }
}

# Test 1: Register a user
Write-Host "[1/7] Testing user registration..."
try {
    $body = @{
        username = "testuser"
        email = "test@example.com"
        password = "testpass123"
        full_name = "Test User"
    } | ConvertTo-Json

    $response = Invoke-API -Method POST -Uri "$BASE_URL/auth/register" -Body $body
    $TOKEN = $response.token
    Write-Host "√ User registered successfully" -ForegroundColor Green
} catch {
    Write-Host "× Registration failed" -ForegroundColor Red
    exit 1
}

# Test 2: Login
Write-Host "[2/7] Testing user login..."
try {
    $body = @{
        username = "testuser"
        password = "testpass123"
    } | ConvertTo-Json

    $response = Invoke-API -Method POST -Uri "$BASE_URL/auth/login" -Body $body
    Write-Host "√ Login successful" -ForegroundColor Green
} catch {
    Write-Host "× Login failed" -ForegroundColor Red
    exit 1
}

# Test 3: Get current user
Write-Host "[3/7] Testing get current user..."
try {
    $headers = @{ Authorization = "Bearer $TOKEN" }
    $response = Invoke-API -Method GET -Uri "$BASE_URL/user" -Headers $headers
    Write-Host "√ Get user successful" -ForegroundColor Green
} catch {
    Write-Host "× Get user failed" -ForegroundColor Red
    exit 1
}

# Test 4: Create repository
Write-Host "[4/7] Testing repository creation..."
try {
    $body = @{
        name = "test-repo"
        description = "Test repository"
        is_private = $false
    } | ConvertTo-Json

    $headers = @{ Authorization = "Bearer $TOKEN" }
    $response = Invoke-API -Method POST -Uri "$BASE_URL/repos" -Body $body -Headers $headers
    Write-Host "√ Repository created" -ForegroundColor Green
} catch {
    Write-Host "× Repository creation failed" -ForegroundColor Red
    exit 1
}

# Test 5: Get repository
Write-Host "[5/7] Testing get repository..."
try {
    $response = Invoke-API -Method GET -Uri "$BASE_URL/repos/testuser/test-repo"
    Write-Host "√ Get repository successful" -ForegroundColor Green
} catch {
    Write-Host "× Get repository failed" -ForegroundColor Red
    exit 1
}

# Test 6: List repositories
Write-Host "[6/7] Testing list repositories..."
try {
    $response = Invoke-API -Method GET -Uri "$BASE_URL/users/testuser/repos"
    Write-Host "√ List repositories successful" -ForegroundColor Green
} catch {
    Write-Host "× List repositories failed" -ForegroundColor Red
    exit 1
}

# Test 7: Delete repository
Write-Host "[7/7] Testing repository deletion..."
try {
    $headers = @{ Authorization = "Bearer $TOKEN" }
    $response = Invoke-API -Method DELETE -Uri "$BASE_URL/repos/testuser/test-repo" -Headers $headers
    Write-Host "√ Repository deleted" -ForegroundColor Green
} catch {
    Write-Host "× Repository deletion failed" -ForegroundColor Red
    exit 1
}

Write-Host ""
Write-Host "======================================"
Write-Host "All API tests passed!"
Write-Host "======================================"
Write-Host ""
