#!/usr/bin/bash
#本人的，你要是运行不了自己改

##有错立即退出
#set -euo pipefail
 #不直接退出是因为要分裂清理旧文件和备份过程，如果备份了但是没清理不要直接退出

#配置
BACKUP_DIR=$(cd "$(dirname "$0")/../vote/backup" && pwd)
MYSQL_USERNAME=$(grep 'username' ../vote/config/database.go | awk -F '"' '{print $2}')
MYSQL_PASSWORD=$(grep 'password' ../vote/config/database.go | awk -F '"' '{print $2}')
MYSQL_HOST=$(grep 'host' ../vote/config/database.go | awk -F '"' '{print $2}')
MYSQL_PORT=$(grep 'port' ../vote/config/database.go | awk -F '"' '{print $2}')
DB_NAME=$(grep 'Dbname' ../vote/config/database.go | awk -F '"' '{print $2}')
REMAIN_DAY=3 #备份三天

#不存在就建一个
mkdir -p "$BACKUP_DIR"

#带上时间戳方便查
BACKUP_FILE="$BACKUP_DIR/${DB_NAME}_$(date +%Y%m%d_%H%M%S).sql.gz"

echo "start to backup $DB_NAME"
#命令行参数
mysqldump -h "$MYSQL_HOST" -P "$MYSQL_PORT" -u "$MYSQL_USERNAME" -p"$MYSQL_PASSWORD" --database "$DB_NAME" | gzip > "$BACKUP_FILE"
#-f检查文件是否存在，du -b可以以字节为单位显示文件大小，cut -f1能提取第一列

if [ -s "$BACKUP_FILE" ]; then
  echo "success to backup"
else
  echo "fail to backup"
  exit 1
fi

echo "clear up out-date data"
find "$BACKUP_DIR" -name "${DB_NAME}_*.sql.gz" -type f -mtime +"$REMAIN_DAY" -delete