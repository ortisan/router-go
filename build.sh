#https://stackoverflow.com/questions/2870992/automatic-exit-from-bash-shell-script-on-error
abort()
{
    echo >&2 '
***************
*** ABORTED ***
***************
'
    echo "An error occurred. Exiting..." >&2
    exit 1
}

trap 'abort' 0

set -e

docker-compose down

cd dummy-app
mvn clean compile package
cd -
docker-compose up --build


trap : 0

echo >&2 '
************
*** DONE ***
************
'
