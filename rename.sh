if [ $# -eq 0 ]
then
  echo "Usage: $0 name_of_go_package"
  exit 1
fi

find . -name "*.go" -exec sed -i '' -e "s/htmx-template/$1/g" {} \;
find . -name "*.templ" -exec sed -i '' -e "s/htmx-template/$1/g" {} \;
sed -i '' -e "s/htmx-template/$1/g" go.mod
