# Reproducible terminal demo for the aqry README.
# Generate it from the repository root with: scripts/generate-demo.sh

Output assets/aqry-demo.gif

Set Shell "bash"
Set FontSize 18
Set Width 1200
Set Height 680
Set Padding 16
Set Framerate 24
Set PlaybackSpeed 1.25
Set TypingSpeed 45ms
Set WaitTimeout 15s
Set Theme "Catppuccin Mocha"

# Use the binary built by scripts/generate-demo.sh and keep setup off-camera.
Hide
Type `export PATH="$PWD/bin:$PATH"`
Enter
Type `export PS1='$ '`
Enter
Type "unset NO_COLOR"
Enter
Type "clear"
Enter
Show

Sleep 800ms

# Start with aqry's script-friendly output.
Type "aqry example.com"
Sleep 350ms
Enter
Wait+Screen /([0-9]{1,3}\.){3}[0-9]{1,3}/
Sleep 1600ms

# Open the interactive DNS workspace.
Ctrl+L
Sleep 500ms
Type "aqry"
Sleep 300ms
Enter
Wait+Screen /DNS Resolver/
Sleep 800ms

# Pick an MX record and Cloudflare, then run the lookup.
Type "gmail.com"
Sleep 500ms
Tab
Sleep 400ms
Enter
Wait+Screen /Select Record Type/
Down@300ms 3
Sleep 500ms
Enter
Sleep 500ms
Tab
Sleep 400ms
Enter
Wait+Screen /DNS Resolver/
Down
Sleep 500ms
Enter
Sleep 500ms
Tab
Sleep 300ms
Enter
Wait+Screen /Resolved successfully/
Sleep 1400ms

# Expand the complete MX result set and finish on the useful view.
Type "a"
Sleep 3200ms
