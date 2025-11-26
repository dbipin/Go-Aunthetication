# Go-Aunthetication


docker exec -it apiserver-db psql -U postgres -d apiserver -c "UPDATE schema_migrations SET version = 3, dirty = false;"

sudo netstat -tulpn | grep :8080

sudo kill -9 77668

docker exec -it apiserver-db psql -U postgres -d apiserver -c "UPDATE users SET password = '\$2a\$10\$XZLa046FX7exGXjWPob0Tu0id6/PNziafmd/vo9HjAbQrKgawqepS' WHERE email = 'admin@gmail.com';"