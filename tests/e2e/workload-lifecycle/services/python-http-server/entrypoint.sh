#!/bin/sh
# entrypoint.sh

# Exit immediately if a command exits with a non-zero status.
set -e

# Run Django database migrations
echo "Running migrations..."
python manage.py makemigrations
python manage.py migrate

# Start the Django development server
echo "Starting server..."
exec python manage.py runserver 0.0.0.0:8000
