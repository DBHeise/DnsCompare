$scriptFolder = Split-Path -Path $MyInvocation.MyCommand.Definition -Parent
$projectRoot = Split-Path $scriptFolder

$workFolder = Join-Path $projectRoot "cmd\DnsCompare"
Start-Process -FilePath (Join-Path $workFolder 'DnsCompare.exe') -WorkingDirectory $workFolder
