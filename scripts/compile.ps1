

###########################
# Folder Variables
$scriptFolder = Split-Path -Path $MyInvocation.MyCommand.Definition -Parent
$projectRoot = Split-Path $scriptFolder


###########################
# Build Information
$buildInfoJson = Join-Path $scriptFolder "builddata.json"
$buildData = (Get-Content $buildInfoJson -raw | ConvertFrom-Json)
$buildData.Version.Build += 1
$buildData.TimeStamp = [DateTime]::Now.ToString("yyyy-MM-ddTHH:mm:ss.fffffffzzz")
$buildStr = $buildData.Version.Major.ToString() + "." + $buildData.Version.Minor.ToString() + "." +  $buildData.Version.Patch.ToString() + "." +  $buildData.Version.Build.ToString()
$buildData | ConvertTo-Json -Compress | Set-Content $buildInfoJson -Encoding Ascii

###########################
# Full program compliation
Write-Progress -Activity "Program Compliation"
$commit = git rev-parse --short HEAD
$buildDate = $buildData.TimeStamp
$ldStr = @"
-X main.commit='$commit' -X main.version='$buildStr' -X main.builddate='$buildDate'
"@
Get-ChildItem -Path (Join-Path $projectRoot "cmd") -Directory | ForEach-Object {
    Push-Location -Path $_.FullName
    Write-Progress -Activity "Program Compliation" -Status $_
    go build -ldflags "$ldStr"
    Write-Progress -Activity "Program Compliation" -Status $_ -Completed
    if ($LASTEXITCODE -ne 0) {
        exit 123453
    }
    Pop-Location
    
}

Write-Progress -Activity "Program Compliation" -Completed
