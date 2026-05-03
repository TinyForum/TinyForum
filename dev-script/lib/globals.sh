# Global variables
OS=$(uname -s)
CURRENT_USER=$(whoami)
DEFAULT_DB_USER=""
DB_USER=""
PACKAGE_MANAGER=""
NEW_DB_USER=""
NEW_DB_PASS=""

LOCAL_IP=$(hostname -I | awk '{print $1}')
if [ -z "$LOCAL_IP" ]; then
    LOCAL_IP=$(ip route get 1 | awk '{print $7; exit}')
fi
BAXKEND_HOST="http://$LOCAL_IP:8080"
FRONTEND_HOST="http://$LOCAL_IP:3000"