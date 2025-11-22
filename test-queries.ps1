$headers = @{
    "x-api-key" = "test"
    "Content-Type" = "application/json"
}

$body = '{"filters":{},"limit":1}'

for ($i = 1; $i -le 10; $i++) {
    Write-Host "---- Request $i ----" -ForegroundColor Cyan

    $start = Get-Date

    try {
        $response = Invoke-RestMethod -Uri "http://localhost:8080/v1/query" `
            -Method Post `
            -Headers $headers `
            -Body $body

        $duration = (Get-Date) - $start
        Write-Host "Status: 200 OK"
        Write-Host "Time: $($duration.TotalMilliseconds) ms"
        Write-Host "Returned rows: $($response.metadata.returned)"
    }
    catch {
        $duration = (Get-Date) - $start
        Write-Host "❌ Error:" $_.Exception.Message -ForegroundColor Red
        Write-Host "Time: $($duration.TotalMilliseconds) ms"
    }

    Start-Sleep -Milliseconds 300
}
