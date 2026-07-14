[CmdletBinding()]
param(
    [string]$InstallDir,
    [string]$Repository,
    [string]$Ref,
    [string]$DownloadBase,
    [string]$ArtifactDir
)

$ErrorActionPreference = "Stop"
$ProgressPreference = "SilentlyContinue"

function Write-Step {
    param([string]$Message)
    Write-Host "==> $Message" -ForegroundColor Cyan
}

function Write-Success {
    param([string]$Message)
    Write-Host "==> $Message" -ForegroundColor Green
}

function Write-AqryWarning {
    param([string]$Message)
    Write-Host "==> $Message" -ForegroundColor Yellow
}

function Get-Setting {
    param(
        [string]$Value,
        [string]$EnvironmentName,
        [string]$DefaultValue
    )

    if (-not [string]::IsNullOrWhiteSpace($Value)) {
        return $Value
    }

    $environmentValue = [Environment]::GetEnvironmentVariable($EnvironmentName)
    if (-not [string]::IsNullOrWhiteSpace($environmentValue)) {
        return $environmentValue
    }

    return $DefaultValue
}

function Test-PathEntry {
    param(
        [string]$PathValue,
        [string]$Directory
    )

    if ([string]::IsNullOrWhiteSpace($PathValue)) {
        return $false
    }

    [char[]]$directorySeparators = @(
        [System.IO.Path]::DirectorySeparatorChar,
        [System.IO.Path]::AltDirectorySeparatorChar
    )
    $normalizedDirectory = $Directory.Trim().TrimEnd($directorySeparators)

    foreach ($entry in ($PathValue -split ";")) {
        $normalizedEntry = $entry.Trim().TrimEnd($directorySeparators)
        if ([string]::Equals($normalizedEntry, $normalizedDirectory, [StringComparison]::OrdinalIgnoreCase)) {
            return $true
        }
    }

    return $false
}

function Get-ExpectedChecksum {
    param(
        [string]$ChecksumPath,
        [string]$ArchiveName
    )

    foreach ($line in Get-Content -LiteralPath $ChecksumPath) {
        if ($line -match "^([A-Fa-f0-9]{64})\s+\*?(.+)$") {
            if ($Matches[2].Trim() -eq $ArchiveName) {
                return $Matches[1].ToLowerInvariant()
            }
        }
    }

    throw "No checksum was published for $ArchiveName"
}

if ($env:OS -ne "Windows_NT") {
    throw "This installer supports Windows only. Use scripts/install.sh on Linux or macOS."
}

$detectedArchitecture = if (-not [string]::IsNullOrWhiteSpace($env:PROCESSOR_ARCHITEW6432)) {
    $env:PROCESSOR_ARCHITEW6432
} else {
    $env:PROCESSOR_ARCHITECTURE
}

if ($detectedArchitecture -ne "AMD64") {
    throw "Unsupported Windows architecture: $detectedArchitecture (supported: amd64)"
}

$Repository = Get-Setting $Repository "AQRY_REPOSITORY" "Satbir6/Aqry"
$Ref = Get-Setting $Ref "AQRY_REF" "main"
$DownloadBase = Get-Setting $DownloadBase "AQRY_DOWNLOAD_BASE" "https://raw.githubusercontent.com/$Repository/$Ref/dist"
$ArtifactDir = Get-Setting $ArtifactDir "AQRY_ARTIFACT_DIR" ""
$InstallDir = Get-Setting $InstallDir "AQRY_INSTALL_DIR" ""

if ([string]::IsNullOrWhiteSpace($InstallDir)) {
    $localAppData = [Environment]::GetFolderPath([System.Environment+SpecialFolder]::LocalApplicationData)
    if ([string]::IsNullOrWhiteSpace($localAppData)) {
        throw "Could not determine the current user's local application data directory"
    }
    $InstallDir = Join-Path $localAppData "Programs\aqry"
}

$InstallDir = [System.IO.Path]::GetFullPath([Environment]::ExpandEnvironmentVariables($InstallDir))
$archiveName = "aqry_Windows_x86_64.zip"
$tempDir = Join-Path ([System.IO.Path]::GetTempPath()) ("aqry-" + [Guid]::NewGuid().ToString("N"))
$archivePath = Join-Path $tempDir $archiveName
$checksumPath = Join-Path $tempDir "SHA256SUMS"
$extractDir = Join-Path $tempDir "extracted"
$installPath = Join-Path $InstallDir "aqry.exe"
$stagedInstallPath = Join-Path $InstallDir (".aqry-install-{0}.exe" -f $PID)

