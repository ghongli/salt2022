#!/usr/bin/env bash

set -euo pipefail

echo -en "
========================================
coding statistics by git log:
  first ls-files to get total lines;
  second find the commit auth name;
  third output everyone commit line num.
"

total_lines=$(git ls-files | xargs cat | wc -l)
echo -e "\nTotal lines: ${total_lines}"
# shellcheck disable=SC2162
git log --format='%aN' | sort -u | while read name
do
  echo -en "$name\t"
  git log --author="$name" --pretty="tformat:" --numstat | \
  awk '{ add += $1; subs += $2; loc += $1 - $2 } END {printf "added lines %s, removed lines %s, total lines %s \n", add, subs, loc}'
done

echo -en "========================================\n\n"

modified=$(git status --porcelain "$@")
if [[ -n "$modified" ]]; then
  git --no-pager diff HEAD "$@"
  untracked=$(echo "${modified}" | grep '???')
  if [[ -n "$untracked" ]]; then
    echo -e "\n\nUNTRACKED FILES:"
    echo "${untracked}" | awk '{print "+++ " $2}'
  fi

  echo -e "\nerror: commit changes to the generated files above"
  exit 1
fi