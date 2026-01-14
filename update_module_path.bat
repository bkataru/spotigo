@echo off
echo Updating module paths from github.com/bkataru-workshop/spotigo to github.com/bkataru/spotigo

REM Update all .go files
for /r %%f in (*.go) do (
    echo Updating %%f
    powershell -Command "(Get-Content %%f) -replace 'github\.com/bkataru-workshop/spotigo', 'github.com/bkataru/spotigo' | Set-Content %%f"
)

REM Update all .md files  
for /r %%f in (*.md) do (
    echo Updating %%f
    powershell -Command "(Get-Content %%f) -replace 'github\.com/bkataru-workshop/spotigo', 'github.com/bkataru/spotigo' | Set-Content %%f"
)

REM Update all .yml and .yaml files
for /r %%f in (*.yml *.yaml) do (
    echo Updating %%f
    powershell -Command "(Get-Content %%f) -replace 'github\.com/bkataru-workshop/spotigo', 'github.com/bkataru/spotigo' | Set-Content %%f"
)

echo Done!