try {
    Write-Step "Detected Windows/x86_64"
    New-Item -ItemType Directory -Path $tempDir -Force | Out-Null

    if (-not [string]::IsNullOrWhiteSpace($ArtifactDir)) {
        $sourceArchive = Join-Path $ArtifactDir $archiveName
        $sourceChecksum = Join-Path $ArtifactDir "SHA256SUMS"
        if (-not (Test-Path -LiteralPath $sourceArchive -PathType Leaf)) {
            throw "Artifact not found: $sourceArchive"
        }
        if (-not (Test-Path -LiteralPath $sourceChecksum -PathType Leaf)) {
            throw "Checksum file not found: $sourceChecksum"
        }

        Write-Step "Using bundled artifact from $ArtifactDir"
        Copy-Item -LiteralPath $sourceArchive -Destination $archivePath
        Copy-Item -LiteralPath $sourceChecksum -Destination $checksumPath
    } else {
        Write-Step "Downloading $archiveName from repository branch $Ref"
        [Net.ServicePointManager]::SecurityProtocol = [Net.ServicePointManager]::SecurityProtocol -bor [Net.SecurityProtocolType]::Tls12
        Invoke-WebRequest -UseBasicParsing -Uri "$DownloadBase/$archiveName" -OutFile $archivePath
        Invoke-WebRequest -UseBasicParsing -Uri "$DownloadBase/SHA256SUMS" -OutFile $checksumPath
    }

    $expectedChecksum = Get-ExpectedChecksum $checksumPath $archiveName
    $actualChecksum = (Get-FileHash -LiteralPath $archivePath -Algorithm SHA256).Hash.ToLowerInvariant()
    if ($actualChecksum -ne $expectedChecksum) {
        throw "Checksum verification failed for $archiveName"
    }
    Write-Success "Checksum verified"

    Write-Step "Extracting aqry.exe"
    Expand-Archive -LiteralPath $archivePath -DestinationPath $extractDir -Force
    $binaryPath = Join-Path $extractDir "aqry.exe"
    if (-not (Test-Path -LiteralPath $binaryPath -PathType Leaf)) {
        throw "The downloaded archive does not contain aqry.exe"
    }

    New-Item -ItemType Directory -Path $InstallDir -Force | Out-Null
    Write-Step "Installing aqry to $installPath"
    Copy-Item -LiteralPath $binaryPath -Destination $stagedInstallPath -Force
    Copy-Item -LiteralPath $stagedInstallPath -Destination $installPath -Force
    Remove-Item -LiteralPath $stagedInstallPath -Force

    Write-Step "Verifying installation"
    $versionOutput = & $installPath --version 2>&1
    if ($LASTEXITCODE -ne 0) {
        throw "Installed binary could not be executed: $versionOutput"
    }
    Write-Success ($versionOutput -join [Environment]::NewLine)
    Write-Success "Installed at $installPath"

    $userPath = [Environment]::GetEnvironmentVariable(
        "Path",
        [System.EnvironmentVariableTarget]::User
    )
    if (-not (Test-PathEntry $userPath $InstallDir)) {
        $newUserPath = if ([string]::IsNullOrWhiteSpace($userPath)) {
            $InstallDir
        } else {
            "$userPath;$InstallDir"
        }
        [Environment]::SetEnvironmentVariable(
            "Path",
            $newUserPath,
            [System.EnvironmentVariableTarget]::User
        )
        Write-Success "Added $InstallDir to your user PATH"
        Write-AqryWarning "Open a new terminal to use the updated PATH everywhere"
    }

    if (-not (Test-PathEntry $env:Path $InstallDir)) {
        $env:Path = "$InstallDir;$env:Path"
    }

    Write-Host ""
    Write-Host "Try aqry now:"
    Write-Host ""
    Write-Host "  aqry example.com"
    Write-Host "  aqry"
    Write-Host ""
} catch {
    Write-Host "error: $($_.Exception.Message)" -ForegroundColor Red
    throw
} finally {
    if (Test-Path -LiteralPath $stagedInstallPath -PathType Leaf) {
        Remove-Item -LiteralPath $stagedInstallPath -Force -ErrorAction SilentlyContinue
    }
    if (Test-Path -LiteralPath $tempDir -PathType Container) {
        Remove-Item -LiteralPath $tempDir -Recurse -Force -ErrorAction SilentlyContinue
    }
}
