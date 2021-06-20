#!/bin/zsh

dirs=(
  "/Users/Guest"
  "/Users/Shared"
  "$HOME/Library"
  "$HOME/Music"
  "$HOME/Dropbox"
  "$HOME/Pictures"
  "$HOME/Applications"
  "$HOME/.Trash"
  "$HOME/.npm"
  "$HOME/.translation"
  "$HOME/.dart"
  "$HOME/.dropbox"
  "$HOME/.sonarlint"
  "$HOME/.vscode"
  "$HOME/.cache"
  "$HOME/.jupyter"
  "$HOME/.flutter-devtools"
  "$HOME/.bash_sessions"
  "$HOME/.zsh_sessions"
  "$HOME/.dlv"
  "$HOME/.emacs.d"
  "$HOME/.oracle_jre_usage"
  "$HOME/.pub-cache"
  "$HOME/.dartServer"
  "*/.git"
  "*/.idea"
)

files=(
  "$HOME/.bash_history"
  "$HOME/.CFUserTextEncoding"
  "$HOME/.DS_Store"
  "$HOME/.flutter"
  "$HOME/.gitconfig"
  "$HOME/.lesshst"
  "$HOME/.profile"
  "$HOME/.python_history"
  "$HOME/.viminfo"
  "$HOME/.wget-hsts"
  "$HOME/.zcompdump"
  "$HOME/.zsh_history"
)

findCmd="find /Users/* -type d \("

doubleIndex=0
for i in "${dirs[@]}"; do
  noPath="-path \"$i\""
  noPath="$noPath -o -path \"$i/*\""
  findCmd="$findCmd $noPath"
  ((doubleIndex++)) || true
  if [[ ${#dirs[@]} != "$doubleIndex" ]]; then
    findCmd="$findCmd -o"
  fi
done

cmd() {
  findCmd="$findCmd \) -prune -o -print $1"
}
cmd ""
#cmd "-maxdepth 1"

for f in "${files[@]}"; do
  findCmd="$findCmd | grep -v \"$f\""
done

#echo "$findCmd"

alias cf='echo $findCmd'
