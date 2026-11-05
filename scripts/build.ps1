param([string]$Go='go',[string]$Output='I:/Projects/tooling/my-news-20261105/build',[switch]$Candidate)
$ErrorActionPreference='Stop'
$root=Split-Path $PSScriptRoot -Parent
function Check($label){if($LASTEXITCODE){throw $label}}
$commit='uncommitted local candidate'
$safe='safe.directory='+$root.Replace('\','/')
if(-not $Candidate){$dirty=& git -c $safe -C $root status --porcelain;if($LASTEXITCODE -or $dirty){throw 'Release builds require clean committed source'};$commit=(& git -c $safe -C $root rev-parse HEAD).Trim();if($commit -notmatch '^[a-f0-9]{40}$'){throw 'Missing source commit'}}
if(Test-Path -LiteralPath $Output){throw 'Use a fresh output directory; packages are immutable'}
New-Item -ItemType Directory -Path $Output | Out-Null
Push-Location $root
try{
 & npm.cmd --prefix web run build;Check 'Frontend build failed'
 if(-not $Candidate){$dirty=& git -c $safe -C $root status --porcelain;if($LASTEXITCODE -or $dirty){throw 'Generation changed release source; review and commit before packaging'}}
 & $Go test -buildvcs=false ./...;Check 'Backend tests failed'
 & $Go vet ./...;Check 'Go vet failed'
 $stage=Join-Path $Output 'package';New-Item -ItemType Directory -Path $stage | Out-Null
 & $Go build -buildvcs=false -trimpath -ldflags '-s -w' -o (Join-Path $stage 'MyNews.exe') ./cmd/my-news;Check 'Windows build failed'
 try{$env:GOOS='linux';$env:GOARCH='amd64';$env:CGO_ENABLED='0';& $Go build -buildvcs=false -trimpath -ldflags '-s -w' -o (Join-Path $Output 'my-news-linux-amd64') ./cmd/my-news;Check 'Linux compile failed'}finally{Remove-Item Env:GOOS,Env:GOARCH,Env:CGO_ENABLED -ErrorAction SilentlyContinue}
 foreach($name in @('LICENSE','NOTICE','README.md','compatibility.json')){Copy-Item -LiteralPath (Join-Path $root $name) -Destination $stage}
 New-Item -ItemType Directory -Path (Join-Path $stage 'docs'),(Join-Path $stage 'web/public')|Out-Null
 foreach($name in @('SECURITY.md','RESEARCH.md','ACCEPTANCE.md')){Copy-Item -LiteralPath (Join-Path $root ('docs/'+$name)) -Destination (Join-Path $stage 'docs')}
 Copy-Item -LiteralPath (Join-Path $root 'web/public/news.svg') -Destination (Join-Path $stage 'web/public')
 $notice=[Text.StringBuilder]::new()
 $mods=& $Go list -buildvcs=false -deps -f '{{if .Module}}{{.Module.Path}}|{{.Module.Version}}|{{.Module.Dir}}{{end}}' ./cmd/my-news|Where-Object {$_}|Sort-Object -Unique;Check 'Dependency inventory failed'
 foreach($line in $mods){$v=$line.Split('|');if($v[0] -eq 'github.com/oneworkspacestudio/my-news'){continue};$licenses=@(Get-ChildItem -LiteralPath $v[2] -File|Where-Object {$_.Name -match '^(LICENSE|LICENCE|COPYING|NOTICE)(\..*)?$'});if(-not $licenses.Count){throw ('Missing licence '+$v[0])};[void]$notice.AppendLine(('===== '+$v[0]+' '+$v[1]+' ====='));foreach($f in $licenses){[void]$notice.AppendLine([IO.File]::ReadAllText($f.FullName))}}
 foreach($name in @('react','react-dom','scheduler','lucide-react','@codedia/parser')){$dir=Join-Path $root ('web/node_modules/'+$name);$pkg=Get-Content -LiteralPath (Join-Path $dir 'package.json') -Raw|ConvertFrom-Json;$licenses=@(Get-ChildItem -LiteralPath $dir -File|Where-Object {$_.Name -match '^(LICENSE|LICENCE)(\..*)?$'});if(-not $licenses.Count){throw ('Missing frontend licence '+$name)};[void]$notice.AppendLine(('===== '+$name+' '+$pkg.version+' ====='));foreach($f in $licenses){[void]$notice.AppendLine([IO.File]::ReadAllText($f.FullName))}}
 [void]$notice.AppendLine('===== Go runtime =====').AppendLine([IO.File]::ReadAllText((Join-Path (& $Go env GOROOT).Trim() 'LICENSE')))
 [IO.File]::WriteAllText((Join-Path $stage 'THIRD-PARTY-LICENSES.txt'),$notice.ToString())
 $sha=(Get-FileHash -LiteralPath (Join-Path $stage 'MyNews.exe') -Algorithm SHA256).Hash.ToLowerInvariant()
 [IO.File]::WriteAllText((Join-Path $stage 'manifest.json'),(@{schema=1;version='0.1.0-preview.1';candidate=[bool]$Candidate;sourceCommit=$commit;executableSHA256=$sha;githubActionsUsed=$false}|ConvertTo-Json))
 Add-Type -AssemblyName System.IO.Compression.FileSystem
 $zip=Join-Path $Output 'my-news-windows-amd64.zip';[IO.Compression.ZipFile]::CreateFromDirectory($stage,$zip,[IO.Compression.CompressionLevel]::Optimal,$false)
 [IO.File]::WriteAllText((Join-Path $Output 'receipt.json'),(@{candidate=[bool]$Candidate;archive=$zip;sha256=(Get-FileHash $zip).Hash.ToLowerInvariant();executableSHA256=$sha;linuxCompiled=$true}|ConvertTo-Json))
 Write-Output ('Built My News: '+$zip)
}finally{Pop-Location}
