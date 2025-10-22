param([string]$Version = "latest")
$Repo="yourname/litepost"
$OS="windows"
$ARCH=if ([Environment]::Is64BitOperatingSystem) { "amd64" } else { "amd64" } # adjust if you ship 386
if ($Version -eq "latest") {
  $Tag = (Invoke-RestMethod "https://api.github.com/repos/$Repo/releases/latest").tag_name
} else { $Tag = $Version }
$File = "litepost-$OS-$ARCH.zip"
$Url = "https://github.com/$Repo/releases/download/$Tag/$File"
$Tmp = New-Item -ItemType Directory -Path ([IO.Path]::GetTempPath() + [IO.Path]::GetRandomFileName())
$Zip = Join-Path $Tmp $File
Invoke-WebRequest -Uri $Url -OutFile $Zip
Expand-Archive -Path $Zip -DestinationPath $Tmp
$Bin = Join-Path $env:LOCALAPPDATA "Litepost"
New-Item -ItemType Directory -Force -Path $Bin | Out-Null
Copy-Item -Force (Join-Path $Tmp "litepost-$OS-$ARCH\litepost.exe") (Join-Path $Bin "litepost.exe")
$Scope = [EnvironmentVariableTarget]::User
$Path = [Environment]::GetEnvironmentVariable("Path", $Scope)
if ($Path -notlike "*$Bin*") {
  [Environment]::SetEnvironmentVariable("Path", "$Path;$Bin", $Scope)
}
Write-Host "Installed to $Bin. Open a new terminal and run: litepost"
