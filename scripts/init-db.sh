set -e

psql -v ON_ERROR_STOP=1 --username "$POSTGRES_USER" <<-EOSQL
    CREATE DATABASE users;
    CREATE DATABASE events;
    CREATE DATABASE bookings;
    CREATE DATABASE notifications;
EOSQL

echo "Databases created: users, events, bookings, notifications"